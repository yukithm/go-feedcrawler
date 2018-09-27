package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/naoina/toml"
)

type Config struct {
	FeedsFile  string `toml:"feeds_file"`
	StatesFile string `toml:"states_file,omitempty"`
	NumWorkers int    `toml:"workers,omitempty"`
}

func loadConfig(name string) (*Config, error) {
	file := findConfigFile(name)
	if file == "" {
		return nil, fmt.Errorf("Cannot find config file: %s", name)
	}
	return loadConfigFile(file)
}

func loadConfigFile(file string) (*Config, error) {
	var config Config
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if err := toml.Unmarshal(buf, &config); err != nil {
		return nil, err
	}

	return &config, nil
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
	if dir, err := executableDir(); err == nil {
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

func executableDir() (string, error) {
	file, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(file), nil
}
