# SnagTag

Scrape social media platforms such as Youtube and Tiktok, powered by Go

[![Buy Me A Coffee](https://img.shields.io/badge/Buy%20Me%20A%20Coffee-Donate-yellow.svg)](https://www.buymeacoffee.com/dohnj0)
 [![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0.en.html) [![Tested on Ubuntu](https://img.shields.io/badge/Tested%20on-Ubuntu-orange.svg)](https://www.ubuntu.com/) [![Beta Version](https://img.shields.io/badge/Stable%20Version-1.1.0--stable-green.svg)](https://github.com/dohnj0e/snagtag/releases/tag/v1.1.0-stable)

## Table of Contents
  - [Table of Contents](#table-of-contents)
    - [Features](#features)
    - [Prerequisites](#prerequisites)
    - [Installation](#installation)
    - [Usage](#usage)
    - [Support](#support)
    - [Author](#author)

## Features
   - **Keyword:** Scrape video titles based on a provided keyword
   - **Hashtag:** Scrape video titles based on a provided hashtag ![new](https://img.shields.io/badge/New-gray)
   - **Scrolling:** Scrape even more data with infinite scrolling ![New](https://img.shields.io/badge/New-gray)
   - **Command-Line Interface:** Easy-to-use command-line interface
   - **Testing:** Test for youtube and tiktok scraping functionality
   - **Errors/Logging:** Robust error handling and logging system
   - **Authentication:** Handling authentication to Tiktok
   - **Captcha:** Waits user to interact with captcha before proceeding
   - **Config:** Configuration file to setup all you need

## Prerequisites
   - Go 1.18 or later
   - Selenium WebDriver
   - Firefox or Chrome Browser
   - Geckodriver (for Firefox) or Chromedriver (for Chrome)
   
## Installation
1. Clone the repository
   - `git clone https://github.com/dohnj0e/snagtag.git`
   - `cd snagtag`

2. Install Selenium WebDriver and browser driver
   
   a. Download Selenium Server Standalone:
      - [Selenium Server Standalone](https://www.selenium.dev/downloads/)
        
   b. Move the downloaded Selenium jar file to a directory:

      - `mv selenium-server-standalone-x.xx.x.jar /path/to/project/bin`
        
   c. Download the WebDriver for your browser:
      - [Geckodriver (for Firefox)](https://github.com/mozilla/geckodriver/releases)
      - [Chrome Driver (for Chrome)](https://sites.google.com/a/chromium.org/chromedriver/) **Recommended**
        
   d. Move the downloaded driver(s) to a directory:
      - `mv geckodriver /usr/local/bin`
      - `mv chromedriver /usr/local/bin`
  
3. Install all dependencies
   - `go mod tidy`

5. Build the project
   - `go build -o snagtag`

## Usage
   - `./snagtag platform youtube --keyword education`
   - `./snagtag platform tiktok --keyword education`

## ENV Variables (for TikTok)
  - `export TIKTOK_USERNAME='your_tiktok_username'`
  - `export TIKTOK_PASSWORD='your_tiktok_password'`

## Support

If you have any questions or need further assistance, feel free to reach out to me via email:

[💌 Contact Me](mailto:dohnj0@proton.me)
  
## Author
  - Dohn Joe (@dohnj0e)
