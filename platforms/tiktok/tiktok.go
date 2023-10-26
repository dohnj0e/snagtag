package tiktok

import (
	"fmt"
	"os"
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
	cfg, err = config.LoadConfig("/path/to/project/config.yaml")

	if err != nil {
		logger.Log.Errorln("Failed to load config file: ", err)
	}

	seleniumPath = cfg.SeleniumPath
	chromeDriverPath = cfg.ChromeDriverPath
	port = cfg.Port
	searchURL = cfg.TiktokSearchURL
}

func waitForUser() {
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
		"chrome_binary":           "/usr/bin/chrome",
		"webdriver.chrome.driver": chromeDriverPath,
	}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))

	if err != nil {
		return nil, err
	}
	return wd, nil
}

func Login(wd selenium.WebDriver) error {
	if err := wd.Get("https://www.tiktok.com/login/phone-or-email/email"); err != nil {
		return err
	}

	usernameField, err := wd.FindElement(selenium.ByCSSSelector, "input[name='username']")
	if err != nil {
		return err
	}

	passwordField, err := wd.FindElement(selenium.ByCSSSelector, "input[placeholder='Password']")
	if err != nil {
		return err
	}

	username := os.Getenv("TIKTOK_USERNAME")
	if err := usernameField.SendKeys(username); err != nil {
		return err
	}

	password := os.Getenv("TIKTOK_PASSWORD")
	if err := passwordField.SendKeys(password); err != nil {
		return err
	}

	loginButton, err := wd.FindElement(selenium.ByCSSSelector, "button[data-e2e='login-button']")
	if err != nil {
		return err
	}

	if err := loginButton.Click(); err != nil {
		return err
	}

	waitForUser()

	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(500 * time.Millisecond)

	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			logger.Log.Errorln("Timed out waiting for login to complete")
		case <-ticker.C:
			currentURL, err := wd.CurrentURL()

			if err != nil {
				return err
			}
			if currentURL == "https://www.tiktok.com/foryou?lang=en" {
				return nil // login successful
			}
		}
	}
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

	err = wd.Get(searchURL + keyword)
	if err != nil {
		return err
	}

	waitForUser()

	elements, err := wd.FindElements(selenium.ByCSSSelector, "span.tiktok-j2a19r-SpanText")
	if err != nil {
		logger.Log.Errorln("Failed to find elements: ", err)
	}

	for index, element := range elements {
		title, err := element.Text()
		if err != nil {
			logger.Log.Errorln("Failed to retrieve element text: ", err)
		} else {
			if title != "" {
				fmt.Printf("%d: %s\n", index, title)
			}
		}
	}
	fmt.Printf("\n")
	logger.Log.Info("Scrape completed successfully")

	return nil
}
