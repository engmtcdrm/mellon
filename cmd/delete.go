package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
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
				form := huh.NewForm(
					huh.NewGroup(
						huh.NewConfirm().
							Title(fmt.Sprintf("Are you sure you want to delete '%s'?", secretName)).
							Description("This action cannot be undone.").
							Value(&confirmDelete),
					),
				)

				err := form.
					WithTheme(app.ThemeMinno()).
					Run()
				if err != nil {
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

		groups := []*huh.Group{
			huh.NewGroup(
				huh.NewSelect[secrets.Secret]().
					Options(options...).
					Title("Available Secrets").
					Description("Choose a secret to delete.").
					Value(&selectedSecret),
			),
		}

		confirmDelete := true
		if !forceDelete {
			confirmDelete = false
			groups = append(groups, huh.NewGroup(
				huh.NewConfirm().
					Title(fmt.Sprintf("Are you sure you want to delete '%s'?", selectedSecret.Name)).
					Description("This action cannot be undone.").
					Value(&confirmDelete),
			))
		}

		form := huh.NewForm(groups...)

		err = form.
			WithTheme(app.ThemeMinno()).
			Run()
		if err != nil {
			return err
		}

		if confirmDelete {
			if err := os.Remove(selectedSecret.Path); err != nil {
				return fmt.Errorf("could not remove secret '%s': %w", selectedSecret.Name, err)
			}
		}

		return nil
	},
}
