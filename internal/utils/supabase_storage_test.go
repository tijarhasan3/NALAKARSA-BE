package utils

import (
	"net/http"
	"net/http/httptest"
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

func TestNormalizeSupabasePublicURLMakesAvatarPathAbsolute(t *testing.T) {
	cfg := &config.Config{SupabaseURL: "https://project.supabase.co"}
	got := NormalizeSupabasePublicURL("/avatars/user.png", "avatars", cfg)
	want := "https://project.supabase.co/storage/v1/object/public/avatars/user.png"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestNormalizeSupabasePublicURLUsesAssetsForUserAvatarPath(t *testing.T) {
	cfg := &config.Config{SupabaseURL: "https://project.supabase.co"}
	got := NormalizeSupabasePublicURL("users/user-id/avatars/photo.png", "assets", cfg)
	want := "https://project.supabase.co/storage/v1/object/public/assets/users/user-id/avatars/photo.png"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestFindLatestPublicObjectURLListsNewestAvatar(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/storage/v1/object/list/assets" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer service-key" {
			t.Fatalf("unexpected authorization header: %s", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"name":"photo.png","created_at":"2026-11-19T10:00:00Z"}]`))
	}))
	defer server.Close()

	cfg := &config.Config{
		SupabaseURL:            server.URL,
		SupabaseServiceRoleKey: "service-key",
	}
	got, err := FindLatestPublicObjectURL("users/user-id/avatars", "assets", cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := server.URL + "/storage/v1/object/public/assets/users/user-id/avatars/photo.png"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
