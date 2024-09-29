package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/engmtcdrm/minno/app"
	"github.com/engmtcdrm/minno/credentials"
	"github.com/engmtcdrm/minno/credentials/prompts"
	"github.com/engmtcdrm/minno/encrypt"
	"github.com/engmtcdrm/minno/env"
	"github.com/engmtcdrm/minno/header"
	pp "github.com/engmtcdrm/minno/utils/prettyprint"
)

func init() {
	updateCmd.Flags().StringVarP(&credName, "cred-name", "n", "", "The credential to update (optional)")

	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update a credential",
	Long:    "Update a credential",
	Example: app.Name + " update",
	RunE: func(cmd *cobra.Command, args []string) error {
		envVars, err := env.GetEnv()
		if err != nil {
			return err
		}

		tomb, err := encrypt.NewTomb(filepath.Join(envVars.AppHomeDir, ".key"))
		if err != nil {
			return err
		}

		header.PrintBanner()

		var cred string
		var selectedCred credentials.Credential
		var form *huh.Form

		if credName == "" {
			options, err := prompts.GetCredOptions()
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
			WithTheme(pp.ThemeMinno()).
			Run()
		if err != nil {
			if err.Error() == "user aborted" {
				fmt.Println("User aborted")
				os.Exit(0)
			} else {
				return err
			}
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
		fmt.Printf("You can run the commmand %s to view the unencrypted credential\n", pp.Greenf("%s view -n %s", app.Name, selectedCred.Name))

		return nil
	},
}
