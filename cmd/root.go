package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/engmtcdrm/mellon/internal/app"
	"github.com/engmtcdrm/mellon/internal/cli/createcmd"
	"github.com/engmtcdrm/mellon/internal/cli/deletecmd"
	"github.com/engmtcdrm/mellon/internal/cli/listcmd"
	"github.com/engmtcdrm/mellon/internal/cli/updatecmd"
	"github.com/engmtcdrm/mellon/internal/cli/viewcmd"
	"github.com/engmtcdrm/mellon/internal/constants"
	"github.com/engmtcdrm/mellon/internal/env"
	"github.com/engmtcdrm/mellon/secrets"
)

var (
	rootCmd = &cobra.Command{
		Use:     app.Name,
		Short:   app.ShortDesc,
		Long:    app.LongDesc,
		Example: app.Name,
		Version: getSemVer(app.Version),
	}
)

func init() {
	env.Init()

	initShellCompletion(env.Instance.Home())
	mkdir(env.Instance.AppHomeDir(), constants.SecureDirMode)
	mkdir(env.Instance.SecretsPath(), constants.SecureDirMode)
	secureFiles(env.Instance.AppHomeDir(), constants.SecureDirMode, constants.SecureFileMode)

	secretFiles, err := secrets.GetSecretFiles(
		env.Instance.KeyPath(),
		env.Instance.SecretsPath(),
		env.Instance.SecretExt(),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rootCmd.SilenceUsage = true
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.AddCommand(createcmd.NewCommand(secretFiles))
	rootCmd.AddCommand(deletecmd.NewCommand(secretFiles))
	rootCmd.AddCommand(listcmd.NewCommand(secretFiles))
	rootCmd.AddCommand(updatecmd.NewCommand(secretFiles))
	rootCmd.AddCommand(viewcmd.NewCommand(secretFiles))
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
