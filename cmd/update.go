package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/engmtcdrm/go-pardon"
	pp "github.com/engmtcdrm/go-prettyprint"
	"github.com/engmtcdrm/mellon/app"
	"github.com/engmtcdrm/mellon/env"
	"github.com/engmtcdrm/mellon/header"
	"github.com/engmtcdrm/mellon/secrets"
	"github.com/engmtcdrm/mellon/secrets/prompts"
)

func init() {
	updateCmd.Flags().StringVarP(
		&secretName,
		"secret",
		"s",
		"",
		"(optional) The name of the secret to update. If -f/--file is provided with this flag, the secret will be updated from the file. If this flag is not provided, you will be prompted to select a secret to update",
	)
	updateCmd.Flags().StringVarP(
		&secretFile,
		"file",
		"f",
		"",
		"(optional) The file containing the unencrypted secret to encrypt",
	)
	updateCmd.Flags().BoolVarP(
		&cleanupFile,
		"cleanup",
		"c",
		false,
		"(optional) Whether to delete the unencrypted secret file after encryption. Defaults to false",
	)

	updateCmd.MarkFlagFilename("file")

	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update a secret",
	Long:    "Update a secret",
	Example: fmt.Sprintf("  %s update", app.Name),
	PreRunE: validateUpdateCreateFlags,
	RunE: func(cmd *cobra.Command, args []string) error {
		var selectedSecret secrets.Secret

		if secretName != "" && secretFile != "" {
			secretPtr := secrets.FindSecretByName(secretName, secretFiles)
			if secretPtr == nil {
				return fmt.Errorf("could not update secret '%s': does not exist", secretName)
			}
			selectedSecret = *secretPtr
			if err := selectedSecret.EncryptFromFile(secretFile, cleanupFile); err != nil {
				return fmt.Errorf("could not encrypt secret from file '%s': %w", secretFile, err)
			}

			return nil
		}

		header.PrintHeader()

		if secretName == "" {
			options, err := prompts.GetSecretOptions(secretFiles, "update")
			if err != nil {
				return err
			}

			promptSelect := pardon.NewSelect(&selectedSecret).
				Title("What secret do you want to update?").
				Options(options...)

			if err := promptSelect.Ask(); err != nil {
				return err
			}

			fmt.Println()
		} else {
			secretPtr := secrets.FindSecretByName(secretName, secretFiles)
			if secretPtr == nil {
				return fmt.Errorf("secret %s does not exist!\n\nUse command %s to create the secret", pp.Red(secretName), pp.Greenf("%s create", env.Instance.ExeCmd()))
			}
			selectedSecret = *secretPtr
		}

		if secretFile == "" {
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
		} else {
			if err := selectedSecret.EncryptFromFile(secretFile, cleanupFile); err != nil {
				return fmt.Errorf("could not encrypt secret from file '%s': %w", secretFile, err)
			}
		}

		fmt.Println(pp.Complete("Secret encrypted and saved"))
		fmt.Println()
		fmt.Printf("You can run the commmand %s to view the unencrypted secret\n", pp.Greenf("%s view -s %s", env.Instance.ExeCmd(), selectedSecret.Name()))

		return nil
	},
}
