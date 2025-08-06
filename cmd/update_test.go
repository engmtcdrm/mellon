package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/engmtcdrm/mellon/env"
)

// TestUpdateCommand_ValidFlags tests the update command with valid flags.
func TestUpdateCommand_ValidFlags(t *testing.T) {
	envVars, err := env.GetEnv()
	if err != nil {
		t.Fatalf("failed to get environment variables: %v", err)
	}
	dir := t.TempDir()
	secretFile := filepath.Join(dir, "secret.txt")
	updateFile := filepath.Join(dir, "update.txt")
	secretName := "testupdatesecret"
	secretContent := "originalsecret"
	updateContent := "updatedsecret"
	secretOut := filepath.Join(envVars.SecretsPath, secretName+envVars.SecretExt)
	// Clean up before test
	os.Remove(secretOut)

	// First create a secret to update
	if err := os.WriteFile(secretFile, []byte(secretContent), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	createCmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	_, err = createCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to create initial secret: %v", err)
	}
	defer os.Remove(secretOut) // Clean up after test

	// Test each permutation of update flags
	cases := [][]string{
		{"--secret", secretName, "--file", updateFile},
		{"--secret", secretName, "-f", updateFile},
		{"-s", secretName, "--file", updateFile},
		{"-s", secretName, "-f", updateFile},
	}

	for _, args := range cases {
		if err := os.WriteFile(updateFile, []byte(updateContent), 0644); err != nil {
			t.Fatalf("failed to write update file: %v", err)
		}

		cmd := exec.Command(testBinary, append([]string{"update"}, args...)...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("expected success, got error: %v, output: %s", err, output)
		}

		if _, err := os.Stat(updateFile); os.IsNotExist(err) {
			t.Errorf("update file should not be deleted")
		}

		if _, err := os.Stat(secretOut); os.IsNotExist(err) {
			t.Errorf("expected output file %s to exist, got error: %v", secretOut, err)
		}
	}
}

func TestUpdateCommand_CleanupFlag(t *testing.T) {
	envVars, err := env.GetEnv()
	if err != nil {
		t.Fatalf("failed to get environment variables: %v", err)
	}
	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	updateFile := filepath.Join(t.TempDir(), "update.txt")
	secretName := "testupdatecleanup"
	secretContent := "originalsecret"
	updateContent := "updatedsecret"
	secretOut := filepath.Join(envVars.SecretsPath, secretName+envVars.SecretExt)
	// Clean up before test
	os.Remove(secretOut)

	// First create a secret to update
	if err := os.WriteFile(secretFile, []byte(secretContent), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	createCmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	_, err = createCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to create initial secret: %v", err)
	}
	defer os.Remove(secretOut) // Clean up after test

	// Test each permutation of flags with --cleanup flag added
	cases := [][]string{
		{"--secret", secretName, "--file", updateFile, "--cleanup"},
		{"--secret", secretName, "--file", updateFile, "-c"},
		{"--secret", secretName, "-f", updateFile, "--cleanup"},
		{"--secret", secretName, "-f", updateFile, "-c"},
		{"-s", secretName, "--file", updateFile, "--cleanup"},
		{"-s", secretName, "--file", updateFile, "-c"},
		{"-s", secretName, "-f", updateFile, "--cleanup"},
		{"-s", secretName, "-f", updateFile, "-c"},
	}

	for _, args := range cases {
		if err := os.WriteFile(updateFile, []byte(updateContent), 0644); err != nil {
			t.Fatalf("failed to write update file: %v", err)
		}

		cmd := exec.Command(testBinary, append([]string{"update"}, args...)...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("expected success, got error: %v, output: %s", err, output)
		}

		if _, err := os.Stat(updateFile); !os.IsNotExist(err) {
			t.Errorf("update file should be deleted with --cleanup")
		}

		if _, err := os.Stat(secretOut); os.IsNotExist(err) {
			t.Errorf("expected output file %s to exist, got error: %v", secretOut, err)
		}
	}
}

func TestUpdateCommand_TildeExpansion(t *testing.T) {
	envVars, err := env.GetEnv()
	if err != nil {
		t.Fatalf("failed to get environment variables: %v", err)
	}
	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	home, _ := os.UserHomeDir()
	updateFile := filepath.Join(home, "updatetilde.txt")
	secretName := "testupdatetilde"
	secretContent := "originalsecret"
	updateContent := "updatedsecret"
	secretOut := filepath.Join(envVars.SecretsPath, secretName+envVars.SecretExt)
	// Clean up before test
	os.Remove(secretOut)

	// First create a secret to update
	if err := os.WriteFile(secretFile, []byte(secretContent), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	createCmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	_, err = createCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to create initial secret: %v", err)
	}
	defer os.Remove(secretOut) // Clean up after test

	if err := os.WriteFile(updateFile, []byte(updateContent), 0644); err != nil {
		t.Fatalf("failed to write update file: %v", err)
	}
	defer os.Remove(updateFile)

	cmd := exec.Command(testBinary, "update", "--secret", secretName, "--file", "~/updatetilde.txt")
	_, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestUpdateCommand_SecretNotExist(t *testing.T) {
	updateFile := filepath.Join(t.TempDir(), "update.txt")
	secretName := "nonexistentsecret"
	updateContent := "updatedsecret"

	if err := os.WriteFile(updateFile, []byte(updateContent), 0644); err != nil {
		t.Fatalf("failed to write update file: %v", err)
	}

	cmd := exec.Command(testBinary, "update", "--secret", secretName, "--file", updateFile)
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Errorf("expected error for non-existent secret, got none")
	}
}

func TestUpdateCommand_FileNotExist(t *testing.T) {
	envVars, err := env.GetEnv()
	if err != nil {
		t.Fatalf("failed to get environment variables: %v", err)
	}
	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	updateFile := filepath.Join(t.TempDir(), "doesnotexist.txt")
	secretName := "testupdatenofile"
	secretContent := "originalsecret"
	secretOut := filepath.Join(envVars.SecretsPath, secretName+envVars.SecretExt)
	// Clean up before test
	os.Remove(secretOut)

	// First create a secret to update
	if err := os.WriteFile(secretFile, []byte(secretContent), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	createCmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	_, err = createCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to create initial secret: %v", err)
	}
	defer os.Remove(secretOut) // Clean up after test

	cmd := exec.Command(testBinary, "update", "--secret", secretName, "--file", updateFile)
	_, err = cmd.CombinedOutput()
	if err == nil {
		t.Errorf("expected error for non-existent update file, got none")
	}
}

func TestUpdateCommand_FileNoReadAccess(t *testing.T) {
	envVars, err := env.GetEnv()
	if err != nil {
		t.Fatalf("failed to get environment variables: %v", err)
	}
	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	updateFile := filepath.Join(t.TempDir(), "noread.txt")
	secretName := "testupdatenoread"
	secretContent := "originalsecret"
	updateContent := "updatedsecret"
	secretOut := filepath.Join(envVars.SecretsPath, secretName+envVars.SecretExt)
	// Clean up before test
	os.Remove(secretOut)

	// First create a secret to update
	if err := os.WriteFile(secretFile, []byte(secretContent), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	createCmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	_, err = createCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to create initial secret: %v", err)
	}
	defer os.Remove(secretOut) // Clean up after test

	if err := os.WriteFile(updateFile, []byte(updateContent), 0000); err != nil {
		t.Fatalf("failed to write update file: %v", err)
	}
	defer os.Chmod(updateFile, 0644)

	cmd := exec.Command(testBinary, "update", "--secret", secretName, "--file", updateFile)
	_, err = cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected error for no read access, got none")
	}
}

func TestUpdateCommand_CleanupNoWriteAccess(t *testing.T) {
	envVars, err := env.GetEnv()
	if err != nil {
		t.Fatalf("failed to get environment variables: %v", err)
	}
	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	dir := t.TempDir()
	updateFile := filepath.Join(dir, "nowrite.txt")
	secretName := "testupdatenowrite"
	secretContent := "originalsecret"
	updateContent := "updatedsecret"
	secretOut := filepath.Join(envVars.SecretsPath, secretName+envVars.SecretExt)
	// Clean up before test
	os.Remove(secretOut)

	// First create a secret to update
	if err := os.WriteFile(secretFile, []byte(secretContent), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	createCmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	_, err = createCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to create initial secret: %v", err)
	}
	defer os.Remove(secretOut) // Clean up after test

	if err := os.WriteFile(updateFile, []byte(updateContent), 0444); err != nil {
		t.Fatalf("failed to write update file: %v", err)
	}

	// Remove write permission from the directory to prevent file deletion
	if err := os.Chmod(dir, 0555); err != nil {
		t.Fatalf("failed to remove write permission from dir: %v", err)
	}
	defer os.Chmod(dir, 0755)

	cmd := exec.Command(testBinary, "update", "--secret", secretName, "--file", updateFile, "--cleanup")
	_, err = cmd.CombinedOutput()
	if err == nil {
		t.Errorf("expected error for no write access to directory, got none")
	}
}

func TestUpdateCommand_CleanupNoReadWriteAccess(t *testing.T) {
	envVars, err := env.GetEnv()
	if err != nil {
		t.Fatalf("failed to get environment variables: %v", err)
	}
	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	updateFile := filepath.Join(t.TempDir(), "noreadwrite.txt")
	secretName := "testupdatenoreadwrite"
	secretContent := "originalsecret"
	updateContent := "updatedsecret"
	secretOut := filepath.Join(envVars.SecretsPath, secretName+envVars.SecretExt)
	// Clean up before test
	os.Remove(secretOut)

	// First create a secret to update
	if err := os.WriteFile(secretFile, []byte(secretContent), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	createCmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	_, err = createCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to create initial secret: %v", err)
	}
	defer os.Remove(secretOut) // Clean up after test

	if err := os.WriteFile(updateFile, []byte(updateContent), 0000); err != nil {
		t.Fatalf("failed to write update file: %v", err)
	}
	defer os.Chmod(updateFile, 0644)

	cmd := exec.Command(testBinary, "update", "--secret", secretName, "--file", updateFile, "--cleanup")
	_, err = cmd.CombinedOutput()
	if err == nil {
		t.Errorf("expected error for no read/write access, got none")
	}
}

// TestUpdateCommand_PreRunValidation tests the PreRunE validation logic
func TestUpdateCommand_PreRunValidation(t *testing.T) {
	// Test cleanup flag without required flags
	cases := []struct {
		name string
		args []string
	}{
		{
			name: "cleanup without secret flag",
			args: []string{"update", "--cleanup", "--file", "somefile.txt"},
		},
		{
			name: "cleanup without file flag",
			args: []string{"update", "--cleanup", "--secret", "somesecret"},
		},
		{
			name: "cleanup without both flags",
			args: []string{"update", "--cleanup"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command(testBinary, tc.args...)
			_, err := cmd.CombinedOutput()
			if err == nil {
				t.Errorf("expected error for %s, got none", tc.name)
			}
		})
	}
}

// TestUpdateCommand_MissingFlags tests missing required flag combinations
func TestUpdateCommand_MissingFlags(t *testing.T) {
	envVars, err := env.GetEnv()
	if err != nil {
		t.Fatalf("failed to get environment variables: %v", err)
	}
	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	updateFile := filepath.Join(t.TempDir(), "update.txt")
	secretName := "testupdatemissing"
	secretContent := "originalsecret"
	updateContent := "updatedsecret"
	secretOut := filepath.Join(envVars.SecretsPath, secretName+envVars.SecretExt)
	// Clean up before test
	os.Remove(secretOut)

	// First create a secret to update
	if err := os.WriteFile(secretFile, []byte(secretContent), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	createCmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	_, err = createCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to create initial secret: %v", err)
	}
	defer os.Remove(secretOut) // Clean up after test

	if err := os.WriteFile(updateFile, []byte(updateContent), 0644); err != nil {
		t.Fatalf("failed to write update file: %v", err)
	}

	// Test cases where only one flag is provided (should work for interactive mode, but we skip those)
	// The update command allows partial flags for interactive mode, so these should not error
	cases := []struct {
		name        string
		args        []string
		shouldError bool
	}{
		{
			name:        "only file flag",
			args:        []string{"--file", updateFile},
			shouldError: false, // Should work - interactive secret selection
		},
		{
			name:        "only secret flag",
			args:        []string{"--secret", secretName},
			shouldError: false, // Should work - interactive content entry
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Skip interactive tests as requested
			t.Skip("Skipping interactive test: " + tc.name)
		})
	}
}

// TestUpdateCommand_InvalidSecretName tests updating with invalid secret names
func TestUpdateCommand_InvalidSecretName(t *testing.T) {
	updateFile := filepath.Join(t.TempDir(), "update.txt")
	secretName := "invalid!name"
	updateContent := "updatedsecret"

	if err := os.WriteFile(updateFile, []byte(updateContent), 0644); err != nil {
		t.Fatalf("failed to write update file: %v", err)
	}

	cmd := exec.Command(testBinary, "update", "--secret", secretName, "--file", updateFile)
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Errorf("expected error for invalid secret name, got none")
	}
}

// TestUpdateCommand_ValidSecretNames tests updating with valid secret name patterns
func TestUpdateCommand_ValidSecretNames(t *testing.T) {
	envVars, err := env.GetEnv()
	if err != nil {
		t.Fatalf("failed to get environment variables: %v", err)
	}

	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	updateFile := filepath.Join(t.TempDir(), "update.txt")
	secretContent := "originalsecret"
	updateContent := "updatedsecret"

	if err := os.WriteFile(secretFile, []byte(secretContent), 0644); err != nil {
		t.Fatalf("failed to write secret file: %v", err)
	}

	if err := os.WriteFile(updateFile, []byte(updateContent), 0644); err != nil {
		t.Fatalf("failed to write update file: %v", err)
	}

	validNames := []string{
		"simple",
		"with_underscores",
		"with-dashes",
		"with123numbers",
		"MixedCase",
		"path/to/secret",
	}

	for _, validName := range validNames {
		t.Run(fmt.Sprintf("valid_name_%s", validName), func(t *testing.T) {
			secretOut := filepath.Join(envVars.SecretsPath, validName+envVars.SecretExt)
			defer os.Remove(secretOut) // Clean up

			// First create the secret
			createCmd := exec.Command(testBinary, "create", "--secret", validName, "--file", secretFile)
			_, err := createCmd.CombinedOutput()
			if err != nil {
				t.Fatalf("failed to create initial secret '%s': %v", validName, err)
			}

			// Then update it
			updateCmd := exec.Command(testBinary, "update", "--secret", validName, "--file", updateFile)
			_, err = updateCmd.CombinedOutput()
			if err != nil {
				t.Errorf("expected success for valid secret name '%s', got error: %v", validName, err)
			}

			// Verify the file still exists
			if _, err := os.Stat(secretOut); os.IsNotExist(err) {
				t.Errorf("expected output file %s to exist for name '%s'", secretOut, validName)
			}
		})
	}
}
