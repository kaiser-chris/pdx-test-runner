package game

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"path/filepath"
	"strings"
)

const userDocumentsVariable = "%USER_DOCUMENTS%"

func replaceDataPlaceholder(path string) (string, error) {
	if strings.Contains(path, userDocumentsVariable) {
		documentsFolderKey, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\User Shell Folders`, registry.QUERY_VALUE)
		if err != nil {
			return "", fmt.Errorf("could not query registry for user directory (%s): %v", userDocumentsVariable, err)
		}
		documentsFolder, _, err := documentsFolderKey.GetStringValue("Personal")
		if err != nil {
			return "", fmt.Errorf("could not query registry for user directory (%s): %v", userDocumentsVariable, err)
		}
		path = strings.ReplaceAll(path, userDocumentsVariable, documentsFolder)
	}
	return filepath.Clean(path), nil
}
