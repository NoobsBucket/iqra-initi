package seo

import (
	"strings"
	"testing"
)

func TestNormalizePagePath(t *testing.T) {
	got, err := normalizePagePath("courses/")
	if err != nil || got != "/courses" {
		t.Fatalf("normalizePagePath() = %q, %v; want /courses", got, err)
	}
	got, err = normalizePagePath("/")
	if err != nil || got != "/" {
		t.Fatalf("normalizePagePath() = %q, %v; want /", got, err)
	}
	if _, err := normalizePagePath(" "); err == nil {
		t.Fatal("empty page path should be rejected")
	}
}

func TestValidatePage(t *testing.T) {
	page := &Page{PagePath: "about/", PageName: "About", MetaTitle: strings.Repeat("a", 70)}
	if err := validatePage(page); err != nil {
		t.Fatalf("valid page should pass validation: %v", err)
	}
	if page.PagePath != "/about" {
		t.Fatalf("page path was not normalized: %q", page.PagePath)
	}

	page.MetaTitle = strings.Repeat("a", 71)
	if err := validatePage(page); err == nil {
		t.Fatal("meta_title over 70 characters should be rejected")
	}
	page.MetaTitle = "Valid title"
	page.JSONLD = []byte(`{"@type":`)
	if err := validatePage(page); err == nil {
		t.Fatal("invalid json_ld should be rejected")
	}
}
