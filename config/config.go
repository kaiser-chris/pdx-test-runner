package config

import (
	"bahmut.de/pdx-test-runner/logging"
	"encoding/json"
	"os"
)

type TestConfig struct {
	GameDirectory  string   `json:"game-directory"`
	ModDirectories []string `json:"mod-directories"`
	IgnoredFiles   []string `json:"ignored-files"`
}

func LoadConfig(path string) (*TestConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logging.Fatal(err)
		}
	}(file)
	decoder := json.NewDecoder(file)
	var config TestConfig
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
