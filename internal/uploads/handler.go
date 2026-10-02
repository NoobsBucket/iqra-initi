package uploads

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Config struct {
	AccountID       string
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
	Endpoint        string
	PublicURL       string
}

func ConfigFromEnv() Config {
	return Config{
		AccountID:       os.Getenv("R2_ACCOUNT_ID"),
		AccessKeyID:     os.Getenv("R2_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("R2_SECRET_ACCESS_KEY"),
		Bucket:          os.Getenv("R2_BUCKET_NAME"),
		Endpoint:        os.Getenv("R2_ENDPOINT"),
		PublicURL:       os.Getenv("R2_PUBLIC_URL"),
	}
}

type Handler struct {
	config    Config
	once      sync.Once
	presigner *s3.PresignClient
	initErr   error
}

func NewHandler(config Config) *Handler { return &Handler{config: config} }

func (h *Handler) Presign(w http.ResponseWriter, r *http.Request) {
	var request struct {
		FileName    string `json:"file_name"`
		Filename    string `json:"filename"`
		ContentType string `json:"content_type"`
	}
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	if err := decoder.Decode(&request); err != nil {
		uploadError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if request.FileName == "" {
		request.FileName = request.Filename
	}
	name := path.Base(strings.ReplaceAll(strings.TrimSpace(request.FileName), "\\", "/"))
	if name == "" || name == "." || name == "/" {
		uploadError(w, "file_name is required", http.StatusBadRequest)
		return
	}
	mediaType, _, err := mime.ParseMediaType(request.ContentType)
	if err != nil || !strings.HasPrefix(strings.ToLower(mediaType), "image/") {
		uploadError(w, "content_type must be an image MIME type", http.StatusBadRequest)
		return
	}

	h.once.Do(h.initialize)
	if h.initErr != nil {
		log.Printf("initialize R2 upload signer: %v", h.initErr)
		uploadError(w, "upload service is not configured", http.StatusServiceUnavailable)
		return
	}

	randomID := make([]byte, 16)
	if _, err := rand.Read(randomID); err != nil {
		log.Printf("generate upload key: %v", err)
		uploadError(w, "failed to prepare upload", http.StatusInternalServerError)
		return
	}
	key := "images/" + hex.EncodeToString(randomID) + "-" + name
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	presigned, err := h.presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(h.config.Bucket),
		Key:         aws.String(key),
		ContentType: aws.String(mediaType),
	}, func(options *s3.PresignOptions) {
		options.Expires = 15 * time.Minute
	})
	if err != nil {
		log.Printf("presign R2 upload: %v", err)
		uploadError(w, "failed to create upload URL", http.StatusInternalServerError)
		return
	}

	segments := strings.Split(key, "/")
	for index := range segments {
		segments[index] = url.PathEscape(segments[index])
	}
	publicURL := strings.TrimRight(h.config.PublicURL, "/") + "/" + strings.Join(segments, "/")
	writeUploadJSON(w, http.StatusOK, map[string]any{
		"upload_url": presigned.URL,
		"file_url":   publicURL,
		"key":        key,
		"method":     http.MethodPut,
		"headers":    map[string]string{"Content-Type": mediaType},
		"expires_in": int((15 * time.Minute).Seconds()),
	})
}

func (h *Handler) initialize() {
	if h.config.AccountID == "" || h.config.AccessKeyID == "" || h.config.SecretAccessKey == "" || h.config.Bucket == "" || h.config.PublicURL == "" {
		h.initErr = fmt.Errorf("R2_ACCOUNT_ID, R2_ACCESS_KEY_ID, R2_SECRET_ACCESS_KEY, R2_BUCKET_NAME, and R2_PUBLIC_URL are required")
		return
	}
	publicURL, err := url.Parse(strings.TrimSpace(h.config.PublicURL))
	if err != nil || publicURL.Hostname() == "" || (publicURL.Scheme != "http" && publicURL.Scheme != "https") {
		h.initErr = fmt.Errorf("R2_PUBLIC_URL must be an absolute HTTP or HTTPS URL")
		return
	}
	endpoint := strings.TrimSpace(h.config.Endpoint)
	if endpoint == "" {
		endpoint = "https://" + h.config.AccountID + ".r2.cloudflarestorage.com"
	}
	awsConfig, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion("auto"),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(h.config.AccessKeyID, h.config.SecretAccessKey, "")),
		awsconfig.WithRequestChecksumCalculation(aws.RequestChecksumCalculationWhenRequired),
		awsconfig.WithResponseChecksumValidation(aws.ResponseChecksumValidationWhenRequired),
	)
	if err != nil {
		h.initErr = err
		return
	}
	client := s3.NewFromConfig(awsConfig, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(endpoint)
		options.UsePathStyle = true
	})
	h.presigner = s3.NewPresignClient(client)
}

func writeUploadJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func uploadError(w http.ResponseWriter, message string, status int) {
	writeUploadJSON(w, status, map[string]string{"error": message})
}
