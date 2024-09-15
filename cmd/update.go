package cmd

import (
	"fmt"

	"github.com/engmtcdrm/minno/utils/env"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update a credential",
	Long:    "Update a credential",
	Example: env.AppNm + " update",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Welcome to update command")

		return nil
	},
}
