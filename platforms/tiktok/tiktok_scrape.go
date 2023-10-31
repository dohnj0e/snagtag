package tiktok

import (
	"fmt"
	"net/url"
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
	searchURL = fmt.Sprintf("https://www.tiktok.com/search/video?q=%s", encodedKeyword)

	const MaxIndex = 125 // do not change this

	logger.Log.Infoln("Initiating scraping for keyword:", keyword)
	fmt.Printf("\n")

	err := wd.Get(searchURL)
	if err != nil {
		return err
	}

	WaitForUser()

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

		time.Sleep(3 * time.Second)

		currScrollPos, err := wd.ExecuteScript("return window.pageYOffset;", nil)
		if err != nil {
			return err
		}

		if prevScrollPos == currScrollPos {
			break
		}

		elements, err := wd.FindElements(selenium.ByCSSSelector, "div.tiktok-1iy6zew-DivContainer")
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

			if title != "" && !existingTitles[title] {
				fmt.Printf("%d: %s\n", index, titleWithoutEmojis)
				existingTitles[title] = true
			}
		}
	}
	return nil
}
