package notifications

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

func (h *handler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	notifications, err := h.service.GetUserNotifications(r.Context(), userID)
	if err != nil {
		jsonError(w, "failed to get notifications", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, notifications)
}

func (h *handler) MarkRead(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.service.MarkRead(r.Context(), id); err != nil {
		jsonError(w, "failed to mark as read", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{"message": "marked as read"})
}

func (h *handler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if err := h.service.MarkAllRead(r.Context(), userID); err != nil {
		jsonError(w, "failed to mark all as read", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{"message": "all marked as read"})
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.service.DeleteNotification(r.Context(), id); err != nil {
		jsonError(w, "failed to delete notification", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{"message": "notification deleted"})
}

func jsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, message string, status int) {
	jsonResponse(w, status, map[string]any{"error": message})
}
