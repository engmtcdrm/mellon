package cmd

import (
	"fmt"

	pp "github.com/engmtcdrm/go-prettyprint"
	"github.com/engmtcdrm/minno/app"
	"github.com/engmtcdrm/minno/header"
	"github.com/spf13/cobra"
)

func init() {
	listCmd.Flags().BoolVarP(
		&print,
		"print",
		"p",
		false,
		"(optional) Whether to print only the names of the secrets without additional information",
	)

	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List available secrets",
	Long:    "List available secrets",
	Example: fmt.Sprintf("  %s list", app.Name),
	RunE: func(cmd *cobra.Command, args []string) error {
		if print {
			for _, secret := range secretFiles {
				fmt.Println(secret.Name)
			}

			return nil
		}

		header.PrintHeader()

		if len(secretFiles) == 0 {
			return fmt.Errorf("no available secrets to list\n\nUse command %s to create a secret", pp.Greenf("%s create", envVars.ExeCmd))
		}

		fmt.Println(pp.Info("Available secrets"))
		fmt.Println()
		for _, secret := range secretFiles {
			fmt.Printf("  - %s\n", pp.Green(secret.Name))
		}

		return nil
	},
}
