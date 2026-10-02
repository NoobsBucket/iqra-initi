package revalidation

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func NotifySEO() {
	frontendURL := strings.TrimRight(strings.TrimSpace(os.Getenv("FRONTEND_URL")), "/")
	secret := os.Getenv("REVALIDATE_SECRET")
	if frontendURL == "" || secret == "" {
		log.Printf("cache revalidation skipped: FRONTEND_URL or REVALIDATE_SECRET is unset")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, frontendURL+"/api/revalidate", bytes.NewBufferString(`{"tag":"seo"}`))
	if err != nil {
		log.Printf("cache revalidation request failed: %v", err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Revalidate-Secret", secret)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("cache revalidation request failed: %v", err)
		return
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		log.Printf("cache revalidation returned status %d", response.StatusCode)
	}
}
