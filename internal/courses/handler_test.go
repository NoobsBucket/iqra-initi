package courses

import (
	"strings"
	"testing"
)

func TestValidateMetadataLimits(t *testing.T) {
	if err := validateMetadata(strings.Repeat("x", 70), strings.Repeat("y", 170)); err != nil {
		t.Fatalf("values at the limits should be allowed: %v", err)
	}
	if err := validateMetadata(strings.Repeat("x", 71), ""); err == nil {
		t.Fatal("meta_title over 70 characters should be rejected")
	}
	if err := validateMetadata("", strings.Repeat("y", 171)); err == nil {
		t.Fatal("meta_description over 170 characters should be rejected")
	}
}
