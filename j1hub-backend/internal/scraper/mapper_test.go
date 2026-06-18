package scraper

import (
	"testing"
)

func TestMapLocationRank(t *testing.T) {
	tests := []struct {
		source     string
		sourceType string
		want       string
	}{
		{"A", "Acadex", "Rank S"},
		{"X", "Acadex", "Rank A"},
		{"Y", "Acadex", "Rank B"},
		{"Z", "Acadex", "Rank C"},
		{"unknown", "Acadex", "Rank D"},
		{"signature", "IHappy", "Rank S"},
		{"prestige", "IHappy", "Rank A"},
		{"super premium", "IHappy", "Rank B"},
		{"premium", "IHappy", "Rank C"},
		{"promotion", "IHappy", "Rank D"},
		{"other", "IHappy", "Rank D"},
	}

	for _, tt := range tests {
		t.Run(tt.sourceType+"_"+tt.source, func(t *testing.T) {
			got := MapLocationRank(tt.source, tt.sourceType)
			if got != tt.want {
				t.Errorf("MapLocationRank(%q, %q) = %q; want %q", tt.source, tt.sourceType, got, tt.want)
			}
		})
	}
}
