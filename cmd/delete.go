package cmd

import (
	"fmt"
	"os"

	"github.com/engmtcdrm/minno/app"
	"github.com/engmtcdrm/minno/secrets"
	"github.com/spf13/cobra"
)

func init() {
	deleteCmd.Flags().StringVarP(
		&secretName,
		"secret",
		"s",
		"",
		"(optional) The name of the secret to update. If -f/--file is provided with this flag, the secret will be updated from the file. If this flag is not provided, you will be prompted to select a secret to update.",
	)

	rootCmd.AddCommand(deleteCmd)
}

func validateDeleteFlags(cmd *cobra.Command, args []string) error {
	return nil
}

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a secret",
	Long:    "Delete a secret",
	Example: fmt.Sprintf("  %s delete", app.Name),
	PreRunE: validateDeleteFlags,
	RunE: func(cmd *cobra.Command, args []string) error {
		var selectedSecret secrets.Secret

		if secretName != "" {
			if err := secrets.ValidateName(secretName); err != nil {
				return fmt.Errorf("%s. The secret name provided was '%s'", err, secretName)
			}

			secretPtr := secrets.FindSecretByName(secretName, secretFiles)
			if secretPtr == nil {
				return fmt.Errorf("could not delete secret '%s': does not exist", secretName)
			}
			selectedSecret = *secretPtr

			if err := os.Remove(selectedSecret.Path); err != nil {
				return fmt.Errorf("could not remove secret '%s': %w", selectedSecret.Name, err)
			}

			return nil
		}

		return nil
	},
}
