package auth

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type handler struct {
	service service
}

func NewHandler(service service) *handler {
	return &handler{service: service}
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		jsonError(w, "name, email and password are required", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 8 {
		jsonError(w, "password must be at least 8 characters", http.StatusBadRequest)
		return
	}

	if err := h.service.RegisterUser(r.Context(), req.Name, req.Email, req.Password); err != nil {
		if errors.Is(err, ErrEmailTaken) {
			jsonError(w, "email already in use", http.StatusConflict)
			return
		}
		log.Printf("register user: %v", err)
		jsonError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusCreated, map[string]any{
		"message": "verification code sent to " + req.Email,
	})
}

func (h *handler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.OTP == "" {
		jsonError(w, "email and otp are required", http.StatusBadRequest)
		return
	}

	user, token, err := h.service.VerifyOTP(r.Context(), req.Email, req.OTP)
	if err != nil {
		if errors.Is(err, ErrInvalidOTP) {
			jsonError(w, "invalid or expired OTP", http.StatusUnauthorized)
			return
		}
		log.Printf("verify OTP: %v", err)
		jsonError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]any{
		"message": "email verified successfully",
		"user":    user,
		"token":   token,
	})
}

func (h *handler) ResendOTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		jsonError(w, "email is required", http.StatusBadRequest)
		return
	}

	if err := h.service.ResendOTP(r.Context(), req.Email); err != nil {
		if errors.Is(err, ErrAlreadyVerified) {
			jsonError(w, "email already verified", http.StatusBadRequest)
			return
		}
		log.Printf("resend OTP: %v", err)
		jsonError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]any{
		"message": "verification code sent to " + req.Email,
	})
}

func (h *handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		jsonError(w, "email is required", http.StatusBadRequest)
		return
	}

	if err := h.service.ForgotPassword(r.Context(), req.Email); err != nil {
		log.Printf("forgot password: %v", err)
	}

	jsonResponse(w, http.StatusOK, map[string]any{
		"message": "if that email exists you will receive a reset code",
	})
}

func (h *handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email       string `json:"email"`
		OTP         string `json:"otp"`
		NewPassword string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.OTP == "" || req.NewPassword == "" {
		jsonError(w, "email, otp and new_password are required", http.StatusBadRequest)
		return
	}

	if len(req.NewPassword) < 8 {
		jsonError(w, "password must be at least 8 characters", http.StatusBadRequest)
		return
	}

	if err := h.service.ResetPassword(r.Context(), req.Email, req.OTP, req.NewPassword); err != nil {
		if errors.Is(err, ErrInvalidOTP) {
			jsonError(w, "invalid or expired code", http.StatusUnauthorized)
			return
		}
		log.Printf("reset password: %v", err)
		jsonError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]any{
		"message": "password reset successfully",
	})
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		jsonError(w, "email and password are required", http.StatusBadRequest)
		return
	}

	user, token, err := h.service.LoginUser(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			jsonError(w, "invalid email or password", http.StatusUnauthorized)
			return
		}
		if errors.Is(err, ErrNotVerified) {
			jsonError(w, "please verify your email first", http.StatusUnauthorized)
			return
		}
		log.Printf("login user: %v", err)
		jsonError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]any{
		"user":  user,
		"token": token,
	})
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {
	// JWT is stateless — client just deletes the token
	jsonResponse(w, http.StatusOK, map[string]any{
		"message": "logged out successfully",
	})
}

// helpers
func jsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, message string, status int) {
	jsonResponse(w, status, map[string]any{
		"error": message,
	})
}
