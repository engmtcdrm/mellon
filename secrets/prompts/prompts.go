package prompts

import (
	"fmt"

	"github.com/charmbracelet/huh"

	pp "github.com/engmtcdrm/go-prettyprint"
	"github.com/engmtcdrm/minno/env"
	"github.com/engmtcdrm/minno/secrets"
)

// GetSecretOptions returns a slice of huh.Options for all available secrets
func GetSecretOptions(secretFiles []secrets.Secret, action string) ([]huh.Option[secrets.Secret], error) {
	envVars, err := env.GetEnv()
	if err != nil {
		return nil, err
	}

	if len(secretFiles) == 0 {
		return nil, fmt.Errorf("no secrets found to %s\n\nPlease run command %s to create a secret", action, pp.Greenf("%s create", envVars.ExeCmd))
	}

	options := []huh.Option[secrets.Secret]{}

	for _, secret := range secretFiles {
		options = append(options, huh.NewOption(secret.Name, secret))
	}

	return options, nil
}
