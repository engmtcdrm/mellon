package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/engmtcdrm/go-pardon"
	pp "github.com/engmtcdrm/go-prettyprint"
	"github.com/engmtcdrm/mellon/app"
	"github.com/engmtcdrm/mellon/env"
	"github.com/engmtcdrm/mellon/header"
	"github.com/engmtcdrm/mellon/secrets"
)

func init() {
	createCmd.Flags().StringVarP(
		&secretName,
		"secret",
		"s",
		"",
		"(optional) The name of the secret to create",
	)
	createCmd.Flags().StringVarP(
		&secretFile,
		"file",
		"f",
		"",
		"(optional) The file containing the plain text secret to encrypt",
	)
	createCmd.Flags().BoolVarP(
		&cleanupFile,
		"cleanup",
		"c",
		false,
		"(optional) Whether to delete the plain text secret file after encryption",
	)

	createCmd.MarkFlagFilename("file")

	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a secret",
	Long:    "Create a secret.\n\nWhen using the flags -s/--secret and -f/--file, the secret will be read from the specified file and encrypted.\n\nIf no flags are provided, an interactive prompt will be used to enter the secret and its name.",
	Example: fmt.Sprintf("  %s create\n  %s create -s my_secret -f /path/to/secret.txt", app.Name, app.Name),
	PreRunE: validateUpdateCreateFlags,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		var newSecret *secrets.Secret

		if secretName != "" && secretFile != "" {
			secretFilePath := filepath.Join(env.Instance.SecretsPath(), secretName+env.Instance.SecretExt())

			newSecret, err := secrets.NewSecret(env.Instance.KeyPath(), secretName, secretFilePath)
			if err != nil {
				return fmt.Errorf("could not create secret: %w", err)
			}

			if secretPtr := secrets.FindSecretByName(newSecret.Name(), secretFiles); secretPtr != nil {
				return errors.New("secret with that name already exists")
			}

			if err := newSecret.EncryptFromFile(secretFile, cleanupFile); err != nil {
				return fmt.Errorf("could not encrypt secret from file '%s': %w", secretFile, err)
			}

			return nil
		}

		header.PrintHeader()

		if secretName == "" {
			promptQuestion := pardon.NewQuestion(&secretName).
				Title("Enter a name for the secret:").
				Validate(validateSecretName)

			if err := promptQuestion.Ask(); err != nil {
				return err
			}

			fmt.Println()
		} else {
			secretPtr := secrets.FindSecretByName(secretName, secretFiles)
			if secretPtr != nil {
				return fmt.Errorf("secret %s already exists", pp.Red(secretName))
			}
		}

		var secret []byte

		if secretFile == "" {
			promptSecret := pardon.NewPassword(&secret).
				Title("Enter a secret to secure:")

			if err := promptSecret.Ask(); err != nil {
				return err
			}

			fmt.Println()
		}

		newSecret, err = secrets.NewSecret(env.Instance.KeyPath(), secretName, filepath.Join(env.Instance.SecretsPath(), secretName+env.Instance.SecretExt()))
		if err != nil {
			return fmt.Errorf("could not create secret: %w", err)
		}

		if secretFile == "" {
			if err := newSecret.Encrypt(secret); err != nil {
				return fmt.Errorf("could not encrypt secret: %w", err)
			}
		} else {
			if err := newSecret.EncryptFromFile(secretFile, cleanupFile); err != nil {
				return fmt.Errorf("could not encrypt secret from file '%s': %w", secretFile, err)
			}
		}

		fmt.Println(pp.Complete("Secret encrypted and saved"))
		fmt.Println()
		fmt.Printf("You can run the commmand %s to view the unencrypted secret\n", pp.Greenf("%s view -s %s", env.Instance.ExeCmd(), secretName))

		return nil
	},
}
