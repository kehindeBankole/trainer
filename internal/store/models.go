package store

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrDuplicateEntry = errors.New("duplicate entry")

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type Exercise struct {
	ID          uuid.UUID
	Name        string
	Description string
	MuscleGroup string
	Category    string
	CreatedAt   time.Time
}
