package tiktok

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	Init()

	code := m.Run()
	os.Exit(code)
}

func TestScrape(t *testing.T) {
	keyword := "travel"
	err := Scrape(keyword)

	if err != nil {
		t.Errorf("Scrape failed: %v", err)
	}
}
