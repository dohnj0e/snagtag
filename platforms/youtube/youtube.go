package youtube

import (
	"fmt"
	"os"
	"strings"

	"github.com/dohnj0e/snagtag/config"
	"github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
)

var (
	cfg     *config.Config
	logger  *logrus.Logger
	service *selenium.Service
	err     error

	seleniumPath     string
	chromeDriverPath string
	port             int
	searchURL        string
)

func Init() {
	logger = logrus.New()
	logger.Out = os.Stdout
	logger.Level = logrus.InfoLevel
	logger.Formatter = &logrus.TextFormatter{
		DisableTimestamp: true,
	}

	cfg, err = config.LoadConfig("config.yaml")

	if err != nil {
		panic(err)
	}

	seleniumPath = cfg.SeleniumPath
	chromeDriverPath = cfg.ChromeDriverPath
	port = cfg.Port
	searchURL = cfg.YoutubeSearchURL
}

func InitWebDriver() (selenium.WebDriver, error) {
	os.Setenv("PATH", os.Getenv("PATH")+":"+chromeDriverPath)

	opts := []selenium.ServiceOption{
		selenium.GeckoDriver(chromeDriverPath),
	}
	service, err = selenium.NewSeleniumService(seleniumPath, port, opts...)

	if err != nil {
		fmt.Printf("Error starting the Selenium service: %v", err)
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

func Scrape(keyword string) error {
	wd, err := InitWebDriver()

	if err != nil {
		return err
	}
	defer wd.Quit()
	defer service.Stop()

	err = wd.Get(searchURL + keyword)
	if err != nil {
		return err
	}

	elements, err := wd.FindElements(selenium.ByCSSSelector, "a#video-title")
	if err != nil {
		return err
	}

	for index, element := range elements {
		title, err := element.Text()
		if err != nil {
			logger.Error("Failed to retrieve element text:", err)
		} else {
			if title != "" && !strings.Contains(title, "Mix -") {
				fmt.Printf("%d: %s\n", index, title)
			}
		}
	}
	fmt.Printf("\n")
	logger.Info("Scrape completed successfully")
	return nil
}
