package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"

	"github.com/engmtcdrm/minno/secrets"
	"github.com/spf13/cobra"
)

// validateUpdateCreateFlags checks if the flags for creating or updating a secret are valid.
func validateUpdateCreateFlags(cmd *cobra.Command, args []string) error {
	if cleanupFile && (secretName == "" || secretFile == "") {
		return errors.New("flag -c/--cleanup can only be used when -s/--secret and -f/--file are provided")
	}

	return nil
}

// validateSecretName checks if the provided secret name is valid.
func validateSecretName(name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}

	if err := secrets.ValidateName(name); err != nil {
		return err
	}

	if secretPtr := secrets.FindSecretByName(name, secretFiles); secretPtr != nil {
		return errors.New("secret with that name already exists")
	}

	return nil
}

// mkdir creates a directory at the specified path with the given mode.
func mkdir(path string, dirMode os.FileMode) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Directory does not exist, create it
		err = os.MkdirAll(path, dirMode)
		if err != nil {
			panic(err)
		}

		// Change permission again to get rid of any sticky bits
		err = os.Chmod(path, dirMode)
		if err != nil {
			panic(err)
		}
	}
}

// secureFiles walks through the given path and sets the permissions
// for directories and files to the specified modes.
func secureFiles(path string, dirMode os.FileMode, secretMode os.FileMode) {
	// Directory exists, make sure directories and files are secure
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return os.Chmod(path, dirMode)
		}
		return os.Chmod(path, secretMode)
	})
	if err != nil {
		panic(err)
	}
}

// getSemVer returns the semantic version of the input string if it
// matches the pattern `vX.Y.Z`. Otherwise, it returns the input string.
func getSemVer(input string) string {
	// Define the regular expression for semantic versioning
	re := regexp.MustCompile(`^v?(\d+\.\d+\.\d+)$`)

	match := re.FindStringSubmatch(input)

	// If there's a match return the semantic version
	if len(match) > 1 {
		return match[1]
	}

	// If no match, return the original input
	return input
}
