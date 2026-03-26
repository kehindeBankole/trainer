package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"workout-trainer/internal/api"
	"workout-trainer/internal/store"
)

type Application struct {
	Logger          *log.Logger
	DB              *pgxpool.Pool
	UserHandler     *api.UserHandler
	ExerciseHandler *api.ExerciseHandler
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("could not create database pool: %w", err)
	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("could not reach database: %w", err)
	}

	logger.Println("connected to database")

	userStore     := store.NewPostgresUserStore(db)
	exerciseStore := store.NewPostgresExerciseStore(db)

	return &Application{
		Logger:          logger,
		DB:              db,
		UserHandler:     api.NewUserHandler(userStore, logger),
		ExerciseHandler: api.NewExerciseHandler(exerciseStore, logger),
	}, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "app is available")
}
