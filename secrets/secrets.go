package secrets

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/engmtcdrm/minno/env"
)

type Secret struct {
	Name string
	Path string
}

// Returns a slice of all available secrets
func GetSecretFiles() ([]Secret, error) {
	envVars, err := env.GetEnv()
	if err != nil {
		return nil, err
	}

	var secretFiles []Secret

	err = filepath.WalkDir(envVars.AppHomeDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".secret" {
			c := Secret{
				Name: strings.TrimSuffix(filepath.Base(path), ".secret"),
				Path: path,
			}
			secretFiles = append(secretFiles, c)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return secretFiles, nil
}

// ResolveSecretName resolves the full path of a secret name
func ResolveSecretName(secretName string) (string, error) {
	if !IsValidName(secretName) {
		return "", errors.New("secret name can only contain alphanumeric, hyphens, and underscores")
	}

	envVars, err := env.GetEnv()
	if err != nil {
		return "", fmt.Errorf("error getting environment variables: %v", err)
	}

	secretName = filepath.Join(envVars.AppHomeDir, secretName)

	if filepath.Ext(secretName) != ".secret" {
		secretName = secretName + ".secret"
	}

	return secretName, nil
}

// IsValidName checks if a string is a valid secret name
func IsValidName(s string) bool {
	var re = regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
	return re.MatchString(s)
}

// IsExists checks if a secret file exists
func IsExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
