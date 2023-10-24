package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	SeleniumPath     string `yaml:"seleniumPath"`
	ChromeDriverPath string `yaml:"chromeDriverPath"`
	Port             int    `yaml:"port"`
	YoutubeSearchURL string `yaml:"youtubeSearchURL"`
	TiktokSearchURL  string `yaml:"tiktokSearchURL"`
}

var config Config

func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &config)

	if err != nil {
		return nil, err
	}
	return &config, nil
}
