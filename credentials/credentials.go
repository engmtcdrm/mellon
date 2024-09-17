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

func GetCredFiles() ([]Credential, error) {
	if env.AppHomeDir == "" {
		err := env.SetAppHomeDir()
		if err != nil {
			return nil, err
		}
	}

	var credFiles []Credential

	err := filepath.WalkDir(env.AppHomeDir, func(path string, d os.DirEntry, err error) error {
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

func ResolveCredName(credName string) string {
	if env.AppHomeDir == "" {
		err := env.SetAppHomeDir()
		if err != nil {
			return ""
		}
	}

	if !filepath.IsAbs(credName) {
		credName = filepath.Join(env.AppHomeDir, credName)
	} else {
		fmt.Println(fmt.Errorf("Credential name cannot be an absolute path").Error())
		os.Exit(1)
	}

	if !strings.HasSuffix(credName, ".cred") {
		credName = credName + ".cred"
	}

	return credName
}
