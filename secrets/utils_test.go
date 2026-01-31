package secrets

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSecretFiles(t *testing.T) {
	keyPath := path.Join(t.TempDir(), "test.key")

	secretExt := ".secret"

	createSecretDir := func() string {
		dirPath := path.Join(t.TempDir(), "secrets")
		err := os.MkdirAll(dirPath, 0755)
		assert.NoError(t, err)
		return dirPath
	}

	t.Run("empty keyPath", func(t *testing.T) {
		secretsPath := createSecretDir()
		defer os.RemoveAll(secretsPath)

		emptySecrets, err := GetSecretFiles("", secretsPath, secretExt)
		assert.NoError(t, err)
		assert.Empty(t, emptySecrets)
	})

	t.Run("empty secretsPath", func(t *testing.T) {
		emptySecrets, err := GetSecretFiles(keyPath, "", secretExt)
		assert.Error(t, err)
		assert.Empty(t, emptySecrets)
	})

	t.Run("empty secretExt", func(t *testing.T) {
		secretsPath := createSecretDir()
		defer os.RemoveAll(secretsPath)

		emptySecrets, err := GetSecretFiles(keyPath, secretsPath, "")
		assert.NoError(t, err)
		assert.Empty(t, emptySecrets)
	})

}

func TestValidateName(t *testing.T) {
	// Valid names
	assert.NoError(t, ValidateName("valid_name"))
	assert.NoError(t, ValidateName("valid-name"))
	assert.NoError(t, ValidateName("validname123"))
	assert.NoError(t, ValidateName("valid/name"))
	assert.NoError(t, ValidateName("valid\\name"))

	// Invalid names
	assert.Error(t, ValidateName(""))
	assert.Error(t, ValidateName("invalid name"))
	assert.Error(t, ValidateName("invalid.name"))
}
