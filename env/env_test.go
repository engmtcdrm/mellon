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

func TestIsInPath(t *testing.T) {
	// Test with a path that should be in PATH
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		t.Skip("PATH environment variable is empty, skipping test")
	}

	pathDirs := filepath.SplitList(pathEnv)
	if len(pathDirs) == 0 {
		t.Skip("No directories in PATH, skipping test")
	}

	// Test with a directory that is in PATH
	testDir := pathDirs[0]
	testExecutable := filepath.Join(testDir, "test-executable")

	result := IsInPath(testExecutable)
	if !result {
		t.Errorf("IsInPath() should return true for path in PATH: %s", testExecutable)
	}

	// Test with a path that should not be in PATH
	nonExistentPath := filepath.Join(os.TempDir(), "non-existent-dir", "executable")
	result = IsInPath(nonExistentPath)
	if result {
		t.Errorf("IsInPath() should return false for path not in PATH: %s", nonExistentPath)
	}
}

func TestIsInPathEdgeCases(t *testing.T) {
	// Test with empty path
	result := IsInPath("")
	if result {
		t.Errorf("IsInPath() should return false for empty path")
	}

	// Test with root path
	result = IsInPath("/")
	// This might be true or false depending on system, just ensure it doesn't panic
	t.Logf("IsInPath('/') returned: %v", result)
}

func TestExpandTilde(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "tilde_only",
			input:    "~",
			expected: home,
		},
		{
			name:     "tilde_with_path",
			input:    "~/documents/file.txt",
			expected: filepath.Join(home, "documents/file.txt"),
		},
		{
			name:     "tilde_with_slash",
			input:    "~/",
			expected: filepath.Join(home, ""),
		},
		{
			name:     "no_tilde",
			input:    "/absolute/path",
			expected: "/absolute/path",
		},
		{
			name:     "relative_path",
			input:    "relative/path",
			expected: "relative/path",
		},
		{
			name:     "empty_string",
			input:    "",
			expected: "",
		},
		{
			name:     "tilde_in_middle",
			input:    "/path/~/file",
			expected: "/path/~/file",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ExpandTilde(tc.input)
			if err != nil {
				t.Errorf("ExpandTilde(%q) returned error: %v", tc.input, err)
				return
			}

			if result != tc.expected {
				t.Errorf("ExpandTilde(%q) = %q, expected %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestExpandTildePathSeparators(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	// Test with different path separators
	testCases := []struct {
		name  string
		input string
	}{
		{"forward_slash", "~/path/to/file"},
		{"multiple_levels", "~/documents/projects/secret.txt"},
		{"single_file", "~/file.txt"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ExpandTilde(tc.input)
			if err != nil {
				t.Errorf("ExpandTilde(%q) returned error: %v", tc.input, err)
				return
			}

			// Should start with home directory
			if !strings.HasPrefix(result, home) {
				t.Errorf("ExpandTilde(%q) result should start with home directory %q, got: %q", tc.input, home, result)
			}

			// Should be a valid path
			if !filepath.IsAbs(result) {
				t.Errorf("ExpandTilde(%q) should return absolute path, got: %q", tc.input, result)
			}
		})
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
