package config

import (
	"encoding/json"
	"fmt"
	"os"

	"bahmut.de/pdx-test-runner/logging"
)

type TestRunnerConfig struct {
	GameDirectory   string   `json:"game-directory"`
	ModDirectories  []string `json:"mod-directories"`
	OutputDirectory string   `json:"output-directory"`
	IgnoredFiles    []string `json:"ignored-files"`
	MoveSaveGames   bool     `json:"move-save-games"`
}

func LoadConfig(path string) (*TestRunnerConfig, error) {
	// Read config
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logging.Fatal(err)
		}
	}(file)

	// Decode json
	decoder := json.NewDecoder(file)
	var config TestRunnerConfig
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Fill optional output parameter
	if config.OutputDirectory == "" {
		config.OutputDirectory = "output"
	}

	return &config, nil
}
