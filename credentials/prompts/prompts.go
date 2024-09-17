package prompts

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/engmtcdrm/minno/credentials"
	"github.com/engmtcdrm/minno/env"
	"github.com/fatih/color"
)

func GetCredOptions(appHomeDir string) ([]huh.Option[credentials.Credential], error) {
	credFiles, err := credentials.GetCredFiles()
	if err != nil {
		return nil, err
	}

	if len(credFiles) == 0 {
		fmt.Println("No credentials found")
		fmt.Println()
		fmt.Printf("Please run command %s to create a credential\n", color.GreenString(env.AppNm+" create"))
		os.Exit(0)
	}

	options := []huh.Option[credentials.Credential]{}

	for _, c := range credFiles {
		options = append(options, huh.NewOption(c.Name, c))
	}

	return options, nil
}
