package updatecmd

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
	testBinary = filepath.Join(os.TempDir(), "mellon-test-bin-updatecmd")
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

// TestUpdateCommand_ValidFlags tests the update command with valid flags.
func TestUpdateCommand_ValidFlags(t *testing.T) {
	env.Init()
	dir := t.TempDir()
	secretFile := filepath.Join(dir, "secret.txt")
	updateFile := filepath.Join(dir, "update.txt")
	secretName := "testupdatesecret"
	secretContent := "originalsecret"
	updateContent := "updatedsecret"
	secretOut := filepath.Join(env.Instance.SecretsPath(), secretName+env.Instance.SecretExt())
	// Clean up before test
	os.Remove(secretOut)

	// First create a secret to update
	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	assert.NoErrorf(t, err, "failed to write temp file: %v", err)

	createCmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	_, err = createCmd.CombinedOutput()
	assert.NoErrorf(t, err, "failed to create initial secret: %v", err)
	defer os.Remove(secretOut) // Clean up after test

	// Test each permutation of update flags
	cases := [][]string{
		{"--secret", secretName, "--file", updateFile},
		{"--secret", secretName, "-f", updateFile},
		{"-s", secretName, "--file", updateFile},
		{"-s", secretName, "-f", updateFile},
	}

	for _, args := range cases {
		err := os.WriteFile(updateFile, []byte(updateContent), 0644)
		assert.NoError(t, err, "failed to write update file")

		cmd := exec.Command(testBinary, append([]string{"update"}, args...)...)
		output, err := cmd.CombinedOutput()
		assert.NoErrorf(t, err, "expected success, got error: %v, output: %s", err, output)

		assert.FileExists(t, updateFile, "update file should not be deleted")

		assert.FileExistsf(t, secretOut, "expected output file %s to exist", secretOut)
	}
}

func TestUpdateCommand_CleanupFlag(t *testing.T) {
	env.Init()
	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	updateFile := filepath.Join(t.TempDir(), "update.txt")
	secretName := "testupdatecleanup"
	secretContent := "originalsecret"
	updateContent := "updatedsecret"
	secretOut := filepath.Join(env.Instance.SecretsPath(), secretName+env.Instance.SecretExt())
	// Clean up before test
	os.Remove(secretOut)

	// First create a secret to update
	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	assert.NoErrorf(t, err, "failed to write temp file: %v", err)

	createCmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	_, err = createCmd.CombinedOutput()
	assert.NoErrorf(t, err, "failed to create initial secret: %v", err)
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
		err := os.WriteFile(updateFile, []byte(updateContent), 0644)
		assert.NoErrorf(t, err, "failed to write update file: %v", err)

		cmd := exec.Command(testBinary, append([]string{"update"}, args...)...)
		output, err := cmd.CombinedOutput()
		assert.NoErrorf(t, err, "expected success, got error: %v, output: %s", err, output)

		assert.NoFileExists(t, updateFile, "update file should be deleted with --cleanup")

		assert.FileExistsf(t, secretOut, "expected output file %s to exist", secretOut)
	}
}

func TestUpdateCommand_TildeExpansion(t *testing.T) {
	env.Init()
	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	home, _ := os.UserHomeDir()
	updateFile := filepath.Join(home, "updatetilde.txt")
	secretName := "testupdatetilde"
	secretContent := "originalsecret"
	updateContent := "updatedsecret"
	secretOut := filepath.Join(env.Instance.SecretsPath(), secretName+env.Instance.SecretExt())
	// Clean up before test
	os.Remove(secretOut)

	// First create a secret to update
	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	assert.NoErrorf(t, err, "failed to write temp file: %v", err)

	createCmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	_, err = createCmd.CombinedOutput()
	assert.NoErrorf(t, err, "failed to create initial secret: %v", err)
	defer os.Remove(secretOut) // Clean up after test

	err = os.WriteFile(updateFile, []byte(updateContent), 0644)
	assert.NoErrorf(t, err, "failed to write update file: %v", err)
	defer os.Remove(updateFile)

	cmd := exec.Command(testBinary, "update", "--secret", secretName, "--file", "~/updatetilde.txt")
	_, err = cmd.CombinedOutput()
	assert.NoErrorf(t, err, "expected success, got error: %v", err)
}

func TestUpdateCommand_SecretNotExist(t *testing.T) {
	updateFile := filepath.Join(t.TempDir(), "update.txt")
	secretName := "nonexistentsecret"
	updateContent := "updatedsecret"

	err := os.WriteFile(updateFile, []byte(updateContent), 0644)
	assert.NoErrorf(t, err, "failed to write update file: %v", err)

	cmd := exec.Command(testBinary, "update", "--secret", secretName, "--file", updateFile)
	_, err = cmd.CombinedOutput()
	assert.Error(t, err, "expected error for non-existent secret, got none")
}

func TestUpdateCommand_FileNotExist(t *testing.T) {
	env.Init()
	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	updateFile := filepath.Join(t.TempDir(), "doesnotexist.txt")
	secretName := "testupdatenofile"
	secretContent := "originalsecret"
	secretOut := filepath.Join(env.Instance.SecretsPath(), secretName+env.Instance.SecretExt())
	// Clean up before test
	os.Remove(secretOut)

	// First create a secret to update
	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	assert.NoErrorf(t, err, "failed to write temp file: %v", err)

	createCmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	_, err = createCmd.CombinedOutput()
	assert.NoErrorf(t, err, "failed to create initial secret: %v", err)
	defer os.Remove(secretOut) // Clean up after test

	cmd := exec.Command(testBinary, "update", "--secret", secretName, "--file", updateFile)
	_, err = cmd.CombinedOutput()
	assert.Error(t, err, "expected error for non-existent update file, got none")
}

func TestUpdateCommand_FileNoReadAccess(t *testing.T) {
	env.Init()
	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	updateFile := filepath.Join(t.TempDir(), "noread.txt")
	secretName := "testupdatenoread"
	secretContent := "originalsecret"
	updateContent := "updatedsecret"
	secretOut := filepath.Join(env.Instance.SecretsPath(), secretName+env.Instance.SecretExt())
	// Clean up before test
	os.Remove(secretOut)

	// First create a secret to update
	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	assert.NoErrorf(t, err, "failed to write temp file: %v", err)

	createCmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	_, err = createCmd.CombinedOutput()
	assert.NoErrorf(t, err, "failed to create initial secret: %v", err)
	defer os.Remove(secretOut) // Clean up after test

	err = os.WriteFile(updateFile, []byte(updateContent), 0000)
	assert.NoErrorf(t, err, "failed to write update file: %v", err)
	defer os.Chmod(updateFile, 0644)

	cmd := exec.Command(testBinary, "update", "--secret", secretName, "--file", updateFile)
	_, err = cmd.CombinedOutput()
	assert.Error(t, err, "expected error for no read access")
}

func TestUpdateCommand_CleanupNoWriteAccess(t *testing.T) {
	env.Init()
	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	dir := t.TempDir()
	updateFile := filepath.Join(dir, "nowrite.txt")
	secretName := "testupdatenowrite"
	secretContent := "originalsecret"
	updateContent := "updatedsecret"
	secretOut := filepath.Join(env.Instance.SecretsPath(), secretName+env.Instance.SecretExt())
	// Clean up before test
	os.Remove(secretOut)

	// First create a secret to update
	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	assert.NoErrorf(t, err, "failed to write temp file: %v", err)

	createCmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	_, err = createCmd.CombinedOutput()
	assert.NoErrorf(t, err, "failed to create initial secret: %v", err)
	defer os.Remove(secretOut) // Clean up after test

	err = os.WriteFile(updateFile, []byte(updateContent), 0444)
	assert.NoErrorf(t, err, "failed to write update file: %v", err)

	// Remove write permission from the directory to prevent file deletion
	err = os.Chmod(dir, 0555)
	assert.NoErrorf(t, err, "failed to remove write permission from dir: %v", err)
	defer os.Chmod(dir, 0755)

	cmd := exec.Command(testBinary, "update", "--secret", secretName, "--file", updateFile, "--cleanup")
	_, err = cmd.CombinedOutput()
	assert.Error(t, err, "expected error for no write access to directory, got none")
}

func TestUpdateCommand_CleanupNoReadWriteAccess(t *testing.T) {
	env.Init()
	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	updateFile := filepath.Join(t.TempDir(), "noreadwrite.txt")
	secretName := "testupdatenoreadwrite"
	secretContent := "originalsecret"
	updateContent := "updatedsecret"
	secretOut := filepath.Join(env.Instance.SecretsPath(), secretName+env.Instance.SecretExt())
	// Clean up before test
	os.Remove(secretOut)

	// First create a secret to update
	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	assert.NoErrorf(t, err, "failed to write temp file: %v", err)

	createCmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	_, err = createCmd.CombinedOutput()
	assert.NoErrorf(t, err, "failed to create initial secret: %v", err)
	defer os.Remove(secretOut) // Clean up after test

	err = os.WriteFile(updateFile, []byte(updateContent), 0000)
	assert.NoErrorf(t, err, "failed to write update file: %v", err)
	defer os.Chmod(updateFile, 0644)

	cmd := exec.Command(testBinary, "update", "--secret", secretName, "--file", updateFile, "--cleanup")
	_, err = cmd.CombinedOutput()
	assert.Error(t, err, "expected error for no read/write access, got none")
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
			assert.Errorf(t, err, "expected error for %s, got none", tc.name)
		})
	}
}

// TestUpdateCommand_MissingFlags tests missing required flag combinations
func TestUpdateCommand_MissingFlags(t *testing.T) {
	env.Init()
	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	updateFile := filepath.Join(t.TempDir(), "update.txt")
	secretName := "testupdatemissing"
	secretContent := "originalsecret"
	updateContent := "updatedsecret"
	secretOut := filepath.Join(env.Instance.SecretsPath(), secretName+env.Instance.SecretExt())
	// Clean up before test
	os.Remove(secretOut)

	// First create a secret to update
	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	assert.NoErrorf(t, err, "failed to write temp file: %v", err)

	createCmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	_, err = createCmd.CombinedOutput()
	assert.NoErrorf(t, err, "failed to create initial secret: %v", err)
	defer os.Remove(secretOut) // Clean up after test

	err = os.WriteFile(updateFile, []byte(updateContent), 0644)
	assert.NoErrorf(t, err, "failed to write update file: %v", err)

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

	err := os.WriteFile(updateFile, []byte(updateContent), 0644)
	assert.NoErrorf(t, err, "failed to write update file: %v", err)

	cmd := exec.Command(testBinary, "update", "--secret", secretName, "--file", updateFile)
	_, err = cmd.CombinedOutput()
	assert.Error(t, err, "expected error for invalid secret name, got none")
}

// TestUpdateCommand_ValidSecretNames tests updating with valid secret name patterns
func TestUpdateCommand_ValidSecretNames(t *testing.T) {
	env.Init()

	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	updateFile := filepath.Join(t.TempDir(), "update.txt")
	secretContent := "originalsecret"
	updateContent := "updatedsecret"

	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	assert.NoErrorf(t, err, "failed to write secret file: %v", err)

	err = os.WriteFile(updateFile, []byte(updateContent), 0644)
	assert.NoErrorf(t, err, "failed to write update file: %v", err)

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
			_, err := createCmd.CombinedOutput()
			assert.NoErrorf(t, err, "failed to create initial secret '%s': %v", validName, err)

			// Then update it
			updateCmd := exec.Command(testBinary, "update", "--secret", validName, "--file", updateFile)
			_, err = updateCmd.CombinedOutput()
			assert.NoErrorf(t, err, "expected success for valid secret name '%s', got error: %v", validName, err)

			// Verify the file still exists
			assert.FileExistsf(t, secretOut, "expected output file %s to exist for name '%s'", secretOut, validName)
		})
	}
}
