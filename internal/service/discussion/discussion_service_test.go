package discussionservice

import "testing"

func TestSplitTagsReturnsArrayForStoredCommaSeparatedTags(t *testing.T) {
	got := splitTags("iot, mesin produksi,  ")
	if len(got) != 2 || got[0] != "iot" || got[1] != "mesin produksi" {
		t.Fatalf("unexpected tags: %#v", got)
	}
}
