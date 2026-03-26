package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"workout-trainer/internal/store"
)

type ExerciseHandler struct {
	store store.ExerciseStore
	log   *log.Logger
}

func NewExerciseHandler(store store.ExerciseStore, logger *log.Logger) *ExerciseHandler {
	return &ExerciseHandler{
		store: store,
		log:   logger,
	}
}

type createExerciseRequest struct {
	Name        string `json:"name"         validate:"required,min=3,max=100"`
	Description string `json:"description"  validate:"required,min=10"`
	MuscleGroup string `json:"muscle_group" validate:"required,oneof=chest back legs shoulders arms core"`
	Category    string `json:"category"     validate:"required,oneof=strength cardio flexibility"`
}

type exerciseResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	MuscleGroup string    `json:"muscle_group"`
	Category    string    `json:"category"`
	CreatedAt   string    `json:"created_at"`
}

func (h *ExerciseHandler) CreateExercise(w http.ResponseWriter, r *http.Request) {
	var req createExerciseRequest
	if err := ReadJSON(r, &req); err != nil {
		ErrorJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := Validate.Struct(req); err != nil {
		ValidationErrorJSON(w, err)
		return
	}

	exercise := &store.Exercise{
		Name:        req.Name,
		Description: req.Description,
		MuscleGroup: req.MuscleGroup,
		Category:    req.Category,
	}

	if err := h.store.CreateExercise(r.Context(), exercise); err != nil {
		h.log.Printf("error creating exercise: %v", err)
		ErrorJSON(w, http.StatusInternalServerError, "could not create exercise")
		return
	}

	WriteJSON(w, http.StatusCreated, toExerciseResponse(exercise))
}

func (h *ExerciseHandler) GetExercise(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		ErrorJSON(w, http.StatusBadRequest, "invalid exercise id")
		return
	}

	exercise, err := h.store.GetExerciseByID(r.Context(), id)
	if err != nil {
		ErrorJSON(w, http.StatusNotFound, "exercise not found")
		return
	}

	WriteJSON(w, http.StatusOK, toExerciseResponse(exercise))
}

func (h *ExerciseHandler) GetAllExercises(w http.ResponseWriter, r *http.Request) {
	exercises, err := h.store.GetAllExercises(r.Context())
	if err != nil {
		h.log.Printf("error fetching exercises: %v", err)
		ErrorJSON(w, http.StatusInternalServerError, "could not fetch exercises")
		return
	}

	var resp []exerciseResponse
	for _, e := range exercises {
		resp = append(resp, toExerciseResponse(e))
	}

	WriteJSON(w, http.StatusOK, resp)
}

func toExerciseResponse(e *store.Exercise) exerciseResponse {
	return exerciseResponse{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		MuscleGroup: e.MuscleGroup,
		Category:    e.Category,
		CreatedAt:   e.CreatedAt.String(),
	}
}
