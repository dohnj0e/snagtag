package youtube

import (
	"fmt"
	"os"
	"path/filepath"

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
	cfg, err = config.LoadConfig("config.yaml")

	seleniumPath, err = filepath.Abs(cfg.SeleniumPath)
	if err != nil {
		logger.Log.Errorln("Failed to get absolute path for selenium: ", err)
	}

	chromeDriverPath, err = filepath.Abs(cfg.ChromeDriverPath)
	if err != nil {
		logger.Log.Errorln("Failed to get absolute path for chrome driver: ", err)
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
		"chrome_binary":           chromeDriverPath,
		"webdriver.chrome.driver": chromeDriverPath,
	}

	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))

	if err != nil {
		return nil, err
	}
	return wd, nil
}

func Scrape(keyword string) error {
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
