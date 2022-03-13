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
	"log"
	"net/http"
	"strings"
)

type Service struct {
	db      *shortDB.ShorterDB
	metrics *metrics.PrometheusService
}

func New(db *database.DB) *Service {
	return &Service{
		db:      shortDB.New(db),
		metrics: metrics.NewPrometheusService(),
	}
}

func (s *Service) Short(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")

	if url == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hasher := md5.New()

	hasher.Write([]byte(url))
	hexStr := hex.EncodeToString(hasher.Sum(nil))

	ub64 := base64.StdEncoding.EncodeToString([]byte(hexStr))

	ipArr := strings.Split(r.RemoteAddr, ":")

	err := s.db.CreateShort(&model.Short{
		Shorted:   ub64[0:12],
		Original:  url,
		CreatedBy: ipArr[0],
	})

	if err != nil {
		log.Println(err)
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
	hash := chi.URLParam(r, "hash")

	if hash == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	short, err := s.db.GetShortByHash(hash)
	if err != nil {
		if errors.Is(err, shortDB.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ipArr := strings.Split(r.RemoteAddr, ":")

	err = s.db.InsertAnalytics(&model.Transition{
		ID:      uuid.New().String(),
		Shorted: hash,
		IP:      ipArr[0],
		UA:      r.UserAgent(),
	})

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, short[0].Original, http.StatusMovedPermanently)
}
