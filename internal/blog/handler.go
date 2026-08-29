package blog

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type handler struct{ service service }

func NewHandler(service service) *handler { return &handler{service: service} }

func (h *handler) GetAll(w http.ResponseWriter, r *http.Request) {
	posts, err := h.service.GetPosts(r.Context())
	if err != nil {
		jsonError(w, "failed to get posts", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, posts)
}
func (h *handler) GetOne(w http.ResponseWriter, r *http.Request) {
	post, err := h.service.GetPost(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		jsonError(w, "post not found", http.StatusNotFound)
		return
	}
	jsonResponse(w, http.StatusOK, post)
}
func (h *handler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	post, err := h.service.GetPostBySlug(r.Context(), chi.URLParam(r, "slug"))
	if err != nil {
		jsonError(w, "post not found", http.StatusNotFound)
		return
	}
	jsonResponse(w, http.StatusOK, post)
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var post Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if post.Title == "" || post.Content == "" {
		jsonError(w, "title and content are required", http.StatusBadRequest)
		return
	}
	created, err := h.service.CreatePost(r.Context(), &post)
	if err != nil {
		log.Printf("create blog post: %v", err)
		jsonError(w, "failed to create post", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusCreated, created)
}
func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	var post Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	post.ID = chi.URLParam(r, "id")
	updated, err := h.service.UpdatePost(r.Context(), &post)
	if err != nil {
		jsonError(w, "failed to update post", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, updated)
}
func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.service.DeletePost(r.Context(), chi.URLParam(r, "id")); err != nil {
		jsonError(w, "failed to delete post", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{"message": "post deleted"})
}
func (h *handler) Publish(w http.ResponseWriter, r *http.Request) {
	if err := h.service.PublishPost(r.Context(), chi.URLParam(r, "id")); err != nil {
		jsonError(w, "failed to publish post", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{"message": "post published"})
}

func (h *handler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	cats, err := h.service.GetBlogCategories(r.Context())
	if err != nil {
		jsonError(w, "failed to get categories", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, cats)
}

func (h *handler) GetOneCategory(w http.ResponseWriter, r *http.Request) {
	cat, err := h.service.GetBlogCategory(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		jsonError(w, "category not found", http.StatusNotFound)
		return
	}
	jsonResponse(w, http.StatusOK, cat)
}

func (h *handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		jsonError(w, "name is required", http.StatusBadRequest)
		return
	}
	cat, err := h.service.CreateBlogCategory(r.Context(), req.Name, req.Description)
	if err != nil {
		jsonError(w, "failed to create category", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusCreated, cat)
}

func (h *handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	cat, err := h.service.UpdateBlogCategory(r.Context(), chi.URLParam(r, "id"), req.Name, req.Description)
	if err != nil {
		jsonError(w, "failed to update category", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, cat)
}

func (h *handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	if err := h.service.DeleteBlogCategory(r.Context(), chi.URLParam(r, "id")); err != nil {
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
