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
	"github.com/engmtcdrm/minno/credentials"
	"github.com/engmtcdrm/minno/env"
	"github.com/engmtcdrm/minno/header"
)

func init() {
	createCmd.Flags().StringVarP(&credName, "cred-name", "n", "", "(optional) The credential to create")
	createCmd.Flags().StringVarP(&rawCredFile, "file", "f", "", "(optional) The file containing the plain text credential to encrypt")

	rootCmd.AddCommand(createCmd)
}

func validateCredName(name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}

	if !credentials.IsValidName(name) {
		return errors.New("name can only be alphanumeric, hyphens, and underscores")
	}

	for _, f := range credFiles {
		if f.Name == name {
			return errors.New("credential with that name already exists")
		}
	}

	return nil
}

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a credential",
	Long:    "Create a credential",
	Example: app.Name + " create",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Enforce that --cred-name and --file must be provided together
		if (credName != "" && rawCredFile == "") || (credName == "" && rawCredFile != "") {
			return errors.New("flags -n/--cred-name and -f/--file must be provided together")
		}

		envVars, err := env.GetEnv()
		if err != nil {
			return err
		}

		tomb, err := entomb.NewTomb(filepath.Join(envVars.AppHomeDir, ".key"))
		if err != nil {
			return err
		}

		var cred string
		var credFile string
		var encTest []byte

		// Interactive mode if no flags are provided
		if credName == "" && rawCredFile == "" {
			header.PrintHeader()

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
						Validate(validateCredName).
						Inline(true),
				),
			)

			err = form.
				WithTheme(app.ThemeMinno()).
				Run()
			if err != nil {
				return err
			}

			encTest, err = tomb.Encrypt([]byte(strings.TrimSpace(cred)))
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
		}

		if err := validateCredName(credName); err != nil {
			return err
		}

		rawCredFile, err := env.ExpandTilde(strings.TrimSpace(rawCredFile))
		if err != nil {
			return err
		}

		credBytes, err := os.ReadFile(rawCredFile)
		if err != nil {
			return fmt.Errorf("could not read file '%s': %w", rawCredFile, err)
		}

		encTest, err = tomb.Encrypt(credBytes)
		credBytes = nil
		if err != nil {
			return err
		}

		credFile = credName

		credFilePath := filepath.Join(envVars.AppHomeDir, credFile+".cred")

		if err = os.WriteFile(credFilePath, encTest, 0600); err != nil {
			return err
		}

		return nil
	},
}
