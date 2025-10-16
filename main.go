package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"bahmut.de/pdx-test-runner/config"
	"bahmut.de/pdx-test-runner/game"
	"bahmut.de/pdx-test-runner/logging"
	"bahmut.de/pdx-test-runner/testing"
)

const (
	FlagConfig        = "config"
	FlagReportIgnored = "report-ignored"
)

func main() {
	configFlag := flag.String(FlagConfig, "test-config.json", "Optional: Path to test config")
	reportIgnored := flag.Bool(FlagReportIgnored, false, "Optional: Enable to list ignored tests")
	flag.Parse()

	configPath, err := filepath.Abs(*configFlag)
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
	testFiles, err := testing.GetTestFiles(settings.ContentPath, testConfig.ModDirectories, settings.GameType)
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

	logging.Info(buildFoundTestsReport(testFiles, reportIgnored != nil && *reportIgnored))

	startTime := time.Now()
	logging.Info("Start running tests")
	results, err := testing.RunTests(settings, testConfig, testFiles)
	if err != nil {
		logging.Fatalf("Could not run tests: %s", err)
		os.Exit(1)
	}
	endTime := time.Now()
	logging.Info("Finished running tests")
	logging.Infof("Running tests took: %v", endTime.Sub(startTime))

	absoluteOutputPath, err := filepath.Abs(results.OutputDirectory)
	if err != nil {
		logging.Infof("Test output: %s", results.OutputDirectory)
	} else {
		logging.Infof("Test output: %s", absoluteOutputPath)
	}
	logging.Info(buildRunTestsReport(results))

	logging.Info("Reactivating all test files")
	err = testing.ActivateTestFiles(testFiles)
	if err != nil {
		logging.Fatalf("Could not activate test files: %s", err)
		os.Exit(1)
	}

}

func buildFoundTestsReport(files []*testing.PdxTestFile, ignored bool) string {
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
		if testFile.Ignored && !ignored {
			continue
		}
		var color string
		if testFile.Ignored {
			color = logging.AnsiFgLightRed
		} else {
			color = logging.AnsiFgBlue
		}
		for _, test := range testFile.Tests {
			if strings.TrimSpace(test.DisplayName) != "" {
				report += fmt.Sprintf(
					"\n - %sTest:%s %s%s%s (%s)",
					logging.AnsiBoldOn, logging.AnsiAllDefault,
					color,
					test.DisplayName,
					logging.AnsiAllDefault,
					test.Name,
				)
			} else {
				report += fmt.Sprintf(
					"\n - %sTest:%s %s%s%s",
					logging.AnsiBoldOn, logging.AnsiAllDefault,
					color,
					test.Name,
					logging.AnsiAllDefault,
				)
			}
			if strings.TrimSpace(test.Description) != "" {
				report += fmt.Sprintf(" :: %s",
					test.Description,
				)
			}
			if strings.TrimSpace(testFile.DisplayName) != "" {
				report += fmt.Sprintf(
					" :: %s%s%s (%s)",
					logging.AnsiBoldOn,
					testFile.DisplayName,
					logging.AnsiAllDefault,
					testFile.Name,
				)
			} else {
				report += fmt.Sprintf(
					" :: %s%s%s",
					logging.AnsiBoldOn,
					testFile.Name,
					logging.AnsiAllDefault,
				)
			}
		}
	}
	return report
}

func buildRunTestsReport(results *testing.ExecutionResults) string {
	countSuccesses := 0
	countFailures := 0
	for _, testResult := range results.TestResults {
		if testResult.Success {
			countSuccesses++
		} else {
			countFailures++
		}
	}
	report := fmt.Sprintf(
		"There were %s%v%s %ssuccessful%s tests and %s%v%s %sfailed%s tests:",
		logging.AnsiBoldOn, countSuccesses, logging.AnsiAllDefault,
		logging.AnsiFgGreen, logging.AnsiAllDefault,
		logging.AnsiBoldOn, countFailures, logging.AnsiAllDefault,
		logging.AnsiFgLightRed, logging.AnsiAllDefault,
	)
	for _, testResult := range results.TestResults {
		if testResult.Success {
			report += fmt.Sprintf(
				"\n - %s%sSuccess:%s ",
				logging.AnsiBoldOn,
				logging.AnsiFgGreen,
				logging.AnsiAllDefault,
			)
		} else {
			report += fmt.Sprintf(
				"\n - %s%sFailure:%s ",
				logging.AnsiBoldOn,
				logging.AnsiFgLightRed,
				logging.AnsiAllDefault,
			)
		}
		if strings.TrimSpace(testResult.Test.DisplayName) != "" {
			report += fmt.Sprintf(
				"%s%s%s (%s)",
				logging.AnsiFgBlue,
				testResult.Test.DisplayName,
				logging.AnsiAllDefault,
				testResult.Test.Name,
			)
		} else {
			report += fmt.Sprintf(
				"%s%s%s",
				logging.AnsiFgBlue,
				testResult.Test.Name,
				logging.AnsiAllDefault,
			)
		}
		if strings.TrimSpace(testResult.Test.Description) != "" {
			report += fmt.Sprintf(" :: %s",
				testResult.Test.Description,
			)
		}
		if strings.TrimSpace(testResult.TestFile.DisplayName) != "" {
			report += fmt.Sprintf(
				" :: %s%s%s (%s)",
				logging.AnsiBoldOn,
				testResult.TestFile.DisplayName,
				logging.AnsiAllDefault,
				testResult.TestFile.Name,
			)
		} else {
			report += fmt.Sprintf(
				" :: %s%s%s",
				logging.AnsiBoldOn,
				testResult.TestFile.Name,
				logging.AnsiAllDefault,
			)
		}
	}

	return report
}
