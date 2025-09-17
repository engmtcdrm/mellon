package main

import (
	"os"

	"github.com/engmtcdrm/mellon/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(99)
	}
}
