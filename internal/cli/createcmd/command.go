package createcmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/engmtcdrm/go-pardon"
	pp "github.com/engmtcdrm/go-prettyprint"
	"github.com/spf13/cobra"

	"github.com/engmtcdrm/mellon/internal/app"
	"github.com/engmtcdrm/mellon/internal/env"
	"github.com/engmtcdrm/mellon/internal/header"
	"github.com/engmtcdrm/mellon/secrets"
)

type cmd struct {
	secretFiles []secrets.Secret // List of secrets available.
	secretName  string           // The name of the secret.
	secretFile  string           // The file containing the plain text secret to encrypt.
	cleanupFile bool             // Whether to delete the plain text secret file after encryption.
}

func NewCommand(secretFile []secrets.Secret) *cobra.Command {
	if secretFile == nil {
		secretFile = []secrets.Secret{}
	}

	c := &cmd{secretFiles: secretFile}

	createCmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a secret",
		Long:    "Create a secret.\n\nWhen using the flags -s/--secret and -f/--file, the secret will be read from the specified file and encrypted.\n\nIf no flags are provided, an interactive prompt will be used to enter the secret and its name.",
		Example: fmt.Sprintf("  %s create\n  %s create -s my_secret -f /path/to/secret.txt", app.Name, app.Name),
		PreRunE: c.validateFlags,
		RunE:    c.runE,
	}

	createCmd.Flags().StringVarP(
		&c.secretName,
		"secret",
		"s",
		"",
		"(optional) The name of the secret to create",
	)
	createCmd.Flags().StringVarP(
		&c.secretFile,
		"file",
		"f",
		"",
		"(optional) The file containing the plain text secret to encrypt",
	)
	createCmd.Flags().BoolVarP(
		&c.cleanupFile,
		"cleanup",
		"c",
		false,
		"(optional) Whether to delete the plain text secret file after encryption",
	)

	createCmd.MarkFlagFilename("file")

	return createCmd
}

// encryptFromFile encrypts a secret, if found, from a file.
func (c *cmd) encryptFromFile() error {
	secretFilePath := filepath.Join(env.Instance.SecretsPath(), c.secretName+env.Instance.SecretExt())

	newSecret, err := secrets.NewSecret(env.Instance.KeyPath(), c.secretName, secretFilePath)
	if err != nil {
		return fmt.Errorf("could not create secret: %w", err)
	}

	if secretPtr := secrets.FindSecretByName(newSecret.Name(), c.secretFiles); secretPtr != nil {
		return errors.New("secret with that name already exists")
	}

	if err := newSecret.EncryptFromFile(c.secretFile, c.cleanupFile); err != nil {
		return fmt.Errorf("could not encrypt secret from file '%s': %w", c.secretFile, err)
	}

	return nil
}

// encryptSecret encrypts the given secret. If the secretFile flag is provided,it encrypts the
// secret from the file. Otherwise, it prompts the user to enter the secret.
func (c *cmd) encryptSecret() error {
	var secret []byte
	promptSecret := pardon.NewPassword(&secret).
		Title("Enter a secret to secure:")

	if err := promptSecret.Ask(); err != nil {
		return err
	}

	fmt.Println()

	newSecret, err := secrets.NewSecret(env.Instance.KeyPath(), c.secretName, filepath.Join(env.Instance.SecretsPath(), c.secretName+env.Instance.SecretExt()))
	if err != nil {
		return fmt.Errorf("could not create secret: %w", err)
	}

	if err := newSecret.Encrypt(secret); err != nil {
		return fmt.Errorf("could not encrypt secret: %w", err)
	}

	return nil
}

// resolveSecret prompts the user to select a secret if the secretName flag is not provided.
// If the secretName flag is provided, it returns the corresponding secret if it exists.
func (c *cmd) resolveSecretName() error {
	if c.secretName == "" {
		promptQuestion := pardon.NewQuestion(&c.secretName).
			Title("Enter a name for the secret:").
			Validate(c.validateSecretName)

		if err := promptQuestion.Ask(); err != nil {
			return err
		}

		fmt.Println()
		return nil
	}

	foundSecret := secrets.FindSecretByName(c.secretName, c.secretFiles)
	if foundSecret != nil {
		return fmt.Errorf("secret %s already exists", pp.Red(c.secretName))
	}

	return nil
}

// runE is the main execution function for the command.
func (c *cmd) runE(cmd *cobra.Command, args []string) error {
	if c.secretName != "" && c.secretFile != "" {
		return c.encryptFromFile()
	}

	header.PrintHeader()

	if err := c.resolveSecretName(); err != nil {
		return err
	}

	if c.secretFile == "" {
		if err := c.encryptSecret(); err != nil {
			return err
		}
	} else {
		if err := c.encryptFromFile(); err != nil {
			return err
		}
	}

	fmt.Println(pp.Complete("Secret encrypted and saved"))
	fmt.Println()
	fmt.Printf("You can run the commmand %s to view the unencrypted secret\n", pp.Greenf("%s view -s %s", env.Instance.ExeCmd(), c.secretName))

	return nil
}

// validateFlags checks if the flags for creating or updating a secret are valid.
func (c *cmd) validateFlags(cmd *cobra.Command, args []string) error {
	if c.cleanupFile && (c.secretName == "" || c.secretFile == "") {
		return errors.New("flag -c/--cleanup can only be used when -s/--secret and -f/--file are provided")
	}

	return nil
}

// validateSecretName checks if the provided secret name is valid.
func (c *cmd) validateSecretName(name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}

	if err := secrets.ValidateName(name); err != nil {
		return err
	}

	if secretPtr := secrets.FindSecretByName(name, c.secretFiles); secretPtr != nil {
		return errors.New("secret with that name already exists")
	}

	return nil
}
