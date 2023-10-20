package youtube

import (
	"testing"
)

// test for Scrape function
func TestScrape(t *testing.T) {
	keyword := "dance"
	err := Scrape(keyword)

	if err != nil {
		t.Errorf("Scrape failed: %v", err)
	}
}
