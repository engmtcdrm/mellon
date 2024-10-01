package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/engmtcdrm/minno/app"
	"github.com/engmtcdrm/minno/credentials"
	"github.com/engmtcdrm/minno/encrypt"
	"github.com/engmtcdrm/minno/env"
	"github.com/engmtcdrm/minno/header"
	pp "github.com/engmtcdrm/minno/utils/prettyprint"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a credential",
	Long:    "Create a credential",
	Example: app.Name + " create",
	RunE: func(cmd *cobra.Command, args []string) error {
		header.PrintHeader()

		envVars, err := env.GetEnv()
		if err != nil {
			return err
		}

		tomb, err := encrypt.NewTomb(filepath.Join(envVars.AppHomeDir, ".key"))
		if err != nil {
			return err
		}

		var cred string
		var credFile string

		credFiles, err := credentials.GetCredFiles()
		if err != nil {
			return err
		}

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Enter a credential to secure").
					Value(&cred).
					EchoMode(huh.EchoModeNone).
					Inline(true),
				huh.NewInput().
					Title("Enter a name for the credential").
					Value(&credFile).
					Validate(func(s string) error {
						if s == "" {
							return errors.New("name cannot be empty")
						}

						if !credentials.IsValidName(s) {
							return errors.New("name can only be alphanumeric, hyphens, and underscores")
						}

						for _, f := range credFiles {
							if f.Name == s {
								return errors.New("credential with that name already exists")
							}
						}

						return nil
					}).
					Inline(true),
			),
		)

		err = form.
			WithTheme(pp.ThemeMinno()).
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

		credFilePath := filepath.Join(envVars.AppHomeDir, credFile+".cred")

		if err = os.WriteFile(credFilePath, encTest, 0600); err != nil {
			return err
		}

		fmt.Println(pp.Complete("Credential saved"))
		fmt.Println()
		fmt.Printf("You can run the commmand %s to view the unencrypted credential\n", pp.Greenf("%s view -n %s", envVars.ExeCmd, credFile))

		return nil
	},
}
