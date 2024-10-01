package main

import (
	"os"

	"github.com/engmtcdrm/minno/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(99)
	}
}
