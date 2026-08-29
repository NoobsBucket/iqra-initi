package lessons

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type handler struct {
	service service
}

func NewHandler(service service) *handler {
	return &handler{service: service}
}

func (h *handler) GetByCourse(w http.ResponseWriter, r *http.Request) {
	courseID := chi.URLParam(r, "courseID")
	lessons, err := h.service.GetCourseLessons(r.Context(), courseID)
	if err != nil {
		jsonError(w, "failed to get lessons", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, lessons)
}

func (h *handler) GetOne(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	lesson, err := h.service.GetLesson(r.Context(), id)
	if err != nil {
		jsonError(w, "lesson not found", http.StatusNotFound)
		return
	}
	jsonResponse(w, http.StatusOK, lesson)
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	courseID := chi.URLParam(r, "courseID")
	var req Lesson
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Title == "" {
		jsonError(w, "title is required", http.StatusBadRequest)
		return
	}
	req.CourseID = courseID
	lesson, err := h.service.CreateLesson(r.Context(), &req)
	if err != nil {
		jsonError(w, "failed to create lesson", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusCreated, lesson)
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req Lesson
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	req.ID = id
	lesson, err := h.service.UpdateLesson(r.Context(), &req)
	if err != nil {
		jsonError(w, "failed to update lesson", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, lesson)
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.service.DeleteLesson(r.Context(), id); err != nil {
		jsonError(w, "failed to delete lesson", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{"message": "lesson deleted"})
}

func (h *handler) Reorder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		Order int `json:"order"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.service.ReorderLesson(r.Context(), id, req.Order); err != nil {
		jsonError(w, "failed to reorder lesson", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{"message": "lesson reordered"})
}

func jsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, message string, status int) {
	jsonResponse(w, status, map[string]any{"error": message})
}
