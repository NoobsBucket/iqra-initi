package users

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type handler struct {
	service service
}

func NewHandler(service service) *handler {
	return &handler{service: service}
}

func (h *handler) GetAll(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	users, err := h.service.GetUsers(r.Context(), search)
	if err != nil {
		jsonError(w, "failed to get users", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{
		"success": true,
		"users":   users,
		"total":   len(users),
	})
}

func (h *handler) GetOne(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, err := h.service.GetUser(r.Context(), id)
	if err != nil {
		jsonError(w, "user not found", http.StatusNotFound)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{
		"success": true,
		"user":    user,
	})
}

func (h *handler) AssignRole(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Role == "" {
		jsonError(w, "role is required", http.StatusBadRequest)
		return
	}

	user, err := h.service.AssignRole(r.Context(), id, req.Role)
	if err != nil {
		if errors.Is(err, ErrInvalidRole) {
			jsonError(w, "invalid role — must be user, admin, instructor or blog_writer", http.StatusBadRequest)
			return
		}
		jsonError(w, "failed to update role", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "role updated to " + req.Role,
		"user":    user,
	})
}

func jsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, message string, status int) {
	jsonResponse(w, status, map[string]any{
		"success": false,
		"error":   message,
	})
}
