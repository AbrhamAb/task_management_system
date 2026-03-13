package routing

import (
	"task-management-backend/internal/handler/middleware"
	"task-management-backend/internal/handler/rest"

	"github.com/go-chi/chi/v5"
)

func RegisterUserRoutes(r chi.Router, handler *rest.Handler) {
	r.Route("/auth", func(auth chi.Router) {
		auth.Post("/register", handler.Register)
		auth.Post("/login", handler.Login)
		auth.With(middleware.Auth(handler.Module().User)).Get("/me", handler.Me)
	})

	r.Route("/tasks", func(tasks chi.Router) {
		tasks.Use(middleware.Auth(handler.Module().User))
		tasks.Get("/", handler.ListTasks)
		tasks.Post("/", handler.CreateTask)
		tasks.Put("/{taskID}", handler.UpdateTask)
		tasks.Delete("/{taskID}", handler.DeleteTask)
		tasks.Patch("/{taskID}/status", handler.UpdateTaskStatus)
	})

}
