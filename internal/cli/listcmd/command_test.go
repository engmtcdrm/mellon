package listcmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/engmtcdrm/mellon/internal/env"
	"github.com/stretchr/testify/assert"
)

var testBinary string

// TestMain builds the CLI binary once for all tests and cleans up after.
func TestMain(m *testing.M) {
	testBinary = filepath.Join(os.TempDir(), "mellon-test-bin-listcmd")
	projectRoot, err := filepath.Abs(filepath.Join("..", "..", ".."))
	if err != nil {
		panic("failed to determine project root: " + err.Error())
	}

	// Build the test binary
	cmd := exec.Command("go", "build", "-o", testBinary, ".")
	cmd.Dir = projectRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic("failed to build test binary: " + err.Error() + "\n" + string(out))
	}

	code := m.Run()

	// Clean up the test binary after tests
	os.Remove(testBinary)
	os.Exit(code)
}

// TestListCommand_NoSecrets tests the list command when no secrets exist.
func TestListCommand_NoSecrets(t *testing.T) {
	// Clear all secrets first
	cmd := exec.Command(testBinary, "delete", "--all", "--force")
	_, _ = cmd.CombinedOutput() // Ignore error if no secrets exist

	// Test list with no secrets
	cmd = exec.Command(testBinary, "list")
	_, err := cmd.CombinedOutput()
	assert.NoError(t, err)
}

// TestListCommand_WithSecrets tests the list command when secrets exist.
func TestListCommand_WithSecrets(t *testing.T) {
	env.Init()

	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	secretContent := "listsecretcontent"

	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	assert.NoError(t, err, "failed to write secret file")

	// Create multiple secrets
	secretNames := []string{"list1", "list2", "list3"}
	secretPaths := make([]string, len(secretNames))

	for i, name := range secretNames {
		secretPaths[i] = filepath.Join(env.Instance.SecretsPath(), name+env.Instance.SecretExt())
		createCmd := exec.Command(testBinary, "create", "--secret", name, "--file", secretFile)
		_, err := createCmd.CombinedOutput()
		assert.NoErrorf(t, err, "failed to create secret '%s': %v", name, err)
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
	assert.NoErrorf(t, err, "expected success for list command, got error: %v", err)

	// Verify all secret names appear in output
	outputStr := string(output)
	for _, name := range secretNames {
		assert.Contains(t, outputStr, name, "expected output to contain secret name '%s', got: %s", name, outputStr)
	}
}

// TestListCommand_PrintFlag tests the list command with --print flag.
func TestListCommand_PrintFlag(t *testing.T) {
	env.Init()

	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	secretContent := "printsecretcontent"

	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	assert.NoError(t, err, "failed to write secret file")

	// Create multiple secrets
	secretNames := []string{"print1", "print2", "print3"}
	secretPaths := make([]string, len(secretNames))

	for i, name := range secretNames {
		secretPaths[i] = filepath.Join(env.Instance.SecretsPath(), name+env.Instance.SecretExt())
		createCmd := exec.Command(testBinary, "create", "--secret", name, "--file", secretFile)
		_, err := createCmd.CombinedOutput()
		assert.NoErrorf(t, err, "failed to create secret '%s': %v", name, err)
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
		assert.NoErrorf(t, err, "expected success for list with args %v, got error: %v", args, err)

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
			assert.True(t, found, "expected to find secret name '%s' as a line in output: %s", name, outputStr)
		}

		// With print flag, output should be minimal (no headers, no decorative text)
		assert.NotContains(t, outputStr, "Available secrets", "print mode should not contain headers")
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
	assert.NoErrorf(t, err, "expected success for list --print with no secrets, got error: %v", err)

	// Output should be empty or minimal
	outputStr := strings.TrimSpace(string(output))
	if len(outputStr) > 0 {
		t.Logf("List --print with no secrets produced output: '%s'", outputStr)
	}
}

// TestListCommand_MixedSecretNames tests the list command with various valid secret name patterns.
func TestListCommand_MixedSecretNames(t *testing.T) {
	env.Init()

	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	secretContent := "mixedsecretcontent"

	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	assert.NoError(t, err, "failed to write secret file")

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
		secretPaths[i] = filepath.Join(env.Instance.SecretsPath(), name+env.Instance.SecretExt())
		createCmd := exec.Command(testBinary, "create", "--secret", name, "--file", secretFile)
		_, err := createCmd.CombinedOutput()
		assert.NoErrorf(t, err, "failed to create secret '%s': %v", name, err)
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
			assert.NoErrorf(t, err, "expected success for list %s, got error: %v", tc.name, err)

			// Verify all secret names appear in output
			outputStr := string(output)
			for _, name := range secretNames {
				assert.Contains(t, outputStr, name, "expected output to contain secret name '%s' in %s", name, tc.name)
			}
		})
	}
}

// TestListCommand_OutputFormat tests the difference between normal and print mode output.
func TestListCommand_OutputFormat(t *testing.T) {
	env.Init()

	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	secretContent := "formatsecretcontent"
	secretName := "formattest"
	secretOut := filepath.Join(env.Instance.SecretsPath(), secretName+env.Instance.SecretExt())

	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	assert.NoError(t, err, "failed to write secret file")

	// Create a secret
	createCmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	_, err = createCmd.CombinedOutput()
	assert.NoErrorf(t, err, "failed to create secret: %v", err)
	defer os.Remove(secretOut)

	// Test normal mode
	normalCmd := exec.Command(testBinary, "list")
	normalOutput, err := normalCmd.CombinedOutput()
	assert.NoErrorf(t, err, "failed to run list in normal mode: %v", err)

	// Test print mode
	printCmd := exec.Command(testBinary, "list", "--print")
	printOutput, err := printCmd.CombinedOutput()
	assert.NoErrorf(t, err, "failed to run list in print mode: %v", err)

	normalStr := string(normalOutput)
	printStr := string(printOutput)

	// Normal mode should have more formatting/headers
	assert.LessOrEqual(t, len(printStr), len(normalStr), "expected normal mode output to be longer than print mode")

	// Print mode should just contain the secret name
	assert.Equalf(t, strings.TrimSpace(printStr), secretName, "expected print mode to output just the secret name '%s', got: '%s'", secretName, strings.TrimSpace(printStr))

	// Normal mode should contain additional formatting
	assert.Contains(t, normalStr, "Available secrets", "expected normal mode to contain header text")
}
