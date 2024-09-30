package credentials

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/engmtcdrm/minno/env"
)

type Credential struct {
	Name string
	Path string
}

// Returns a slice of all available credentials
func GetCredFiles() ([]Credential, error) {
	envVars, err := env.GetEnv()
	if err != nil {
		return nil, err
	}

	var credFiles []Credential

	err = filepath.WalkDir(envVars.AppHomeDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".cred" {
			c := Credential{
				Name: strings.TrimSuffix(filepath.Base(path), ".cred"),
				Path: path,
			}
			credFiles = append(credFiles, c)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return credFiles, nil
}

func ResolveCredName(credName string) (string, error) {
	envVars, err := env.GetEnv()
	if err != nil {
		return "", err
	}

	if !filepath.IsAbs(credName) {
		credName = filepath.Join(envVars.AppHomeDir, credName)
	} else {
		fmt.Println(fmt.Errorf("Credential name cannot be an absolute path").Error())
		os.Exit(99)
	}

	if filepath.Ext(credName) != ".cred" {
		credName = credName + ".cred"
	}

	return credName, nil
}
