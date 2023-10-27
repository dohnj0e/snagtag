package youtube

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/dohnj0e/snagtag/config"
	"github.com/dohnj0e/snagtag/logger"
	"github.com/tebeka/selenium"
)

var (
	cfg     *config.Config
	service *selenium.Service
	err     error

	seleniumPath     string
	chromeDriverPath string
	port             int
	searchURL        string
)

func Init() {
	cfg, err = config.LoadConfig("/path/to/project/config.yaml") // absolute path

	if err != nil {
		logger.Log.Errorln("Failed to load config file: ", err)
	}

	seleniumPath = cfg.SeleniumPath
	chromeDriverPath = cfg.ChromeDriverPath
	port = cfg.Port
	searchURL = cfg.YoutubeSearchURL
}

func InitWebDriver() (selenium.WebDriver, error) {
	os.Setenv("PATH", os.Getenv("PATH")+":"+chromeDriverPath)

	opts := []selenium.ServiceOption{}
	service, err = selenium.NewSeleniumService(seleniumPath, port, opts...)

	if err != nil {
		logger.Log.Errorln("Error starting the Selenium service: ", err)
		return nil, err
	}

	caps := selenium.Capabilities{
		"browserName":             "chrome",
		"chrome_binary":           "/usr/bin/chrome",
		"webdriver.chrome.driver": chromeDriverPath,
	}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))

	if err != nil {
		return nil, err
	}
	return wd, nil
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

	const MaxIndex = 100

	err := wd.Get(searchURL)
	if err != nil {
		return err
	}

	closeButton, err := wd.FindElement(selenium.ByCSSSelector, "button.yt-spec-button-shape-next.yt-spec-button-shape-next--filled.yt-spec-button-shape-next--mono.yt-spec-button-shape-next--size-m")
	if err != nil {
		logger.Log.Error("Failed to find the 'accept all' button: ", err)
		return err
	}

	err = closeButton.Click()
	if err != nil {
		logger.Log.Error("Failed to click the 'accept all' button: ", err)
		return err
	}

	for {
		time.Sleep(1 * time.Second)

		err = ScrollIncrementally(wd, 1000)
		if err != nil {
			return err
		}

		prevScrollPos, err := wd.ExecuteScript("return window.pageYOffset;", nil)
		if err != nil {
			return err
		}

		_, err = wd.ExecuteScript("window.scroll(0, 25000);", nil)
		if err != nil {
			logger.Log.Errorln("Failed to scroll: ", err)
			return err
		}

		time.Sleep(1 * time.Second)

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
			if err != nil {
				logger.Log.Errorln("Failed to retrieve element text: ", err)
				continue
			}

			if title != "" &&
				!strings.Contains(title, "Mix") &&
				!strings.Contains(title, "Playlist") &&
				!strings.Contains(title, "Greatest Hits") &&
				!existingTitles[title] {
				fmt.Printf("%d: %s\n", index, title)
				existingTitles[title] = true
			}
		}
	}
	return nil
}

func Scrape(keyword string) error {
	encodedKeyword := url.QueryEscape(keyword)
	searchURL = fmt.Sprintf("https://www.youtube.com/results?search_query=%s", encodedKeyword)

	wd, err := InitWebDriver()
	if err != nil {
		return err
	}

	defer wd.Quit()
	defer service.Stop()

	err = ScrollAndScrape(wd, keyword)
	if err != nil {
		return err
	}

	return nil
}
