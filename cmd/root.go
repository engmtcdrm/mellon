package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/engmtcdrm/mellon/app"
	"github.com/engmtcdrm/mellon/env"
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

	secretName  string // The name of the secret to create/view/update/delete
	secretFile  string // The file containing the plain text secret to encrypt
	cleanupFile bool   // Whether to delete the raw secret file after encryption
	forceDelete bool   // Whether to force overwrite an existing secret file (only used with delete command)
	deleteAll   bool   // Whether to delete all secrets (only used with delete command)
	output      string // The file to write decrypted secret to (only used with view command)
	print       bool   // Whether to print only the names of the secrets without additional information (only used with list command)

	secretFiles []secrets.Secret // List of secrets available in the app
	envVars     *env.Env         // Environment variables for the app

	// Modes for files and directories
	dirMode    os.FileMode = 0700 // Default directory mode for app home directory as well as output of secret directories
	secretMode os.FileMode = 0600 // Default file mode for secret files
)

// Execute executes the root command.
func Execute() error {
	rootCmd.SilenceUsage = true
	return rootCmd.ExecuteContext(context.Background())
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	cobra.OnInitialize(configInit)
}

func configInit() {
	var err error

	envVars, err = env.GetEnv()
	if err != nil {
		panic(err)
	}

	mkdir(envVars.AppHomeDir, dirMode)
	mkdir(envVars.SecretsPath, dirMode)
	secureFiles(envVars.AppHomeDir, dirMode, secretMode)

	secretFiles, err = secrets.GetSecretFiles()
	if err != nil {
		panic(err)
	}
}
