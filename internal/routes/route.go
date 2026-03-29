package routes

import (
	"workout-trainer/internal/app"
	"workout-trainer/internal/middleware"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	r.Get("/health", app.HealthCheck)

	r.Route("/api/v1", func(r chi.Router) {

		r.Route("/users", func(r chi.Router) {
			r.Post("/register", app.UserHandler.RegisterUser)
			r.Get("/", app.UserHandler.GetAllUsers)
			r.Post("/login", app.UserHandler.Login)
			r.Get("/{id}", app.UserHandler.GetUser)
		})

		r.Route("/exercises", func(r chi.Router) {
			r.Get("/", app.ExerciseHandler.GetAllExercises)
			r.Get("/{id}", app.ExerciseHandler.GetExercise)

			r.Group(func(r chi.Router) {
				r.Use(middleware.Authenticate(app.Authenticator))
				r.Post("/", app.ExerciseHandler.CreateExercise)
			})
		})

	})

	return r
}
