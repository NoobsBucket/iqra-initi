package lessons

import "testing"

func TestValidThumbnailURL(t *testing.T) {
	valid := "https://cdn.example.com/images/lesson.webp"
	invalid := "/images/lesson.webp"

	if !validThumbnailURL(nil) {
		t.Fatal("nil thumbnail URL should be allowed")
	}
	if !validThumbnailURL(&valid) {
		t.Fatal("absolute HTTPS URL should be allowed")
	}
	if validThumbnailURL(&invalid) {
		t.Fatal("relative URL should be rejected")
	}
}
