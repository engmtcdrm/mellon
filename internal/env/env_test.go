package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvFields(t *testing.T) {
	Init()
	assert.NotNil(t, Instance, "Instance should not be nil")
	assert.NotEmpty(t, Instance.Home(), "Home field should not be empty")

	expectedHome, _ := os.UserHomeDir()
	assert.Equalf(t, expectedHome, Instance.Home(), "Home field should be user home directory, got: %s, expected: %s", Instance.Home(), expectedHome)

	assert.NotEmpty(t, Instance.AppHomeDir(), "AppHomeDir field should not be empty")
	assert.Truef(t, strings.HasPrefix(Instance.AppHomeDir(), Instance.Home()), "AppHomeDir should be within Home directory, got: %s, Home: %s", Instance.AppHomeDir(), Instance.Home())
	assert.NotEmpty(t, Instance.SecretExt(), "SecretExt field should not be empty")
	assert.Equalf(t, ".thurin", Instance.SecretExt(), "SecretExt should be '.thurin', got: %s", Instance.SecretExt())
	assert.NotEmpty(t, Instance.KeyPath(), "KeyPath field should not be empty")
	assert.Truef(t, strings.HasPrefix(Instance.KeyPath(), Instance.AppHomeDir()), "KeyPath should be within AppHomeDir, got: %s, AppHomeDir: %s", Instance.KeyPath(), Instance.AppHomeDir())

	expectedKeyPath := filepath.Join(Instance.AppHomeDir(), ".key")
	assert.Equalf(t, expectedKeyPath, Instance.KeyPath(), "KeyPath should be %s, got: %s", expectedKeyPath, Instance.KeyPath())
	assert.NotEmpty(t, Instance.SecretsPath(), "SecretsPath field should not be empty")
	assert.Truef(t, strings.HasPrefix(Instance.SecretsPath(), Instance.AppHomeDir()), "SecretsPath should be within AppHomeDir, got: %s, AppHomeDir: %s", Instance.SecretsPath(), Instance.AppHomeDir())

	expectedSecretsPath := filepath.Join(Instance.AppHomeDir(), Instance.SecretExt())
	assert.Equalf(t, expectedSecretsPath, Instance.SecretsPath(), "SecretsPath should be %s, got: %s", expectedSecretsPath, Instance.SecretsPath())

	assert.NotEmpty(t, Instance.ExeCmd(), "ExeCmd field should not be empty")
}

func TestEnvSingleton(t *testing.T) {
	// Reset the singleton for this test by creating a new test
	// Note: We can't actually reset the singleton in production code due to sync.Once
	// But we can test that multiple calls return the same instance

	instances := make([]*Env, 10)

	// Call GetEnv multiple times
	for i := 0; i < 10; i++ {
		Init()
		instances[i] = Instance
	}

	// Verify all instances are the same
	firstInstance := instances[0]
	for i, instance := range instances {
		assert.Equalf(t, firstInstance, instance, "Instance %d is different from first instance (singleton violation)", i)
	}
}
