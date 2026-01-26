package listcmd

import (
	"fmt"

	pp "github.com/engmtcdrm/go-prettyprint"
	"github.com/spf13/cobra"

	"github.com/engmtcdrm/mellon/internal/app"
	"github.com/engmtcdrm/mellon/internal/env"
	"github.com/engmtcdrm/mellon/internal/header"
	"github.com/engmtcdrm/mellon/secrets"
)

var (
	secretFiles []secrets.Secret // List of available secrets
	print       bool             // Whether to print only the names of the secrets without additional information
)

func NewCommand(secretFilesList []secrets.Secret) *cobra.Command {
	secretFiles = secretFilesList

	listCmd := &cobra.Command{
		Use:     "list",
		Short:   "List available secrets",
		Long:    "List available secrets",
		Example: fmt.Sprintf("  %s list", app.Name),
		RunE: func(cmd *cobra.Command, args []string) error {
			if print {
				return printSecrets()
			}

			return listSecrets()
		},
	}

	listCmd.Flags().BoolVarP(
		&print,
		"print",
		"p",
		false,
		"(optional) Whether to print only the names of the secrets without additional information",
	)

	return listCmd
}

func printSecrets() error {
	for _, secret := range secretFiles {
		fmt.Println(secret.Name())
	}
	return nil
}

func listSecrets() error {
	header.PrintHeader()

	if len(secretFiles) == 0 {
		fmt.Printf("No available secrets to list\n\nUse command %s to create a secret", pp.Greenf("%s create", env.Instance.ExeCmd()))
		return nil
	}

	fmt.Println(pp.Info("Available secrets"))
	fmt.Println()

	for _, secret := range secretFiles {
		fmt.Printf("  - %s\n", pp.Green(secret.Name()))
	}

	return nil
}
