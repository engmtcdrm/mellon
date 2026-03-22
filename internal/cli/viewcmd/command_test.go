package viewcmd

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
	testBinary = filepath.Join(os.TempDir(), "mellon-test-bin-viewcmd")
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

// TestViewCommand_ValidFlags tests the view command with valid flags.
func TestViewCommand_ValidFlags(t *testing.T) {
	env.Init()

	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	secretName := "testviewsecret"
	secretContent := "supersecretcontent"
	secretOut := filepath.Join(env.Instance.SecretsPath(), secretName+env.Instance.SecretExt())

	// Clean up before test
	os.Remove(secretOut)
	defer os.Remove(secretOut)

	// Create secret file
	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	assert.NoErrorf(t, err, "failed to write secret file: %v", err)

	// First create a secret to view
	createCmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	_, err = createCmd.CombinedOutput()
	assert.NoErrorf(t, err, "failed to create initial secret: %v", err)

	// Test each permutation of view flags
	cases := [][]string{
		{"--secret", secretName},
		{"-s", secretName},
	}

	for _, args := range cases {
		cmd := exec.Command(testBinary, append([]string{"view"}, args...)...)
		output, err := cmd.CombinedOutput()
		assert.NoErrorf(t, err, "expected success for args %v, got error: %v, output: %s", args, err, output)

		// Verify the output contains the secret content
		outputStr := string(output)
		assert.Containsf(t, outputStr, secretContent, "expected output to contain secret content '%s', got: %s", secretContent, outputStr)
	}
}

func TestViewCommand_OutputFlag(t *testing.T) {
	env.Init()

	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	outputFile := filepath.Join(t.TempDir(), "output.txt")
	secretName := "testviewoutput"
	secretContent := "outputsecretcontent"
	secretOut := filepath.Join(env.Instance.SecretsPath(), secretName+env.Instance.SecretExt())

	// Clean up before test
	os.Remove(secretOut)
	defer os.Remove(secretOut)

	// Create secret file
	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	assert.NoErrorf(t, err, "failed to write secret file: %v", err)

	// First create a secret to view
	createCmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	_, err = createCmd.CombinedOutput()
	assert.NoErrorf(t, err, "failed to create initial secret: %v", err)

	// Test output flag variations
	cases := [][]string{
		{"--secret", secretName, "--output", outputFile},
		{"-s", secretName, "-o", outputFile},
	}

	for _, args := range cases {
		// Remove output file if it exists
		os.Remove(outputFile)

		cmd := exec.Command(testBinary, append([]string{"view"}, args...)...)
		_, err := cmd.CombinedOutput()
		assert.NoErrorf(t, err, "expected success for args %v, got error: %v", args, err)

		// Verify the output file was created and contains the secret
		outputContent, err := os.ReadFile(outputFile)
		assert.NoErrorf(t, err, "failed to read output file: %v", err)

		assert.Equalf(t, secretContent, string(outputContent), "expected output file content '%s', got '%s'", secretContent, string(outputContent))

		// Verify file permissions
		info, err := os.Stat(outputFile)
		assert.NoErrorf(t, err, "failed to stat output file: %v", err)

		assert.Equalf(t, os.FileMode(0600), info.Mode().Perm(), "expected output file mode 0600, got %o", info.Mode().Perm())
	}
}

func TestViewCommand_SecretNotExist(t *testing.T) {
	secretName := "nonexistentviewsecret"

	cmd := exec.Command(testBinary, "view", "--secret", secretName)
	_, err := cmd.CombinedOutput()
	assert.Error(t, err, "expected error for non-existent secret, got none")
}

func TestViewCommand_InvalidSecretName(t *testing.T) {
	secretName := "invalid!name"

	cmd := exec.Command(testBinary, "view", "--secret", secretName)
	_, err := cmd.CombinedOutput()
	assert.Error(t, err, "expected error for invalid secret name, got none")
}

func TestViewCommand_ValidSecretNames(t *testing.T) {
	env.Init()

	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	secretContent := "viewsecretcontent"

	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	assert.NoErrorf(t, err, "failed to write secret file: %v", err)

	validNames := []string{
		"simple",
		"with_underscores",
		"with-dashes",
		"with123numbers",
		"path/to/secret",
	}

	for _, validName := range validNames {
		t.Run(fmt.Sprintf("valid_name_%s", validName), func(t *testing.T) {
			secretOut := filepath.Join(env.Instance.SecretsPath(), validName+env.Instance.SecretExt())
			defer os.Remove(secretOut)

			// Create the secret
			createCmd := exec.Command(testBinary, "create", "--secret", validName, "--file", secretFile)
			_, err := createCmd.CombinedOutput()
			assert.NoErrorf(t, err, "failed to create secret with valid name '%s': %v", validName, err)

			// View the secret
			viewCmd := exec.Command(testBinary, "view", "--secret", validName)
			output, err := viewCmd.CombinedOutput()
			assert.NoErrorf(t, err, "failed to view secret with valid name '%s': %v, output: %s", validName, err, output)

			// Verify content
			outputStr := string(output)
			assert.Containsf(t, outputStr, secretContent, "expected output to contain secret content for name '%s'", validName)
		})
	}
}

func TestViewCommand_PreRunValidation(t *testing.T) {
	// Test output flag without secret flag
	outputFile := filepath.Join(t.TempDir(), "output.txt")

	cmd := exec.Command(testBinary, "view", "--output", outputFile)
	_, err := cmd.CombinedOutput()
	assert.Error(t, err, "expected error for output flag without secret flag, got none")
}

func TestViewCommand_OutputDirectoryCreation(t *testing.T) {
	env.Init()

	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	outputDir := filepath.Join(t.TempDir(), "nested", "output")
	outputFile := filepath.Join(outputDir, "secret.txt")
	secretName := "testviewdircreate"
	secretContent := "dircreatecontent"
	secretOut := filepath.Join(env.Instance.SecretsPath(), secretName+env.Instance.SecretExt())

	// Clean up before test
	os.Remove(secretOut)
	defer os.Remove(secretOut)

	// Create secret file
	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	assert.NoErrorf(t, err, "failed to write secret file: %v", err)

	// First create a secret to view
	createCmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	_, err = createCmd.CombinedOutput()
	assert.NoErrorf(t, err, "failed to create initial secret: %v", err)

	// Test that output directory is created
	cmd := exec.Command(testBinary, "view", "--secret", secretName, "--output", outputFile)
	_, err = cmd.CombinedOutput()
	assert.NoErrorf(t, err, "expected success for nested output directory, got error: %v", err)

	// Verify the output file was created
	assert.FileExistsf(t, outputFile, "expected output file to be created at %s", outputFile)

	// Verify directory was created with correct permissions
	info, err := os.Stat(outputDir)
	assert.NoErrorf(t, err, "failed to stat output directory: %v", err)
	assert.Equalf(t, os.FileMode(0700), info.Mode().Perm(), "expected output directory mode 0700, got %o", info.Mode().Perm())
}

func TestViewCommand_NoFlags(t *testing.T) {
	// Test view command without any flags (should enter interactive mode, but we skip this)
	t.Skip("Skipping interactive test: view without flags")
}
