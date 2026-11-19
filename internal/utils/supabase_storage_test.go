package utils

import (
	"testing"

	"nalakarsa/internal/config"
)

func TestParseSupabaseSignedURLSupportsArrayResponse(t *testing.T) {
	got, err := parseSupabaseSignedURL([]byte(`[{"signedURL":"https://example.test/signed"}]`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "https://example.test/signed" {
		t.Fatalf("unexpected URL: %s", got)
	}
}

func TestParseSupabaseSignedURLSupportsWrappedResponse(t *testing.T) {
	got, err := parseSupabaseSignedURL([]byte(`{"signedURLs":[{"signedURL":"https://example.test/signed"}]}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "https://example.test/signed" {
		t.Fatalf("unexpected URL: %s", got)
	}
}

func TestParseSupabaseSignedURLSupportsSingleObjectResponse(t *testing.T) {
	got, err := parseSupabaseSignedURL([]byte(`{"signedURL":"https://example.test/signed"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "https://example.test/signed" {
		t.Fatalf("unexpected URL: %s", got)
	}
}

func TestNormalizeSupabaseStorageURLMakesRelativeSignedURLAbsolute(t *testing.T) {
	cfg := &config.Config{SupabaseURL: "https://project.supabase.co"}
	got := NormalizeSupabaseStorageURL("/object/sign/assets/users/file.png?token=test", cfg)
	want := "https://project.supabase.co/storage/v1/object/sign/assets/users/file.png?token=test"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
