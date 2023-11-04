package tiktok

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dohnj0e/snagtag/logger"
	"github.com/tebeka/selenium"
)

func WaitForUser() {
	logger.Log.Warnln("Please solve the captcha, then press Enter to continue...")
	var input string
	fmt.Scanln(&input)
}

func WaitForElement(wd selenium.WebDriver, selector string, timeout time.Duration) (selenium.WebElement, error) {
	var element selenium.WebElement
	var err error

	endTime := time.Now().Add(timeout)

	for time.Now().Before(endTime) {
		element, err = wd.FindElement(selenium.ByCSSSelector, selector)

		if err == nil {
			return element, nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return nil, errors.New("timed out waiting for element: " + selector)
}

func Login(wd selenium.WebDriver) error {
	if err := wd.Get("https://www.tiktok.com/login/phone-or-email/email"); err != nil {
		return err
	}

	time.Sleep(2 * time.Second)

	usernameField, err := wd.FindElement(selenium.ByCSSSelector, "input[name='username']")
	if err != nil {
		return err
	}

	passwordField, err := wd.FindElement(selenium.ByCSSSelector, "input[type='password']")
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

	time.Sleep(2 * time.Second)

	loginButton, err := WaitForElement(wd, "button[type='submit']", 10*time.Second)
	if err != nil {
		return err
	}

	if err := loginButton.Click(); err != nil {
		return err
	}

	fmt.Printf("\n")
	WaitForUser()

	timeout := time.After(60 * time.Second)
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
				return nil
			}
		}
	}
}
