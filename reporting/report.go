package reporting

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"bahmut.de/pdx-test-runner/game"
	"bahmut.de/pdx-test-runner/testing"
)

func WriteReport(results *testing.ExecutionResults, testFiles []*testing.PdxTestFile, settings *game.LauncherSettings) error {
	builder := strings.Builder{}

	builder.WriteString("# Test Run - ")
	builder.WriteString(results.StartTime.Format(time.DateTime))
	builder.WriteString("\n")
	builder.WriteString("\n")
	builder.WriteString("## General")
	builder.WriteString("\n")
	builder.WriteString("\n")
	builder.WriteString("**Game:** ")
	switch settings.GameType {
	case game.Victoria3:
		builder.WriteString("Victoria 3\n")
	case game.CrusaderKings3:
		builder.WriteString("Crusader Kings 3\n")
	default:
		builder.WriteString("Unknown\n")
	}
	builder.WriteString("\n")
	builder.WriteString("**Start Time:** ")
	builder.WriteString(results.StartTime.Format(time.DateTime))
	builder.WriteString("\n")
	builder.WriteString("\n")
	builder.WriteString("**End Time:** ")
	builder.WriteString(results.EndTime.Format(time.DateTime))
	builder.WriteString("\n")
	builder.WriteString("\n")
	builder.WriteString("**Duration:** ")
	builder.WriteString(results.Duration.String())
	builder.WriteString("\n")
	builder.WriteString("\n")
	builder.WriteString("## Found Test Files & Tests\n\n")
	builder.WriteString("| Active | Test | Description | File |\n")
	builder.WriteString("|---|---|---|---|\n")
	for _, file := range testFiles {
		for _, test := range file.Tests {
			builder.WriteString("| ")
			if file.Ignored {
				builder.WriteString("❌")
			} else {
				builder.WriteString("✅")
			}
			builder.WriteString(" | ")
			if test.DisplayName != "" {
				builder.WriteString(test.DisplayName)
				builder.WriteString(" (")
				builder.WriteString(test.Name)
				builder.WriteString(")")
			} else {
				builder.WriteString(test.Name)
			}
			builder.WriteString(" | ")
			if test.Description != "" {
				builder.WriteString(test.Description)
			} else {
				builder.WriteString(" - ")
			}
			builder.WriteString(" | ")
			if file.DisplayName != "" {
				builder.WriteString(file.DisplayName)
				builder.WriteString(" (")
				builder.WriteString(file.Name)
				builder.WriteString(")")
			} else {
				builder.WriteString(file.Name)
			}
			builder.WriteString(" |\n")
		}
	}
	builder.WriteString("\n")
	builder.WriteString("## Test Results\n\n")
	builder.WriteString("| Success | Test | Date | Description | File |\n")
	builder.WriteString("|---|---|---|---|---|\n")
	for _, result := range results.TestResults {
		builder.WriteString("| ")
		if result.Success {
			builder.WriteString("✅")
		} else {
			builder.WriteString("❌")
		}
		builder.WriteString(" | ")
		if result.Test.DisplayName != "" {
			builder.WriteString(result.Test.DisplayName)
			builder.WriteString(" (")
			builder.WriteString(result.Test.Name)
			builder.WriteString(" )")
		} else {
			builder.WriteString(result.Test.Name)
		}
		builder.WriteString(" | ")
		builder.WriteString(result.Date)
		builder.WriteString(" | ")
		if result.Test.Description != "" {
			builder.WriteString(result.Test.Description)
		} else {
			builder.WriteString(" - ")
		}
		builder.WriteString(" | ")
		if result.TestFile.DisplayName != "" {
			builder.WriteString(result.TestFile.DisplayName)
			builder.WriteString(" (")
			builder.WriteString(result.TestFile.Name)
			builder.WriteString(" )")
		} else {
			builder.WriteString(result.TestFile.Name)
		}
		builder.WriteString(" |\n")
	}

	reportFile := filepath.Join(results.OutputDirectory, "report.md")
	err := os.WriteFile(reportFile, []byte(builder.String()), os.ModePerm)
	if err != nil {
		return fmt.Errorf("error writing report: %v", err)
	}

	return nil
}
