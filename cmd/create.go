package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/engmtcdrm/go-pardon"
	pp "github.com/engmtcdrm/go-prettyprint"
	"github.com/engmtcdrm/mellon/app"
	"github.com/engmtcdrm/mellon/header"
	"github.com/engmtcdrm/mellon/secrets"
)

func init() {
	createCmd.Flags().StringVarP(
		&secretName,
		"secret",
		"s",
		"",
		"(optional) The name of the secret to create. If this is provided then -f/--file must also be provided",
	)
	createCmd.Flags().StringVarP(
		&secretFile,
		"file",
		"f",
		"",
		"(optional) The file containing the plain text secret to encrypt. If this is provided then -s/--secret must also be provided",
	)
	createCmd.Flags().BoolVarP(
		&cleanupFile,
		"cleanup",
		"c",
		false,
		"(optional) Whether to delete the plain text secret file after encryption",
	)

	createCmd.MarkFlagsRequiredTogether("file", "secret")
	createCmd.MarkFlagFilename("file")

	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a secret",
	Long:    "Create a secret.\n\nWhen using the flags -s/--secret and -f/--file, the secret will be read from the specified file and encrypted.\n\nIf no flags are provided, an interactive prompt will be used to enter the secret and its name.",
	Example: fmt.Sprintf("  %s create\n  %s create -n my_secret -f /path/to/secret.txt", app.Name, app.Name),
	PreRunE: validateUpdateCreateFlags,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		var newSecret *secrets.Secret

		if secretName != "" && secretFile != "" {
			secretFilePath := filepath.Join(envVars.SecretsPath, secretName+envVars.SecretExt)

			newSecret, err = secrets.NewSecret(envVars.KeyPath, secretName, secretFilePath)
			if err != nil {
				return fmt.Errorf("could not create secret: %w", err)
			}

			if secretPtr := secrets.FindSecretByName(newSecret.Name, secretFiles); secretPtr != nil {
				return errors.New("secret with that name already exists")
			}

			if err := newSecret.EncryptFromFile(secretFile, cleanupFile); err != nil {
				return fmt.Errorf("could not encrypt secret from file '%s': %w", secretFile, err)
			}

			return nil
		}

		header.PrintHeader()

		var secret []byte

		promptSecret := pardon.NewPassword().
			Title("Enter a secret to secure:").
			Value(&secret)

		if err := promptSecret.Ask(); err != nil {
			return err
		}

		fmt.Println()

		if secretName == "" {
			promptQuestion := pardon.NewQuestion().
				Title("Enter a name for the secret:").
				Value(&secretName).
				Validate(validateSecretName)

			if err := promptQuestion.Ask(); err != nil {
				return err
			}

		} else {
			secretPtr := secrets.FindSecretByName(secretName, secretFiles)
			if secretPtr != nil {
				return fmt.Errorf("secret %s already exists", pp.Red(secretName))
			}
		}

		fmt.Println()

		newSecret, err = secrets.NewSecret(envVars.KeyPath, secretName, filepath.Join(envVars.SecretsPath, secretName+envVars.SecretExt))
		if err != nil {
			return fmt.Errorf("could not create secret: %w", err)
		}

		if err := newSecret.Encrypt(secret); err != nil {
			return fmt.Errorf("could not encrypt secret: %w", err)
		}

		fmt.Println(pp.Complete("Secret encrypted and saved"))
		fmt.Println()
		fmt.Printf("You can run the commmand %s to view the unencrypted secret\n", pp.Greenf("%s view -s %s", envVars.ExeCmd, secretName))

		return nil
	},
}
