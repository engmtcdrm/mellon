package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/engmtcdrm/minno/app"
	"github.com/engmtcdrm/minno/credentials"
	"github.com/engmtcdrm/minno/credentials/prompts"
	"github.com/engmtcdrm/minno/encrypt"
	"github.com/engmtcdrm/minno/env"
	"github.com/engmtcdrm/minno/utils/header"
	pp "github.com/engmtcdrm/minno/utils/prettyprint"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
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
		tomb, err := encrypt.NewTomb(filepath.Join(env.AppHomeDir, ".key"))
		if err != nil {
			return err
		}

		header.PrintBanner()

		var cred string
		var selectedCredFile credentials.Credential
		var form *huh.Form

		if credName == "" {
			options, err := prompts.GetCredOptions(env.AppHomeDir)
			if err != nil {
				return err
			}

			form = huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[credentials.Credential]().
						Options(options...).
						Title("Available Credentials").
						Description("Choose a credential to update.").
						Value(&selectedCredFile),
					huh.NewInput().
						Title("Enter the updated credential").
						Value(&cred).
						EchoMode(huh.EchoModeNone).
						Inline(true),
				),
			)
		} else {
			selectedCredFile.Name = filepath.Base(credName)
			selectedCredFile.Path = credentials.ResolveCredName(credName)

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

		if err = os.WriteFile(selectedCredFile.Path, encTest, 0600); err != nil {
			return err
		}

		fmt.Println()
		fmt.Println(pp.Complete("Credential saved"))
		fmt.Println()
		fmt.Printf("Please run the commmand %s to view the unencrypted credential\n", color.GreenString(app.Name+" view -f "+selectedCredFile.Name))

		return nil
	},
}
