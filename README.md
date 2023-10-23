# SnagTag [![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/dohnj0)

Tailored for scraping data from social media platforms, powered by Go.

[![Buy Me A Coffee](https://img.shields.io/badge/Buy%20Me%20A%20Coffee-Donate-yellow.svg)](https://www.buymeacoffee.com/dohnj0) [![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0.en.html) [![Tested on Ubuntu](https://img.shields.io/badge/Tested%20on-Ubuntu-orange.svg)](https://www.ubuntu.com/) [![Beta Version](https://img.shields.io/badge/Beta%20Version-0.1.0--beta-red.svg)](https://github.com/dohnj0/snagtag/releases/tag/v0.1.0-beta)

## Table of Contents
  - [Table of Contents](#table-of-contents)
    - [Features](#features)
    - [Prerequisites](#prerequisites)
    - [Installation](#installation)
    - [Usage](#usage)
    - [Support](#support)
    - [Author](#author)

## Features
   - **Scraping:** Scrape video titles based on a provided keyword
   - **Command-Line Interface:** Easy-to-use command-line interface
   - **Testing:** Test for youtube and tiktok scraping functionality
   - **Errors/Logging:** Robust error handling and logging system
   - **Authentication:** Handling authentication to Tiktok

## Prerequisites
   - Go 1.15 or later
   - Selenium WebDriver
   - Firefox or Chrome Browser
   - Geckodriver (for Firefox) or Chromedriver (for Chrome)
   
## Installation
   1. Clone the repository
      - `git clone https://github.com/dohnj0/snagtag.git`
      - `cd snagtag`
     
   2. Install selenium webdriver
      - Download the webdriver for your browser:
        - [Geckodriver (for Firefox)](https://github.com/mozilla/geckodriver/releases)
        - [Chrome Driver (for Chrome)](https://sites.google.com/a/chromium.org/chromedriver/) **Recommended**
       
      - Move the downloaded driver(s) to a directory, for example:
        - `mv geckodriver /usr/local/bin`
        - `mv selenium /path/to/project/bin`
          
   3. Build the project
       - `go build -o snagtag`

## Usage
   - `./snagtag platform youtube --keyword education` (for YouTube)

   - `./snagtag platform tiktok --keyword education` (for TikTok)

## Support

If you have any questions or need further assistance, feel free to reach out to me via email:

[💌 Contact Me](mailto:dohnj0@proton.me)
  
## Author
  - Dohn Joe (@dohnj0)
