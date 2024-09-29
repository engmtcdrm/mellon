package prompts

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"

	"github.com/engmtcdrm/minno/app"
	"github.com/engmtcdrm/minno/credentials"
	pp "github.com/engmtcdrm/minno/utils/prettyprint"
)

func GetCredOptions() ([]huh.Option[credentials.Credential], error) {
	credFiles, err := credentials.GetCredFiles()
	if err != nil {
		return nil, err
	}

	if len(credFiles) == 0 {
		fmt.Println("No credentials found")
		fmt.Println()
		fmt.Printf("Please run command %s to create a credential\n", pp.Greenf("%s create", app.Name))
		os.Exit(0)
	}

	options := []huh.Option[credentials.Credential]{}

	for _, cred := range credFiles {
		options = append(options, huh.NewOption(cred.Name, cred))
	}

	return options, nil
}
