package youtube

import (
	"fmt"
	"log"
	"os"

	"github.com/tebeka/selenium"
)

var (
	logger  = log.New(os.Stdout, "DEBUG: ", log.Lshortfile)
	service *selenium.Service
)

func InitWebDriver() (selenium.WebDriver, error) {
	var err error

	const (
		seleniumPath    = "/path/to/project/bin/selenium-server-standalone-3.141.59.jar" // absolute path to selenium
		geckoDriverPath = "/path/to/project/bin/geckodriver"                             // absolute path to geckodriver (firefox)
	)

	// add geckodriver path to the system PATH
	os.Setenv("PATH", os.Getenv("PATH")+":"+geckoDriverPath)

	// port number
	port := 4444

	opts := []selenium.ServiceOption{
		selenium.GeckoDriver(geckoDriverPath),
		selenium.Output(os.Stderr),
	}
	service, err = selenium.NewSeleniumService(seleniumPath, port, opts...)

	if err != nil {
		fmt.Printf("Error starting the Selenium service: %v", err)
		return nil, err
	}

	caps := selenium.Capabilities{
		"browserName":            "firefox",
		"firefox_binary":         "/usr/bin/firefox",
		"webdriver.gecko.driver": geckoDriverPath,
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

	err = wd.Get("https://youtube.com/results?search_query=" + keyword)
	if err != nil {
		return err
	}

	elements, err := wd.FindElements(selenium.ByCSSSelector, "span#video-title")
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
