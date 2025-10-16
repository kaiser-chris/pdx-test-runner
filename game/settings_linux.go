package game

import (
	"path/filepath"
	"strings"
)

const userDocumentsVariable = "$LINUX_DATA_HOME"

func replaceDataPlaceholder(path string) (string, error) {
	result := strings.ReplaceAll(path, userDocumentsVariable, "~/.local/share")
	result = filepath.Clean(result)
	return result, nil
}
