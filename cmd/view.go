package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/engmtcdrm/go-entomb"
	pp "github.com/engmtcdrm/go-prettyprint"
	"github.com/engmtcdrm/minno/app"
	"github.com/engmtcdrm/minno/header"
	"github.com/engmtcdrm/minno/secrets"
	"github.com/engmtcdrm/minno/secrets/prompts"
)

func init() {
	viewCmd.Flags().StringVarP(
		&secretName,
		"secret",
		"s",
		"",
		"(optional) The name of the secret to view. Only names containing alphanumeric, hyphens, and underscores are allowed.",
	)
	viewCmd.Flags().StringVarP(
		&output,
		"output",
		"o",
		"",
		"(optional) File to write decrypted secret to. Defaults to outputting to stdout. This only works with the option -s/--secret",
	)

	rootCmd.AddCommand(viewCmd)
}

var viewCmd = &cobra.Command{
	Use:     "view",
	Short:   "View a secret",
	Long:    "View a secret",
	Example: fmt.Sprintf("  %s view\n  %s view -s awesome-secret", app.Name, app.Name),
	RunE: func(cmd *cobra.Command, args []string) error {
		tomb, err := entomb.NewTomb(envVars.KeyPath)
		if err != nil {
			return err
		}

		var selectedSecretFile secrets.Secret

		if secretName == "" {
			header.PrintHeader()

			options, err := prompts.GetSecretOptions(secretFiles, "view")
			if err != nil {
				return err
			}

			form := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[secrets.Secret]().
						Options(options...).
						Title("Available Secrets").
						Description("Choose a secret to view.").
						Value(&selectedSecretFile),
				),
			)

			err = form.
				WithTheme(app.ThemeMinno()).
				Run()
			if err != nil {
				return err
			}

			data, err := os.ReadFile(selectedSecretFile.Path)
			if err != nil {
				return errors.New("failed to read secret. Encrypted secret may be corrupted")
			}

			fmt.Println(pp.Complete("Secret read"))

			secret, err := tomb.Decrypt(data)
			data = nil
			if err != nil {
				return errors.New("failed to decrypt secret. Encrypted secret may be corrupted")
			}

			fmt.Println(pp.Complete("Secret decrypted"))
			fmt.Println()
			fmt.Println(pp.Info("The secret is " + pp.Green(string(secret))))
		} else {
			if err := secrets.ValidateName(secretName); err != nil {
				return fmt.Errorf("%s\n\nThe secret name provided was %s", err, pp.Red(secretName))
			}

			secretPtr := secrets.FindSecretByName(secretName, secretFiles)
			if secretPtr == nil {
				return fmt.Errorf("secret %s does not exist!\n\nUse command %s to create the secret", pp.Red(secretName), pp.Greenf("%s create", envVars.ExeCmd))
			}

			selectedSecretFile = *secretPtr

			data, err := os.ReadFile(selectedSecretFile.Path)
			if err != nil {
				return errors.New("failed to read secret. Encrypted secret may be corrupted")
			}

			secret, err := tomb.Decrypt(data)
			data = nil
			if err != nil {
				return errors.New("failed to decrypt secret. Encrypted secret may be corrupted")
			}

			if output == "" {
				fmt.Print(string(secret))
			} else {
				outputDir := filepath.Dir(output)
				if _, err := os.Stat(outputDir); os.IsNotExist(err) {
					err = os.MkdirAll(outputDir, dirMode)
					if err != nil {
						return errors.New("failed to create output directory for output file")
					}
				}

				err = os.WriteFile(output, secret, secretMode)
				if err != nil {
					return errors.New("failed to write secret to output file")
				}
			}
			secret = nil
		}

		return nil
	},
}
