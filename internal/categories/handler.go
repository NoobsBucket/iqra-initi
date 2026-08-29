package categories

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
	categories, err := h.service.GetCategories(r.Context())
	if err != nil {
		jsonError(w, "failed to get categories", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, categories)
}

func (h *handler) GetOne(w http.ResponseWriter, r *http.Request) {
	category, err := h.service.GetCategory(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		jsonError(w, "category not found", http.StatusNotFound)
		return
	}
	jsonResponse(w, http.StatusOK, category)
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		ImageURL    string `json:"image_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if request.Name == "" {
		jsonError(w, "name is required", http.StatusBadRequest)
		return
	}
	category, err := h.service.CreateCategory(r.Context(), request.Name, request.Description, request.ImageURL)
	if err != nil {
		jsonError(w, "failed to create category", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusCreated, category)
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		ImageURL    string `json:"image_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	category, err := h.service.UpdateCategory(
		r.Context(), chi.URLParam(r, "id"), request.Name, request.Description, request.ImageURL,
	)
	if err != nil {
		jsonError(w, "failed to update category", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, category)
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.service.DeleteCategory(r.Context(), chi.URLParam(r, "id")); err != nil {
		jsonError(w, "failed to delete category", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{"message": "category deleted"})
}

func jsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, message string, status int) {
	jsonResponse(w, status, map[string]any{"error": message})
}
