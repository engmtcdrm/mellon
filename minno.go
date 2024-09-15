package main

import (
	"fmt"
	"os"

	"github.com/engmtcdrm/minno/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(99)
	}
}
