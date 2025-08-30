package prompts

import (
	"fmt"

	"github.com/engmtcdrm/go-pardon"
	pp "github.com/engmtcdrm/go-prettyprint"
	"github.com/engmtcdrm/mellon/env"
	"github.com/engmtcdrm/mellon/secrets"
)

// GetSecretOptions returns a list of options for selecting a secret from the provided list of secret files.
func GetSecretOptions(secretFiles []secrets.Secret, action string) ([]pardon.Option[secrets.Secret], error) {
	if err := env.Init(); err != nil {
		return nil, err
	}

	if len(secretFiles) == 0 {
		return nil, fmt.Errorf(
			"no secrets found to %s\n\nPlease run command %s to create a secret",
			action,
			pp.Greenf("%s create", env.Instance.ExeCmd()),
		)
	}

	options := []pardon.Option[secrets.Secret]{}

	for _, secret := range secretFiles {
		options = append(options, pardon.NewOption(secret.Name(), secret))
	}

	return options, nil
}
