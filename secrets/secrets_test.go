package secrets

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSecret(t *testing.T) {
	t.Run("empty secret keyPath", func(t *testing.T) {
		_, err := NewSecret("", "value", "value")
		assert.Error(t, err)
	})

	t.Run("invalid secret name", func(t *testing.T) {
		_, err := NewSecret(path.Join(t.TempDir(), "temp.key"), "value!:", "value")
		assert.Error(t, err)
	})
}
