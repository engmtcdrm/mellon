package deletecmd

import (
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

const confirmationWord = "NAVAER" // The word the user must type to confirm deletion of all secrets.

type cmd struct {
	secretFiles []secrets.Secret // List of available secrets.
	secretName  string           // The name of the secret.
	forceDelete bool             // Whether to delete without confirmation.
	deleteAll   bool             // Whether to delete all secrets.
}

func NewCommand(secretFiles []secrets.Secret) *cobra.Command {
	if secretFiles == nil {
		secretFiles = []secrets.Secret{}
	}

	c := &cmd{secretFiles: secretFiles}

	deleteCmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a secret",
		Long:    "Delete a secret",
		Example: fmt.Sprintf("  %s delete", app.Name),
		RunE:    c.runE,
	}

	deleteCmd.Flags().StringVarP(
		&c.secretName,
		"secret",
		"s",
		"",
		"(optional) The name of the secret to delete",
	)
	deleteCmd.Flags().BoolVarP(
		&c.forceDelete,
		"force",
		"f",
		false,
		"(optional) Whether to force delete the secrets without confirmation",
	)
	deleteCmd.Flags().BoolVar(
		&c.deleteAll,
		"all",
		false,
		"(optional) Whether to delete all secrets",
	)

	deleteCmd.MarkFlagsMutuallyExclusive("secret", "all")
	deleteCmd.RegisterFlagCompletionFunc("secret", c.secretFlagCompletion)

	return deleteCmd
}

// deleteAllSecrets deletes all secrets. If forceDelete is false, it will prompt the user for confirmation.
func (c *cmd) deleteAllSecrets() error {
	deleteConfirmation := confirmationWord

	if !c.forceDelete {
		var confirmDelete bool
		promptInitConfirm := pardon.NewConfirm(&confirmDelete).
			Title(fmt.Sprintf("Are you sure you want to delete ALL secrets? %s", pp.Red("There is no going back.")))

		if err := promptInitConfirm.Ask(); err != nil {
			return err
		}

		if !confirmDelete {
			fmt.Println()
			fmt.Println(pp.Fail("Aborted deleting all secrets"))
			return nil
		}

		fmt.Println()

		var deleteConfirmation string
		promptFinalConfirm := pardon.NewQuestion(&deleteConfirmation).
			Title(fmt.Sprintf("To confirm, type %s:", pp.Red(confirmationWord))).
			Icon("")
		if err := promptFinalConfirm.Ask(); err != nil {
			return err
		}

		fmt.Println()
	}

	if deleteConfirmation == confirmationWord {
		for _, secret := range c.secretFiles {
			if err := secrets.RemoveSecret(env.Instance.SecretsPath(), secret); err != nil {
				return fmt.Errorf("could not remove secret '%s': %w", secret.Name(), err)
			}
		}

		if !c.forceDelete {
			fmt.Println(pp.Complete("All secrets deleted successfully"))
		}
	} else {
		fmt.Println(pp.Fail("Aborted deleting all secrets"))
	}

	return nil
}

// deleteSecret deletes a single secret. If confirmDelete is false, it will not delete the secret.
func (c *cmd) deleteSecret(secret *secrets.Secret, confirmDelete bool) error {
	if confirmDelete {
		if err := secrets.RemoveSecret(env.Instance.SecretsPath(), *secret); err != nil {
			return fmt.Errorf("could not remove secret '%s': %w", secret.Name(), err)
		}

		if !c.forceDelete {
			fmt.Println(pp.Complete("Secret deleted successfully"))
		}

		return nil
	}

	fmt.Println(pp.Fail("Aborted deleting secret"))

	return nil
}

// findAndDeleteSecret finds a secret by name and deletes it if the user confirms the deletion.
func (c *cmd) findAndDeleteSecret() error {
	foundSecret := secrets.FindSecretByName(c.secretName, c.secretFiles)
	if foundSecret == nil {
		return fmt.Errorf("could not delete secret '%s': does not exist", c.secretName)
	}

	confirmDelete := true
	if !c.forceDelete {
		confirmDelete = false
		promptConfirm := pardon.NewConfirm(&confirmDelete).
			Title(fmt.Sprintf("Are you sure you want to delete %s?", pp.Red(c.secretName)))

		if err := promptConfirm.Ask(); err != nil {
			return err
		}

		fmt.Println()
	}

	return c.deleteSecret(foundSecret, confirmDelete)
}

// secretFlagCompletion provides shell completion for the -s/--secret flag.
func (c *cmd) secretFlagCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var secretNames []string

	for _, secret := range c.secretFiles {
		secretNames = append(secretNames, secret.Name())
	}
	return secretNames, cobra.ShellCompDirectiveNoFileComp
}

// selectAndDeleteSecret prompts the user to select a secret and deletes it if the user confirms the deletion.
func (c *cmd) selectAndDeleteSecret() error {
	header.PrintHeader()

	options, err := prompts.GetSecretOptions(c.secretFiles, "delete", env.Instance.ExeCmd())
	if err != nil {
		return err
	}

	var selectedSecret secrets.Secret
	promptSelect := pardon.NewSelect(&selectedSecret).
		Title("What secret do you want to delete?").
		Options(options...)
	if err := promptSelect.Ask(); err != nil {
		return err
	}

	confirmDelete := true
	if !c.forceDelete {
		fmt.Println()

		confirmDelete = false
		promptConfirm := pardon.NewConfirm(&confirmDelete).
			Title(fmt.Sprintf("Are you sure you want to delete %s?", pp.Red(selectedSecret.Name())))
		if err := promptConfirm.Ask(); err != nil {
			return err
		}
	}

	fmt.Println()

	return c.deleteSecret(&selectedSecret, confirmDelete)
}

// runE is the main execution function for the command.
func (c *cmd) runE(cmd *cobra.Command, args []string) error {
	if !c.forceDelete {
		header.PrintHeader()
	}

	if c.deleteAll {
		return c.deleteAllSecrets()
	}

	if c.secretName != "" {
		return c.findAndDeleteSecret()
	}

	return c.selectAndDeleteSecret()
}
