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

const (
	grayColor  = "\033[37m"
	cyanColor  = "\033[0;36m"
	resetColor = "\033[0m"
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
			(r >= 0x1FA00 && r <= 0x1FA6F) ||
			(r == 0x2665) {
			return true
		}
	}
	return false
}

func RemoveEmojis(s string) string {
	var result string
	var clusters []string

	gr := uniseg.NewGraphemes(s)
	for gr.Next() {
		clusters = append(clusters, gr.Str())
	}

	for i, cluster := range clusters {
		if cluster == "#" {
			if i < len(clusters)-1 && !strings.ContainsAny(clusters[i+1], "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789") {
				continue
			}
		} else if IsEmoji(cluster) {
			continue
		}
		result += cluster
	}
	return result
}

func PrintWithColor(text string, colorCode string) {
	fmt.Printf("%s%s%s", colorCode, text, resetColor)
}

func PrintTextWithHashtags(index int, text string) {
	words := strings.Fields(text)
	for i, word := range words {
		if strings.HasPrefix(word, "#") {
			PrintWithColor(word, cyanColor)
		} else {
			PrintWithColor(word, grayColor)
		}
		if i < len(words)-1 {
			fmt.Print(" ")
		}
	}
	fmt.Println()
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

	const MaxIndex = 125 // do not change this

	fmt.Printf("\n")
	logger.Log.Infoln("Initiating scraping for keyword:", keyword)

	err := wd.Get(searchURL)
	if err != nil {
		return err
	}

	time.Sleep(1 * time.Second)
	closeButton, err := wd.FindElement(selenium.ByCSSSelector, "button[aria-label='Accept the use of cookies and other data for the purposes described']")

	if err != nil {
		logger.Log.Warn("Failed to find the 'accept all' button: ", err)
		fmt.Printf("\n")
	} else {
		err = closeButton.Click()
		if err != nil {
			logger.Log.Error("Failed to click the 'accept all' button: ", err)
			fmt.Printf("\n")
			return err
		}
	}
	fmt.Printf("\n")

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

			if titleWithoutEmojis != "" &&
				!strings.Contains(title, "Mix") &&
				!strings.Contains(title, "Playlist") &&
				!strings.Contains(title, "Greatest Hits") &&
				!existingTitles[title] {
				PrintTextWithHashtags(index, fmt.Sprintf("%d: %s", index, titleWithoutEmojis))
				existingTitles[title] = true
			}
		}
	}
	return nil
}
