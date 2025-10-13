package main

import (
	"bahmut.de/pdx-test-runner/config"
	"bahmut.de/pdx-test-runner/game"
	"bahmut.de/pdx-test-runner/logging"
	"bahmut.de/pdx-test-runner/testing"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	FlagConfig = "config"
)

func main() {
	configFlag := flag.String(FlagConfig, "test-config.json", "Optional: Path to test config")
	flag.Parse()

	var configPath string
	if configFlag == nil || *configFlag == "" {
		configPath = "test-config.json"
	} else {
		configPath = *configFlag
	}
	configPath, err := filepath.Abs(configPath)
	if err != nil {
		logging.Fatalf("Provided config file path is invalid: %s", err)
		os.Exit(1)
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		logging.Fatalf("Config file does not exist: %s", configPath)
		os.Exit(1)
	}

	logging.Info("Loading Runner Config")
	testConfig, err := config.LoadConfig(configPath)
	if err != nil {
		logging.Fatalf("Could not load config file: %s", err)
		os.Exit(1)
	}

	logging.Info("Loading Game Settings")
	settings, err := game.GetLauncherSettings(testConfig.GameDirectory)
	if err != nil {
		logging.Fatalf("Could not load game launcher settings: %s", err)
		os.Exit(1)
	}

	logging.Info("Reading Tests")
	testFiles, err := testing.GetTestFiles(settings.ContentPath, testConfig.ModDirectories, game.Victoria3)
	if err != nil {
		logging.Fatalf("Could not parse tests: %s", err)
		os.Exit(1)
	}

	logging.Info("Deactivating ignored test files")
	err = testing.DeactivateTestFiles(testFiles, testConfig.IgnoredFiles)
	if err != nil {
		logging.Fatalf("Could not deactivate ignored test files: %s", err)
		os.Exit(1)
	}

	logging.Info(buildFoundTestsReport(testFiles, true))
	logging.Info(buildFoundTestsReport(testFiles, false))

	logging.Info("Reactivating all test files")
	err = testing.ActivateTestFiles(testFiles)
	if err != nil {
		logging.Fatalf("Could not activate test files: %s", err)
		os.Exit(1)
	}

}

func buildFoundTestsReport(files []*testing.PdxTestFile, ignored bool) string {
	var color string
	if ignored {
		color = logging.AnsiFgLightRed
	} else {
		color = logging.AnsiFgBlue
	}

	countFiles := 0
	countTests := 0
	for _, testFile := range files {
		if testFile.Ignored == !ignored {
			continue
		}
		countFiles++
		countTests += len(testFile.Tests)
	}

	var report string
	if ignored {
		report = "Ignored"
	} else {
		report = "Found"
	}
	report += fmt.Sprintf(
		" %s%v%s Tests in %s%v%s Files:",
		logging.AnsiBoldOn, countTests, logging.AnsiAllDefault,
		logging.AnsiBoldOn, countFiles, logging.AnsiAllDefault,
	)

	for _, testFile := range files {
		if testFile.Ignored == !ignored {
			continue
		}
		report += fmt.Sprintf(
			"\n%sFile:%s %s%s%s",
			logging.AnsiBoldOn, logging.AnsiAllDefault,
			color,
			testFile.Name,
			logging.AnsiAllDefault,
		)
		for _, test := range testFile.Tests {
			report += fmt.Sprintf(
				"\n - %sTest:%s %s%s%s",
				logging.AnsiBoldOn, logging.AnsiAllDefault,
				color,
				test.Name,
				logging.AnsiAllDefault,
			)
			if strings.TrimSpace(test.Comment) != "" {
				report += fmt.Sprintf(" :: %s%s%s",
					logging.AnsiFgGreen,
					test.Comment,
					logging.AnsiAllDefault,
				)
			}
		}
	}
	return report
}
