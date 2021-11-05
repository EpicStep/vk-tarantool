package shorter

import (
	"github.com/EpicStep/vk-tarantool/internal/jsonutil"
	v1 "github.com/EpicStep/vk-tarantool/pkg/api/v1"
	ua "github.com/mileusna/useragent"
	"net/http"
)

func (s *Service) Analytics(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")

	t, err := s.db.GetTransitionByShort(hash)
	if err != nil {
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
