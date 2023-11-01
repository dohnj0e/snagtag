package tiktok

import (
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

func Login(wd selenium.WebDriver) error {
	if err := wd.Get("https://www.tiktok.com/login/phone-or-email/email"); err != nil {
		return err
	}

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

	loginButton, err := wd.FindElement(selenium.ByCSSSelector, "button[type='submit']")
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
