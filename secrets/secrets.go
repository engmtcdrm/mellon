package secrets

import (
	"errors"
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

	err = filepath.WalkDir(envVars.SecretsPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == envVars.SecretExt && !d.IsDir() {
			c := Secret{
				Name: strings.TrimSuffix(filepath.Base(path), envVars.SecretExt),
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

// ValidateName checks if a string is a valid secret name
func ValidateName(s string) error {
	var re = regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)

	if re.MatchString(s) {
		return nil
	}

	return errors.New("invalid secret name: Secret name can only contain alphanumeric, hyphens, and underscores")
}

// FindSecretByName searches for a secret by its name in the provided slice of secrets.
func FindSecretByName(name string, secretFiles []Secret) *Secret {
	for _, f := range secretFiles {
		if f.Name == name {
			return &f
		}
	}

	return nil
}
