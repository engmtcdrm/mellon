package cmd

import (
	"errors"
	"fmt"
	"os"

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
	viewCmd.Flags().StringVarP(&credName, "cred-name", "n", "", "The name of the credential to view. Only names containing alphanumeric, hyphens, and underscores are allowed. (optional)")

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

		tomb, err := encrypt.NewTomb(envVars.KeyPath)
		if err != nil {
			return err
		}

		if credName == "" {
			header.PrintHeader()

			var selectedCredFile credentials.Credential

			options, err := prompts.GetCredOptions()
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
				return err
			}

			data, err := os.ReadFile(selectedCredFile.Path)
			if err != nil {
				return errors.New("failed to read credential. Encrypted credential may be corrupted")
			}

			fmt.Println(pp.Complete("Credential read"))

			cred, err := tomb.Decrypt(data)
			data = nil
			if err != nil {
				return errors.New("failed to decrypt credential. Encrypted credential may be corrupted")
			}

			fmt.Println(pp.Complete("Credential decrypted"))
			fmt.Println()
			fmt.Println(pp.Info("The credential is " + pp.Green(string(cred))))
		} else {
			credName, err = credentials.ResolveCredName(credName)
			if err != nil {
				return err
			}

			if !credentials.IsExists(credName) {
				return errors.New("credential name does not exist")
			}

			data, err := os.ReadFile(credName)
			if err != nil {
				return errors.New("failed to read credential. Encrypted credential may be corrupted")
			}

			cred, err := tomb.Decrypt(data)
			data = nil
			if err != nil {
				return errors.New("failed to decrypt credential. Encrypted credential may be corrupted")
			}

			fmt.Print(string(cred))
			cred = nil
		}

		return nil
	},
}
