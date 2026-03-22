package updatecmd

import (
	"errors"
	"fmt"

	"github.com/engmtcdrm/go-pardon"
	pp "github.com/engmtcdrm/go-prettyprint"
	"github.com/spf13/cobra"

	"github.com/engmtcdrm/mellon/internal/app"
	"github.com/engmtcdrm/mellon/internal/env"
	"github.com/engmtcdrm/mellon/internal/header"
	"github.com/engmtcdrm/mellon/secrets"
	"github.com/engmtcdrm/mellon/secrets/prompts"
)

type cmd struct {
	secretFiles []secrets.Secret // List of available secrets.
	secretName  string           // The name of the secret.
	secretFile  string           // The file containing the unencrypted secret to encrypt.
	cleanupFile bool             // Whether to delete the unencrypted secret file after encryption.
}

func NewCommand(secretFiles []secrets.Secret) *cobra.Command {
	if secretFiles == nil {
		secretFiles = []secrets.Secret{}
	}

	c := &cmd{secretFiles: secretFiles}

	updateCmd := &cobra.Command{
		Use:     "update",
		Short:   "Update a secret",
		Long:    "Update a secret",
		Example: fmt.Sprintf("  %s update", app.Name),
		PreRunE: c.validateUpdateCreateFlags,
		RunE:    c.runE,
	}

	updateCmd.Flags().StringVarP(
		&c.secretName,
		"secret",
		"s",
		"",
		"(optional) The name of the secret to update. If -f/--file is provided with this flag, the secret will be updated from the file. If this flag is not provided, you will be prompted to select a secret to update",
	)
	updateCmd.Flags().StringVarP(
		&c.secretFile,
		"file",
		"f",
		"",
		"(optional) The file containing the unencrypted secret to encrypt",
	)
	updateCmd.Flags().BoolVarP(
		&c.cleanupFile,
		"cleanup",
		"c",
		false,
		"(optional) Whether to delete the unencrypted secret file after encryption. Defaults to false",
	)

	updateCmd.MarkFlagFilename("file")
	updateCmd.RegisterFlagCompletionFunc("secret", c.secretFlagCompletion)

	return updateCmd
}

// encryptFromFile encrypts a secret, if found, from a file.
func (c *cmd) encryptFromFile() error {
	foundSecret := secrets.FindSecretByName(c.secretName, c.secretFiles)
	if foundSecret == nil {
		return fmt.Errorf("could not update secret '%s': does not exist", c.secretName)
	}

	if err := foundSecret.EncryptFromFile(c.secretFile, c.cleanupFile); err != nil {
		return fmt.Errorf("could not encrypt secret from file '%s': %w", c.secretFile, err)
	}

	return nil
}

// encryptSecret encrypts the given secret. If the secretFile flag is provided,it encrypts the
// secret from the file. Otherwise, it prompts the user to enter the secret.
func (c *cmd) encryptSecret(selectedSecret *secrets.Secret) error {
	if c.secretFile == "" {
		var secret []byte

		promptSecret := pardon.NewPassword(&secret).
			Title("Enter the updated secret:")

		if err := promptSecret.Ask(); err != nil {
			return err
		}

		if err := selectedSecret.Encrypt(secret); err != nil {
			return fmt.Errorf("could not encrypt secret: %w", err)
		}

		fmt.Println()
		return nil
	}

	if err := selectedSecret.EncryptFromFile(c.secretFile, c.cleanupFile); err != nil {
		return fmt.Errorf("could not encrypt secret from file '%s': %w", c.secretFile, err)
	}

	return nil
}

// resolveSecret prompts the user to select a secret if the secretName flag is not provided.
// If the secretName flag is provided, it returns the corresponding secret if it exists.
func (c *cmd) resolveSecret() (*secrets.Secret, error) {

	if c.secretName == "" {
		var selectedSecret secrets.Secret

		options, err := prompts.GetSecretOptions(c.secretFiles, "update", env.Instance.ExeCmd())
		if err != nil {
			return nil, err
		}

		promptSelect := pardon.NewSelect(&selectedSecret).
			Title("What secret do you want to update?").
			Options(options...)

		if err := promptSelect.Ask(); err != nil {
			return nil, err
		}

		fmt.Println()
		return &selectedSecret, nil
	}

	foundSecret := secrets.FindSecretByName(c.secretName, c.secretFiles)
	if foundSecret == nil {
		return nil, fmt.Errorf("secret %s does not exist!\n\nUse command %s to create the secret", pp.Red(c.secretName), pp.Greenf("%s create", env.Instance.ExeCmd()))
	}

	return foundSecret, nil
}

// secretFlagCompletion provides shell completion for the -s/--secret flag.
func (c *cmd) secretFlagCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var secretNames []string

	for _, secret := range c.secretFiles {
		secretNames = append(secretNames, secret.Name())
	}

	return secretNames, cobra.ShellCompDirectiveNoFileComp
}

// runE is the main execution function for the command.
func (c *cmd) runE(cmd *cobra.Command, args []string) error {
	if c.secretName != "" && c.secretFile != "" {
		return c.encryptFromFile()
	}

	header.PrintHeader()

	selectedSecret, err := c.resolveSecret()
	if err != nil {
		return err
	}

	if err := c.encryptSecret(selectedSecret); err != nil {
		return err
	}

	fmt.Println(pp.Complete("Secret encrypted and saved"))
	fmt.Println()
	fmt.Printf("You can run the commmand %s to view the unencrypted secret\n", pp.Greenf("%s view -s %s", env.Instance.ExeCmd(), selectedSecret.Name()))

	return nil
}

// validateUpdateCreateFlags checks if the flags for creating or updating a secret are valid.
func (c *cmd) validateUpdateCreateFlags(cmd *cobra.Command, args []string) error {
	if c.cleanupFile && (c.secretName == "" || c.secretFile == "") {
		return errors.New("flag -c/--cleanup can only be used when -s/--secret and -f/--file are provided")
	}

	return nil
}
