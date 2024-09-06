package cmd

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/engmtcdrm/minno/utils/env"
	"github.com/engmtcdrm/minno/utils/header"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	debug bool

	rootCmd = &cobra.Command{
		Use:     env.AppNm,
		Short:   env.AppShortDesc,
		Long:    env.AppLongDesc,
		Version: env.AppVersion,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	cobra.OnInitialize(header.PrintBanner)
	cobra.OnInitialize(configInit)

	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug mode")
}

func configInit() {
	appDir, err := env.AppHomeDir()
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		// Directory does not exist, create it
		err = os.MkdirAll(appDir, 0700)
		if err != nil {
			panic(err)
		}

		// Change permission again to get rid of any sticky bits
		err = os.Chmod(appDir, 0700)
		if err != nil {
			panic(err)
		}
	} else {
		// Directory exists, make sure directories and files are secure
		err = filepath.Walk(appDir, func(path string, info os.FileInfo, err error) error {
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
