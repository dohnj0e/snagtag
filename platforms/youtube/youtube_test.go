package youtube

import (
	"testing"
)

func TestScrape(t *testing.T) {
	keyword := "travel"
	err := Scrape(keyword)

	if err != nil {
		t.Errorf("Scrape failed: %v", err)
	}
}
