package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetEnv(t *testing.T) {
	// Test that GetEnv returns a valid Env instance
	env, err := GetEnv()
	if err != nil {
		t.Fatalf("GetEnv() returned error: %v", err)
	}

	if env == nil {
		t.Fatalf("GetEnv() returned nil")
	}

	// Test that subsequent calls return the same instance (singleton)
	env2, err := GetEnv()
	if err != nil {
		t.Fatalf("Second GetEnv() call returned error: %v", err)
	}

	if env != env2 {
		t.Errorf("GetEnv() should return the same instance (singleton pattern)")
	}
}

func TestEnvFields(t *testing.T) {
	env, err := GetEnv()
	if err != nil {
		t.Fatalf("GetEnv() returned error: %v", err)
	}

	// Test Home field
	if env.Home == "" {
		t.Errorf("Home field should not be empty")
	}

	expectedHome, _ := os.UserHomeDir()
	if env.Home != expectedHome {
		t.Errorf("Home field should be user home directory, got: %s, expected: %s", env.Home, expectedHome)
	}

	// Test AppHomeDir field
	if env.AppHomeDir == "" {
		t.Errorf("AppHomeDir field should not be empty")
	}

	if !strings.HasPrefix(env.AppHomeDir, env.Home) {
		t.Errorf("AppHomeDir should be within Home directory")
	}

	// Test SecretExt field
	if env.SecretExt == "" {
		t.Errorf("SecretExt field should not be empty")
	}

	if env.SecretExt != ".thurin" {
		t.Errorf("SecretExt should be '.thurin', got: %s", env.SecretExt)
	}

	// Test KeyPath field
	if env.KeyPath == "" {
		t.Errorf("KeyPath field should not be empty")
	}

	if !strings.HasPrefix(env.KeyPath, env.AppHomeDir) {
		t.Errorf("KeyPath should be within AppHomeDir")
	}

	expectedKeyPath := filepath.Join(env.AppHomeDir, ".key")
	if env.KeyPath != expectedKeyPath {
		t.Errorf("KeyPath should be %s, got: %s", expectedKeyPath, env.KeyPath)
	}

	// Test SecretsPath field
	if env.SecretsPath == "" {
		t.Errorf("SecretsPath field should not be empty")
	}

	if !strings.HasPrefix(env.SecretsPath, env.AppHomeDir) {
		t.Errorf("SecretsPath should be within AppHomeDir")
	}

	expectedSecretsPath := filepath.Join(env.AppHomeDir, env.SecretExt)
	if env.SecretsPath != expectedSecretsPath {
		t.Errorf("SecretsPath should be %s, got: %s", expectedSecretsPath, env.SecretsPath)
	}

	// Test ExeCmd field
	if env.ExeCmd == "" {
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
		env, err := GetEnv()
		if err != nil {
			t.Fatalf("GetEnv() call %d returned error: %v", i, err)
		}
		instances[i] = env
	}

	// Verify all instances are the same
	firstInstance := instances[0]
	for i, instance := range instances {
		if instance != firstInstance {
			t.Errorf("Instance %d is different from first instance (singleton violation)", i)
		}
	}
}
