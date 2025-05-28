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
	"github.com/engmtcdrm/minno/credentials"
	"github.com/engmtcdrm/minno/credentials/prompts"
	"github.com/engmtcdrm/minno/header"
)

func init() {
	updateCmd.Flags().StringVarP(
		&credName,
		"cred-name",
		"n",
		"",
		"(optional) The credential to update",
	)

	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update a credential",
	Long:    "Update a credential",
	Example: fmt.Sprintf("  %s update", app.Name),
	RunE: func(cmd *cobra.Command, args []string) error {
		tomb, err := entomb.NewTomb(filepath.Join(envVars.AppHomeDir, ".key"))
		if err != nil {
			return err
		}

		header.PrintHeader()

		var cred string
		var selectedCred credentials.Credential
		var form *huh.Form

		if credName == "" {
			options, err := prompts.GetCredOptions(credFiles)
			if err != nil {
				return err
			}

			form = huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[credentials.Credential]().
						Options(options...).
						Title("Available Credentials").
						Description("Choose a credential to update.").
						Value(&selectedCred),
					huh.NewInput().
						Title("Enter the updated credential").
						Value(&cred).
						EchoMode(huh.EchoModeNone).
						Inline(true),
				),
			)
		} else {
			selectedCred.Name = filepath.Base(credName)
			selectedCred.Path, err = credentials.ResolveCredName(credName)
			if err != nil {
				return err
			}

			if !credentials.IsExists(selectedCred.Path) {
				return fmt.Errorf("credential %s does not exist!\n\nUse command %s to create the credential", pp.Red(selectedCred.Name), pp.Greenf("%s create", envVars.ExeCmd))
			}

			form = huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Enter the updated credential").
						Value(&cred).
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

		encTest, err := tomb.Encrypt([]byte(strings.TrimSpace(cred)))
		cred = ""
		if err != nil {
			return err
		}

		fmt.Println(pp.Complete("Credential encrypted"))

		if err = os.WriteFile(selectedCred.Path, encTest, 0600); err != nil {
			return err
		}
		fmt.Println(pp.Complete("Credential saved"))
		fmt.Println()
		fmt.Printf("You can run the commmand %s to view the unencrypted credential\n", pp.Greenf("%s view -n %s", envVars.ExeCmd, selectedCred.Name))

		return nil
	},
}
