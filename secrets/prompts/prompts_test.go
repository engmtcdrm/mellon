package prompts_test

import (
	"path"
	"testing"

	"github.com/engmtcdrm/mellon/secrets"
	"github.com/engmtcdrm/mellon/secrets/prompts"
	"github.com/stretchr/testify/assert"
)

func TestGetSecretOptions(t *testing.T) {
	keyPath := path.Join(t.TempDir(), "temp.key")

	t.Run("no secrets available", func(t *testing.T) {
		_, err := prompts.GetSecretOptions([]secrets.Secret{}, "use", "mellon")
		assert.Error(t, err)
	})

	t.Run("secrets available", func(t *testing.T) {
		secret1, err := secrets.NewSecret(keyPath, "secret1", path.Join(t.TempDir(), "secret1"))
		assert.NoError(t, err)
		secret2, err := secrets.NewSecret(keyPath, "secret2", path.Join(t.TempDir(), "secret2"))
		assert.NoError(t, err)

		var secrets []secrets.Secret
		secrets = append(secrets, *secret1, *secret2)

		options, err := prompts.GetSecretOptions(secrets, "use", "mellon")
		assert.NoError(t, err)
		assert.Len(t, options, 2)
		assert.Equal(t, secrets[0], options[0].Value)
		assert.Equal(t, secrets[1], options[1].Value)
	})
}
