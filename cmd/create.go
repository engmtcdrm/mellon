package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/engmtcdrm/minno/credentials"
	"github.com/engmtcdrm/minno/encrypt"
	"github.com/engmtcdrm/minno/env"
	"github.com/engmtcdrm/minno/utils/header"
	pp "github.com/engmtcdrm/minno/utils/prettyprint"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a credential",
	Long:    "Create a credential",
	Example: env.AppNm + " create",
	RunE: func(cmd *cobra.Command, args []string) error {
		header.PrintBanner()

		tomb, err := encrypt.NewTomb(filepath.Join(env.AppHomeDir, ".key"))
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
					Title("Enter a filename for the credential").
					Value(&credFile).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("filename cannot be empty")
						}

						if strings.HasSuffix(s, ".cred") {
							return fmt.Errorf("filename cannot end with .cred")
						}

						var re = regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
						if !re.MatchString(s) {
							return fmt.Errorf("filename must be alphanumeric and can contain hyphens and underscores")
						}

						for _, f := range credFiles {
							if f.Name == s {
								return fmt.Errorf("credential file with that name already exists")
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

		credFilePath := filepath.Join(env.AppHomeDir, credFile+".cred")

		if err = os.WriteFile(credFilePath, encTest, 0600); err != nil {
			return err
		}

		fmt.Println()
		fmt.Println(pp.Complete("Credential saved"))
		fmt.Println()
		fmt.Printf("Please run the commmand %s to view the unencrypted credential\n", color.GreenString(env.AppNm+" view -f "+credFile))

		return nil
	},
}
