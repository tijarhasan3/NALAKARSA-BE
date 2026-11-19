package discussionrepository

import "testing"

func TestDiscussionOrderFollowsSortContract(t *testing.T) {
	tests := []struct {
		name string
		sort string
		want string
	}{
		{name: "newest", sort: "newest", want: "discussions.created_at desc"},
		{name: "oldest", sort: "oldest", want: "discussions.created_at asc"},
		{name: "popular", sort: "popular", want: "discussions.upvote_count desc"},
		{name: "empty defaults to newest", sort: "", want: "discussions.created_at desc"},
		{name: "unknown defaults to newest", sort: "unsupported", want: "discussions.created_at desc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := discussionOrder(tt.sort); got != tt.want {
				t.Fatalf("discussionOrder(%q) = %q, want %q", tt.sort, got, tt.want)
			}
		})
	}
}
