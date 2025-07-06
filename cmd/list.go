package cmd

import (
	"fmt"

	"github.com/engmtcdrm/minno/app"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List available secrets",
	Long:    "List available secrets",
	Example: fmt.Sprintf("  %s list", app.Name),
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, secret := range secretFiles {
			fmt.Println(secret.Name)
		}

		return nil
	},
}
