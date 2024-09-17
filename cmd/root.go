package cmd

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/engmtcdrm/minno/env"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	debug bool

	rootCmd = &cobra.Command{
		Use:     env.AppNm,
		Short:   env.AppShortDesc,
		Long:    env.AppLongDesc,
		Example: env.AppNm,
		Version: env.AppVersion,
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
	err := env.SetAppHomeDir()
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(env.AppHomeDir); os.IsNotExist(err) {
		// Directory does not exist, create it
		err = os.MkdirAll(env.AppHomeDir, 0700)
		if err != nil {
			panic(err)
		}

		// Change permission again to get rid of any sticky bits
		err = os.Chmod(env.AppHomeDir, 0700)
		if err != nil {
			panic(err)
		}
	} else {
		// Directory exists, make sure directories and files are secure
		err = filepath.Walk(env.AppHomeDir, func(path string, info os.FileInfo, err error) error {
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
