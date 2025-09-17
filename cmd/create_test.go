package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/engmtcdrm/mellon/env"
)

var testBinary string

// TestMain builds the CLI binary once for all tests and cleans up after.
func TestMain(m *testing.M) {
	testBinary = filepath.Join(os.TempDir(), "mellon-test-bin")
	projectRoot, err := filepath.Abs(filepath.Join(".."))
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

// TestCreateCommand_ValidFlags tests the create command with valid flags.
func TestCreateCommand_ValidFlags(t *testing.T) {
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
		if err := os.WriteFile(secretFile, []byte(secretContent), 0644); err != nil {
			t.Fatalf("failed to write temp file: %v", err)
		}

		cmd := exec.Command(testBinary, append([]string{"create"}, args...)...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("expected success, got error: %v, output: %s", err, output)
		}

		if _, err := os.Stat(secretFile); os.IsNotExist(err) {
			t.Errorf("file should not be deleted")
		}

		if _, err := os.Stat(secretOut); os.IsNotExist(err) {
			t.Errorf("expected output file %s to exist, got error: %v", secretOut, err)
		} else {
			os.Remove(secretOut)
		}
	}
}

func TestCreateCommand_CleanupFlag(t *testing.T) {
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
		if err := os.WriteFile(secretFile, []byte(secretContent), 0644); err != nil {
			t.Fatalf("failed to write temp file: %v", err)
		}

		cmd := exec.Command(testBinary, append([]string{"create"}, args...)...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("expected success, got error: %v, output: %s", err, output)
		}

		if _, err := os.Stat(secretFile); !os.IsNotExist(err) {
			t.Errorf("file should be deleted with --cleanup")
		}

		if _, err := os.Stat(secretOut); os.IsNotExist(err) {
			t.Errorf("expected output file %s to exist, got error: %v", secretOut, err)
		} else {
			// Clean up after test
			os.Remove(secretOut)
		}
	}
}

func TestCreateCommand_Permission0600(t *testing.T) {
	env.Init()
	secretFile := filepath.Join(os.TempDir(), "secret.txt")
	secretName := "testperm"
	secretContent := "supersecret"
	secretOut := filepath.Join(env.Instance.SecretsPath(), secretName+env.Instance.SecretExt())
	// Clean up before test
	os.Remove(secretOut)

	if err := os.WriteFile(secretFile, []byte(secretContent), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	cmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	_, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	info, err := os.Stat(secretOut)
	if err != nil {
		t.Fatalf("expected output file, got error: %v", err)
	}

	if info.Mode().Perm() != 0600 {
		t.Fatalf("expected file mode 0600, got %v", info.Mode().Perm())
	}
}

func TestCreateCommand_TildeExpansion(t *testing.T) {
	env.Init()
	home, _ := os.UserHomeDir()
	secretFile := filepath.Join(home, "secrettilde.txt")
	secretName := "testtilde"
	secretContent := "supersecret"
	secretOut := filepath.Join(env.Instance.SecretsPath(), secretName+env.Instance.SecretExt())
	// Clean up before test
	os.Remove(secretOut)

	if err := os.WriteFile(secretFile, []byte(secretContent), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	defer os.Remove(secretFile)

	cmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", "~/secrettilde.txt")
	_, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestCreateCommand_FileNotExist(t *testing.T) {
	secretFile := filepath.Join(t.TempDir(), "doesnotexist.txt")
	secretName := "testnofile"

	cmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Errorf("expected error for non-existent file, got none")
	}
}

func TestCreateCommand_FileNoReadAccess(t *testing.T) {
	secretFile := filepath.Join(t.TempDir(), "noread.txt")
	secretName := "testnoread"
	secretContent := "supersecret"

	if err := os.WriteFile(secretFile, []byte(secretContent), 0000); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	defer os.Chmod(secretFile, 0644)

	cmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected error for no read access, got none")
	}
}

func TestCreateCommand_CleanupNoWriteAccess(t *testing.T) {
	dir := t.TempDir()
	secretFile := filepath.Join(dir, "nowrite.txt")
	secretName := "testnowrite"
	secretContent := "supersecret"

	if err := os.WriteFile(secretFile, []byte(secretContent), 0444); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	// Remove write permission from the directory to prevent file deletion
	if err := os.Chmod(dir, 0555); err != nil {
		t.Fatalf("failed to remove write permission from dir: %v", err)
	}
	defer os.Chmod(dir, 0755)

	cmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile, "--cleanup")
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Errorf("expected error for no write access to directory, got none")
	}
}

func TestCreateCommand_CleanupNoReadWriteAccess(t *testing.T) {
	secretFile := filepath.Join(t.TempDir(), "noreadwrite.txt")
	secretName := "testnoreadwrite"
	secretContent := "supersecret"

	if err := os.WriteFile(secretFile, []byte(secretContent), 0000); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	defer os.Chmod(secretFile, 0644)

	cmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile, "--cleanup")
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Errorf("expected error for no read/write access, got none")
	}
}

func TestCreateCommand_AlreadyExists(t *testing.T) {
	env.Init()
	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	secretName := "testexists"
	secretContent := "supersecret"
	secretOut := filepath.Join(env.Instance.SecretsPath(), secretName+env.Instance.SecretExt())
	// Clean up before test
	os.Remove(secretOut)

	if err := os.WriteFile(secretFile, []byte(secretContent), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	// First create
	cmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	_, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	// Try again with same name
	cmd = exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	_, err = cmd.CombinedOutput()
	if err == nil {
		t.Errorf("expected error for existing secret, got none")
	}
	// Clean up after test
	os.Remove(secretOut)
}

func TestCreateCommand_InvalidSecretName(t *testing.T) {
	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	secretName := "invalid!name"
	secretContent := "supersecret"

	if err := os.WriteFile(secretFile, []byte(secretContent), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	cmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Errorf("expected error for invalid secret, got none")
	}
}

// TestCreateCommand_PreRunValidation tests the PreRunE validation logic
func TestCreateCommand_PreRunValidation(t *testing.T) {
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
			_, err := cmd.CombinedOutput()
			if err == nil {
				t.Errorf("expected error for %s, got none", tc.name)
			}
		})
	}
}

// TestCreateCommand_SecretNameValidation tests various invalid secret name patterns
func TestCreateCommand_SecretNameValidation(t *testing.T) {
	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	secretContent := "supersecret"

	if err := os.WriteFile(secretFile, []byte(secretContent), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

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
			_, err := cmd.CombinedOutput()
			if err == nil {
				t.Errorf("expected error for invalid secret name '%s', got none", invalidName)
			}
		})
	}
}

// TestCreateCommand_ValidSecretNames tests valid secret name patterns
func TestCreateCommand_ValidSecretNames(t *testing.T) {
	env.Init()

	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	secretContent := "supersecret"

	if err := os.WriteFile(secretFile, []byte(secretContent), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
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
			secretOut := filepath.Join(env.Instance.SecretsPath(), validName+env.Instance.SecretExt())
			defer os.Remove(secretOut) // Clean up

			cmd := exec.Command(testBinary, "create", "--secret", validName, "--file", secretFile)
			_, err := cmd.CombinedOutput()
			if err != nil {
				t.Errorf("expected success for valid secret name '%s', got error: %v", validName, err)
			}

			// Verify the file was created
			if _, err := os.Stat(secretOut); os.IsNotExist(err) {
				t.Errorf("expected output file %s to exist for name '%s'", secretOut, validName)
			}
		})
	}
}

// TestCreateCommand_Force tests if there's a force flag (like in other commands)
func TestCreateCommand_ForceOverwrite(t *testing.T) {
	env.Init()

	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	secretName := "testforce"
	secretContent := "supersecret"
	secretOut := filepath.Join(env.Instance.SecretsPath(), secretName+env.Instance.SecretExt())
	defer os.Remove(secretOut) // Clean up

	if err := os.WriteFile(secretFile, []byte(secretContent), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	// First create
	cmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	_, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected success for first create, got error: %v", err)
	}

	// Try to create again with --force flag (if it exists)
	cmd = exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile, "--force")
	output, err := cmd.CombinedOutput()

	// This test will help identify if a force flag should be added
	if err != nil {
		t.Logf("Force flag not implemented (expected): %v, output: %s", err, output)
	} else {
		t.Logf("Force flag appears to be implemented")
	}
}
