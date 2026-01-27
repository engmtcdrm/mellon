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

const confirmationWord = "NAVAER"

var (
	secretFiles []secrets.Secret // List of available secrets.
	secretName  string           // The name of the secret.
	forceDelete bool             // Whether to delete without confirmation.
	deleteAll   bool             // Whether to delete all secrets.
)

func NewCommand(secretFilesList []secrets.Secret) *cobra.Command {
	secretFiles = secretFilesList

	deleteCmd := &cobra.Command{
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
				finalDelete := confirmationWord
				if !forceDelete {
					confirmDelete := false
					promptConfirm2 := pardon.NewConfirm(&confirmDelete).
						Title(fmt.Sprintf("Are you sure you want to delete ALL secrets? %s", pp.Red("There is no going back.")))

					if err := promptConfirm2.Ask(); err != nil {
						return err
					}

					if !confirmDelete {
						fmt.Println()
						fmt.Println(pp.Fail("Aborted deleting all secrets"))
						return nil
					}

					fmt.Println()

					finalDelete = ""
					promptConfirm := pardon.NewQuestion(&finalDelete).
						Title(fmt.Sprintf("To confirm, type %s:", pp.Red(confirmationWord))).
						Icon("")
					if err := promptConfirm.Ask(); err != nil {
						return err
					}

					fmt.Println()
				}

				if finalDelete == confirmationWord {
					for _, secret := range secretFiles {
						if err := secrets.RemoveSecret(env.Instance.SecretsPath(), secret); err != nil {
							return fmt.Errorf("could not remove secret '%s': %w", secret.Name(), err)
						}
					}

					if !forceDelete {
						fmt.Println(pp.Complete("All secrets deleted successfully"))
					}
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
					promptConfirm := pardon.NewConfirm(&confirmDelete).
						Title(fmt.Sprintf("Are you sure you want to delete %s?", pp.Red(secretName)))

					if err := promptConfirm.Ask(); err != nil {
						return err
					}

					fmt.Println()
				}

				if confirmDelete {
					if err := secrets.RemoveSecret(env.Instance.SecretsPath(), selectedSecret); err != nil {
						return fmt.Errorf("could not remove secret '%s': %w", selectedSecret.Name(), err)
					}

					if !forceDelete {
						fmt.Println(pp.Complete("Secret deleted successfully"))
					}
				} else {
					fmt.Println(pp.Fail("Aborted deleting secret"))
				}

				return nil
			}

			header.PrintHeader()

			options, err := prompts.GetSecretOptions(secretFiles, "delete", env.Instance.ExeCmd())
			if err != nil {
				return err
			}

			promptSelect := pardon.NewSelect(&selectedSecret).
				Title("What secret do you want to delete?").
				Options(options...)

			if err := promptSelect.Ask(); err != nil {
				return err
			}

			confirmDelete := true
			if !forceDelete {
				fmt.Println()

				confirmDelete = false
				promptConfirm := pardon.NewConfirm(&confirmDelete).
					Title(fmt.Sprintf("Are you sure you want to delete %s?", pp.Red(selectedSecret.Name())))

				if err := promptConfirm.Ask(); err != nil {
					return err
				}
			}

			fmt.Println()

			if confirmDelete {
				if err := secrets.RemoveSecret(env.Instance.SecretsPath(), selectedSecret); err != nil {
					return fmt.Errorf("could not remove secret '%s': %w", selectedSecret.Name(), err)
				}

				fmt.Println(pp.Complete("Secret deleted successfully"))
			} else {

				fmt.Println(pp.Fail("Aborted deleting secret"))
			}

			return nil
		},
	}

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
	deleteCmd.RegisterFlagCompletionFunc("secret", secretFlagCompletion)

	return deleteCmd
}

// secretFlagCompletion provides shell completion for the -s/--secret flag.
func secretFlagCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var secretNames []string

	for _, secret := range secretFiles {
		secretNames = append(secretNames, secret.Name())
	}
	return secretNames, cobra.ShellCompDirectiveNoFileComp
}
