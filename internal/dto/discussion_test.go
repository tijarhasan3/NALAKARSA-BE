package dto

import (
	"encoding/json"
	"testing"
)

func TestStringOrArrayAcceptsOptionalTags(t *testing.T) {
	for _, input := range []string{`{"tags":null}`, `{}`, `{"tags":[]}`} {
		var request CreateDiscussionRequest
		if err := json.Unmarshal([]byte(input), &request); err != nil {
			t.Fatalf("input %s should be accepted: %v", input, err)
		}
		if request.Tags != "" {
			t.Fatalf("input %s should produce empty tags, got %q", input, request.Tags)
		}
	}
}
