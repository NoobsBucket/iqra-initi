package main

import (
	"net/http"

	dicebear "github.com/dicebear/dicebear-go/v10"
	"github.com/dicebear/styles/v10"
	"github.com/go-chi/chi/v5"
)

func (app *application) avatarHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := app.authStore.GetUserByID(r.Context(), id)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	style, err := dicebear.NewStyle([]byte(styles.Clay))
	if err != nil {
		http.Error(w, "failed to generate avatar", http.StatusInternalServerError)
		return
	}

	avatar, err := dicebear.NewAvatar(style, map[string]any{
		"seed": user.Email,
		"backgroundColor": []any{
			"ff2e63",
			"00c2a8",
			"ffb300",
			"3d5afe",
			"8e24aa",
			"00e676",
		},
	})
	if err != nil {
		http.Error(w, "failed to generate avatar", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=604800")
	w.Write([]byte(avatar.SVG()))
}
