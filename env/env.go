package env

import (
	"os"
	"path/filepath"

	"github.com/engmtcdrm/minno/app"
)

var (
	AppHomeDir string
)

// Returns the home directory for the application.
func SetAppHomeDir() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	AppHomeDir = filepath.Join(home, "."+app.Name)

	return nil
}
