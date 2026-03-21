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

type cmd struct {
	secretFiles []secrets.Secret // List of available secrets.
	print       bool             // Whether to print only the names of the secrets.
}

func NewCommand(secretFiles []secrets.Secret) *cobra.Command {
	if secretFiles == nil {
		secretFiles = []secrets.Secret{}
	}

	c := &cmd{secretFiles: secretFiles}

	listCmd := &cobra.Command{
		Use:     "list",
		Short:   "List available secrets",
		Long:    "List available secrets",
		Example: fmt.Sprintf("  %s list", app.Name),
		RunE:    c.runE,
	}

	listCmd.Flags().BoolVarP(
		&c.print,
		"print",
		"p",
		false,
		"(optional) Whether to print only the names of the secrets",
	)

	return listCmd
}

// listSecrets lists all available secrets.
func (c *cmd) listSecrets() error {
	header.PrintHeader()

	if len(c.secretFiles) == 0 {
		fmt.Printf("No available secrets to list\n\nUse command %s to create a secret", pp.Greenf("%s create", env.Instance.ExeCmd()))
		return nil
	}

	fmt.Println(pp.Info("Available secrets"))
	fmt.Println()

	for _, secret := range c.secretFiles {
		fmt.Printf("  - %s\n", pp.Green(secret.Name()))
	}

	return nil
}

// printSecrets prints only the names of the secrets.
func (c *cmd) printSecrets() error {
	for _, secret := range c.secretFiles {
		fmt.Println(secret.Name())
	}

	return nil
}

// runE is the main execution function for the command.
func (c *cmd) runE(cmd *cobra.Command, args []string) error {
	if c.print {
		return c.printSecrets()
	}

	return c.listSecrets()
}
