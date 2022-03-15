package shorter

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"github.com/EpicStep/vk-tarantool/internal/jsonutil"
	shortDB "github.com/EpicStep/vk-tarantool/internal/shorter/database"
	"github.com/EpicStep/vk-tarantool/internal/shorter/model"
	v1 "github.com/EpicStep/vk-tarantool/pkg/api/v1"
	"github.com/EpicStep/vk-tarantool/pkg/database"
	"github.com/EpicStep/vk-tarantool/pkg/metrics"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strings"
)

type Service struct {
	db      *shortDB.ShorterDB
	metrics *metrics.PrometheusService
	tracer  trace.Tracer
}

func New(db *database.DB, tr trace.Tracer) *Service {
	return &Service{
		db:      shortDB.New(db),
		metrics: metrics.NewPrometheusService(),
		tracer:  tr,
	}
}

func (s *Service) Short(w http.ResponseWriter, r *http.Request) {
	spanCtx, span := s.tracer.Start(r.Context(), "http.api.short")
	defer span.End()

	url := r.URL.Query().Get("url")

	if url == "" {
		span.End()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hasher := md5.New()

	hasher.Write([]byte(url))
	hexStr := hex.EncodeToString(hasher.Sum(nil))

	ub64 := base64.StdEncoding.EncodeToString([]byte(hexStr))

	ipArr := strings.Split(r.RemoteAddr, ":")

	err := s.db.CreateShort(spanCtx, &model.Short{
		Shorted:   ub64[0:12],
		Original:  url,
		CreatedBy: ipArr[0],
	})
	if err != nil {
		if errors.Is(err, shortDB.ErrAlreadyExists) {
			span.RecordError(err)
			span.End()
			w.WriteHeader(http.StatusConflict)
			return
		}

		span.RecordError(err)
		span.End()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.metrics.TotalRequests.With(prometheus.Labels{
		"method": r.Method,
		"path":   r.RequestURI,
	}).Inc()

	jsonutil.MarshalResponse(w, http.StatusOK, v1.ShorterResponse{
		Shorted: "http://localhost:8182/" + ub64[0:12],
	})
}

func (s *Service) Transition(w http.ResponseWriter, r *http.Request) {
	spanCtx, span := s.tracer.Start(r.Context(), "http.api.transition")
	defer span.End()

	hash := chi.URLParam(r, "hash")

	if hash == "" {
		span.End()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	short, err := s.db.GetShortByHash(spanCtx, hash)
	if err != nil {
		if errors.Is(err, shortDB.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		span.RecordError(err)
		span.End()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ipArr := strings.Split(r.RemoteAddr, ":")

	err = s.db.InsertAnalytics(spanCtx, &model.Transition{
		ID:      uuid.New().String(),
		Shorted: hash,
		IP:      ipArr[0],
		UA:      r.UserAgent(),
	})

	if err != nil {
		span.RecordError(err)
		span.End()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, short[0].Original, http.StatusMovedPermanently)
}
