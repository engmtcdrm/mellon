package secrets

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/engmtcdrm/go-entomb"
	"github.com/engmtcdrm/mellon/env"
)

// Secret represents a secret value stored in the system.
type Secret struct {
	name string
	path string
	tomb *entomb.Tomb
}

// NewSecret creates a new secret with the given key path, name, and path.
func NewSecret(keyPath string, name string, path string) (*Secret, error) {
	if keyPath == "" {
		return nil, errors.New("key path cannot be empty")
	}

	if err := ValidateName(name); err != nil {
		return nil, fmt.Errorf("%s. The secret name provided was '%s'", err, name)
	}

	if path == "" {
		return nil, errors.New("path cannot be empty")
	}

	tomb, err := entomb.NewTomb(keyPath, true, true)
	if err != nil {
		return nil, fmt.Errorf("could not create tomb: %w", err)
	}

	return &Secret{
		name: name,
		path: path,
		tomb: tomb,
	}, nil
}

// Name returns the name of the secret.
func (s *Secret) Name() string {
	return s.name
}

// Path returns the path of the secret.
func (s *Secret) Path() string {
	return s.path
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

	err = os.MkdirAll(filepath.Dir(s.path), dirMode)
	if err != nil {
		return fmt.Errorf("could not create directory for secret '%s': %w", s.name, err)
	}

	if err = os.WriteFile(s.path, encSecret, secretMode); err != nil {
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

	err = os.MkdirAll(filepath.Dir(s.path), dirMode)
	if err != nil {
		return fmt.Errorf("could not create directory for secret '%s': %w", s.name, err)
	}

	if err = os.WriteFile(s.path, encSecret, secretMode); err != nil {
		return err
	}

	return nil
}

// Decrypt reads the encrypted secret from the file and decrypts it.
func (s *Secret) Decrypt() ([]byte, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsPermission(err) {
			return nil, fmt.Errorf("failed to read secret '%s': permission denied", s.name)
		}

		if os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to read secret '%s': secret does not exist", s.name)
		}

		return nil, err
	}

	secret, err := s.tomb.Decrypt(data)
	data = nil
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt secret '%s'. Encrypted secret may be corrupted", s.name)
	}

	return secret, nil
}
