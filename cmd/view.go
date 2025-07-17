package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/engmtcdrm/go-pardon"
	pp "github.com/engmtcdrm/go-prettyprint"
	"github.com/engmtcdrm/minno/app"
	"github.com/engmtcdrm/minno/header"
	"github.com/engmtcdrm/minno/secrets"
	"github.com/engmtcdrm/minno/secrets/prompts"
)

func init() {
	viewCmd.Flags().StringVarP(
		&secretName,
		"secret",
		"s",
		"",
		"(optional) The name of the secret to view. Only names containing alphanumeric, hyphens, and underscores are allowed",
	)
	viewCmd.Flags().StringVarP(
		&output,
		"output",
		"o",
		"",
		"(optional) File to write decrypted secret to. Defaults to outputting to stdout. This only works with the option -s/--secret",
	)

	rootCmd.AddCommand(viewCmd)
}

func validateViewFlags(cmd *cobra.Command, args []string) error {
	if output != "" && secretName == "" {
		return errors.New("flag -o/--output can only be used when -s/--secret is provided")
	}

	return nil
}

var viewCmd = &cobra.Command{
	Use:     "view",
	Short:   "View a secret",
	Long:    "View a secret",
	Example: fmt.Sprintf("  %s view\n  %s view -s awesome-secret", app.Name, app.Name),
	PreRunE: validateViewFlags,
	RunE: func(cmd *cobra.Command, args []string) error {
		var selectedSecretFile secrets.Secret

		if secretName == "" {
			header.PrintHeader()

			options, err := prompts.GetSecretOptions(secretFiles, "view")
			if err != nil {
				return err
			}

			promptSelect := pardon.NewSelect[secrets.Secret]().
				Options(options...).
				Title("What secret do you want to view?").
				Value(&selectedSecretFile).
				SelectFunc(
					func(s string) string {
						return pp.Yellow(s)
					}).
				Icon(pp.Cyan(pardon.Icons.QuestionMark))
			if err := promptSelect.Ask(); err != nil {
				return err
			}

			secret, err := selectedSecretFile.Decrypt()
			if err != nil {
				return errors.New("failed to decrypt secret. Encrypted secret may be corrupted")
			}

			fmt.Println(pp.Complete("Secret decrypted"))
			fmt.Println()
			fmt.Println(pp.Info("The secret is " + pp.Green(string(secret))))

			return nil
		}

		secretPtr := secrets.FindSecretByName(secretName, secretFiles)
		if secretPtr == nil {
			return fmt.Errorf("failed to read secret '%s': secret does not exist", secretName)
		}

		selectedSecretFile = *secretPtr

		secret, err := selectedSecretFile.Decrypt()
		if err != nil {
			return fmt.Errorf("failed to decrypt secret '%s'. Encrypted secret may be corrupted", secretName)
		}

		if output == "" {
			fmt.Print(string(secret))
		} else {
			outputDir := filepath.Dir(output)
			if _, err := os.Stat(outputDir); os.IsNotExist(err) {
				err = os.MkdirAll(outputDir, dirMode)
				if err != nil {
					return fmt.Errorf("failed to create output directory for output file '%s'", output)
				}
			}

			err = os.WriteFile(output, secret, secretMode)
			if err != nil {
				return fmt.Errorf("failed to write secret to output file '%s'", output)
			}
		}
		secret = nil

		return nil
	},
}
