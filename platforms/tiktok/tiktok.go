package tiktok

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tebeka/selenium"
)

var (
	logger  = log.New(os.Stdout, "DEBUG: ", log.Lshortfile)
	service *selenium.Service
)

func InitWebDriver() (selenium.WebDriver, error) {
	var err error

	const (
		seleniumPath     = "/home/ian/Documents/DEV/projects/go/snagtag/bin/selenium-server-standalone-3.141.59.jar" // absolute path to selenium
		chromeDriverPath = "/home/ian/Documents/DEV/projects/go/snagtag/bin/chromedriver"                            // absolute path to chromedriver (chrome)
	)

	// add geckodriver path to system PATH
	os.Setenv("PATH", os.Getenv("PATH")+":"+chromeDriverPath)

	// port number
	port := 4445

	opts := []selenium.ServiceOption{
		selenium.GeckoDriver(chromeDriverPath),
		selenium.Output(os.Stderr),
	}
	service, err = selenium.NewSeleniumService(seleniumPath, port, opts...)

	if err != nil {
		fmt.Printf("Error starting the Selenium service: %v", err)
		return nil, err
	}

	caps := selenium.Capabilities{
		"browserName":            "chrome",
		"chrome_binary":          "/usr/bin/chrome",
		"webdriver.gecko.driver": chromeDriverPath,
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

	// tiktok credentials
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

	// wait for url to change
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return errors.New("timed out waiting for login to complete")
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

	err = wd.Get("https://www.tiktok.com/search/video?q=" + keyword)
	if err != nil {
		return err
	}

	elements, err := wd.FindElements(selenium.ByCSSSelector, "span.tiktok-j2a19r-SpanText")
	if err != nil {
		return err
	}

	for index, element := range elements {
		title, err := element.Text()
		if err != nil {
			logger.Println("Failed to retrieve element text:", err)
		} else {
			fmt.Printf("Video title %d: %s\n", index, title)
		}
	}
	return nil
}
