package routing

import (
	"net/http"

	"task-management-backend/internal/handler/middleware"
	"task-management-backend/internal/handler/rest"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(handler *rest.Handler) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.CORS)
	r.Use(chimiddleware.StripSlashes)

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	r.Route("/api", func(api chi.Router) {
		RegisterUserRoutes(api, handler)
	})

	return r
}
