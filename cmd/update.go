package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	updateCmd.Flags().StringVarP(
		&secretName,
		"secret",
		"s",
		"",
		"(optional) The secret to update",
	)

	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update a secret",
	Long:    "Update a secret",
	Example: fmt.Sprintf("  %s update", app.Name),
	RunE: func(cmd *cobra.Command, args []string) error {
		tomb, err := entomb.NewTomb(filepath.Join(envVars.AppHomeDir, ".key"))
		if err != nil {
			return err
		}

		header.PrintHeader()

		var secret string
		var selectedSecret secrets.Secret
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
			selectedSecret.Name = filepath.Base(secretName)
			selectedSecret.Path, err = secrets.ResolveSecretName(secretName)
			if err != nil {
				return err
			}

			if !secrets.IsExists(selectedSecret.Path) {
				return fmt.Errorf("secret %s does not exist!\n\nUse command %s to create the secret", pp.Red(selectedSecret.Name), pp.Greenf("%s create", envVars.ExeCmd))
			}

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

		encTest, err := tomb.Encrypt([]byte(strings.TrimSpace(secret)))
		secret = ""
		if err != nil {
			return err
		}

		fmt.Println(pp.Complete("Secret encrypted"))

		if err = os.WriteFile(selectedSecret.Path, encTest, 0600); err != nil {
			return err
		}
		fmt.Println(pp.Complete("Secret saved"))
		fmt.Println()
		fmt.Printf("You can run the commmand %s to view the unencrypted secret\n", pp.Greenf("%s view -n %s", envVars.ExeCmd, selectedSecret.Name))

		return nil
	},
}
