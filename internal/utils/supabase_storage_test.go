package utils

import "testing"

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
