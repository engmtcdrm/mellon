package viewcmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/engmtcdrm/go-pardon"
	pp "github.com/engmtcdrm/go-prettyprint"
	"github.com/spf13/cobra"

	"github.com/engmtcdrm/mellon/internal/app"
	"github.com/engmtcdrm/mellon/internal/constants"
	"github.com/engmtcdrm/mellon/internal/env"
	"github.com/engmtcdrm/mellon/internal/header"
	"github.com/engmtcdrm/mellon/secrets"
	"github.com/engmtcdrm/mellon/secrets/prompts"
)

type cmd struct {
	secretFiles []secrets.Secret // List of available secrets.
	secretName  string           // The name of the secret.
	output      string           // The file to write decrypted secret to.
}

func NewCommand(secretFiles []secrets.Secret) *cobra.Command {
	if secretFiles == nil {
		secretFiles = []secrets.Secret{}
	}

	c := &cmd{secretFiles: secretFiles}

	viewCmd := &cobra.Command{
		Use:     "view",
		Short:   "View a secret",
		Long:    "View a secret",
		Example: fmt.Sprintf("  %s view\n  %s view -s awesome-secret", app.Name, app.Name),
		PreRunE: c.validateViewFlags,
		RunE:    c.runE,
	}

	viewCmd.Flags().StringVarP(
		&c.secretName,
		"secret",
		"s",
		"",
		"(optional) The name of the secret to view. Only names containing alphanumeric, hyphens, and underscores are allowed",
	)
	viewCmd.Flags().StringVarP(
		&c.output,
		"output",
		"o",
		"",
		"(optional) File to write decrypted secret to. Defaults to outputting to stdout. This only works with the option -s/--secret",
	)

	viewCmd.RegisterFlagCompletionFunc("secret", c.secretFlagCompletion)

	return viewCmd
}

// promptViewSecret prompts the user to select a secret to view.
func (c *cmd) promptViewSecret() error {
	var selectedSecretFile secrets.Secret
	header.PrintHeader()

	options, err := prompts.GetSecretOptions(c.secretFiles, "view", env.Instance.ExeCmd())
	if err != nil {
		return err
	}

	promptSelect := pardon.NewSelect(&selectedSecretFile).
		Options(options...).
		Title("What secret do you want to view?")

	if err := promptSelect.Ask(); err != nil {
		return err
	}

	secret, err := selectedSecretFile.Decrypt()
	if err != nil {
		return errors.New("failed to decrypt secret. Encrypted secret may be corrupted")
	}

	fmt.Println()
	fmt.Println(pp.Complete("Secret decrypted"))
	fmt.Println()
	fmt.Println(pp.Info("The secret is " + pp.Green(string(secret))))

	return nil
}

// runE is the main execution function for the command.
func (c *cmd) runE(cmd *cobra.Command, args []string) error {
	if c.secretName == "" {
		return c.promptViewSecret()
	}

	return c.viewSecret()
}

// secretFlagCompletion provides shell completion for the -s/--secret flag.
func (c *cmd) secretFlagCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var secretNames []string

	for _, secret := range c.secretFiles {
		secretNames = append(secretNames, secret.Name())
	}
	return secretNames, cobra.ShellCompDirectiveNoFileComp
}

// validateViewFlags checks if the flags for viewing a secret are valid.
func (c *cmd) validateViewFlags(cmd *cobra.Command, args []string) error {
	if c.output != "" && c.secretName == "" {
		return errors.New("flag -o/--output can only be used when -s/--secret is provided")
	}

	return nil
}

// viewSecret decrypts and displays the secret specified by the secretName flag.
// If the output flag is provided, the decrypted secret is written to the specified file.
func (c *cmd) viewSecret() error {
	secretPtr := secrets.FindSecretByName(c.secretName, c.secretFiles)
	if secretPtr == nil {
		return fmt.Errorf("failed to read secret '%s': secret does not exist", c.secretName)
	}

	secret, err := secretPtr.Decrypt()
	if err != nil {
		return fmt.Errorf("failed to decrypt secret '%s'. Encrypted secret may be corrupted", c.secretName)
	}

	if c.output == "" {
		fmt.Print(string(secret))
	} else {
		outputDir := filepath.Dir(c.output)
		if _, err := os.Stat(outputDir); os.IsNotExist(err) {
			err = os.MkdirAll(outputDir, constants.SecureDirMode)
			if err != nil {
				return fmt.Errorf("failed to create output directory for output file '%s'", c.output)
			}
		}

		err = os.WriteFile(c.output, secret, constants.SecureFileMode)
		if err != nil {
			return fmt.Errorf("failed to write secret to output file '%s'", c.output)
		}
	}

	secret = nil

	return nil
}
