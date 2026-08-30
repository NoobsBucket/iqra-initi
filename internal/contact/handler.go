package contact

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

func (h *handler) Send(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Phone   string `json:"phone"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.Email == "" || req.Message == "" {
		jsonError(w, "name, email and message are required", http.StatusBadRequest)
		return
	}

	msg, err := h.service.SendMessage(r.Context(), req.Name, req.Email, req.Phone, req.Subject, req.Message)
	if err != nil {
		jsonError(w, "failed to send message", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusCreated, map[string]any{
		"success": true,
		"message": "message sent successfully",
		"data":    msg,
	})
}

func (h *handler) GetAll(w http.ResponseWriter, r *http.Request) {
	messages, err := h.service.GetMessages(r.Context())
	if err != nil {
		jsonError(w, "failed to get messages", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{
		"success":  true,
		"messages": messages,
		"total":    len(messages),
	})
}

func (h *handler) MarkRead(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.service.MarkRead(r.Context(), id); err != nil {
		jsonError(w, "failed to mark as read", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "marked as read",
	})
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.service.DeleteMessage(r.Context(), id); err != nil {
		jsonError(w, "failed to delete message", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "message deleted",
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
