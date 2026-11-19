package discussionservice

import (
	"testing"

	"github.com/google/uuid"
	"nalakarsa/internal/model"
)

func TestDiscussionAvatarURLKeepsStoredAvatar(t *testing.T) {
	user := model.User{AvatarURL: "https://example.com/avatar.png"}
	if got := discussionAvatarURL(user); got != user.AvatarURL {
		t.Fatalf("got %q, want stored avatar %q", got, user.AvatarURL)
	}
}

func TestDiscussionAvatarURLProvidesFallbackWhenEmpty(t *testing.T) {
	user := model.User{ID: uuid.MustParse("11111111-1111-1111-1111-111111111111"), FullName: "Demo User"}
	got := discussionAvatarURL(user)
	if got == "" || got[:len("https://api.dicebear.com/")] != "https://api.dicebear.com/" {
		t.Fatalf("expected DiceBear fallback, got %q", got)
	}
}
