package youtube

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// create and send http request to youtube
func FetchYoutubeSearchResults(keyword string) (*http.Response, error) {
	url := "https://youtube.com/results?search_query=" + keyword
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}
	return resp, nil
}

// print html content of response
func PrintHTML(resp *http.Response) error {
	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)

	return nil
}

// parse html and extract data
func ParseYoutubeSearchResults(resp *http.Response) error {
	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return err
	}

	doc.Find(".video-title").Each(func(index int, element *goquery.Selection) {
		title := element.Text()
		fmt.Println(title)
	})
	return nil
}

// initialize scraping
func Scrape(keyword string) error {
	resp, err := FetchYoutubeSearchResults(keyword)

	if err != nil {
		return err
	}

	// ensure body is closed after reading from it
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("Failed to close response body: %v\n", closeErr)
		}
	}()

	err = PrintHTML(resp)

	if err != nil {
		return err
	}
	return ParseYoutubeSearchResults(resp)
}
