package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	assert.Truef(t, result, "IsInPath() should return true for path in PATH: %s", testExecutable)

	// Test with a path that should not be in PATH
	nonExistentPath := filepath.Join(os.TempDir(), "non-existent-dir", "executable")
	result = IsInPath(nonExistentPath)
	assert.Falsef(t, result, "IsInPath() should return false for path not in PATH: %s", nonExistentPath)
}

func TestIsInPathEdgeCases(t *testing.T) {
	// Test with empty path
	result := IsInPath("")
	assert.False(t, result, "IsInPath() should return false for empty path")

	// Test with root path
	result = IsInPath("/")
	// This might be true or false depending on system, just ensure it doesn't panic
	t.Logf("IsInPath('/') returned: %v", result)
}

func TestExpandTilde(t *testing.T) {
	home, err := os.UserHomeDir()
	assert.NoErrorf(t, err, "Failed to get user home directory: %v", err)

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
			assert.NoErrorf(t, err, "ExpandTilde(%q) returned error: %v", tc.input, err)
			assert.Equalf(t, tc.expected, result, "ExpandTilde(%q) = %q, expected %q", tc.input, result, tc.expected)
		})
	}
}

func TestExpandTildePathSeparators(t *testing.T) {
	home, err := os.UserHomeDir()
	assert.NoErrorf(t, err, "Failed to get user home directory: %v", err)

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
			assert.NoErrorf(t, err, "ExpandTilde(%q) returned error: %v", tc.input, err)
			assert.Truef(t, strings.HasPrefix(result, home), "ExpandTilde(%q) result should start with home directory %q, got: %q", tc.input, home, result)
			assert.Truef(t, filepath.IsAbs(result), "ExpandTilde(%q) should return absolute path, got: %q", tc.input, result)
		})
	}
}
