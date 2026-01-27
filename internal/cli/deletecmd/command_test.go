package deletecmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/engmtcdrm/mellon/internal/env"
	"github.com/stretchr/testify/assert"
)

var testBinary string

// TestMain builds the CLI binary once for all tests and cleans up after.
func TestMain(m *testing.M) {
	testBinary = filepath.Join(os.TempDir(), "mellon-test-bin-deletecmd")
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

// TestDeleteCommand_ValidFlags tests the delete command with valid flags.
func TestDeleteCommand_ValidFlags(t *testing.T) {
	env.Init()
	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	secretName := "testdeletesecret"
	secretContent := "supersecret"
	secretOut := filepath.Join(env.Instance.SecretsPath(), secretName+env.Instance.SecretExt())
	// Clean up before test
	os.Remove(secretOut)

	// First create a secret to delete
	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	assert.NoErrorf(t, err, "failed to write temp file: %v", err)

	createCmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	output, err := createCmd.CombinedOutput()
	assert.NoErrorf(t, err, "failed to create initial secret: %v, output: %s", err, output)

	// Verify the secret exists before deletion
	_, err = os.Stat(secretOut)
	assert.NoErrorf(t, err, "secret should exist before deletion")

	// Test each permutation of delete flags with force to avoid prompts
	cases := [][]string{
		{"--secret", secretName, "--force"},
		{"-s", secretName, "--force"},
		{"--secret", secretName, "-f"},
		{"-s", secretName, "-f"},
	}

	for i, args := range cases {
		// Recreate the secret for each test case (except first)
		if i > 0 {
			createCmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
			output, err = createCmd.CombinedOutput()
			assert.NoErrorf(t, err, "failed to recreate secret for test case %d: %v, output: %s", i, err, output)
		}

		cmd := exec.Command(testBinary, append([]string{"delete"}, args...)...)
		output, err := cmd.CombinedOutput()
		assert.NoErrorf(t, err, "expected success for args %v, got error: %v, output: %s", args, err, output)

		// Verify the secret was deleted
		_, err = os.Stat(secretOut)
		assert.NoFileExistsf(t, secretOut, "secret should be deleted after delete command")
	}
}

func TestDeleteCommand_SecretNotExist(t *testing.T) {
	secretName := "nonexistentsecret"

	cmd := exec.Command(testBinary, "delete", "--secret", secretName, "--force")
	output, err := cmd.CombinedOutput()
	assert.Errorf(t, err, "expected error for non-existent secret, got none, output: %s", output)
}

func TestDeleteCommand_InvalidSecretName(t *testing.T) {
	secretName := "invalid!name"

	cmd := exec.Command(testBinary, "delete", "--secret", secretName, "--force")
	output, err := cmd.CombinedOutput()
	assert.Errorf(t, err, "expected error for invalid secret name, got none, output: %s", output)
}

func TestDeleteCommand_ValidSecretNames(t *testing.T) {
	env.Init()
	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	secretContent := "supersecret"

	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	assert.NoErrorf(t, err, "failed to write secret file: %v", err)

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
			secretOut := filepath.Join(env.Instance.SecretsPath(), validName+env.Instance.SecretExt())
			defer os.Remove(secretOut) // Clean up

			// First create the secret
			createCmd := exec.Command(testBinary, "create", "--secret", validName, "--file", secretFile)
			output, err := createCmd.CombinedOutput()
			assert.NoErrorf(t, err, "failed to create initial secret '%s': %v, output: %s", validName, err, output)

			// Verify it exists
			_, err = os.Stat(secretOut)
			assert.NoErrorf(t, err, "expected secret '%s' to exist before deletion, got error: %v", validName, err)
			assert.FileExistsf(t, secretOut, "secret should exist before deletion")

			// Then delete it
			deleteCmd := exec.Command(testBinary, "delete", "--secret", validName, "--force")
			output, err = deleteCmd.CombinedOutput()
			assert.NoErrorf(t, err, "failed to delete secret '%s': %v, output: %s", validName, err, output)

			// Verify it was deleted
			_, err = os.Stat(secretOut)
			assert.Errorf(t, err, "expected secret '%s' to be deleted, got error: %v", validName, err)
			assert.NoFileExistsf(t, secretOut, "secret should be deleted for name '%s'", validName)
		})
	}
}

func TestDeleteCommand_AllFlag(t *testing.T) {
	env.Init()
	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	secretContent := "supersecret"

	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	assert.NoErrorf(t, err, "failed to write secret file: %v", err)

	// Create multiple secrets
	secretNames := []string{"test1", "test2", "test3"}
	secretPaths := make([]string, len(secretNames))

	for i, name := range secretNames {
		secretPaths[i] = filepath.Join(env.Instance.SecretsPath(), name+env.Instance.SecretExt())
		createCmd := exec.Command(testBinary, "create", "--secret", name, "--file", secretFile)
		output, err := createCmd.CombinedOutput()
		assert.NoErrorf(t, err, "failed to create secret '%s': %v, output: %s", name, err, output)

		// Verify it exists
		_, err = os.Stat(secretPaths[i])
		assert.NoErrorf(t, err, "secret '%s' should exist before deletion", name)
		assert.FileExistsf(t, secretPaths[i], "secret '%s' should exist before deletion", name)
	}

	// Delete all secrets with force flag
	cmd := exec.Command(testBinary, "delete", "--all", "--force")
	output, err := cmd.CombinedOutput()
	assert.NoErrorf(t, err, "expected success for delete all, got error: %v, output: %s", err, output)

	// Verify all secrets were deleted
	for i, name := range secretNames {
		_, err := os.Stat(secretPaths[i])
		assert.Errorf(t, err, "expected secret '%s' to be deleted, got error: %v", name, err)
		assert.NoFileExistsf(t, secretPaths[i], "secret '%s' should be deleted", name)
	}
}

func TestDeleteCommand_AllFlagEmptySecrets(t *testing.T) {
	// Test delete all when no secrets exist
	cmd := exec.Command(testBinary, "delete", "--all", "--force")
	output, err := cmd.CombinedOutput()
	assert.NoErrorf(t, err, "expected success for delete all with no secrets, got error: %v, output: %s", err, output)
}

func TestDeleteCommand_MutuallyExclusiveFlags(t *testing.T) {
	env.Init()
	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	secretName := "testmutex"
	secretContent := "supersecret"
	secretOut := filepath.Join(env.Instance.SecretsPath(), secretName+env.Instance.SecretExt())

	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	assert.NoErrorf(t, err, "failed to write secret file: %v", err)

	// Create a secret first
	createCmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	output, err := createCmd.CombinedOutput()
	assert.NoErrorf(t, err, "failed to create secret '%s': %v, output: %s", secretName, err, output)
	defer os.Remove(secretOut)

	// Test that --secret and --all are mutually exclusive
	cmd := exec.Command(testBinary, "delete", "--secret", secretName, "--all", "--force")
	output, err = cmd.CombinedOutput()
	assert.Errorf(t, err, "expected error for mutually exclusive flags --secret and --all, got none, output: %s", output)
}

func TestDeleteCommand_NoFlags(t *testing.T) {
	// Test delete command without any flags (should enter interactive mode, but we skip this)
	t.Skip("Skipping interactive test: delete without flags")
}

func TestDeleteCommand_ForceFlag(t *testing.T) {
	env.Init()
	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	secretName := "testforce"
	secretContent := "supersecret"
	secretOut := filepath.Join(env.Instance.SecretsPath(), secretName+env.Instance.SecretExt())

	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	assert.NoErrorf(t, err, "failed to write secret file: %v", err)

	// Create a secret first
	createCmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	output, err := createCmd.CombinedOutput()
	assert.NoErrorf(t, err, "failed to create initial secret: %v, output: %s", err, output)

	// Test force flag (both short and long form)
	cases := [][]string{
		{"--secret", secretName, "--force"},
		{"-s", secretName, "-f"},
	}

	for i, args := range cases {
		// Recreate the secret for second test case
		if i > 0 {
			createCmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
			output, err = createCmd.CombinedOutput()
			assert.NoErrorf(t, err, "failed to recreate secret for test case %d: %v, output: %s", i, err, output)
		}

		cmd := exec.Command(testBinary, append([]string{"delete"}, args...)...)
		output, err := cmd.CombinedOutput()
		assert.NoErrorf(t, err, "expected success for force delete with args %v, got error: %v, output: %s", args, err, output)

		// Verify the secret was deleted
		_, err = os.Stat(secretOut)
		assert.Errorf(t, err, "expected secret '%s' to be deleted with force flag, got error: %v", secretName, err)
		assert.NoFileExistsf(t, secretOut, "secret '%s' should be deleted with force flag", secretName)
	}
}

func TestDeleteCommand_WithoutForceFlag(t *testing.T) {
	// Test delete without force flag (should enter interactive mode, but we skip this)
	t.Skip("Skipping interactive test: delete without force flag")
}

func TestDeleteCommand_SilentMode(t *testing.T) {
	env.Init()
	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	secretName := "testsilent"
	secretContent := "supersecret"
	secretOut := filepath.Join(env.Instance.SecretsPath(), secretName+env.Instance.SecretExt())

	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	assert.NoErrorf(t, err, "failed to write secret file: %v", err)

	// Create a secret first
	createCmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	output, err := createCmd.CombinedOutput()
	assert.NoErrorf(t, err, "failed to create initial secret: %v, output: %s", err, output)
	defer os.Remove(secretOut)

	// Test that force flag suppresses output
	cmd := exec.Command(testBinary, "delete", "--secret", secretName, "--force")
	output, err = cmd.CombinedOutput()
	assert.NoErrorf(t, err, "expected success for silent delete, got error: %v, output: %s", err, output)

	// With force flag, output should be minimal (no success messages)
	outputStr := string(output)
	assert.Lessf(t, len(outputStr), 100, "Force delete output (should be minimal): %s", outputStr)

	// Verify the secret was deleted
	_, err = os.Stat(secretOut)
	assert.Errorf(t, err, "expected secret '%s' to be deleted with force flag, got error: %v", secretName, err)
	assert.NoFileExistsf(t, secretOut, "secret '%s' should be deleted with force flag", secretName)
}

func TestDeleteCommand_AllFlagSilentMode(t *testing.T) {
	env.Init()
	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	secretContent := "supersecret"

	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	assert.NoErrorf(t, err, "failed to write secret file: %v", err)

	// Create multiple secrets
	secretNames := []string{"silent1", "silent2"}
	secretPaths := make([]string, len(secretNames))

	for i, name := range secretNames {
		secretPaths[i] = filepath.Join(env.Instance.SecretsPath(), name+env.Instance.SecretExt())
		createCmd := exec.Command(testBinary, "create", "--secret", name, "--file", secretFile)
		output, err := createCmd.CombinedOutput()
		assert.NoErrorf(t, err, "failed to create secret '%s': %v, output: %s", name, err, output)
	}

	// Delete all secrets with force flag (silent mode)
	cmd := exec.Command(testBinary, "delete", "--all", "--force")
	output, err := cmd.CombinedOutput()
	assert.NoErrorf(t, err, "expected success for silent delete all, got error: %v, output: %s", err, output)

	// With force flag, output should be minimal
	outputStr := string(output)
	assert.Lessf(t, len(outputStr), 100, "Force delete all output (should be minimal): %s", outputStr)

	// Verify all secrets were deleted
	for i, name := range secretNames {
		_, err := os.Stat(secretPaths[i])
		assert.Errorf(t, err, "expected secret '%s' to be deleted with force flag, got error: %v", name, err)
		assert.NoFileExistsf(t, secretPaths[i], "secret '%s' should be deleted with force flag", name)
	}
}

func TestDeleteCommand_EdgeCases(t *testing.T) {
	// Test edge case: empty secret name
	cmd := exec.Command(testBinary, "delete", "--secret", "", "--force")
	output, err := cmd.CombinedOutput()
	assert.Errorf(t, err, "expected error for empty secret name, got none, output: %s", output)

	// Test edge case: only force flag - should enter interactive mode
	t.Skip("Skipping interactive test: delete with only force flag")
}
