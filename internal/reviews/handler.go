package reviews

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type handler struct{ service service }

func NewHandler(service service) *handler { return &handler{service: service} }

func (h *handler) GetByCourse(w http.ResponseWriter, r *http.Request) {
	reviews, err := h.service.GetCourseReviews(r.Context(), chi.URLParam(r, "courseID"))
	if err != nil {
		jsonError(w, "failed to get reviews", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, reviews)
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID     string  `json:"user_id"`
		ReviewText string  `json:"review_text"`
		Rating     float64 `json:"rating"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	review, err := h.service.CreateReview(r.Context(), request.UserID, chi.URLParam(r, "courseID"), request.ReviewText, request.Rating)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonResponse(w, http.StatusCreated, review)
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ReviewText string  `json:"review_text"`
		Rating     float64 `json:"rating"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	review, err := h.service.UpdateReview(r.Context(), chi.URLParam(r, "id"), request.ReviewText, request.Rating)
	if err != nil {
		jsonError(w, "failed to update review", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, review)
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.service.DeleteReview(r.Context(), chi.URLParam(r, "id")); err != nil {
		jsonError(w, "failed to delete review", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{"message": "review deleted"})
}

func jsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
func jsonError(w http.ResponseWriter, message string, status int) {
	jsonResponse(w, status, map[string]any{"error": message})
}
