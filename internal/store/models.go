package store

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrDuplicateEntry = errors.New("duplicate entry")

type User struct {
	ID        uuid.UUID
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
}

type Exercise struct {
	ID          uuid.UUID
	Name        string
	Description string
	MuscleGroup string
	Category    string
	CreatedAt   time.Time
}
