package cmd

import (
	"fmt"

	"github.com/engmtcdrm/minno/utils/env"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(viewCmd)
}

var viewCmd = &cobra.Command{
	Use:     "view",
	Short:   "View a credential",
	Long:    "View a credential",
	Example: env.AppNm + " view",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Welcome to view command")

		return nil
	},
}

// decTest, err := tomb.Decrypt(encTest)
// 		if err != nil {
// 			return err
// 		}

// 		fmt.Printf("Decoded value: \"%s\"\n", string(decTest))

// 		return nil
