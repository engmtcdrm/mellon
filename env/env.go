package env

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/engmtcdrm/mellon/app"
)

var (
	Instance  *Env // Singleton instance of Env
	once      sync.Once
	secretExt = ".thurin" // The file extension for secret files.
)

type Env struct {
	home        string // User's home directory.
	appHomeDir  string // The directory in the user's home directory where the app stores its data.
	keyPath     string // The path to the encryption key file.
	secretsPath string // The path to the directory where secrets are stored.
	secretExt   string // The file extension for secret files.
	exeCmd      string // The command to run the executable. If the executable is in the PATH environment variable, this will be the executable name.
}

// Home returns the home directory of the user.
func (e *Env) Home() string {
	return e.home
}

// AppHomeDir returns the app home directory.
func (e *Env) AppHomeDir() string {
	return e.appHomeDir
}

// KeyPath returns the encryption key path.
func (e *Env) KeyPath() string {
	return e.keyPath
}

// SecretsPath returns the secrets path.
func (e *Env) SecretsPath() string {
	return e.secretsPath
}

// SecretExt returns the secret file extension.
func (e *Env) SecretExt() string {
	return e.secretExt
}

// ExeCmd returns the executable command.
func (e *Env) ExeCmd() string {
	return e.exeCmd
}

// Init initializes the environment variables.
func Init() {
	once.Do(func() {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}

		Instance = &Env{
			home:       home,
			appHomeDir: filepath.Join(home, app.DotName),
			secretExt:  secretExt,
		}

		executablePath, err := os.Executable()
		if err != nil {
			panic(err)
		}

		// If executable is in path, use the base name, i.e. the executable name
		if IsInPath(executablePath) {
			Instance.exeCmd = filepath.Base(executablePath)
		} else {
			Instance.exeCmd = executablePath
		}

		Instance.keyPath = filepath.Join(Instance.appHomeDir, ".key")
		Instance.secretsPath = filepath.Join(Instance.appHomeDir, Instance.secretExt)
	})
}
