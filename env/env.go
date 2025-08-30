package env

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/engmtcdrm/mellon/app"
)

type Env struct {
	Home        string // Home is the user's home directory.
	AppHomeDir  string // AppHomeDir is the directory in the user's home directory where the app stores its data.
	KeyPath     string // KeyPath is the path to the key file.
	SecretsPath string // SecretsPath is the path to the directory where secrets are stored.
	SecretExt   string
	ExeCmd      string // ExeCmd is the command to run the executable. If the executable is in the PATH environment variable, this will be the executable name.
}

var (
	Instance *Env // Singleton instance of Env
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

		Instance = &Env{
			Home:       home,
			AppHomeDir: filepath.Join(home, app.DotName),
			SecretExt:  ".thurin",
		}

		executablePath, e := os.Executable()
		if e != nil {
			err = e
			return
		}

		executableName := filepath.Base(executablePath)

		if IsInPath(executablePath) {
			Instance.ExeCmd = executableName
		} else {
			Instance.ExeCmd = executablePath
		}

		Instance.KeyPath = filepath.Join(Instance.AppHomeDir, ".key")
		Instance.SecretsPath = filepath.Join(Instance.AppHomeDir, Instance.SecretExt)
	})

	return Instance, err
}
