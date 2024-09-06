package header

import (
	"fmt"
	"strings"

	"github.com/engmtcdrm/minno/utils/env"
	"github.com/fatih/color"
)

func PrintBanner() {
	fmt.Println(color.MagentaString("       _"))
	fmt.Println(color.MagentaString(" _____|_|___ ___ ___"))
	fmt.Println(color.RedString("|     | |   |   | . |"))
	fmt.Println(color.RedString("|_|_|_|_|_|_|_|_|___| ") + color.GreenString("v"+env.AppVersion))
	fmt.Println(env.AppLongDesc)
	fmt.Println(color.CyanString(env.RepoUrl))
	fmt.Println(strings.Repeat("-", max(len(env.AppLongDesc), len(env.RepoUrl))))
	fmt.Println()
}
