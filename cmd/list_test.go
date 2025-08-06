package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/engmtcdrm/mellon/env"
)

// TestListCommand_NoSecrets tests the list command when no secrets exist.
func TestListCommand_NoSecrets(t *testing.T) {
	// Clear all secrets first
	cmd := exec.Command(testBinary, "delete", "--all", "--force")
	_, _ = cmd.CombinedOutput() // Ignore error if no secrets exist

	// Test list with no secrets
	cmd = exec.Command(testBinary, "list")
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Errorf("expected error when no secrets exist, got none")
	}
}

// TestListCommand_WithSecrets tests the list command when secrets exist.
func TestListCommand_WithSecrets(t *testing.T) {
	envVars, err := env.GetEnv()
	if err != nil {
		t.Fatalf("failed to get environment variables: %v", err)
	}

	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	secretContent := "listsecretcontent"

	if err := os.WriteFile(secretFile, []byte(secretContent), 0644); err != nil {
		t.Fatalf("failed to write secret file: %v", err)
	}

	// Create multiple secrets
	secretNames := []string{"list1", "list2", "list3"}
	secretPaths := make([]string, len(secretNames))

	for i, name := range secretNames {
		secretPaths[i] = filepath.Join(envVars.SecretsPath, name+envVars.SecretExt)
		createCmd := exec.Command(testBinary, "create", "--secret", name, "--file", secretFile)
		_, err := createCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("failed to create secret '%s': %v", name, err)
		}
	}

	// Clean up after test
	defer func() {
		for _, path := range secretPaths {
			os.Remove(path)
		}
	}()

	// Test list command
	cmd := exec.Command(testBinary, "list")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("expected success for list command, got error: %v", err)
	}

	// Verify all secret names appear in output
	outputStr := string(output)
	for _, name := range secretNames {
		if !strings.Contains(outputStr, name) {
			t.Errorf("expected output to contain secret name '%s', got: %s", name, outputStr)
		}
	}
}

// TestListCommand_PrintFlag tests the list command with --print flag.
func TestListCommand_PrintFlag(t *testing.T) {
	envVars, err := env.GetEnv()
	if err != nil {
		t.Fatalf("failed to get environment variables: %v", err)
	}

	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	secretContent := "printsecretcontent"

	if err := os.WriteFile(secretFile, []byte(secretContent), 0644); err != nil {
		t.Fatalf("failed to write secret file: %v", err)
	}

	// Create multiple secrets
	secretNames := []string{"print1", "print2", "print3"}
	secretPaths := make([]string, len(secretNames))

	for i, name := range secretNames {
		secretPaths[i] = filepath.Join(envVars.SecretsPath, name+envVars.SecretExt)
		createCmd := exec.Command(testBinary, "create", "--secret", name, "--file", secretFile)
		_, err := createCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("failed to create secret '%s': %v", name, err)
		}
	}

	// Clean up after test
	defer func() {
		for _, path := range secretPaths {
			os.Remove(path)
		}
	}()

	// Test print flag variations
	cases := [][]string{
		{"--print"},
		{"-p"},
	}

	for _, args := range cases {
		cmd := exec.Command(testBinary, append([]string{"list"}, args...)...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("expected success for list with args %v, got error: %v", args, err)
			continue
		}

		// With print flag, output should only contain secret names (minimal formatting)
		outputStr := string(output)
		lines := strings.Split(strings.TrimSpace(outputStr), "\n")

		// Verify all secret names appear as individual lines
		for _, name := range secretNames {
			found := false
			for _, line := range lines {
				if strings.TrimSpace(line) == name {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected to find secret name '%s' as a line in output: %s", name, outputStr)
			}
		}

		// With print flag, output should be minimal (no headers, no decorative text)
		if strings.Contains(outputStr, "Available secrets") {
			t.Errorf("print mode should not contain headers, got: %s", outputStr)
		}
	}
}

// TestListCommand_EmptySecretsWithPrintFlag tests the list command with --print flag when no secrets exist.
func TestListCommand_EmptySecretsWithPrintFlag(t *testing.T) {
	// Clear all secrets first
	cmd := exec.Command(testBinary, "delete", "--all", "--force")
	_, _ = cmd.CombinedOutput() // Ignore error if no secrets exist

	// Test list with print flag and no secrets
	cmd = exec.Command(testBinary, "list", "--print")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("expected success for list --print with no secrets, got error: %v", err)
	}

	// Output should be empty or minimal
	outputStr := strings.TrimSpace(string(output))
	if len(outputStr) > 0 {
		t.Logf("List --print with no secrets produced output: '%s'", outputStr)
	}
}

// TestListCommand_MixedSecretNames tests the list command with various valid secret name patterns.
func TestListCommand_MixedSecretNames(t *testing.T) {
	envVars, err := env.GetEnv()
	if err != nil {
		t.Fatalf("failed to get environment variables: %v", err)
	}

	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	secretContent := "mixedsecretcontent"

	if err := os.WriteFile(secretFile, []byte(secretContent), 0644); err != nil {
		t.Fatalf("failed to write secret file: %v", err)
	}

	// Create secrets with various valid name patterns
	secretNames := []string{
		"simple",
		"with_underscores",
		"with-dashes",
		"with123numbers",
		"MixedCase",
		"path/to/secret",
		"another\\path\\secret",
	}
	secretPaths := make([]string, len(secretNames))

	for i, name := range secretNames {
		secretPaths[i] = filepath.Join(envVars.SecretsPath, name+envVars.SecretExt)
		createCmd := exec.Command(testBinary, "create", "--secret", name, "--file", secretFile)
		_, err := createCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("failed to create secret '%s': %v", name, err)
		}
	}

	// Clean up after test
	defer func() {
		for _, path := range secretPaths {
			os.Remove(path)
		}
	}()

	// Test both normal and print modes
	testCases := []struct {
		name string
		args []string
	}{
		{"normal_mode", []string{}},
		{"print_mode", []string{"--print"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command(testBinary, append([]string{"list"}, tc.args...)...)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Errorf("expected success for list %s, got error: %v", tc.name, err)
				return
			}

			// Verify all secret names appear in output
			outputStr := string(output)
			for _, name := range secretNames {
				if !strings.Contains(outputStr, name) {
					t.Errorf("expected output to contain secret name '%s' in %s, got: %s", name, tc.name, outputStr)
				}
			}
		})
	}
}

// TestListCommand_OutputFormat tests the difference between normal and print mode output.
func TestListCommand_OutputFormat(t *testing.T) {
	envVars, err := env.GetEnv()
	if err != nil {
		t.Fatalf("failed to get environment variables: %v", err)
	}

	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	secretContent := "formatsecretcontent"
	secretName := "formattest"
	secretOut := filepath.Join(envVars.SecretsPath, secretName+envVars.SecretExt)

	if err := os.WriteFile(secretFile, []byte(secretContent), 0644); err != nil {
		t.Fatalf("failed to write secret file: %v", err)
	}

	// Create a secret
	createCmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	_, err = createCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to create secret: %v", err)
	}
	defer os.Remove(secretOut)

	// Test normal mode
	normalCmd := exec.Command(testBinary, "list")
	normalOutput, err := normalCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to run list in normal mode: %v", err)
	}

	// Test print mode
	printCmd := exec.Command(testBinary, "list", "--print")
	printOutput, err := printCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to run list in print mode: %v", err)
	}

	normalStr := string(normalOutput)
	printStr := string(printOutput)

	// Normal mode should have more formatting/headers
	if len(normalStr) <= len(printStr) {
		t.Errorf("expected normal mode output to be longer than print mode")
	}

	// Print mode should just contain the secret name
	if strings.TrimSpace(printStr) != secretName {
		t.Errorf("expected print mode to output just the secret name '%s', got: '%s'", secretName, strings.TrimSpace(printStr))
	}

	// Normal mode should contain additional formatting
	if !strings.Contains(normalStr, "Available secrets") {
		t.Errorf("expected normal mode to contain header text")
	}
}
