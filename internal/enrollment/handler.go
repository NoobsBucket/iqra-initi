package enrollment

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type handler struct {
	service service
}

func NewHandler(service service) *handler {
	return &handler{service: service}
}

func (h *handler) Enroll(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID   string `json:"user_id"`
		CourseID string `json:"course_id"`
		FullName string `json:"full_name"`
		Whatsapp string `json:"whatsapp"`
		Notes    string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.UserID == "" || req.CourseID == "" || req.FullName == "" || req.Whatsapp == "" {
		jsonError(w, "user_id, course_id, full_name and whatsapp are required", http.StatusBadRequest)
		return
	}

	enrollment, err := h.service.EnrollUser(r.Context(), req.UserID, req.CourseID, req.FullName, req.Whatsapp, req.Notes)
	if err != nil {
		if errors.Is(err, ErrAlreadyEnrolled) {
			jsonError(w, "already enrolled in this course", http.StatusConflict)
			return
		}
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusCreated, map[string]any{
		"success":    true,
		"message":    "enrolled successfully",
		"enrollment": enrollment,
	})
}

func (h *handler) GetAll(w http.ResponseWriter, r *http.Request) {
	enrollments, err := h.service.GetAllEnrollments(r.Context())
	if err != nil {
		log.Printf("get all enrollments error: %v", err)
		jsonError(w, "failed to get enrollments", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{
		"success":     true,
		"enrollments": enrollments,
		"total":       len(enrollments),
	})
}

func (h *handler) GetUserEnrollments(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	enrollments, err := h.service.GetUserEnrollments(r.Context(), userID)
	if err != nil {
		log.Printf("get user enrollments error: %v", err)
		jsonError(w, "failed to get enrollments", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{
		"success":     true,
		"enrollments": enrollments,
		"total":       len(enrollments),
	})
}

func (h *handler) GetCourseEnrollments(w http.ResponseWriter, r *http.Request) {
	courseID := chi.URLParam(r, "courseID")
	enrollments, err := h.service.GetCourseEnrollments(r.Context(), courseID)
	if err != nil {
		log.Printf("get course enrollments error: %v", err)
		jsonError(w, "failed to get enrollments", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{
		"success":     true,
		"enrollments": enrollments,
		"total":       len(enrollments),
	})
}

func (h *handler) UpdateProgress(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		Progress int `json:"progress"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.service.UpdateProgress(r.Context(), id, req.Progress); err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "progress updated successfully",
	})
}

func (h *handler) UpdatePaymentStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Status == "" {
		jsonError(w, "status is required", http.StatusBadRequest)
		return
	}
	if err := h.service.UpdatePaymentStatus(r.Context(), id, req.Status); err != nil {
		jsonError(w, "failed to update payment status", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "payment status updated to " + req.Status,
	})
}

func (h *handler) Complete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.service.CompleteEnrollment(r.Context(), id); err != nil {
		jsonError(w, "failed to complete enrollment", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "course marked as completed",
	})
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.service.DeleteEnrollment(r.Context(), id); err != nil {
		jsonError(w, "failed to delete enrollment", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "enrollment deleted successfully",
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
