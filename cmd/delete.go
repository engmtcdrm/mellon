package cmd

import (
	"fmt"
	"os"

	"github.com/engmtcdrm/go-pardon"
	pp "github.com/engmtcdrm/go-prettyprint"
	"github.com/engmtcdrm/minno/app"
	"github.com/engmtcdrm/minno/header"
	"github.com/engmtcdrm/minno/secrets"
	"github.com/engmtcdrm/minno/secrets/prompts"
	"github.com/spf13/cobra"
)

func init() {
	deleteCmd.Flags().StringVarP(
		&secretName,
		"secret",
		"s",
		"",
		"(optional) The name of the secret to delete",
	)
	deleteCmd.Flags().BoolVarP(
		&forceDelete,
		"force",
		"f",
		false,
		"(optional) Whether to force delete the secrets without confirmation",
	)
	deleteCmd.Flags().BoolVarP(
		&deleteAll,
		"all",
		"a",
		false,
		"(optional) Whether to delete all secrets",
	)

	deleteCmd.MarkFlagsMutuallyExclusive("secret", "all")

	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a secret",
	Long:    "Delete a secret",
	Example: fmt.Sprintf("  %s delete", app.Name),
	RunE: func(cmd *cobra.Command, args []string) error {
		var selectedSecret secrets.Secret

		// TODO: Add logic to delete all secrets if deleteAll is true.

		if secretName != "" {
			secretPtr := secrets.FindSecretByName(secretName, secretFiles)
			if secretPtr == nil {
				return fmt.Errorf("could not delete secret '%s': does not exist", secretName)
			}
			selectedSecret = *secretPtr

			confirmDelete := true
			if !forceDelete {
				confirmDelete = false
				promptConfirm := pardon.NewConfirm().
					Title(fmt.Sprintf("Are you sure you want to delete '%s'?", secretName)).
					Value(&confirmDelete)
				if err := promptConfirm.Ask(); err != nil {
					return err
				}
			}

			if confirmDelete {
				if err := os.Remove(selectedSecret.Path); err != nil {
					return fmt.Errorf("could not remove secret '%s': %w", selectedSecret.Name, err)
				}
			}

			return nil
		}

		header.PrintHeader()

		options, err := prompts.GetSecretOptions(secretFiles, "delete")
		if err != nil {
			return err
		}

		promptSelect := pardon.NewSelect[secrets.Secret]().
			Title("What secret do you want to delete?").
			Options(options...).
			Value(&selectedSecret).
			SelectFunc(
				func(s string) string {
					return pp.Yellow(s)
				}).
			Icon(pp.Cyan(pardon.Icons.QuestionMark))

		if err := promptSelect.Ask(); err != nil {
			return err
		}

		confirmDelete := true
		if !forceDelete {
			confirmDelete = false
			promptConfirm := pardon.NewConfirm().
				Title(fmt.Sprintf("Are you sure you want to delete %s?", pp.Cyan(selectedSecret.Name))).
				Icon(pp.Cyan(pardon.Icons.QuestionMark)).
				Value(&confirmDelete)
			if err := promptConfirm.Ask(); err != nil {
				return err
			}
		}

		if confirmDelete {
			if err := os.Remove(selectedSecret.Path); err != nil {
				return fmt.Errorf("could not remove secret '%s': %w", selectedSecret.Name, err)
			}
		}

		return nil
	},
}
