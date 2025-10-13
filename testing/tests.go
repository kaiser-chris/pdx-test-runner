package testing

import (
	"bahmut.de/pdx-test-runner/game"
	"bahmut.de/pdx-test-runner/logging"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const activeTestSuffix = ".txt"

var regexLastDate = regexp.MustCompile(`^last_date\s*=\s*(.*)$`)
var regexTest = regexp.MustCompile(`(?m)(?:^\s*#\s*(?P<comment>.+)\s*)*\s*(?P<test>[a-zA-Z_\-0-9]+)\s*=\s*{\s+(?:acceptable_fail_rate|success|fail)`)

var victoria3IgnoreList = []string{
	"test.txt",
}

type PdxTestFile struct {
	Ignored  bool
	Name     string
	Path     string
	LastDate string
	Tests    []*PdxTest
}
type PdxTest struct {
	Name    string
	Comment string
}

func GetTestFiles(gamePath string, modPaths []string, gameType game.Type) ([]*PdxTestFile, error) {
	testFiles := make([]*PdxTestFile, 0)
	ignoreList := getIgnoreList(gameType)

	// Add base game tests
	baseGameTests, err := parseTestDirectory(filepath.Join(gamePath, "tools", "scripted_tests"), ignoreList)
	if err != nil {
		return nil, err
	}
	testFiles = append(testFiles, baseGameTests...)

	// Add mod tests
	for _, modPath := range modPaths {
		modTests, err := parseTestDirectory(filepath.Join(modPath, "tools", "scripted_tests"), ignoreList)
		if err != nil {
			return nil, err
		}
		testFiles = mergeTestFiles(testFiles, modTests)
	}

	// Remove empty tests
	results := make([]*PdxTestFile, 0)
	for _, testFile := range testFiles {
		if len(testFile.Tests) == 0 {
			continue
		}
		results = append(results, testFile)
	}

	return results, nil
}

func parseTestDirectory(directory string, gameIgnoreList []string) ([]*PdxTestFile, error) {
	tests := make([]*PdxTestFile, 0)
	err := filepath.WalkDir(directory, func(filepath string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), activeTestSuffix) && !strings.HasSuffix(info.Name(), ignoreSuffix) {
			return nil
		}
		for _, ignore := range gameIgnoreList {
			if strings.HasSuffix(filepath, ignore) {
				return nil
			}
		}
		testsInFile, err := parseTestFile(filepath)
		if err != nil {
			return err
		}
		tests = append(tests, testsInFile)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tests, nil
}

func parseTestFile(file string) (*PdxTestFile, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	matches := regexTest.FindAllStringSubmatch(string(content), -1)
	tests := make([]*PdxTest, len(matches))

	lastDate := regexLastDate.FindString(string(content))

	for i, match := range matches {
		tests[i] = &PdxTest{
			Name:    match[2],
			Comment: match[1],
		}
	}

	if len(tests) <= 0 {
		logging.Debugf("No tests found in file %s", file)
	}

	testFile := &PdxTestFile{
		Ignored:  strings.Contains(file, ignoreSuffix),
		Name:     filepath.Base(file),
		Path:     file,
		LastDate: lastDate,
		Tests:    tests,
	}
	return testFile, nil
}

func getIgnoreList(gameType game.Type) []string {
	switch gameType {
	case game.Victoria3:
		return victoria3IgnoreList
	default:
		return make([]string, 0)
	}
}

func mergeTestFiles(existingTests []*PdxTestFile, newTests []*PdxTestFile) []*PdxTestFile {
	mergedFiles := make([]*PdxTestFile, len(existingTests))
	copy(mergedFiles, existingTests)
	for _, newFile := range newTests {
		found := false
		for exitingIndex, existingFile := range existingTests {
			if existingFile.Name == newFile.Name {
				mergedFiles[exitingIndex] = newFile
				found = true
				break
			}
		}
		if !found {
			mergedFiles = append(mergedFiles, newFile)
		}
	}
	return mergedFiles
}
