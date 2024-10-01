package credentials

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/engmtcdrm/minno/env"
)

type Credential struct {
	Name string
	Path string
}

// Returns a slice of all available credentials
func GetCredFiles() ([]Credential, error) {
	envVars, err := env.GetEnv()
	if err != nil {
		return nil, err
	}

	var credFiles []Credential

	err = filepath.WalkDir(envVars.AppHomeDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".cred" {
			c := Credential{
				Name: strings.TrimSuffix(filepath.Base(path), ".cred"),
				Path: path,
			}
			credFiles = append(credFiles, c)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return credFiles, nil
}

// ResolveCredName resolves the full path of a credential name
func ResolveCredName(credName string) (string, error) {
	if !IsValidName(credName) {
		return "", errors.New("credential name can only contain alphanumeric, hyphens, and underscores")
	}

	envVars, err := env.GetEnv()
	if err != nil {
		return "", fmt.Errorf("error getting environment variables: %v", err)
	}

	credName = filepath.Join(envVars.AppHomeDir, credName)

	if filepath.Ext(credName) != ".cred" {
		credName = credName + ".cred"
	}

	return credName, nil
}

// IsValidName checks if a string is a valid credential name
func IsValidName(s string) bool {
	var re = regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
	return re.MatchString(s)
}

// IsExists checks if a credential file exists
func IsExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
