package env

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/engmtcdrm/minno/app"
)

type Env struct {
	// Home is the user's home directory.
	Home string
	// AppHomeDir is the directory in the user's home directory where the app
	// stores its data.
	AppHomeDir string
	// KeyPath is the path to the key file.
	KeyPath string
	// ExeCmd is the command to run the executable. If the executable is in the
	// PATH environment variable, this will be the executable name.
	ExeCmd string
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

// IsInPath checks if the directory of the executable is in the PATH
// environment variable
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
