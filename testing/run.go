package testing

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"bahmut.de/pdx-test-runner/game"
)

func RunTests(settings *game.LauncherSettings) error {
	resultPath := filepath.Join(settings.DataPath, "tests.txt")
	if _, err := os.Stat(resultPath); err == nil {
		err = os.Remove(resultPath)
		if err != nil {
			return fmt.Errorf("old test results could not be deleted: %v", err)
		}
	}

	binary := exec.Command(settings.ExecPath, "-nographics", "-handsoff", "-scripted_tests")
	err := binary.Start()
	if err != nil {
		return fmt.Errorf("error starting game: %v", err)
	}

	done := false
	for !done {
		time.Sleep(30 * time.Second)
		if _, err := os.Stat(resultPath); os.IsNotExist(err) {
			continue
		}
		content, err := os.ReadFile(resultPath)
		if err != nil {
			return fmt.Errorf("could not check test results: %v", err)
		}
		// when tests are finished they are logged with [ OK ] or [ FAIL ]
		if strings.ContainsAny(string(content), "[]") {
			done = true
		}
	}

	err = binary.Process.Kill()
	if err != nil {
		return fmt.Errorf("error stopping game: %v", err)
	}

	return nil
}
