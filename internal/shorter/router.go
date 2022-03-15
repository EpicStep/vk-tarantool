package shorter

import "github.com/go-chi/chi/v5"

// Routes add new routes to chi Router.
func (s *Service) Routes(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		s.apiV1(r)
	})
}

func (s *Service) apiV1(r chi.Router) {
	r.Route("/v1", func(r chi.Router) {
		r.Post("/set", s.Short)
		r.Get("/analytics", s.Analytics)
		r.Get("/{hash}", s.Transition)
	})
}
