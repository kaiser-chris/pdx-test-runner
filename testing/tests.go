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

var regexTestDisplayName = regexp.MustCompile(`^### name\s*=\s*(?P<name>.*)`)
var regexLastDate = regexp.MustCompile(`last_date\s*=\s*(.*)`)
var regexTest = regexp.MustCompile(`(?m)(?:^\s*###\s*name\s*=\s*(?P<name>.+)\s*)*\s*(?:^\s*###\s*desc\s*=\s*(?P<desc>.+)\s*)*\s*(?P<test>[a-zA-Z_\-0-9]+)\s*=\s*{\s+(?:acceptable_fail_rate|success|fail)`)

// Base game files that should not be parsed
var baseIgnoreList = map[game.Type][]string{
	game.Victoria3: {
		"test.txt",
	},
	game.CrusaderKings3: {},
}

type PdxTestFile struct {
	Ignored     bool
	Name        string
	DisplayName string
	Path        string
	LastDate    string
	Tests       []*PdxTest
}
type PdxTest struct {
	Name        string
	DisplayName string
	Description string
}

func GetTestFiles(gamePath string, modPaths []string, gameType game.Type) ([]*PdxTestFile, error) {
	testFiles := make([]*PdxTestFile, 0)
	ignoreList := baseIgnoreList[gameType]

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

	lastDate := regexLastDate.FindStringSubmatch(string(content))
	displayName := regexTestDisplayName.FindStringSubmatch(string(content))

	for i, match := range matches {
		tests[i] = &PdxTest{
			Name:        match[3],
			DisplayName: match[1],
			Description: match[2],
		}
	}

	if len(tests) <= 0 {
		logging.Debugf("No tests found in file %s", file)
	}

	testFile := &PdxTestFile{
		Ignored: strings.HasSuffix(file, ignoreSuffix),
		Name:    filepath.Base(file),
		Path:    file,
		Tests:   tests,
	}
	if lastDate != nil {
		testFile.LastDate = lastDate[1]
	}
	if displayName != nil {
		testFile.DisplayName = displayName[1]
	}

	return testFile, nil
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
