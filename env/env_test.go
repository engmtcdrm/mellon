package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnvFields(t *testing.T) {
	Init()

	// Test Home field
	if Instance.Home() == "" {
		t.Errorf("Home field should not be empty")
	}

	expectedHome, _ := os.UserHomeDir()
	if Instance.Home() != expectedHome {
		t.Errorf("Home field should be user home directory, got: %s, expected: %s", Instance.Home(), expectedHome)
	}

	// Test AppHomeDir field
	if Instance.AppHomeDir() == "" {
		t.Errorf("AppHomeDir field should not be empty")
	}

	if !strings.HasPrefix(Instance.AppHomeDir(), Instance.Home()) {
		t.Errorf("AppHomeDir should be within Home directory")
	}

	// Test SecretExt field
	if Instance.SecretExt() == "" {
		t.Errorf("SecretExt field should not be empty")
	}

	if Instance.SecretExt() != ".thurin" {
		t.Errorf("SecretExt should be '.thurin', got: %s", Instance.SecretExt())
	}

	// Test KeyPath field
	if Instance.KeyPath() == "" {
		t.Errorf("KeyPath field should not be empty")
	}

	if !strings.HasPrefix(Instance.KeyPath(), Instance.AppHomeDir()) {
		t.Errorf("KeyPath should be within AppHomeDir")
	}

	expectedKeyPath := filepath.Join(Instance.AppHomeDir(), ".key")
	if Instance.KeyPath() != expectedKeyPath {
		t.Errorf("KeyPath should be %s, got: %s", expectedKeyPath, Instance.KeyPath())
	}

	// Test SecretsPath field
	if Instance.SecretsPath() == "" {
		t.Errorf("SecretsPath field should not be empty")
	}

	if !strings.HasPrefix(Instance.SecretsPath(), Instance.AppHomeDir()) {
		t.Errorf("SecretsPath should be within AppHomeDir")
	}

	expectedSecretsPath := filepath.Join(Instance.AppHomeDir(), Instance.SecretExt())
	if Instance.SecretsPath() != expectedSecretsPath {
		t.Errorf("SecretsPath should be %s, got: %s", expectedSecretsPath, Instance.SecretsPath())
	}

	// Test ExeCmd field
	if Instance.ExeCmd() == "" {
		t.Errorf("ExeCmd field should not be empty")
	}
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
		if instance != firstInstance {
			t.Errorf("Instance %d is different from first instance (singleton violation)", i)
		}
	}
}
