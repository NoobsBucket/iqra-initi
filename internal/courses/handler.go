package courses

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"unicode/utf8"

	"github.com/NoobsBucket/iqra-initi/internal/revalidation"
	"github.com/go-chi/chi/v5"
)

type handler struct {
	service service
}

func NewHandler(service service) *handler {
	return &handler{service: service}
}

func (h *handler) GetAll(w http.ResponseWriter, r *http.Request) {
	courses, err := h.service.GetCourses(r.Context())
	if err != nil {
		jsonError(w, "failed to get courses", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, courses)
}

func (h *handler) GetOne(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	course, err := h.service.GetCourse(r.Context(), id)
	if err != nil {
		jsonError(w, "course not found", http.StatusNotFound)
		return
	}
	jsonResponse(w, http.StatusOK, course)
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title           string   `json:"title"`
		Description     string   `json:"description"`
		ImageURL        string   `json:"image_url"`
		MetaTitle       string   `json:"meta_title"`
		MetaDescription string   `json:"meta_description"`
		MetaKeywords    string   `json:"meta_keywords"`
		Price           float64  `json:"price"`
		DiscountPrice   float64  `json:"discount_price"`
		Currency        string   `json:"currency"`
		CreatedBy       string   `json:"created_by"`
		CategoryIDs     []string `json:"category_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Title == "" {
		jsonError(w, "title is required", http.StatusBadRequest)
		return
	}
	if err := validateMetadata(req.MetaTitle, req.MetaDescription); err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Currency == "" {
		req.Currency = "USD"
	}
	course, err := h.service.CreateCourse(r.Context(), req.CreatedBy, req.Title, req.Description, req.ImageURL, req.Currency, req.MetaTitle, req.MetaDescription, req.MetaKeywords, req.Price, req.DiscountPrice, req.CategoryIDs)
	if err != nil {
		log.Println("create course error:", err)
		jsonError(w, "failed to create course", http.StatusInternalServerError)
		return
	}
	revalidation.NotifySEO()
	jsonResponse(w, http.StatusCreated, course)
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		Title           string  `json:"title"`
		Description     string  `json:"description"`
		ImageURL        string  `json:"image_url"`
		MetaTitle       string  `json:"meta_title"`
		MetaDescription string  `json:"meta_description"`
		MetaKeywords    string  `json:"meta_keywords"`
		Price           float64 `json:"price"`
		DiscountPrice   float64 `json:"discount_price"`
		Currency        string  `json:"currency"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := validateMetadata(req.MetaTitle, req.MetaDescription); err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
	course, err := h.service.UpdateCourse(r.Context(), id, req.Title, req.Description, req.ImageURL, req.Currency, req.MetaTitle, req.MetaDescription, req.MetaKeywords, req.Price, req.DiscountPrice)
	if err != nil {
		log.Println("update course error:", err)
		jsonError(w, "failed to update course", http.StatusInternalServerError)
		return
	}
	revalidation.NotifySEO()
	jsonResponse(w, http.StatusOK, course)
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.service.DeleteCourse(r.Context(), id); err != nil {
		jsonError(w, "failed to delete course", http.StatusInternalServerError)
		return
	}
	revalidation.NotifySEO()
	jsonResponse(w, http.StatusOK, map[string]any{"message": "course deleted"})
}

func (h *handler) Publish(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.service.PublishCourse(r.Context(), id); err != nil {
		jsonError(w, "failed to publish course", http.StatusInternalServerError)
		return
	}
	revalidation.NotifySEO()
	jsonResponse(w, http.StatusOK, map[string]any{"message": "course published"})
}

func jsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, message string, status int) {
	jsonResponse(w, status, map[string]any{"error": message})
}

func validateMetadata(title, description string) error {
	if utf8.RuneCountInString(title) > 70 {
		return fmt.Errorf("meta_title must be at most 70 characters")
	}
	if utf8.RuneCountInString(description) > 170 {
		return fmt.Errorf("meta_description must be at most 170 characters")
	}
	return nil
}
