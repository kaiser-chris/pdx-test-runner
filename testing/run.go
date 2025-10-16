package testing

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"bahmut.de/pdx-test-runner/config"
	"bahmut.de/pdx-test-runner/game"
	"bahmut.de/pdx-test-runner/logging"
)

const failTestPrefix = "TEST_FAIL_"
const resultFileName = "tests.txt"
const saveGameDirectoryName = "save games"

const testResultSuccess = "OK"

var regexTestResult = regexp.MustCompile(`(?m)^\[\s(OK|FAIL) ]\s(.*)\s\(\s(.*)\s\)`)

var saveGameSuffixes = map[game.Type]string{
	game.Victoria3:      ".v3",
	game.CrusaderKings3: ".ck3",
}

type ExecutionResults struct {
	OutputDirectory string
	TestResults     []*TestResult
}

type TestResult struct {
	Success  bool
	Date     string
	Test     *PdxTest
	TestFile *PdxTestFile
}

func RunTests(settings *game.LauncherSettings, config *config.TestRunnerConfig, testFiles []*PdxTestFile) (*ExecutionResults, error) {
	resultFile := filepath.Join(settings.DataPath, resultFileName)

	// Delete old test results
	err := deleteTestResults(resultFile)
	if err != nil {
		return nil, err
	}

	err = runGame(settings.ExecPath, resultFile)
	if err != nil {
		return nil, err
	}

	results, err := collectTestResults(resultFile, settings, config, testFiles)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func deleteTestResults(resultFile string) error {
	if _, err := os.Stat(resultFile); err == nil {
		err = os.Remove(resultFile)
		if err != nil {
			return fmt.Errorf("old test results could not be deleted: %v", err)
		}
	}
	return nil
}

func runGame(gameBinary, resultFile string) error {
	binary := exec.Command(gameBinary, "-nographics", "-handsoff", "-scripted_tests")
	err := binary.Start()
	if err != nil {
		return fmt.Errorf("error starting game: %v", err)
	}

	done := false
	for !done {
		time.Sleep(30 * time.Second)
		if _, err := os.Stat(resultFile); os.IsNotExist(err) {
			continue
		}
		content, err := os.ReadFile(resultFile)
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

func collectTestResults(resultFile string, settings *game.LauncherSettings, config *config.TestRunnerConfig, testFiles []*PdxTestFile) (*ExecutionResults, error) {
	saveDirectory := filepath.Join(settings.DataPath, saveGameDirectoryName)

	// create output directory
	runOutputDirectory := filepath.Join(config.OutputDirectory, time.Now().Format("2006-01-02_15_04_05"))
	if _, err := os.Stat(runOutputDirectory); os.IsNotExist(err) {
		err = os.MkdirAll(runOutputDirectory, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("error creating test result output directory: %v", err)
		}
	}
	if _, err := os.Stat(saveDirectory); os.IsNotExist(err) {
		return nil, fmt.Errorf("save game directory does not exist: %s", saveDirectory)
	}
	if _, err := os.Stat(resultFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("test result file does not exist: %s", resultFile)
	}

	// Move test result file to output directory
	content, err := os.ReadFile(resultFile)
	if err != nil {
		return nil, fmt.Errorf("could not read test result file: %v", err)
	}
	output := filepath.Join(runOutputDirectory, filepath.Base(resultFile))
	err = os.WriteFile(output, content, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("could not copy test result file to output directory: %v", err)
	}

	// Move test fail save games to output directory
	err = filepath.WalkDir(saveDirectory, func(file string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), saveGameSuffixes[settings.GameType]) {
			// Ignore non save game files
			return nil
		}
		if !strings.HasPrefix(info.Name(), failTestPrefix) {
			// We only care about test fail save games
			return nil
		}
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("could not read test result save game: %v", err)
		}
		output := filepath.Join(runOutputDirectory, info.Name())
		err = os.WriteFile(output, content, os.ModePerm)
		if err != nil {
			return fmt.Errorf("could not write test result save game to output directory: %v", err)
		}
		if config.MoveSaveGames {
			err = os.Remove(file)
			if err != nil {
				return fmt.Errorf("could not remove test result save game: %v", err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	testResults, err := parseTestResults(resultFile, testFiles)
	if err != nil {
		return nil, err
	}

	return &ExecutionResults{
		OutputDirectory: runOutputDirectory,
		TestResults:     testResults,
	}, nil
}

func parseTestResults(resultFile string, testFiles []*PdxTestFile) ([]*TestResult, error) {
	file, err := os.Open(resultFile)
	if err != nil {
		return nil, fmt.Errorf("could not open test results file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logging.Fatalf("could not close test results file: %v", err)
		}
	}(file)

	results := make([]*TestResult, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := regexTestResult.FindStringSubmatch(line)
		if matches == nil {
			continue
		}
		testFile, test := getTestFileAndTestByName(matches[2], testFiles)
		if testFile == nil || test == nil {
			logging.Errorf("Could not match test result (%s) to parsed tests: %s", matches[2], line)
			continue
		}
		testResult := &TestResult{}
		if matches[1] == testResultSuccess {
			testResult.Success = true
		}
		testResult.Test = test
		testResult.TestFile = testFile
		testResult.Date = matches[3]
		results = append(results, testResult)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("could not read test results file: %v", err)
	}

	return results, nil
}

func getTestFileAndTestByName(name string, testFiles []*PdxTestFile) (*PdxTestFile, *PdxTest) {
	for _, file := range testFiles {
		for _, test := range file.Tests {
			if test.Name == name {
				return file, test
			}
		}
	}
	return nil, nil
}
