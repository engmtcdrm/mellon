package env

import (
	"os"
	"path/filepath"
)

const (
	AppNm        = "minno"
	AppVersion   = "0.0.1"
	AppLongDesc  = "A CLI tool for securing and obtaining credentials"
	AppShortDesc = "A CLI tool for securing and obtaining credentials"
	RepoUrl      = "https://github.com/engmtcdrm/minno"
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

	AppHomeDir = filepath.Join(home, "."+AppNm)

	return nil
}
