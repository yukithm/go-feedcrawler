package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kardianos/osext"
	"github.com/naoina/toml"
	"github.com/yukithm/go-feedcrawler"
)

type Config struct {
	FeedCrawler struct {
		StateFile  string `toml:"state_file,omitempty"`
		NumWorkers int    `toml:"workers,omitempty"`
	} `toml:"feedcrawler"`
	Feed feedcrawler.Feeds
}

func loadConfig() (Config, error) {
	file := findConfigFile(configFile)
	if file == "" {
		return Config{}, fmt.Errorf("Cannot find config file: %s", configFile)
	}
	return loadConfigFile(file)
}

func loadConfigFile(file string) (Config, error) {
	var config Config
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return config, err
	}

	if err := toml.Unmarshal(buf, &config); err != nil {
		return config, err
	}

	return config, nil
}

func findConfigFile(filename string) string {
	// current directory
	if dir, err := os.Getwd(); err == nil {
		file := filepath.Join(dir, filename)
		if fileExists(file) {
			return file
		}
	}

	// executable directory
	if dir, err := osext.ExecutableFolder(); err == nil {
		file := filepath.Join(dir, filename)
		if fileExists(file) {
			return file
		}
	}

	return ""
}

func fileExists(file string) bool {
	if _, err := os.Stat(file); err != nil {
		return false
	}
	return true
}
