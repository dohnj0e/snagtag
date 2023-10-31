package youtube

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/dohnj0e/snagtag/logger"
	"github.com/rivo/uniseg"
	"github.com/tebeka/selenium"
)

func IsEmoji(s string) bool {
	for _, r := range s {
		if (r >= 0x1F600 && r <= 0x1F64F) ||
			(r >= 0x1F300 && r <= 0x1F5FF) ||
			(r >= 0x1F680 && r <= 0x1F6FF) ||
			(r >= 0x1F700 && r <= 0x1F77F) ||
			(r >= 0x1F780 && r <= 0x1F7FF) ||
			(r >= 0x1F800 && r <= 0x1F8FF) ||
			(r >= 0x1F900 && r <= 0x1F9FF) ||
			(r >= 0x1FA00 && r <= 0x1FA6F) {
			return true
		}
	}
	return false
}

func RemoveEmojis(s string) string {
	var result string
	gr := uniseg.NewGraphemes(s)
	for gr.Next() {
		cluster := gr.Str()
		if !IsEmoji(cluster) {
			result += cluster
		}
	}
	return result
}

func ScrollIncrementally(wd selenium.WebDriver, amount int) error {
	script := fmt.Sprintf("window.scroll(0, %d);", amount)
	_, err := wd.ExecuteScript(script, nil)
	return err
}

func ScrollAndScrape(wd selenium.WebDriver, keyword string) error {
	existingTitles := map[string]bool{}
	encodedKeyword := url.QueryEscape(keyword)
	searchURL = fmt.Sprintf("https://www.youtube.com/results?search_query=%s", encodedKeyword)

	const MaxIndex = 125

	fmt.Printf("\n")
	logger.Log.Infoln("Initiating scraping for keyword:", keyword)

	err := wd.Get(searchURL)
	if err != nil {
		return err
	}

	fmt.Printf("\n")
	closeButton, err := wd.FindElement(selenium.ByCSSSelector, "button.yt-spec-button-shape-next.yt-spec-button-shape-next--filled.yt-spec-button-shape-next--mono.yt-spec-button-shape-next--size-m")

	if err != nil {
		logger.Log.Warn("Failed to find the 'accept all' button: ", err)
	} else {
		err = closeButton.Click()
		if err != nil {
			logger.Log.Error("Failed to click the 'accept all' button: ", err)
			return err
		}
	}

	for {
		time.Sleep(3 * time.Second)

		err = ScrollIncrementally(wd, 500)
		if err != nil {
			return err
		}

		prevScrollPos, err := wd.ExecuteScript("return window.pageYOffset;", nil)
		if err != nil {
			return err
		}

		_, err = wd.ExecuteScript("window.scroll(0, 45000);", nil)
		if err != nil {
			logger.Log.Errorln("Failed to scroll: ", err)
			return err
		}

		time.Sleep(5 * time.Second)

		currScrollPos, err := wd.ExecuteScript("return window.pageYOffset;", nil)
		if err != nil {
			return err
		}

		if prevScrollPos == currScrollPos {
			break
		}

		elements, err := wd.FindElements(selenium.ByCSSSelector, "a#video-title")
		if err != nil {
			logger.Log.Errorln("Failed to find elements: ", err)
			return err
		}

		for index, element := range elements {
			if index >= MaxIndex {
				fmt.Printf("\n")
				logger.Log.Info("Scrape completed successfully")
				return nil
			}

			title, err := element.Text()
			titleWithoutEmojis := RemoveEmojis(title)

			if err != nil {
				logger.Log.Errorln("Failed to retrieve element text: ", err)
				continue
			}

			if title != "" &&
				!strings.Contains(title, "Mix") &&
				!strings.Contains(title, "Playlist") &&
				!strings.Contains(title, "Greatest Hits") &&
				!existingTitles[title] {
				fmt.Printf("%d: %s\n", index, titleWithoutEmojis)
				existingTitles[title] = true
			}
		}
	}
	return nil
}
