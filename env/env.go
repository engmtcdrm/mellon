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
	// AppHomeDir is the directory in the user's home directory where the app stores its data.
	AppHomeDir string
	// KeyPath is the path to the key file.
	KeyPath string
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

		instance.KeyPath = filepath.Join(instance.AppHomeDir, ".key")
	})

	return instance, err
}
