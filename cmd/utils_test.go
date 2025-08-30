package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/engmtcdrm/mellon/env"
	"github.com/engmtcdrm/mellon/secrets"
	"github.com/spf13/cobra"
)

// TestValidateUpdateCreateFlags tests the validateUpdateCreateFlags function
func TestValidateUpdateCreateFlags(t *testing.T) {
	// Save original values
	origCleanupFile := cleanupFile
	origSecretName := secretName
	origSecretFile := secretFile
	defer func() {
		cleanupFile = origCleanupFile
		secretName = origSecretName
		secretFile = origSecretFile
	}()

	cmd := &cobra.Command{}

	tests := []struct {
		name        string
		cleanupFile bool
		secretName  string
		secretFile  string
		expectError bool
	}{
		{
			name:        "cleanup with both secret and file",
			cleanupFile: true,
			secretName:  "testsecret",
			secretFile:  "testfile.txt",
			expectError: false,
		},
		{
			name:        "cleanup without secret name",
			cleanupFile: true,
			secretName:  "",
			secretFile:  "testfile.txt",
			expectError: true,
		},
		{
			name:        "cleanup without secret file",
			cleanupFile: true,
			secretName:  "testsecret",
			secretFile:  "",
			expectError: true,
		},
		{
			name:        "cleanup without both",
			cleanupFile: true,
			secretName:  "",
			secretFile:  "",
			expectError: true,
		},
		{
			name:        "no cleanup flag",
			cleanupFile: false,
			secretName:  "",
			secretFile:  "",
			expectError: false,
		},
		{
			name:        "no cleanup with partial flags",
			cleanupFile: false,
			secretName:  "testsecret",
			secretFile:  "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanupFile = tt.cleanupFile
			secretName = tt.secretName
			secretFile = tt.secretFile

			err := validateUpdateCreateFlags(cmd, []string{})
			if tt.expectError && err == nil {
				t.Errorf("expected error for %s, got none", tt.name)
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected no error for %s, got: %v", tt.name, err)
			}
		})
	}
}

// TestValidateSecretName tests the validateSecretName function
func TestValidateSecretName(t *testing.T) {
	secretFiles = []secrets.Secret{}

	newSecret, err := secrets.NewSecret(env.Instance.KeyPath(), "existing_secret", "existing_secret")
	if err != nil {
		t.Fatalf("failed to create new secret: %v", err)
	}

	secretFiles = append(secretFiles, *newSecret)

	tests := []struct {
		name        string
		secretName  string
		expectError bool
	}{
		{
			name:        "valid new secret name",
			secretName:  "new_secret",
			expectError: false,
		},
		{
			name:        "empty secret name",
			secretName:  "",
			expectError: true,
		},
		{
			name:        "existing secret name",
			secretName:  "existing_secret",
			expectError: true,
		},
		{
			name:        "invalid characters - dots",
			secretName:  "invalid.name",
			expectError: true,
		},
		{
			name:        "invalid characters - spaces",
			secretName:  "invalid name",
			expectError: true,
		},
		{
			name:        "invalid characters - special chars",
			secretName:  "invalid@name!",
			expectError: true,
		},
		{
			name:        "valid name with underscores",
			secretName:  "valid_name_123",
			expectError: false,
		},
		{
			name:        "valid name with hyphens",
			secretName:  "valid-name-123",
			expectError: false,
		},
		{
			name:        "valid name with forward slashes",
			secretName:  "valid/name/123",
			expectError: false,
		},
		{
			name:        "valid name with backslashes",
			secretName:  "valid\\name\\123",
			expectError: false,
		},
		{
			name:        "valid alphanumeric only",
			secretName:  "valid123",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSecretName(tt.secretName)
			if tt.expectError && err == nil {
				t.Errorf("expected error for %s, got none", tt.name)
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected no error for %s, got: %v", tt.name, err)
			}
		})
	}
}

// TestMkdir tests the mkdir function
func TestMkdir(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name     string
		path     string
		mode     os.FileMode
		existing bool
	}{
		{
			name:     "create new directory",
			path:     filepath.Join(tempDir, "newdir"),
			mode:     0700,
			existing: false,
		},
		{
			name:     "create nested directory",
			path:     filepath.Join(tempDir, "nested", "deep", "dir"),
			mode:     0750,
			existing: false,
		},
		{
			name:     "existing directory",
			path:     tempDir, // tempDir already exists
			mode:     0755,
			existing: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create directory if it should exist beforehand
			if tt.existing && tt.path != tempDir {
				if err := os.MkdirAll(tt.path, 0755); err != nil {
					t.Fatalf("failed to create existing directory: %v", err)
				}
			}

			// Test mkdir function
			mkdir(tt.path, tt.mode)

			// Verify directory exists
			info, err := os.Stat(tt.path)
			if err != nil {
				t.Errorf("directory %s should exist after mkdir: %v", tt.path, err)
				return
			}

			if !info.IsDir() {
				t.Errorf("%s should be a directory", tt.path)
			}

			// Verify permissions (only for new directories, as tempDir might have different perms)
			if !tt.existing || tt.path != tempDir {
				if info.Mode().Perm() != tt.mode {
					t.Errorf("expected mode %o, got %o", tt.mode, info.Mode().Perm())
				}
			}
		})
	}
}

// TestSecureFiles tests the secureFiles function
func TestSecureFiles(t *testing.T) {
	tempDir := t.TempDir()

	// Create test structure
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdirectory: %v", err)
	}

	file1 := filepath.Join(tempDir, "file1.txt")
	file2 := filepath.Join(subDir, "file2.txt")

	if err := os.WriteFile(file1, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create file1: %v", err)
	}
	if err := os.WriteFile(file2, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create file2: %v", err)
	}

	// Test secureFiles
	dirMode := os.FileMode(0700)
	fileMode := os.FileMode(0600)

	secureFiles(tempDir, dirMode, fileMode)

	// Check root directory permissions
	info, err := os.Stat(tempDir)
	if err != nil {
		t.Fatalf("failed to stat root directory: %v", err)
	}
	if info.Mode().Perm() != dirMode {
		t.Errorf("root directory: expected mode %o, got %o", dirMode, info.Mode().Perm())
	}

	// Check subdirectory permissions
	info, err = os.Stat(subDir)
	if err != nil {
		t.Fatalf("failed to stat subdirectory: %v", err)
	}
	if info.Mode().Perm() != dirMode {
		t.Errorf("subdirectory: expected mode %o, got %o", dirMode, info.Mode().Perm())
	}

	// Check file permissions
	info, err = os.Stat(file1)
	if err != nil {
		t.Fatalf("failed to stat file1: %v", err)
	}
	if info.Mode().Perm() != fileMode {
		t.Errorf("file1: expected mode %o, got %o", fileMode, info.Mode().Perm())
	}

	info, err = os.Stat(file2)
	if err != nil {
		t.Fatalf("failed to stat file2: %v", err)
	}
	if info.Mode().Perm() != fileMode {
		t.Errorf("file2: expected mode %o, got %o", fileMode, info.Mode().Perm())
	}
}

// TestGetSemVer tests the getSemVer function
func TestGetSemVer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid semver with v prefix",
			input:    "v1.2.3",
			expected: "1.2.3",
		},
		{
			name:     "valid semver without v prefix",
			input:    "1.2.3",
			expected: "1.2.3",
		},
		{
			name:     "valid semver with leading zeros",
			input:    "v01.02.03",
			expected: "01.02.03",
		},
		{
			name:     "valid semver with large numbers",
			input:    "v123.456.789",
			expected: "123.456.789",
		},
		{
			name:     "invalid semver - too many parts",
			input:    "v1.2.3.4",
			expected: "v1.2.3.4",
		},
		{
			name:     "invalid semver - too few parts",
			input:    "v1.2",
			expected: "v1.2",
		},
		{
			name:     "invalid semver - non-numeric",
			input:    "v1.2.a",
			expected: "v1.2.a",
		},
		{
			name:     "invalid semver - with text",
			input:    "v1.2.3-beta",
			expected: "v1.2.3-beta",
		},
		{
			name:     "non-semver string",
			input:    "main",
			expected: "main",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "just v",
			input:    "v",
			expected: "v",
		},
		{
			name:     "double v prefix",
			input:    "vv1.2.3",
			expected: "vv1.2.3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getSemVer(tt.input)
			if result != tt.expected {
				t.Errorf("getSemVer(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}
