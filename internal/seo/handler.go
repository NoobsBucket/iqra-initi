package seo

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/NoobsBucket/iqra-initi/internal/auth"
	"github.com/NoobsBucket/iqra-initi/internal/revalidation"
	"github.com/go-chi/chi/v5"
)

type handler struct{ service service }

func NewHandler(service service) *handler { return &handler{service: service} }

func (h *handler) GetPublic(w http.ResponseWriter, r *http.Request) {
	pagePath, err := normalizePagePath(r.URL.Query().Get("path"))
	if err != nil {
		writeError(w, "path is required", http.StatusBadRequest)
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=60, s-maxage=300")
	page, err := h.service.GetPage(r.Context(), pagePath)
	if errors.Is(err, ErrNotFound) {
		writeError(w, "page SEO not found", http.StatusNotFound)
		return
	}
	if err != nil {
		writeError(w, "failed to get page SEO", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, page)
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	pages, err := h.service.GetPages(r.Context())
	if err != nil {
		writeError(w, "failed to get page SEO", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, pages)
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var page Page
	if err := json.NewDecoder(r.Body).Decode(&page); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := validatePage(&page); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	page.UpdatedBy = auth.UserIDFromContext(r.Context())
	created, err := h.service.CreatePage(r.Context(), &page)
	if errors.Is(err, ErrDuplicate) {
		writeError(w, "page_path already exists", http.StatusConflict)
		return
	}
	if err != nil {
		writeError(w, "failed to create page SEO", http.StatusInternalServerError)
		return
	}
	revalidation.NotifySEO()
	writeJSON(w, http.StatusCreated, created)
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	var page Page
	if err := json.NewDecoder(r.Body).Decode(&page); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := validatePage(&page); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	page.ID = chi.URLParam(r, "id")
	page.UpdatedBy = auth.UserIDFromContext(r.Context())
	updated, err := h.service.UpdatePage(r.Context(), &page)
	if errors.Is(err, ErrDuplicate) {
		writeError(w, "page_path already exists", http.StatusConflict)
		return
	}
	if errors.Is(err, ErrNotFound) {
		writeError(w, "page SEO not found", http.StatusNotFound)
		return
	}
	if err != nil {
		writeError(w, "failed to update page SEO", http.StatusInternalServerError)
		return
	}
	revalidation.NotifySEO()
	writeJSON(w, http.StatusOK, updated)
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	err := h.service.DeletePage(r.Context(), chi.URLParam(r, "id"))
	if errors.Is(err, ErrNotFound) {
		writeError(w, "page SEO not found", http.StatusNotFound)
		return
	}
	if err != nil {
		writeError(w, "failed to delete page SEO", http.StatusInternalServerError)
		return
	}
	revalidation.NotifySEO()
	writeJSON(w, http.StatusOK, map[string]string{"message": "page SEO deleted"})
}

func validatePage(page *Page) error {
	pagePath, err := normalizePagePath(page.PagePath)
	if err != nil {
		return err
	}
	page.PagePath = pagePath
	page.PageName = strings.TrimSpace(page.PageName)
	if page.PageName == "" {
		return errors.New("page_name is required")
	}
	if utf8.RuneCountInString(page.MetaTitle) > 70 {
		return errors.New("meta_title must be at most 70 characters")
	}
	if utf8.RuneCountInString(page.MetaDescription) > 170 {
		return errors.New("meta_description must be at most 170 characters")
	}
	if utf8.RuneCountInString(page.OGTitle) > 100 {
		return errors.New("og_title must be at most 100 characters")
	}
	if utf8.RuneCountInString(page.OGDescription) > 200 {
		return errors.New("og_description must be at most 200 characters")
	}
	if len(page.JSONLD) > 0 && !json.Valid(page.JSONLD) {
		return errors.New("json_ld must be valid JSON")
	}
	if strings.TrimSpace(page.Robots) == "" {
		page.Robots = "index,follow"
	}
	return nil
}

func normalizePagePath(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", errors.New("page_path is required")
	}
	if !strings.HasPrefix(value, "/") {
		value = "/" + value
	}
	if value != "/" {
		value = strings.TrimRight(value, "/")
		if value == "" {
			value = "/"
		}
	}
	return value, nil
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, message string, status int) {
	writeJSON(w, status, map[string]string{"error": message})
}
