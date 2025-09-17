package secrets

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const reValidName = `^[\w\/\\\-]+$`

var (
	dirMode    os.FileMode = 0700 // Default directory mode for secrets
	secretMode os.FileMode = 0600 // Default file mode for secret files
)

// Returns a slice of all available secrets
func GetSecretFiles(keyPath, secretsPath, secretExt string) ([]Secret, error) {
	var secretFiles []Secret

	err := filepath.WalkDir(secretsPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(secretsPath, path)
		if err != nil {
			return err
		}

		if filepath.Ext(path) == secretExt && !d.IsDir() {
			c, err := NewSecret(keyPath, strings.TrimSuffix(relPath, secretExt), path)
			if err != nil {
				return err
			}
			secretFiles = append(secretFiles, *c)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return secretFiles, nil
}

// RemoveSecret removes a secret from the specified secrets path
func RemoveSecret(secretsPath string, secret Secret) error {
	if secret.Path() == "" {
		return errors.New("secret path cannot be empty")
	}

	if err := os.Remove(secret.Path()); err != nil {
		return fmt.Errorf("could not remove secret '%s': %w", secret.name, err)
	}

	// Ignore trying to delete the secrets directory itself
	if secret.Path() == secretsPath {
		return nil
	}

	dirIsEmpty, err := isDirEmpty(filepath.Dir(secret.Path()))
	if err != nil {
		return fmt.Errorf("could not check if directory is empty: %w", err)
	}

	if dirIsEmpty {
		if err := os.Remove(filepath.Dir(secret.Path())); err != nil {
			return fmt.Errorf("could not remove directory '%s': %w", filepath.Dir(secret.Path()), err)
		}
	}

	return nil
}

// ValidateName checks if a string is a valid secret name
func ValidateName(s string) error {
	var re = regexp.MustCompile(reValidName)

	if re.MatchString(s) {
		return nil
	}

	return errors.New("invalid secret name: Secret name can only contain alphanumeric, hyphens, underscores, and slashes")
}

// FindSecretByName searches for a secret by its name in the provided slice of secrets.
func FindSecretByName(name string, secretFiles []Secret) *Secret {
	for _, f := range secretFiles {
		if f.name == name {
			return &f
		}
	}

	return nil
}

// ClearSecret overwrites the contents of a byte slice with zeros clearing sensitive data from memory.
func ClearSecret(s *[]byte) {
	if s != nil {
		for i := range *s {
			(*s)[i] = 0
		}
	}
}

// isDirEmpty checks if a directory is empty.
func isDirEmpty(dirPath string) (bool, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
}

// trimSpaceBytes trims leading and trailing ASCII whitespace from a byte slice in-place.
// Returns a subslice of the original slice, so the underlying array is not copied.
func trimSpaceBytes(b *[]byte) []byte {
	start := 0
	end := len(*b)

	// Trim leading spaces
	for start < end && ((*b)[start] == ' ' || (*b)[start] == '\t' || (*b)[start] == '\n' || (*b)[start] == '\r') {
		start++
	}
	// Trim trailing spaces
	for end > start && ((*b)[end-1] == ' ' || (*b)[end-1] == '\t' || (*b)[end-1] == '\n' || (*b)[end-1] == '\r') {
		end--
	}
	return (*b)[start:end]
}
