package discussionhandler

import "testing"

func TestIsMyTopicsScopeSupportsContractVariants(t *testing.T) {
	tests := []struct {
		name  string
		scope string
		mine  string
		want  bool
	}{
		{name: "recommended scope", scope: "my_topics", want: true},
		{name: "scope alias", scope: "mine", want: true},
		{name: "boolean alias", mine: "true", want: true},
		{name: "all", scope: "all", want: false},
		{name: "empty", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isMyTopicsScope(tt.scope, tt.mine); got != tt.want {
				t.Fatalf("isMyTopicsScope(%q, %q) = %v, want %v", tt.scope, tt.mine, got, tt.want)
			}
		})
	}
}
