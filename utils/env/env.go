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

// Returns the home directory for the application.
func AppHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, "."+AppNm), nil
}
