package secrets

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/engmtcdrm/go-entomb"
	"github.com/engmtcdrm/mellon/env"
)

type Secret struct {
	Name string
	Path string
	tomb *entomb.Tomb
}

var (
	dirMode    os.FileMode = 0700 // Default directory mode for secrets
	secretMode os.FileMode = 0600 // Default file mode for secret files
)

func NewSecret(keyPath string, name string, path string) (*Secret, error) {
	if keyPath == "" {
		return nil, errors.New("key path cannot be empty")
	}

	tomb, err := entomb.NewTomb(keyPath)
	if err != nil {
		return nil, fmt.Errorf("could not create tomb: %w", err)
	}

	if err := ValidateName(name); err != nil {
		return nil, fmt.Errorf("%s. The secret name provided was '%s'", err, name)
	}

	if path == "" {
		return nil, errors.New("path cannot be empty")
	}

	return &Secret{Name: name, Path: path, tomb: tomb}, nil
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

		relPath, err := filepath.Rel(envVars.SecretsPath, path)
		if err != nil {
			return err
		}

		if filepath.Ext(path) == envVars.SecretExt && !d.IsDir() {
			c, err := NewSecret(envVars.KeyPath, strings.TrimSuffix(relPath, envVars.SecretExt), path)
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

func isDirEmpty(dirPath string) (bool, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
}

func RemoveSecret(secret Secret) error {
	if secret.Path == "" {
		return errors.New("secret path cannot be empty")
	}

	if err := os.Remove(secret.Path); err != nil {
		return fmt.Errorf("could not remove secret '%s': %w", secret.Name, err)
	}

	envVars, err := env.GetEnv()
	if err != nil {
		return err
	}

	// Ignore trying to delete the secrets directory itself
	if secret.Path == envVars.SecretsPath {
		return nil
	}

	dirIsEmpty, err := isDirEmpty(filepath.Dir(secret.Path))
	if err != nil {
		return fmt.Errorf("could not check if directory is empty: %w", err)
	}

	if dirIsEmpty {
		if err := os.Remove(filepath.Dir(secret.Path)); err != nil {
			return fmt.Errorf("could not remove directory '%s': %w", filepath.Dir(secret.Path), err)
		}
	}

	return nil
}

// ValidateName checks if a string is a valid secret name
func ValidateName(s string) error {
	var re = regexp.MustCompile(`^[\w\/\\\-]+$`)

	if re.MatchString(s) {
		return nil
	}

	return errors.New("invalid secret name: Secret name can only contain alphanumeric, hyphens, underscores, and slashes")
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

// trimSpaceBytes trims leading and trailing ASCII whitespace from a byte slice in-place.
// Returns a subslice of the original slice, so the underlying array is not copied.
func trimSpaceBytes(b []byte) []byte {
	start := 0
	end := len(b)

	// Trim leading spaces
	for start < end && (b[start] == ' ' || b[start] == '\t' || b[start] == '\n' || b[start] == '\r') {
		start++
	}
	// Trim trailing spaces
	for end > start && (b[end-1] == ' ' || b[end-1] == '\t' || b[end-1] == '\n' || b[end-1] == '\r') {
		end--
	}
	return b[start:end]
}

// EncryptFromFile reads a secret from a file, trims leading and trailing whitespace
// and encrypts it before writing it to the secret's path.
func (s *Secret) EncryptFromFile(file string, cleanup bool) error {
	var encSecret []byte

	rawFile, err := env.ExpandTilde(strings.TrimSpace(file))
	if err != nil {
		return err
	}

	secretBytes, err := os.ReadFile(rawFile)
	if err != nil {
		return fmt.Errorf("could not read file '%s': %w", rawFile, err)
	}

	secretBytes = trimSpaceBytes(secretBytes)
	encSecret, err = s.tomb.Encrypt(secretBytes)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(s.Path), dirMode)
	if err != nil {
		return fmt.Errorf("could not create directory for secret '%s': %w", s.Name, err)
	}

	if err = os.WriteFile(s.Path, encSecret, secretMode); err != nil {
		return err
	}

	if cleanup {
		if err = os.Remove(rawFile); err != nil {
			return fmt.Errorf("could not remove file '%s': %w", rawFile, err)
		}
	}

	return nil
}

// Encrypt encrypts a secret and writes it to the secret's path.
// The secret is trimmed of leading and trailing whitespace before encryption.
func (s *Secret) Encrypt(secret []byte) error {
	encSecret, err := s.tomb.Encrypt(trimSpaceBytes(secret))
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(s.Path), dirMode)
	if err != nil {
		return fmt.Errorf("could not create directory for secret '%s': %w", s.Name, err)
	}

	if err = os.WriteFile(s.Path, encSecret, secretMode); err != nil {
		return err
	}

	return nil
}

// Decrypt reads the encrypted secret from the file and decrypts it.
func (s *Secret) Decrypt() ([]byte, error) {
	data, err := os.ReadFile(s.Path)
	if err != nil {
		if os.IsPermission(err) {
			return nil, fmt.Errorf("failed to read secret '%s': permission denied", s.Name)
		}

		if os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to read secret '%s': secret does not exist", s.Name)
		}

		return nil, err
	}

	secret, err := s.tomb.Decrypt(data)
	data = nil
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt secret '%s'. Encrypted secret may be corrupted", s.Name)
	}

	return secret, nil
}
