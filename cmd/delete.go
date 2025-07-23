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
	deleteCmd.Flags().BoolVar(
		&deleteAll,
		"all",
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

		if !forceDelete {
			header.PrintHeader()
		}

		if deleteAll {
			confirmDelete := "NAVAER"
			if !forceDelete {
				confirmDelete2 := false
				promptConfirm2 := pardon.NewConfirm().
					Title(fmt.Sprintf("Are you sure you want to delete ALL secrets? %s", pp.Red("There is no going back."))).
					Value(&confirmDelete2)
				if err := promptConfirm2.Ask(); err != nil {
					return err
				}

				if !confirmDelete2 {
					fmt.Println()
					fmt.Println(pp.Fail("Aborted deleting all secrets"))
					return nil
				}

				fmt.Println()

				confirmDelete = ""
				promptConfirm := pardon.NewQuestion().
					Title(fmt.Sprintf("To confirm, type %s:", pp.Red("NAVAER"))).
					Icon("").
					Value(&confirmDelete)
				if err := promptConfirm.Ask(); err != nil {
					return err
				}

				fmt.Println()
			}

			if confirmDelete == "NAVAER" {
				secrets, err := secrets.GetSecretFiles()
				if err != nil {
					return fmt.Errorf("could not retrieve secrets: %w", err)
				}

				for _, secret := range secrets {
					if err := os.Remove(secret.Path); err != nil {
						return fmt.Errorf("could not remove secret '%s': %w", secret.Name, err)
					}
				}

				fmt.Println(pp.Complete("All secrets deleted successfully"))
			} else {
				fmt.Println(pp.Fail("Aborted deleting all secrets"))
			}

			return nil
		}

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
					Title(fmt.Sprintf("Are you sure you want to delete %s?", pp.Red(secretName))).
					Value(&confirmDelete)
				if err := promptConfirm.Ask(); err != nil {
					return err
				}

				fmt.Println()
			}

			if confirmDelete {
				if err := os.Remove(selectedSecret.Path); err != nil {
					return fmt.Errorf("could not remove secret '%s': %w", selectedSecret.Name, err)
				}

				fmt.Println(pp.Complete("All secrets deleted successfully"))
			} else {
				fmt.Println(pp.Fail("Aborted deleting secret"))
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
			Value(&selectedSecret)

		if err := promptSelect.Ask(); err != nil {
			return err
		}

		confirmDelete := true
		if !forceDelete {
			fmt.Println()

			confirmDelete = false
			promptConfirm := pardon.NewConfirm().
				Title(fmt.Sprintf("Are you sure you want to delete %s?", pp.Red(selectedSecret.Name))).
				Value(&confirmDelete)
			if err := promptConfirm.Ask(); err != nil {
				return err
			}
		}

		fmt.Println()

		if confirmDelete {
			if err := os.Remove(selectedSecret.Path); err != nil {
				return fmt.Errorf("could not remove secret '%s': %w", selectedSecret.Name, err)
			}

			fmt.Println(pp.Complete("All secrets deleted successfully"))
		} else {

			fmt.Println(pp.Fail("Aborted deleting secret"))
		}

		return nil
	},
}
