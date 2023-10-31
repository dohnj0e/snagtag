package tiktok

import (
	"fmt"
	"os"

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
	searchURL = cfg.TiktokSearchURL
}

func WaitForUser() {
	logger.Log.Warnln("Please solve the captcha, then press Enter to continue...")
	var input string
	fmt.Scanln(&input)
}

func InitWebDriver() (selenium.WebDriver, error) {
	os.Setenv("PATH", os.Getenv("PATH")+":"+chromeDriverPath)

	opts := []selenium.ServiceOption{}
	service, err = selenium.NewSeleniumService(seleniumPath, port, opts...)

	if err != nil {
		logger.Log.Errorln("Error starting the Selenium service: ", err)
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

	if err := Login(wd); err != nil {
		return err
	}

	err = ScrollAndScrape(wd, keyword)
	if err != nil {
		return err
	}
	return nil
}
