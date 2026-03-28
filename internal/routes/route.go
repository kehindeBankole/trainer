package routes

import (
	"workout-trainer/internal/app"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", app.HealthCheck)

	r.Route("/api/v1", func(r chi.Router) {

		r.Route("/users", func(r chi.Router) {
			r.Post("/register", app.UserHandler.RegisterUser)
			r.Get("/", app.UserHandler.GetAllUsers)
			r.Get("/{id}", app.UserHandler.GetUser)
		})

		r.Route("/exercises", func(r chi.Router) {
			r.Post("/", app.ExerciseHandler.CreateExercise)
			r.Get("/", app.ExerciseHandler.GetAllExercises)
			r.Get("/{id}", app.ExerciseHandler.GetExercise)
		})

	})

	return r
}
