package shorter

import (
	"errors"
	"github.com/EpicStep/vk-tarantool/internal/jsonutil"
	"github.com/EpicStep/vk-tarantool/internal/shorter/database"
	v1 "github.com/EpicStep/vk-tarantool/pkg/api/v1"
	ua "github.com/mileusna/useragent"
	"net/http"
)

func (s *Service) Analytics(w http.ResponseWriter, r *http.Request) {
	spanCtx, span := s.tracer.Start(r.Context(), "http.api.analytics")
	defer span.End()

	hash := r.URL.Query().Get("hash")

	if hash == "" {
		span.End()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	t, err := s.db.GetTransitionByShort(spanCtx, hash)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			span.End()
			w.WriteHeader(http.StatusNotFound)
			return
		}

		span.RecordError(err)
		span.End()

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	os := make(map[string]int64)
	browser := make(map[string]int64)

	for _, v := range t {
		parsed := ua.Parse(v.UA)

		os[parsed.OS] += 1
		browser[parsed.Name] += 1
	}

	jsonutil.MarshalResponse(w, http.StatusOK, jsonutil.NewSuccessfulResponse(v1.AnalyticsResponse{
		Views:   len(t),
		OS:      os,
		Browser: browser,
	}))
}
