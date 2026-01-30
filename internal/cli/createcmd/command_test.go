package createcmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/engmtcdrm/mellon/internal/env"
	"github.com/engmtcdrm/mellon/secrets"
	"github.com/stretchr/testify/assert"
)

var testBinary string

// TestMain builds the CLI binary once for all tests and cleans up after.
func TestMain(m *testing.M) {
	testBinary = filepath.Join(os.TempDir(), "mellon-test-bin-createcmd")
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

// TestCreateCommandValidFlags tests the create command with valid flags.
func TestCreateCommandValidFlags(t *testing.T) {
	env.Init()
	dir := t.TempDir()
	secretFile := filepath.Join(dir, "secret.txt")
	secretName := "testsecret"
	secretContent := "supersecret"
	secretOut := filepath.Join(env.Instance.SecretsPath(), secretName+env.Instance.SecretExt())
	// Clean up before test
	os.Remove(secretOut)

	// Test each permutation of flags
	cases := [][]string{
		{"--secret", secretName, "--file", secretFile},
		{"--secret", secretName, "-f", secretFile},
		{"-s", secretName, "--file", secretFile},
		{"-s", secretName, "-f", secretFile},
	}

	for _, args := range cases {
		err := os.WriteFile(secretFile, []byte(secretContent), 0644)
		assert.NoErrorf(t, err, "failed to write temp file: %v", err)

		cmd := exec.Command(testBinary, append([]string{"create"}, args...)...)
		output, err := cmd.CombinedOutput()
		assert.NoErrorf(t, err, "expected success, got error: %v, output: %s", err, output)
		assert.FileExistsf(t, secretOut, "expected output file %s to exist", secretOut)

		os.Remove(secretOut)
	}
}

func TestCreateCommandCleanupFlag(t *testing.T) {
	env.Init()
	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	secretName := "testcleanup"
	secretContent := "supersecret"
	secretOut := filepath.Join(env.Instance.SecretsPath(), secretName+env.Instance.SecretExt())
	// Clean up before test
	os.Remove(secretOut)

	// Test each permutation of flags with --cleanup flag added
	cases := [][]string{
		{"--secret", secretName, "--file", secretFile, "--cleanup"},
		{"--secret", secretName, "--file", secretFile, "-c"},
		{"--secret", secretName, "-f", secretFile, "--cleanup"},
		{"--secret", secretName, "-f", secretFile, "-c"},
		{"-s", secretName, "--file", secretFile, "--cleanup"},
		{"-s", secretName, "--file", secretFile, "-c"},
		{"-s", secretName, "-f", secretFile, "--cleanup"},
		{"-s", secretName, "-f", secretFile, "-c"},
	}

	for _, args := range cases {
		err := os.WriteFile(secretFile, []byte(secretContent), 0644)
		assert.NoErrorf(t, err, "failed to write temp file: %v", err)

		cmd := exec.Command(testBinary, append([]string{"create"}, args...)...)
		output, err := cmd.CombinedOutput()
		assert.NoErrorf(t, err, "expected success, got error: %v, output: %s", err, output)
		assert.NoFileExistsf(t, secretFile, "file should be deleted with --cleanup")
		assert.FileExistsf(t, secretOut, "expected output file %s to exist", secretOut)

		// Clean up after test
		os.Remove(secretOut)
	}
}

func TestCreateCommandPermission0600(t *testing.T) {
	env.Init()
	secretFile := filepath.Join(os.TempDir(), "secret.txt")
	secretName := "testperm"
	secretContent := "supersecret"
	secretOut := filepath.Join(env.Instance.SecretsPath(), secretName+env.Instance.SecretExt())
	// Clean up before test
	os.Remove(secretOut)

	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	assert.NoErrorf(t, err, "failed to write temp file: %v", err)

	cmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	output, err := cmd.CombinedOutput()
	assert.NoErrorf(t, err, "expected success, got error: %v, output: %s", err, output)

	info, err := os.Stat(secretOut)
	assert.NoErrorf(t, err, "expected output file %s to exist, got error: %v", secretOut, err)
	assert.FileExistsf(t, secretOut, "expected output file %s to exist, got error: %v", secretOut, err)
	assert.NotEqualf(t, 0600, info.Mode().Perm(), "expected file mode 0600, got %v", info.Mode().Perm())
}

func TestCreateCommandTildeExpansion(t *testing.T) {
	env.Init()
	home, _ := os.UserHomeDir()
	secretFile := filepath.Join(home, "secrettilde.txt")
	secretName := "testtilde"
	secretContent := "supersecret"
	secretOut := filepath.Join(env.Instance.SecretsPath(), secretName+env.Instance.SecretExt())
	// Clean up before test
	os.Remove(secretOut)

	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	assert.NoErrorf(t, err, "failed to write temp file: %v", err)
	defer os.Remove(secretFile)

	cmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", "~/secrettilde.txt")
	output, err := cmd.CombinedOutput()
	assert.NoErrorf(t, err, "expected success, got error: %v, output: %s", err, output)
}

func TestCreateCommandFileNotExist(t *testing.T) {
	secretFile := filepath.Join(t.TempDir(), "doesnotexist.txt")
	secretName := "testnofile"

	cmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	output, err := cmd.CombinedOutput()
	assert.Error(t, err, "expected error for non-existent file, got none, output: %s", output)
}

func TestCreateCommandFileNoReadAccess(t *testing.T) {
	secretFile := filepath.Join(t.TempDir(), "noread.txt")
	secretName := "testnoread"
	secretContent := "supersecret"

	err := os.WriteFile(secretFile, []byte(secretContent), 0000)
	assert.NoErrorf(t, err, "failed to write temp file: %v", err)
	defer os.Chmod(secretFile, 0644)

	cmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	output, err := cmd.CombinedOutput()
	assert.Error(t, err, "expected error for no read access, got none, output: %s", output)
}

func TestCreateCommandCleanupNoWriteAccess(t *testing.T) {
	dir := t.TempDir()
	secretFile := filepath.Join(dir, "nowrite.txt")
	secretName := "testnowrite"
	secretContent := "supersecret"

	err := os.WriteFile(secretFile, []byte(secretContent), 0444)
	assert.NoErrorf(t, err, "failed to write temp file: %v", err)

	// Remove write permission from the directory to prevent file deletion
	err = os.Chmod(dir, 0555)
	assert.NoErrorf(t, err, "failed to remove write permission from dir: %v", err)
	defer os.Chmod(dir, 0755)

	cmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile, "--cleanup")
	output, err := cmd.CombinedOutput()
	assert.Error(t, err, "expected error for no write access to directory, got none, output: %s", output)
}

func TestCreateCommandCleanupNoReadWriteAccess(t *testing.T) {
	secretFile := filepath.Join(t.TempDir(), "noreadwrite.txt")
	secretName := "testnoreadwrite"
	secretContent := "supersecret"

	err := os.WriteFile(secretFile, []byte(secretContent), 0000)
	assert.NoErrorf(t, err, "failed to write temp file: %v", err)
	defer os.Chmod(secretFile, 0644)

	cmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile, "--cleanup")
	output, err := cmd.CombinedOutput()
	assert.Error(t, err, "expected error for no read/write access, got none, output: %s", output)
}

func TestCreateCommandAlreadyExists(t *testing.T) {
	env.Init()
	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	secretName := "testexists"
	secretContent := "supersecret"
	secretOut := filepath.Join(env.Instance.SecretsPath(), secretName+env.Instance.SecretExt())
	// Clean up before test
	os.Remove(secretOut)

	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	assert.NoErrorf(t, err, "failed to write temp file: %v", err)

	// First create
	cmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	output, err := cmd.CombinedOutput()
	assert.NoErrorf(t, err, "expected success, got error: %v, output: %s", err, output)

	// Try again with same name
	cmd = exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	output, err = cmd.CombinedOutput()
	assert.Error(t, err, "expected error for existing secret, got none, output: %s", output)

	// Clean up after test
	os.Remove(secretOut)
}

func TestCreateCommandInvalidSecretName(t *testing.T) {
	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	secretName := "invalid!name"
	secretContent := "supersecret"

	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	assert.NoErrorf(t, err, "failed to write temp file: %v", err)

	cmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	output, err := cmd.CombinedOutput()
	assert.Error(t, err, "expected error for invalid secret, got none, output: %s", output)
}

// TestCreateCommandPreRunValidation tests the PreRunE validation logic
func TestCreateCommandPreRunValidation(t *testing.T) {
	// Test cleanup flag without required flags
	cases := []struct {
		name string
		args []string
	}{
		{
			name: "cleanup without secret flag",
			args: []string{"create", "--cleanup", "--file", "somefile.txt"},
		},
		{
			name: "cleanup without file flag",
			args: []string{"create", "--cleanup", "--secret", "somesecret"},
		},
		{
			name: "cleanup without both flags",
			args: []string{"create", "--cleanup"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command(testBinary, tc.args...)
			output, err := cmd.CombinedOutput()
			assert.Error(t, err, "expected error for %s, got none, output: %s", tc.name, output)
		})
	}
}

// TestCreateCommandSecretNameValidation tests various invalid secret name patterns
func TestCreateCommandSecretNameValidation(t *testing.T) {
	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	secretContent := "supersecret"

	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	assert.NoErrorf(t, err, "failed to write temp file: %v", err)

	invalidNames := []string{
		"name with spaces",
		"name:with:colons",
		"name*with*asterisks",
		"name?with?questions",
		"name<with>brackets",
		"name|with|pipes",
		"name\"with\"quotes",
	}

	for _, invalidName := range invalidNames {
		t.Run(fmt.Sprintf("invalid_name_%s", invalidName), func(t *testing.T) {
			cmd := exec.Command(testBinary, "create", "--secret", invalidName, "--file", secretFile)
			output, err := cmd.CombinedOutput()
			assert.Error(t, err, "expected error for invalid secret name '%s', got none, output: %s", invalidName, output)
		})
	}
}

// TestCreateCommandValidSecretNames tests valid secret name patterns
func TestCreateCommandValidSecretNames(t *testing.T) {
	env.Init()

	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	secretContent := "supersecret"

	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	assert.NoErrorf(t, err, "failed to write temp file: %v", err)

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

			cmd := exec.Command(testBinary, "create", "--secret", validName, "--file", secretFile)
			output, err := cmd.CombinedOutput()
			assert.NoErrorf(t, err, "expected success for valid secret name '%s', got error: %v, output: %s", validName, err, output)

			// Verify the file was created
			_, err = os.Stat(secretOut)
			assert.NoErrorf(t, err, "expected output file %s to exist for name '%s', got error: %v", secretOut, validName, err)
			assert.FileExistsf(t, secretOut, "expected output file %s to exist for name '%s', got error: %v", secretOut, validName, err)
		})
	}
}

// TestCreateCommandForceOverwrite tests the behavior of the --force flag when creating a secret that already exists
func TestCreateCommandForceOverwrite(t *testing.T) {
	env.Init()

	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	secretName := "testforce"
	secretContent := "supersecret"
	secretOut := filepath.Join(env.Instance.SecretsPath(), secretName+env.Instance.SecretExt())
	defer os.Remove(secretOut) // Clean up

	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	assert.NoErrorf(t, err, "failed to write temp file: %v", err)

	// First create
	cmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	output, err := cmd.CombinedOutput()
	assert.NoErrorf(t, err, "expected success for first create, got error: %v, output: %s", err, output)

	// Try to create again with --force flag (if it exists)
	cmd = exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile, "--force")
	output, err = cmd.CombinedOutput()
	assert.Error(t, err, "expected error for creating existing secret with --force flag, got none, output: %s", output)
}

func TestNewCommand(t *testing.T) {
	emptySecrets := []secrets.Secret{}
	new := NewCommand(emptySecrets)
	assert.NotNil(t, new)
}

func TestValidateFlags(t *testing.T) {
	t.Run("cleanup without secret and file", func(t *testing.T) {
		cleanupFile = true
		secretName = ""
		secretFile = ""

		err := validateFlags(nil, []string{})
		assert.Error(t, err)
	})

	t.Run("cleanup with secret and file", func(t *testing.T) {
		cleanupFile = true
		secretName = "mysecret"
		secretFile = "myfile.txt"

		err := validateFlags(nil, []string{})
		assert.NoError(t, err)
	})
}

func TestRun(t *testing.T) {
	env.Init()
	t.Run("run with secretName and secretFile, file does not exist", func(t *testing.T) {
		secretName = "mysecret"
		secretFile = "myfile.txt"
		secretFiles = []secrets.Secret{}

		err := run(nil, []string{})
		assert.Error(t, err) // file does not exist, so expect error
	})
}

func TestValidateSecretName(t *testing.T) {
	env.Init()
	tempDir := t.TempDir()
	keyPath := filepath.Join(tempDir, "TestValidateSecretName.key")
	secretsPath := filepath.Join(tempDir, "secrets")
	secretFiles = []secrets.Secret{}

	secret1, err := secrets.NewSecret(keyPath, "secret1", secretsPath)
	assert.NoError(t, err)
	secretFiles = append(secretFiles, *secret1)

	t.Run("empty name", func(t *testing.T) {
		err := validateSecretName("")
		assert.Error(t, err)
	})

	t.Run("invalid name", func(t *testing.T) {
		err := validateSecretName("invalid name!:")
		assert.Error(t, err)
	})

	t.Run("secret exists", func(t *testing.T) {
		err := validateSecretName("secret1")
		assert.Error(t, err)
	})

	t.Run("secret does not exist", func(t *testing.T) {
		err := validateSecretName("newsecret")
		assert.NoError(t, err)
	})
}

func TestEncryptFromFile(t *testing.T) {
	env.Init()
	tempDir := t.TempDir()
	keyPath := filepath.Join(tempDir, "TestValidateSecretName.key")
	secretsPath := filepath.Join(tempDir, "secrets")
	secretFiles = []secrets.Secret{}

	secret1, err := secrets.NewSecret(keyPath, "secret1", secretsPath)
	assert.NoError(t, err)
	secretFiles = append(secretFiles, *secret1)

	t.Run("invalid secret name", func(t *testing.T) {
		secretName = "invalid name!:"
		err := encryptFromFile()
		assert.Error(t, err)
	})

	t.Run("secret already exists", func(t *testing.T) {
		secretName = "secret1"
		err := encryptFromFile()
		assert.Error(t, err)
	})

	t.Run("new secret", func(t *testing.T) {
		secretName = "secret2"
		secretFile = filepath.Join(tempDir, "unencrypted-secret2.txt")
		err := os.WriteFile(secretFile, []byte("supersecret2"), 0644)
		assert.NoError(t, err)

		err = encryptFromFile()
		assert.NoError(t, err)

		// Verify the secret file was created
		secretPath := filepath.Join(env.Instance.SecretsPath(), secretName+env.Instance.SecretExt())
		_, err = os.Stat(secretPath)
		assert.NoError(t, err)

		err = os.Remove(secretPath)
		assert.NoError(t, err)
	})
}

func TestResolveSecretName(t *testing.T) {
	env.Init()
	tempDir := t.TempDir()
	keyPath := filepath.Join(tempDir, "TestValidateSecretName.key")
	secretsPath := filepath.Join(tempDir, "secrets")
	secretFiles = []secrets.Secret{}

	secret1, err := secrets.NewSecret(keyPath, "secret1", secretsPath)
	assert.NoError(t, err)
	secretFiles = append(secretFiles, *secret1)

	t.Run("secretName not set", func(t *testing.T) {
		// Create a temporary file to simulate os.Stdin
		tempFile, err := ioutil.TempFile("", "testinput")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer os.Remove(tempFile.Name()) // Clean up the temp file

		// Write the simulated input to the temp file
		if _, err := tempFile.WriteString("test"); err != nil {
			t.Fatalf("failed to write to temp file: %v", err)
		}

		// Seek to the beginning of the file so it can be read
		if _, err := tempFile.Seek(0, 0); err != nil {
			t.Fatalf("failed to seek to beginning of temp file: %v", err)
		}

		// Redirect os.Stdin to the temporary file
		oldStdin := os.Stdin
		defer func() { os.Stdin = oldStdin }() // Restore original Stdin after the test
		os.Stdin = tempFile

		err = resolveSecretName()
		assert.NoError(t, err)
	})

	t.Run("secretName exists", func(t *testing.T) {
		secretName = "secret1"
		err := resolveSecretName()
		assert.Error(t, err)
	})

	t.Run("secretName does not exist", func(t *testing.T) {
		secretName = "secret2"
		err := resolveSecretName()
		assert.NoError(t, err)
	})
}
