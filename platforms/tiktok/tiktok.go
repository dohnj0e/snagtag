package tiktok

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// create and send http request to tiktok
func FetchTiktokSearchResults(keyword string) (*http.Response, error) {
	url := "https://www.tiktok.com/search?q=" + keyword
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}
	return resp, nil
}

// parse html and extract data
func ParseTiktokSearchResults(resp *http.Response) error {
	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return err
	}

	// replace with correct css selector
	doc.Find(".video-title").Each(func(index int, element *goquery.Selection) {
		title := element.Text()
		fmt.Println(title)
	})
	return nil
}

// initialize scraping
func Scrape(keyword string) error {
	resp, err := FetchTiktokSearchResults(keyword)

	if err != nil {
		return err
	}

	// ensure body is closed after reading from it
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("Failed to close response body: %v\n", closeErr)
		}
	}()

	err = ParseTiktokSearchResults(resp)

	if err != nil {
		return err
	}
	return nil
}
