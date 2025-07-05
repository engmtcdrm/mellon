package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/engmtcdrm/go-entomb"
	pp "github.com/engmtcdrm/go-prettyprint"
	"github.com/engmtcdrm/minno/app"
	"github.com/engmtcdrm/minno/env"
	"github.com/engmtcdrm/minno/header"
	"github.com/engmtcdrm/minno/secrets"
	"github.com/engmtcdrm/minno/secrets/prompts"
)

func init() {
	updateCmd.Flags().StringVarP(
		&secretName,
		"secret",
		"s",
		"",
		"(optional) The name of the secret to update. If -f/--file is provided with this flag, the secret will be updated from the file. If this flag is not provided, you will be prompted to select a secret to update.",
	)
	updateCmd.Flags().StringVarP(
		&rawSecretFile,
		"file",
		"f",
		"",
		"(optional) The file containing the unencrypted secret to encrypt. If this is provided then -s/--secret must also be provided",
	)
	updateCmd.Flags().BoolVarP(
		&cleanupFile,
		"cleanup",
		"c",
		false,
		"(optional) Whether to delete the unencrypted secret file after encryption. Defaults to false",
	)

	rootCmd.AddCommand(updateCmd)
}

func validateUpdateFlags(cmd *cobra.Command, args []string) error {
	// Make sure both flags are provided if one is used
	if rawSecretFile != "" && secretName == "" {
		return errors.New("flag -f/--file must be provided with -s/--secret")
	}

	if cleanupFile && (secretName == "" || rawSecretFile == "") {
		return errors.New("flag -c/--cleanup can only be used when -s/--secret and -f/--file are provided")
	}

	return nil
}

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update a secret",
	Long:    "Update a secret",
	Example: fmt.Sprintf("  %s update", app.Name),
	PreRunE: validateUpdateFlags,
	RunE: func(cmd *cobra.Command, args []string) error {
		tomb, err := entomb.NewTomb(envVars.KeyPath)
		if err != nil {
			return err
		}

		var secret string
		var selectedSecret secrets.Secret
		var encSecret []byte

		if secretName != "" && rawSecretFile != "" {
			if err := secrets.ValidateName(secretName); err != nil {
				return fmt.Errorf("%s. The secret name provided was '%s'", err, secretName)
			}

			secretPtr := secrets.FindSecretByName(secretName, secretFiles)
			if secretPtr == nil {
				return fmt.Errorf("could not update secret '%s': does not exist", secretName)
			}
			selectedSecret = *secretPtr

			rawSecretFile, err := env.ExpandTilde(strings.TrimSpace(rawSecretFile))
			if err != nil {
				return err
			}

			secretBytes, err := os.ReadFile(rawSecretFile)
			if err != nil {
				return fmt.Errorf("could not read file '%s': %w", rawSecretFile, err)
			}

			secret = strings.TrimSpace(string(secretBytes))
			secretBytes = nil

			encSecret, err = tomb.Encrypt([]byte(secret))
			if err != nil {
				return err
			}

			if err = os.WriteFile(selectedSecret.Path, encSecret, secretMode); err != nil {
				return err
			}

			if cleanupFile {
				if err = os.Remove(rawSecretFile); err != nil {
					return fmt.Errorf("could not remove file '%s': %w", rawSecretFile, err)
				}
			}

			return nil
		}

		header.PrintHeader()

		var form *huh.Form

		if secretName == "" {
			options, err := prompts.GetSecretOptions(secretFiles, "update")
			if err != nil {
				return err
			}

			form = huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[secrets.Secret]().
						Options(options...).
						Title("Available Secrets").
						Description("Choose a secret to update.").
						Value(&selectedSecret),
					huh.NewInput().
						Title("Enter the updated secret").
						Value(&secret).
						EchoMode(huh.EchoModeNone).
						Inline(true),
				),
			)
		} else {
			if err := secrets.ValidateName(secretName); err != nil {
				return fmt.Errorf("%s\n\nThe secret name provided was %s", err, pp.Red(secretName))
			}

			secretPtr := secrets.FindSecretByName(secretName, secretFiles)
			if secretPtr == nil {
				return fmt.Errorf("secret %s does not exist!\n\nUse command %s to create the secret", pp.Red(secretName), pp.Greenf("%s create", envVars.ExeCmd))
			}
			selectedSecret = *secretPtr

			form = huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Enter the updated secret").
						Value(&secret).
						EchoMode(huh.EchoModeNone).
						Inline(true),
				),
			)
		}

		err = form.
			WithTheme(app.ThemeMinno()).
			Run()
		if err != nil {
			return err
		}

		encSecret, err = tomb.Encrypt([]byte(strings.TrimSpace(secret)))
		secret = ""
		if err != nil {
			return err
		}

		fmt.Println(pp.Complete("Secret encrypted"))

		if err = os.WriteFile(selectedSecret.Path, encSecret, secretMode); err != nil {
			return err
		}
		fmt.Println(pp.Complete("Secret saved"))
		fmt.Println()
		fmt.Printf("You can run the commmand %s to view the unencrypted secret\n", pp.Greenf("%s view -s %s", envVars.ExeCmd, selectedSecret.Name))

		return nil
	},
}
