package shorter

import "github.com/go-chi/chi/v5"

// Routes add new routes to chi Router.
func (s *Service) Routes(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		r.Post("/set", s.Short)
		r.Get("/{hash}", s.Transition)
	})
}