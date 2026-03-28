package api

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"workout-trainer/internal/store"
)

type UserHandler struct {
	store  store.UserStore
	logger *log.Logger
}

func NewUserHandler(store store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{store: store, logger: logger}
}

type registerRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type userResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt string    `json:"created_at"`
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := ReadJSON(r, &req); err != nil {
		ErrorJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := Validate.Struct(req); err != nil {
		ValidationErrorJSON(w, err)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Printf("error hashing password: %v", err)
		ErrorJSON(w, http.StatusInternalServerError, "internal server error")
		return
	}

	user := &store.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hash),
	}

	if err := h.store.CreateUser(r.Context(), user); err != nil {
		if errors.Is(err, store.ErrDuplicateEntry) {
			ErrorJSON(w, http.StatusConflict, "username or email already exists")
			return
		}
		h.logger.Printf("error creating user: %v", err)
		ErrorJSON(w, http.StatusInternalServerError, "could not create user")
		return
	}

	WriteJSON(w, http.StatusCreated, userResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
	})
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.store.GetAllUsers(r.Context())
	if err != nil {
		h.logger.Printf("error fetching users: %v", err)
		ErrorJSON(w, http.StatusInternalServerError, "could not fetch users")
		return
	}

	resp := make([]userResponse, len(users))
	for i, u := range users {
		resp[i] = userResponse{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			CreatedAt: u.CreatedAt.String(),
		}
	}

	WriteJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		ErrorJSON(w, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.store.GetUserByID(r.Context(), id)
	if err != nil {
		ErrorJSON(w, http.StatusNotFound, "user not found")
		return
	}

	WriteJSON(w, http.StatusOK, userResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
	})
}
