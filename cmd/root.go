package cmd

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/engmtcdrm/minno/app"
	"github.com/engmtcdrm/minno/env"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	debug bool

	rootCmd = &cobra.Command{
		Use:     app.Name,
		Short:   app.ShortDesc,
		Long:    app.LongDesc,
		Example: app.Name,
		Version: app.Version,
	}

	credName string
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.ExecuteContext(context.Background())
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	cobra.OnInitialize(configInit)

	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug mode")
}

func configInit() {
	envVars, err := env.GetEnv()
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(envVars.AppHomeDir); os.IsNotExist(err) {
		// Directory does not exist, create it
		err = os.MkdirAll(envVars.AppHomeDir, 0700)
		if err != nil {
			panic(err)
		}

		// Change permission again to get rid of any sticky bits
		err = os.Chmod(envVars.AppHomeDir, 0700)
		if err != nil {
			panic(err)
		}
	} else {
		// Directory exists, make sure directories and files are secure
		err = filepath.Walk(envVars.AppHomeDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return os.Chmod(path, 0700)
			}
			return os.Chmod(path, 0600)
		})
		if err != nil {
			panic(err)
		}
	}

	if debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug(color.CyanString("Debug mode enabled"))
	}
}
