package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"workout-trainer/internal/api"
	"workout-trainer/internal/auth"
	"workout-trainer/internal/store"
)

type Application struct {
	Logger          *log.Logger
	DB              *pgxpool.Pool
	Authenticator   *auth.JWTAuthenticator
	UserHandler     *api.UserHandler
	ExerciseHandler *api.ExerciseHandler
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is not set")
	}

	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("could not create database pool: %w", err)
	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("could not reach database: %w", err)
	}

	logger.Println("connected to database")

	authenticator := auth.NewJWTAuthenticator(jwtSecret)
	userStore      := store.NewPostgresUserStore(db)
	exerciseStore  := store.NewPostgresExerciseStore(db)

	return &Application{
		Logger:          logger,
		DB:              db,
		Authenticator:   authenticator,
		UserHandler:     api.NewUserHandler(userStore, logger, authenticator),
		ExerciseHandler: api.NewExerciseHandler(exerciseStore, logger),
	}, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "app is available")
}
