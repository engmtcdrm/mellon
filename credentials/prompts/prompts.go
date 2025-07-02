package prompts

import (
	"fmt"

	"github.com/charmbracelet/huh"

	pp "github.com/engmtcdrm/go-prettyprint"
	"github.com/engmtcdrm/minno/credentials"
	"github.com/engmtcdrm/minno/env"
)

// GetCredOptions returns a slice of huh.Options for all available credentials
func GetCredOptions(credFiles []credentials.Credential, action string) ([]huh.Option[credentials.Credential], error) {
	envVars, err := env.GetEnv()
	if err != nil {
		return nil, err
	}

	if len(credFiles) == 0 {
		return nil, fmt.Errorf("no credentials found to %s\n\nPlease run command %s to create a credential", action, pp.Greenf("%s create", envVars.ExeCmd))
	}

	options := []huh.Option[credentials.Credential]{}

	for _, cred := range credFiles {
		options = append(options, huh.NewOption(cred.Name, cred))
	}

	return options, nil
}
