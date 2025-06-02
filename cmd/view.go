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
	"github.com/engmtcdrm/minno/credentials"
	"github.com/engmtcdrm/minno/credentials/prompts"
	"github.com/engmtcdrm/minno/header"
)

func init() {
	viewCmd.Flags().StringVarP(
		&credName,
		"cred-name",
		"n",
		"",
		"(optional) The name of the credential to view. Only names containing alphanumeric, hyphens, and underscores are allowed.",
	)
	viewCmd.Flags().StringVarP(
		&output,
		"output",
		"o",
		"",
		"(optional) File to write decrypted credential to. Defaults to outputting to stdout. This only works with the option -n, --cred-name",
	)

	rootCmd.AddCommand(viewCmd)
}

var viewCmd = &cobra.Command{
	Use:     "view",
	Short:   "View a credential",
	Long:    "View a credential",
	Example: fmt.Sprintf("  %s view\n  %s view -n awesome-cred", app.Name, app.Name),
	RunE: func(cmd *cobra.Command, args []string) error {
		tomb, err := entomb.NewTomb(envVars.KeyPath)
		if err != nil {
			return err
		}

		if credName == "" {
			header.PrintHeader()

			var selectedCredFile credentials.Credential

			options, err := prompts.GetCredOptions(credFiles, "view")
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
				WithTheme(app.ThemeMinno()).
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

			if output == "" {
				fmt.Print(string(cred))
			} else {
				outputDir := filepath.Dir(output)
				if _, err := os.Stat(outputDir); os.IsNotExist(err) {
					err = os.MkdirAll(outputDir, 0700)
					if err != nil {
						return errors.New("failed to create output directory for output file")
					}
				}

				err = os.WriteFile(output, cred, 0600)
				if err != nil {
					return errors.New("failed to write credential to output file")
				}
			}
			cred = nil
		}

		return nil
	},
}
