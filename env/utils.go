package env

import (
	"os"
	"path/filepath"
	"strings"
)

// IsInPath checks if the directory of the executable is in the PATH environment variable
func IsInPath(executablePath string) bool {
	executableDir := filepath.Dir(executablePath)
	pathEnv := os.Getenv("PATH")
	for _, dir := range filepath.SplitList(pathEnv) {
		if dir == executableDir {
			return true
		}
	}
	return false
}

// ExpandTilde expands the tilde (~) in the given path to the user's home directory.
func ExpandTilde(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, path[1:]), nil
	}

	return path, nil
}
