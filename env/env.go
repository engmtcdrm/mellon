package env

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/engmtcdrm/minno/app"
)

type Env struct {
	Home       string // Home is the user's home directory.
	AppHomeDir string // AppHomeDir is the directory in the user's home directory where the app stores its data.
	KeyPath    string // KeyPath is the path to the key file.
	ExeCmd     string // ExeCmd is the command to run the executable. If the executable is in the PATH environment variable, this will be the executable name.
}

var (
	instance *Env
	once     sync.Once
)

// GetEnv returns the singleton instance of Env.
func GetEnv() (*Env, error) {
	var err error

	once.Do(func() {
		home, e := os.UserHomeDir()
		if e != nil {
			err = e
			return
		}

		instance = &Env{
			Home:       home,
			AppHomeDir: filepath.Join(home, app.DotName),
		}

		executablePath, e := os.Executable()
		if e != nil {
			err = e
			return
		}

		executableName := filepath.Base(executablePath)

		if IsInPath(executablePath) {
			instance.ExeCmd = executableName
		} else {
			instance.ExeCmd = executablePath
		}

		instance.KeyPath = filepath.Join(instance.AppHomeDir, ".key")
	})

	return instance, err
}

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
