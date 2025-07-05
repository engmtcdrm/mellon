package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/engmtcdrm/go-entomb"
	pp "github.com/engmtcdrm/go-prettyprint"
	"github.com/engmtcdrm/minno/app"
	"github.com/engmtcdrm/minno/env"
	"github.com/engmtcdrm/minno/header"
	"github.com/engmtcdrm/minno/secrets"
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
		&rawSecretFile,
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
		"(optional) Whether to delete the plain text secret file after encryption. Defaults to false",
	)

	rootCmd.AddCommand(createCmd)
}

func validateCreateFlags(cmd *cobra.Command, args []string) error {
	// Make sure both flags are provided if one is used
	if rawSecretFile != "" && secretName == "" {
		return errors.New("flag -f/--file must be provided with -s/--secret")
	}

	if cleanupFile && (secretName == "" || rawSecretFile == "") {
		return errors.New("flag -c/--cleanup can only be used when -s/--secret and -f/--file are provided")
	}

	return nil
}

func validateSecretName(name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}

	if err := secrets.ValidateName(name); err != nil {
		return err
	}

	if secretPtr := secrets.FindSecretByName(name, secretFiles); secretPtr != nil {
		return errors.New("secret with that name already exists")
	}

	return nil
}

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a secret",
	Long:    "Create a secret.\n\nWhen using the flags -s/--secret and -f/--file, the secret will be read from the specified file and encrypted.\n\nIf no flags are provided, an interactive prompt will be used to enter the secret and its name.",
	Example: fmt.Sprintf("  %s create\n  %s create -n my_secret -f /path/to/secret.txt", app.Name, app.Name),
	PreRunE: validateCreateFlags,
	RunE: func(cmd *cobra.Command, args []string) error {
		tomb, err := entomb.NewTomb(envVars.KeyPath)
		if err != nil {
			return err
		}

		var secret string
		var secretFile string
		var encSecret []byte

		if secretName != "" && rawSecretFile != "" {
			if err := validateSecretName(secretName); err != nil {
				return err
			}

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

			secretFile = secretName
			secretFilePath := filepath.Join(envVars.SecretsPath, secretFile+envVars.SecretExt)

			if err = os.WriteFile(secretFilePath, encSecret, secretMode); err != nil {
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

		// Interactive mode if no flags are provided
		if secretName == "" {
			form = huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Enter a secret to secure").
						Value(&secret).
						EchoMode(huh.EchoModeNone).
						Inline(true),
					huh.NewInput().
						Title("Enter a name for the secret").
						Value(&secretFile).
						Validate(validateSecretName).
						Inline(true),
				),
			)

		} else {
			if err := secrets.ValidateName(secretName); err != nil {
				return fmt.Errorf("%s\n\nThe secret name provided was %s", err, pp.Red(secretName))
			}

			secretPtr := secrets.FindSecretByName(secretName, secretFiles)
			if secretPtr != nil {
				return fmt.Errorf("secret %s already exists", pp.Red(secretName))
			}

			secretFile = secretName

			form = huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Enter a secret to secure").
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

		secretFilePath := filepath.Join(envVars.SecretsPath, secretFile+envVars.SecretExt)

		if err = os.WriteFile(secretFilePath, encSecret, secretMode); err != nil {
			return err
		}

		fmt.Println(pp.Complete("Secret saved"))
		fmt.Println()
		fmt.Printf("You can run the commmand %s to view the unencrypted secret\n", pp.Greenf("%s view -s %s", envVars.ExeCmd, secretFile))

		return nil
	},
}
