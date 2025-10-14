package game

import (
	"path/filepath"
)

func replaceDataPlaceholder(path string) (string, error) {
	return filepath.Clean(path), nil
}
