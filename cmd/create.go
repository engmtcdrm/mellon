package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/engmtcdrm/minno/utils/encrypt"
	"github.com/engmtcdrm/minno/utils/env"
	pp "github.com/engmtcdrm/minno/utils/prettyprint"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	appHomeDir string
	tomb       *encrypt.Tomb
)

func initApp() error {
	var err error
	appHomeDir, err = env.AppHomeDir()

	if err != nil {
		return err
	}

	tomb, err = encrypt.NewTomb(filepath.Join(appHomeDir, ".key"))

	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a set of credentials",
	Long:    "Create a set of credentials",
	Example: env.AppNm + " create",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := initApp(); err != nil {
			return err
		}

		var cred string
		var credFile string
		var credFiles []string

		filepath.WalkDir(appHomeDir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() && strings.HasSuffix(path, ".cred") {
				path = filepath.Base(path)
				path = strings.Replace(path, ".cred", "", 1)
				credFiles = append(credFiles, path)
			}

			return nil
		})

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
							return fmt.Errorf("Filename cannot be empty")
						}

						if strings.HasSuffix(s, ".cred") {
							return fmt.Errorf("Filename cannot end with .cred")
						}

						var re = regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
						if !re.MatchString(s) {
							return fmt.Errorf("Filename must be alphanumeric and can contain hyphens and underscores")
						}

						for _, f := range credFiles {
							if f == s {
								return fmt.Errorf("Credential file with that name already exists")
							}
						}

						return nil
					}).
					Inline(true),
			),
		)

		err := form.
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

		credFilePath := filepath.Join(appHomeDir, credFile+".cred")

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
