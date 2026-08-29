package settings

import (
	"encoding/json"
	"net/http"
)

type handler struct {
	service service
}

func NewHandler(service service) *handler {
	return &handler{service: service}
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	settings, err := h.service.GetSettings(r.Context())
	if err != nil {
		jsonError(w, "failed to get settings", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, settings)
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	var req Settings
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	updated, err := h.service.UpdateSettings(r.Context(), &req)
	if err != nil {
		jsonError(w, "failed to update settings", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, updated)
}

func jsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, message string, status int) {
	jsonResponse(w, status, map[string]any{"error": message})
}
