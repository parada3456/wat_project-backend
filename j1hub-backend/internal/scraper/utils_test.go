package scraper

import (
	"testing"
)

func TestGenerateJobID(t *testing.T) {
	url1 := "https://acadex.com/job/123"
	url2 := "https://acadex.com/job/123"
	url3 := "https://ihappy.com/job/456"

	id1 := GenerateJobID(url1)
	id2 := GenerateJobID(url2)
	id3 := GenerateJobID(url3)

	if id1 == "" {
		t.Error("GenerateJobID returned empty string")
	}

	if id1 != id2 {
		t.Errorf("GenerateJobID(%q) != GenerateJobID(%q); got %q and %q", url1, url2, id1, id2)
	}

	if id1 == id3 {
		t.Errorf("GenerateJobID(%q) == GenerateJobID(%q); both got %q", url1, url3, id1)
	}
}
