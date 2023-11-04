package config

import (
	"io/ioutil"
	"path/filepath"

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
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	config.SeleniumPath, err = filepath.Abs(config.SeleniumPath)
	if err != nil {
		return nil, err
	}

	config.ChromeDriverPath, err = filepath.Abs(config.ChromeDriverPath)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
