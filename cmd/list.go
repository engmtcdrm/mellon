package cmd

import (
	"fmt"

	"github.com/engmtcdrm/minno/app"
	"github.com/engmtcdrm/minno/header"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

// func validateListFlags(cmd *cobra.Command, args []string) error {
// 	return nil
// }

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List available secrets",
	Long:    "List available secrets",
	Example: fmt.Sprintf("  %s list", app.Name),
	// PreRunE: validateListFlags,
	RunE: func(cmd *cobra.Command, args []string) error {
		header.PrintHeader()

		if len(secretFiles) == 0 {
			fmt.Println("No secrets found.")

			return nil
		}

		fmt.Printf("%d available secret(s):\n\n", len(secretFiles))

		for _, secret := range secretFiles {
			fmt.Printf("- %s\n", secret.Name)
		}

		return nil
	},
}
