package testing

import (
	"fmt"
	"os"
	"strings"
)

const ignoreSuffix = ".ignore"

func DeactivateTestFiles(testFiles []*PdxTestFile, ignoreFiles []string) error {
	for _, testFile := range testFiles {
		deactivated := false
		for _, ignoreFile := range ignoreFiles {
			if ignoreFile == testFile.Name {
				err := deactivateTestFile(testFile)
				if err != nil {
					return err
				}
				deactivated = true
				break
			}
		}
		if !deactivated {
			err := activateTestFile(testFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ActivateTestFiles(testFiles []*PdxTestFile) error {
	for _, testFile := range testFiles {
		err := activateTestFile(testFile)
		if err != nil {
			return err
		}
	}
	return nil
}

func deactivateTestFile(file *PdxTestFile) error {
	if !strings.HasSuffix(file.Path, ignoreSuffix) {
		deactivatedName := file.Path + ignoreSuffix
		err := os.Rename(file.Path, deactivatedName)
		if err != nil {
			return fmt.Errorf("could not deactivate test file (%s): %v", file, err)
		}
		file.Path = deactivatedName
		file.Ignored = true
	}
	return nil
}

func activateTestFile(file *PdxTestFile) error {
	if strings.HasSuffix(file.Path, ignoreSuffix) {
		activatedName, _ := strings.CutSuffix(file.Path, ignoreSuffix)
		err := os.Rename(file.Path, activatedName)
		if err != nil {
			return fmt.Errorf("could not activate test file (%s): %v", file, err)
		}
		file.Path = activatedName
		file.Ignored = false
	}
	return nil
}
