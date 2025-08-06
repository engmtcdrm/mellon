package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/engmtcdrm/mellon/env"
)

// TestViewCommand_ValidFlags tests the view command with valid flags.
func TestViewCommand_ValidFlags(t *testing.T) {
	envVars, err := env.GetEnv()
	if err != nil {
		t.Fatalf("failed to get environment variables: %v", err)
	}

	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	secretName := "testviewsecret"
	secretContent := "supersecretcontent"
	secretOut := filepath.Join(envVars.SecretsPath, secretName+envVars.SecretExt)

	// Clean up before test
	os.Remove(secretOut)
	defer os.Remove(secretOut)

	// Create secret file
	if err := os.WriteFile(secretFile, []byte(secretContent), 0644); err != nil {
		t.Fatalf("failed to write secret file: %v", err)
	}

	// First create a secret to view
	createCmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	_, err = createCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to create initial secret: %v", err)
	}

	// Test each permutation of view flags
	cases := [][]string{
		{"--secret", secretName},
		{"-s", secretName},
	}

	for _, args := range cases {
		cmd := exec.Command(testBinary, append([]string{"view"}, args...)...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("expected success for args %v, got error: %v, output: %s", args, err, output)
			continue
		}

		// Verify the output contains the secret content
		outputStr := string(output)
		if !strings.Contains(outputStr, secretContent) {
			t.Errorf("expected output to contain secret content '%s', got: %s", secretContent, outputStr)
		}
	}
}

func TestViewCommand_OutputFlag(t *testing.T) {
	envVars, err := env.GetEnv()
	if err != nil {
		t.Fatalf("failed to get environment variables: %v", err)
	}

	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	outputFile := filepath.Join(t.TempDir(), "output.txt")
	secretName := "testviewoutput"
	secretContent := "outputsecretcontent"
	secretOut := filepath.Join(envVars.SecretsPath, secretName+envVars.SecretExt)

	// Clean up before test
	os.Remove(secretOut)
	defer os.Remove(secretOut)

	// Create secret file
	if err := os.WriteFile(secretFile, []byte(secretContent), 0644); err != nil {
		t.Fatalf("failed to write secret file: %v", err)
	}

	// First create a secret to view
	createCmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	_, err = createCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to create initial secret: %v", err)
	}

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
		if err != nil {
			t.Errorf("expected success for args %v, got error: %v", args, err)
			continue
		}

		// Verify the output file was created and contains the secret
		outputContent, err := os.ReadFile(outputFile)
		if err != nil {
			t.Errorf("failed to read output file: %v", err)
			continue
		}

		if string(outputContent) != secretContent {
			t.Errorf("expected output file content '%s', got '%s'", secretContent, string(outputContent))
		}

		// Verify file permissions
		info, err := os.Stat(outputFile)
		if err != nil {
			t.Errorf("failed to stat output file: %v", err)
			continue
		}

		if info.Mode().Perm() != 0600 {
			t.Errorf("expected output file mode 0600, got %o", info.Mode().Perm())
		}
	}
}

func TestViewCommand_SecretNotExist(t *testing.T) {
	secretName := "nonexistentviewsecret"

	cmd := exec.Command(testBinary, "view", "--secret", secretName)
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Errorf("expected error for non-existent secret, got none")
	}
}

func TestViewCommand_InvalidSecretName(t *testing.T) {
	secretName := "invalid!name"

	cmd := exec.Command(testBinary, "view", "--secret", secretName)
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Errorf("expected error for invalid secret name, got none")
	}
}

func TestViewCommand_ValidSecretNames(t *testing.T) {
	envVars, err := env.GetEnv()
	if err != nil {
		t.Fatalf("failed to get environment variables: %v", err)
	}

	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	secretContent := "viewsecretcontent"

	if err := os.WriteFile(secretFile, []byte(secretContent), 0644); err != nil {
		t.Fatalf("failed to write secret file: %v", err)
	}

	validNames := []string{
		"simple",
		"with_underscores",
		"with-dashes",
		"with123numbers",
		"path/to/secret",
	}

	for _, validName := range validNames {
		t.Run(fmt.Sprintf("valid_name_%s", validName), func(t *testing.T) {
			secretOut := filepath.Join(envVars.SecretsPath, validName+envVars.SecretExt)
			defer os.Remove(secretOut)

			// Create the secret
			createCmd := exec.Command(testBinary, "create", "--secret", validName, "--file", secretFile)
			_, err := createCmd.CombinedOutput()
			if err != nil {
				t.Fatalf("failed to create secret '%s': %v", validName, err)
			}

			// View the secret
			viewCmd := exec.Command(testBinary, "view", "--secret", validName)
			output, err := viewCmd.CombinedOutput()
			if err != nil {
				t.Errorf("expected success for valid secret name '%s', got error: %v", validName, err)
				return
			}

			// Verify content
			outputStr := string(output)
			if !strings.Contains(outputStr, secretContent) {
				t.Errorf("expected output to contain secret content for name '%s'", validName)
			}
		})
	}
}

func TestViewCommand_PreRunValidation(t *testing.T) {
	// Test output flag without secret flag
	outputFile := filepath.Join(t.TempDir(), "output.txt")

	cmd := exec.Command(testBinary, "view", "--output", outputFile)
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Errorf("expected error for output flag without secret flag, got none")
	}
}

func TestViewCommand_OutputDirectoryCreation(t *testing.T) {
	envVars, err := env.GetEnv()
	if err != nil {
		t.Fatalf("failed to get environment variables: %v", err)
	}

	secretFile := filepath.Join(t.TempDir(), "secret.txt")
	outputDir := filepath.Join(t.TempDir(), "nested", "output")
	outputFile := filepath.Join(outputDir, "secret.txt")
	secretName := "testviewdircreate"
	secretContent := "dircreatecontent"
	secretOut := filepath.Join(envVars.SecretsPath, secretName+envVars.SecretExt)

	// Clean up before test
	os.Remove(secretOut)
	defer os.Remove(secretOut)

	// Create secret file
	if err := os.WriteFile(secretFile, []byte(secretContent), 0644); err != nil {
		t.Fatalf("failed to write secret file: %v", err)
	}

	// First create a secret to view
	createCmd := exec.Command(testBinary, "create", "--secret", secretName, "--file", secretFile)
	_, err = createCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to create initial secret: %v", err)
	}

	// Test that output directory is created
	cmd := exec.Command(testBinary, "view", "--secret", secretName, "--output", outputFile)
	_, err = cmd.CombinedOutput()
	if err != nil {
		t.Errorf("expected success for nested output directory, got error: %v", err)
	}

	// Verify the output file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("expected output file to be created")
	}

	// Verify directory was created with correct permissions
	info, err := os.Stat(outputDir)
	if err != nil {
		t.Errorf("failed to stat output directory: %v", err)
	} else if info.Mode().Perm() != 0700 {
		t.Errorf("expected output directory mode 0700, got %o", info.Mode().Perm())
	}
}

func TestViewCommand_NoFlags(t *testing.T) {
	// Test view command without any flags (should enter interactive mode, but we skip this)
	t.Skip("Skipping interactive test: view without flags")
}
