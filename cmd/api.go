package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/NoobsBucket/iqra-initi/internal/auth"
	"github.com/NoobsBucket/iqra-initi/internal/blog"
	"github.com/NoobsBucket/iqra-initi/internal/categories"
	"github.com/NoobsBucket/iqra-initi/internal/contact"
	"github.com/NoobsBucket/iqra-initi/internal/courses"
	"github.com/NoobsBucket/iqra-initi/internal/enrollment"
	"github.com/NoobsBucket/iqra-initi/internal/lessons"
	"github.com/NoobsBucket/iqra-initi/internal/mailer"
	"github.com/NoobsBucket/iqra-initi/internal/notifications"
	"github.com/NoobsBucket/iqra-initi/internal/products"
	"github.com/NoobsBucket/iqra-initi/internal/reviews"
	"github.com/NoobsBucket/iqra-initi/internal/settings"
	"github.com/NoobsBucket/iqra-initi/internal/users"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type routeRateLimiter struct {
	mu        sync.Mutex
	limit     int
	window    time.Duration
	requests  map[string][]time.Time
}

func newRouteRateLimiter(limit int, window time.Duration) *routeRateLimiter {
	return &routeRateLimiter{
		limit:    limit,
		window:   window,
		requests: make(map[string][]time.Time),
	}
}

func (rl *routeRateLimiter) allow(key string) bool {
	now := time.Now()
	cutoff := now.Add(-rl.window)

	rl.mu.Lock()
	defer rl.mu.Unlock()

	entries := rl.requests[key]
	filtered := entries[:0]
	for _, ts := range entries {
		if ts.After(cutoff) {
			filtered = append(filtered, ts)
		}
	}

	if len(filtered) >= rl.limit {
		rl.requests[key] = filtered
		return false
	}

	filtered = append(filtered, now)
	rl.requests[key] = filtered
	return true
}

func hashString(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func extractIP(r *http.Request) string {
	for _, header := range []string{"X-Forwarded-For", "X-Real-IP", "CF-Connecting-IP", "X-Client-IP"} {
		if value := strings.TrimSpace(r.Header.Get(header)); value != "" {
			return strings.Split(value, ",")[0]
		}
	}
	return strings.TrimSpace(r.RemoteAddr)
}

func extractBrowserID(r *http.Request) string {
	if value := strings.TrimSpace(r.Header.Get("X-Browser-ID")); value != "" {
		return value
	}
	return hashString(r.UserAgent())
}

func extractDeviceFingerprint(r *http.Request) string {
	if value := strings.TrimSpace(r.Header.Get("X-Device-Fingerprint")); value != "" {
		return value
	}
	raw := strings.Join([]string{
		r.UserAgent(),
		r.Header.Get("Accept-Language"),
		r.Header.Get("Accept-Encoding"),
		r.Header.Get("Sec-CH-UA"),
	}, "|")
	return hashString(raw)
}

func rateLimitMiddleware(limit int, window time.Duration) func(http.Handler) http.Handler {
	limiter := newRouteRateLimiter(limit, window)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := extractIP(r)
			browserID := extractBrowserID(r)
			deviceFingerprint := extractDeviceFingerprint(r)
			key := fmt.Sprintf("%s|%s|%s|%s|%s", r.Method, r.URL.Path, ip, browserID, deviceFingerprint)

			if !limiter.allow(key) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = w.Write([]byte(`{"error":"rate limit exceeded"}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr) // pick one ClientIPFrom* based on your infra, see below
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(rateLimitMiddleware(60, time.Minute))
	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			ip = forwarded
		} else if realip := r.Header.Get("X-Real-IP"); realip != "" {
			ip = realip
		}
		log.Printf("Client IP: %s", ip)

		userAgent := r.UserAgent()
		log.Printf("User Agent: %s", userAgent)
		reqId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("Request ID: %s", reqId)

		w.Write([]byte("Ahaan go for now , yoo yoo hi"))

	})
	// stores
	authStore := auth.NewStore(app.db)
	categoryStore := categories.NewStore(app.db)
	courseStore := courses.NewStore(app.db)
	reviewStore := reviews.NewStore(app.db)
	blogStore := blog.NewStore(app.db)
	lessonStore := lessons.NewStore(app.db)
	enrollmentStore := enrollment.NewStore(app.db)
	userStore := users.NewStore(app.db)
	contactStore := contact.NewStore(app.db)
	notifStore := notifications.NewStore(app.db)
	settingsStore := settings.NewStore(app.db)

	// services
	productsService := products.NewService()
	productsHandler := products.NewHandler(productsService)
	authService := auth.NewService(authStore, app.config.jwt.secret, app.config.jwt.exp, app.mailer)
	categoryService := categories.NewService(categoryStore)
	courseService := courses.NewService(courseStore)
	reviewService := reviews.NewService(reviewStore)
	blogService := blog.NewService(blogStore)
	lessonService := lessons.NewService(lessonStore)
	enrollmentService := enrollment.NewService(enrollmentStore)
	userService := users.NewService(userStore)
	contactService := contact.NewService(contactStore)
	notifService := notifications.NewService(notifStore)
	settingsService := settings.NewService(settingsStore)

	// handlers
	authHandler := auth.NewHandler(authService)
	categoryHandler := categories.NewHandler(categoryService)
	courseHandler := courses.NewHandler(courseService)
	reviewHandler := reviews.NewHandler(reviewService)
	blogHandler := blog.NewHandler(blogService)
	lessonHandler := lessons.NewHandler(lessonService)
	enrollmentHandler := enrollment.NewHandler(enrollmentService)
	userHandler := users.NewHandler(userService)
	contactHandler := contact.NewHandler(contactService)
	notifHandler := notifications.NewHandler(notifService)
	settingsHandler := settings.NewHandler(settingsService)

	// routes
	r.Get("/products", productsHandler.GetProducts)
	r.Get("/v1/users/{id}/avatar", app.avatarHandler)
	r.Route("/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
			r.Post("/logout", authHandler.Logout)
			r.Post("/verify-otp", authHandler.VerifyOTP)
			r.Post("/resend-otp", authHandler.ResendOTP)
			r.Post("/forgot-password", authHandler.ForgotPassword)
			r.Post("/reset-password", authHandler.ResetPassword)
			r.Post("/google", authHandler.GoogleAuth)
		})
		r.Route("/categories", func(r chi.Router) {
			r.Get("/", categoryHandler.GetAll)
			r.Get("/{id}", categoryHandler.GetOne)
			r.Post("/", categoryHandler.Create)
			r.Patch("/{id}", categoryHandler.Update)
			r.Delete("/{id}", categoryHandler.Delete)
		})
		r.Route("/courses", func(r chi.Router) {
			r.Get("/", courseHandler.GetAll)
			r.Get("/{id}", courseHandler.GetOne)
			r.Post("/", courseHandler.Create)
			r.Patch("/{id}", courseHandler.Update)
			r.Delete("/{id}", courseHandler.Delete)
			r.Patch("/{id}/publish", courseHandler.Publish)
			r.Get("/{courseID}/reviews", reviewHandler.GetByCourse)
			r.Post("/{courseID}/reviews", reviewHandler.Create)
			r.Get("/{courseID}/lessons", lessonHandler.GetByCourse)
			r.Post("/{courseID}/lessons", lessonHandler.Create)
			r.Get("/{courseID}/enrollments", enrollmentHandler.GetCourseEnrollments)
		})
		r.Route("/lessons", func(r chi.Router) {
			r.Get("/{id}", lessonHandler.GetOne)
			r.Patch("/{id}", lessonHandler.Update)
			r.Delete("/{id}", lessonHandler.Delete)
			r.Patch("/{id}/reorder", lessonHandler.Reorder)
		})
		r.Route("/reviews", func(r chi.Router) {
			r.Patch("/{id}", reviewHandler.Update)
			r.Delete("/{id}", reviewHandler.Delete)
		})
		r.Route("/enrollments", func(r chi.Router) {
			r.Get("/", enrollmentHandler.GetAll)
			r.Post("/", enrollmentHandler.Enroll)
			r.Get("/user/{userID}", enrollmentHandler.GetUserEnrollments)
			r.Patch("/{id}/progress", enrollmentHandler.UpdateProgress)
			r.Patch("/{id}/payment", enrollmentHandler.UpdatePaymentStatus)
			r.Patch("/{id}/complete", enrollmentHandler.Complete)
			r.Delete("/{id}", enrollmentHandler.Delete)
		})
		r.Route("/users", func(r chi.Router) {
			r.Get("/", userHandler.GetAll)
			r.Get("/{id}", userHandler.GetOne)
			r.Patch("/{id}/role", userHandler.AssignRole)
		})
		r.Route("/contact", func(r chi.Router) {
			r.Post("/", contactHandler.Send)
			r.Get("/", contactHandler.GetAll)
			r.Patch("/{id}/read", contactHandler.MarkRead)
			r.Delete("/{id}", contactHandler.Delete)
		})
		r.Route("/blog", func(r chi.Router) {
			r.Get("/categories", blogHandler.GetAllCategories)
			r.Post("/categories", blogHandler.CreateCategory)
			r.Get("/categories/{id}", blogHandler.GetOneCategory)
			r.Patch("/categories/{id}", blogHandler.UpdateCategory)
			r.Delete("/categories/{id}", blogHandler.DeleteCategory)

			r.Get("/", blogHandler.GetAll)
			r.Get("/slug/{slug}", blogHandler.GetBySlug)
			r.Get("/{id}", blogHandler.GetOne)
			r.Post("/", blogHandler.Create)
			r.Patch("/{id}", blogHandler.Update)
			r.Delete("/{id}", blogHandler.Delete)
			r.Patch("/{id}/publish", blogHandler.Publish)
		})
		r.Route("/notifications", func(r chi.Router) {
			r.Get("/user/{userID}", notifHandler.GetAll)
			r.Patch("/{id}/read", notifHandler.MarkRead)
			r.Patch("/user/{userID}/read-all", notifHandler.MarkAllRead)
			r.Delete("/{id}", notifHandler.Delete)
		})
		r.Route("/settings", func(r chi.Router) {
			r.Get("/", settingsHandler.Get)
			r.Patch("/", settingsHandler.Update)
		})
	})

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 60,
		IdleTimeout:  time.Minute,
	}
	log.Printf("server started at %s", app.config.addr)
	return srv.ListenAndServe()
}

type application struct {
	config    config
	db        *pgxpool.Pool
	mailer    *mailer.Mailer
	authStore auth.Store
}

type config struct {
	addr   string
	env    string
	db     dbConfig
	jwt    jwtConfig
	mailer mailerConfig
}

type mailerConfig struct {
	apiKey string
	from   string
}

type dbConfig struct {
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  time.Duration
}

type jwtConfig struct {
	secret string
	exp    time.Duration
}
