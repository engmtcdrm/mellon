package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/engmtcdrm/minno/app"
	"github.com/engmtcdrm/minno/credentials"
	"github.com/engmtcdrm/minno/credentials/prompts"
	"github.com/engmtcdrm/minno/encrypt"
	"github.com/engmtcdrm/minno/env"
	"github.com/engmtcdrm/minno/header"
	pp "github.com/engmtcdrm/minno/utils/prettyprint"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	viewCmd.Flags().StringVarP(&credName, "cred-name", "n", "", "The name of the credential to view (optional)")

	rootCmd.AddCommand(viewCmd)
}

var viewCmd = &cobra.Command{
	Use:     "view",
	Short:   "View a credential",
	Long:    "View a credential",
	Example: app.Name + " view",
	RunE: func(cmd *cobra.Command, args []string) error {
		envVars, err := env.GetEnv()
		if err != nil {
			return err
		}

		tomb, err := encrypt.NewTomb(filepath.Join(envVars.AppHomeDir, ".key"))
		if err != nil {
			return err
		}

		if credName == "" {
			header.PrintBanner()

			var selectedCredFile credentials.Credential

			options, err := prompts.GetCredOptions(envVars.AppHomeDir)
			if err != nil {
				return err
			}

			form := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[credentials.Credential]().
						Options(options...).
						Title("Available Credentials").
						Description("Choose a credential to update.").
						Value(&selectedCredFile),
				),
			)

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

			data, err := os.ReadFile(selectedCredFile.Path)
			if err != nil {
				return err
			}

			fmt.Println(pp.Complete("Credential read"))

			cred, err := tomb.Decrypt(data)
			data = nil
			if err != nil {
				fmt.Println(pp.Fail("Failed to decrypt credential. Encrypted credential may be corrupted"))
				os.Exit(99)
			}

			fmt.Println(pp.Complete("Credential decrypted"))
			fmt.Println()
			fmt.Println(pp.Info("The credential is " + color.CyanString(string(cred))))
		} else {
			credName, err = credentials.ResolveCredName(credName)
			if err != nil {
				return err
			}

			data, err := os.ReadFile(credName)
			if err != nil {
				return err
			}

			cred, err := tomb.Decrypt(data)
			data = nil
			if err != nil {
				return err
			}

			fmt.Print(string(cred))
			cred = nil
		}

		return nil
	},
}
