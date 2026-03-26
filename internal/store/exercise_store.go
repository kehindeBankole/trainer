package store

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ExerciseStore interface {
	CreateExercise(ctx context.Context, exercise *Exercise) error
	GetExerciseByID(ctx context.Context, id uuid.UUID) (*Exercise, error)
	GetAllExercises(ctx context.Context) ([]*Exercise, error)
}

type PostgresExerciseStore struct {
	db *pgxpool.Pool
}

func NewPostgresExerciseStore(db *pgxpool.Pool) *PostgresExerciseStore {
	return &PostgresExerciseStore{db: db}
}

func (s *PostgresExerciseStore) CreateExercise(ctx context.Context, exercise *Exercise) error {
	query := `
		INSERT INTO exercises (name, description, muscle_group, category)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	return s.db.QueryRow(ctx, query,
		exercise.Name,
		exercise.Description,
		exercise.MuscleGroup,
		exercise.Category,
	).Scan(&exercise.ID, &exercise.CreatedAt)
}

func (s *PostgresExerciseStore) GetExerciseByID(ctx context.Context, id uuid.UUID) (*Exercise, error) {
	query := `
		SELECT id, name, description, muscle_group, category, created_at
		FROM exercises
		WHERE id = $1`

	exercise := &Exercise{}
	err := s.db.QueryRow(ctx, query, id).Scan(
		&exercise.ID,
		&exercise.Name,
		&exercise.Description,
		&exercise.MuscleGroup,
		&exercise.Category,
		&exercise.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("exercise not found")
		}
		return nil, err
	}
	return exercise, nil
}

func (s *PostgresExerciseStore) GetAllExercises(ctx context.Context) ([]*Exercise, error) {
	query := `
		SELECT id, name, description, muscle_group, category, created_at
		FROM exercises
		ORDER BY name ASC`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []*Exercise
	for rows.Next() {
		exercise := &Exercise{}
		err := rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.Description,
			&exercise.MuscleGroup,
			&exercise.Category,
			&exercise.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}
	return exercises, nil
}
