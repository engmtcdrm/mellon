package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var testBinary string

// TestMain builds the CLI binary once for all tests and cleans up after.
func TestMain(m *testing.M) {
	testBinary = filepath.Join(os.TempDir(), "minno-test-bin")
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
	dir := t.TempDir()
	credFile := filepath.Join(dir, "cred.txt")
	credName := "testcred"
	credContent := "supersecret"
	home, _ := os.UserHomeDir()
	credOut := filepath.Join(home, ".minno", credName+".cred")
	// Clean up before test
	os.Remove(credOut)

	// Test each permutation of flags
	cases := [][]string{
		{"--cred-name", credName, "--file", credFile},
		{"--cred-name", credName, "-f", credFile},
		{"-n", credName, "--file", credFile},
		{"-n", credName, "-f", credFile},
	}

	for _, args := range cases {
		if err := os.WriteFile(credFile, []byte(credContent), 0644); err != nil {
			t.Fatalf("failed to write temp file: %v", err)
		}

		cmd := exec.Command(testBinary, append([]string{"create"}, args...)...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("expected success, got error: %v, output: %s", err, output)
		}

		if _, err := os.Stat(credFile); os.IsNotExist(err) {
			t.Errorf("file should not be deleted")
		}

		if _, err := os.Stat(credOut); os.IsNotExist(err) {
			t.Errorf("expected output file %s to exist, got error: %v", credOut, err)
		} else {
			os.Remove(credOut)
		}
	}
}

func TestCreateCommand_CleanupFlag(t *testing.T) {
	credFile := filepath.Join(t.TempDir(), "cred.txt")
	credName := "testcleanup"
	credContent := "supersecret"
	home, _ := os.UserHomeDir()
	credOut := filepath.Join(home, ".minno", credName+".cred")
	// Clean up before test
	os.Remove(credOut)

	// Test each permutation of flags with --cleanup flag added
	cases := [][]string{
		{"--cred-name", credName, "--file", credFile, "--cleanup"},
		{"--cred-name", credName, "--file", credFile, "-c"},
		{"--cred-name", credName, "-f", credFile, "--cleanup"},
		{"--cred-name", credName, "-f", credFile, "-c"},
		{"-n", credName, "--file", credFile, "--cleanup"},
		{"-n", credName, "--file", credFile, "-c"},
		{"-n", credName, "-f", credFile, "--cleanup"},
		{"-n", credName, "-f", credFile, "-c"},
	}

	for _, args := range cases {
		if err := os.WriteFile(credFile, []byte(credContent), 0644); err != nil {
			t.Fatalf("failed to write temp file: %v", err)
		}

		cmd := exec.Command(testBinary, append([]string{"create"}, args...)...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("expected success, got error: %v, output: %s", err, output)
		}

		if _, err := os.Stat(credFile); !os.IsNotExist(err) {
			t.Errorf("file should be deleted with --cleanup")
		}

		if _, err := os.Stat(credOut); os.IsNotExist(err) {
			t.Errorf("expected output file %s to exist, got error: %v", credOut, err)
		} else {
			// Clean up after test
			os.Remove(credOut)
		}
	}
}

func TestCreateCommand_Permission0600(t *testing.T) {
	credFile := filepath.Join(os.TempDir(), "cred.txt")
	credName := "testperm"
	credContent := "supersecret"
	home, _ := os.UserHomeDir()
	credOut := filepath.Join(home, ".minno", credName+".cred")
	// Clean up before test
	os.Remove(credOut)

	if err := os.WriteFile(credFile, []byte(credContent), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	cmd := exec.Command(testBinary, "create", "--cred-name", credName, "--file", credFile)
	_, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	info, err := os.Stat(credOut)
	if err != nil {
		t.Fatalf("expected output file, got error: %v", err)
	}

	if info.Mode().Perm() != 0600 {
		t.Fatalf("expected file mode 0600, got %v", info.Mode().Perm())
	}
}

func TestCreateCommand_TildeExpansion(t *testing.T) {
	home, _ := os.UserHomeDir()
	credFile := filepath.Join(home, "credtilde.txt")
	credName := "testtilde"
	credContent := "supersecret"
	credOut := filepath.Join(home, ".minno", credName+".cred")
	// Clean up before test
	os.Remove(credOut)

	if err := os.WriteFile(credFile, []byte(credContent), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	defer os.Remove(credFile)

	cmd := exec.Command(testBinary, "create", "--cred-name", credName, "--file", "~/credtilde.txt")
	_, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestCreateCommand_MissingFlags(t *testing.T) {
	credFile := filepath.Join(t.TempDir(), "cred.txt")
	credName := "testmissing"

	cases := [][]string{
		{"--cred-name", credName},
		{"--file", credFile},
		{"-n", credName},
		{"-f", credFile},
	}

	for _, args := range cases {
		cmd := exec.Command(testBinary, append([]string{"create"}, args...)...)
		_, err := cmd.CombinedOutput()
		if err == nil {
			t.Errorf("expected error for args %v, got none", args)
		}
	}
}

func TestCreateCommand_FileNotExist(t *testing.T) {
	credFile := filepath.Join(t.TempDir(), "doesnotexist.txt")
	credName := "testnofile"

	cmd := exec.Command(testBinary, "create", "--cred-name", credName, "--file", credFile)
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Errorf("expected error for non-existent file, got none")
	}
}

func TestCreateCommand_FileNoReadAccess(t *testing.T) {
	credFile := filepath.Join(t.TempDir(), "noread.txt")
	credName := "testnoread"
	credContent := "supersecret"

	if err := os.WriteFile(credFile, []byte(credContent), 0000); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	defer os.Chmod(credFile, 0644)

	cmd := exec.Command(testBinary, "create", "--cred-name", credName, "--file", credFile)
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected error for no read access, got none")
	}
}

func TestCreateCommand_CleanupNoWriteAccess(t *testing.T) {
	dir := t.TempDir()
	credFile := filepath.Join(dir, "nowrite.txt")
	credName := "testnowrite"
	credContent := "supersecret"

	if err := os.WriteFile(credFile, []byte(credContent), 0444); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	// Remove write permission from the directory to prevent file deletion
	if err := os.Chmod(dir, 0555); err != nil {
		t.Fatalf("failed to remove write permission from dir: %v", err)
	}
	defer os.Chmod(dir, 0755)

	cmd := exec.Command(testBinary, "create", "--cred-name", credName, "--file", credFile, "--cleanup")
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Errorf("expected error for no write access to directory, got none")
	}
}

func TestCreateCommand_CleanupNoReadWriteAccess(t *testing.T) {
	credFile := filepath.Join(t.TempDir(), "noreadwrite.txt")
	credName := "testnoreadwrite"
	credContent := "supersecret"

	if err := os.WriteFile(credFile, []byte(credContent), 0000); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	defer os.Chmod(credFile, 0644)

	cmd := exec.Command(testBinary, "create", "--cred-name", credName, "--file", credFile, "--cleanup")
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Errorf("expected error for no read/write access, got none")
	}
}

func TestCreateCommand_AlreadyExists(t *testing.T) {
	credFile := filepath.Join(t.TempDir(), "cred.txt")
	credName := "testexists"
	credContent := "supersecret"
	home, _ := os.UserHomeDir()
	credOut := filepath.Join(home, ".minno", credName+".cred")
	// Clean up before test
	os.Remove(credOut)

	if err := os.WriteFile(credFile, []byte(credContent), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	// First create
	cmd := exec.Command(testBinary, "create", "--cred-name", credName, "--file", credFile)
	_, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	// Try again with same name
	cmd = exec.Command(testBinary, "create", "--cred-name", credName, "--file", credFile)
	_, err = cmd.CombinedOutput()
	if err == nil {
		t.Errorf("expected error for existing cred-name, got none")
	}
	// Clean up after test
	os.Remove(credOut)
}

func TestCreateCommand_InvalidCredName(t *testing.T) {
	credFile := filepath.Join(t.TempDir(), "cred.txt")
	credName := "invalid!name"
	credContent := "supersecret"

	if err := os.WriteFile(credFile, []byte(credContent), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	cmd := exec.Command(testBinary, "create", "--cred-name", credName, "--file", credFile)
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Errorf("expected error for invalid cred-name, got none")
	}
}
