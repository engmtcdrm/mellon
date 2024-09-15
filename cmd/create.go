package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/engmtcdrm/minno/utils/encrypt"
	"github.com/engmtcdrm/minno/utils/env"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a set of credentials",
	Long:    "Create a set of credentials",
	Example: env.AppNm + " create",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Welcome to create command")

		appHomeDir, err := env.AppHomeDir()

		if err != nil {
			return err
		}

		tomb, err := encrypt.NewTomb(filepath.Join(appHomeDir, ".key"))
		if err != nil {
			return err
		}

		fmt.Print("Enter a credential: ")
		reader := bufio.NewReader(os.Stdin)
		cred, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		// Remove any trailing newline or spaces
		cred = strings.TrimSpace(cred)

		fmt.Printf("Value to encode: \"%s\"\n", cred)

		encTest, err := tomb.Encrypt([]byte(cred))
		if err != nil {
			return err
		}

		decTest, err := tomb.Decrypt(encTest)
		if err != nil {
			return err
		}

		fmt.Printf("Decoded value: \"%s\"\n", string(decTest))

		return nil
	},
}
