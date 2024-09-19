package header

import (
	"fmt"
	"strings"

	"github.com/engmtcdrm/minno/app"
	"github.com/fatih/color"
)

func PrintBanner() {
	fmt.Println(color.MagentaString("       _"))
	fmt.Println(color.MagentaString(" _____|_|___ ___ ___"))
	fmt.Println(color.RedString("|     | |   |   | . |"))
	fmt.Println(color.RedString("|_|_|_|_|_|_|_|_|___| ") + color.GreenString("v"+app.Version))
	fmt.Println(app.LongDesc)
	fmt.Println(color.CyanString(app.RepoUrl))
	fmt.Println(strings.Repeat("-", max(len(app.LongDesc), len(app.RepoUrl))))
	fmt.Println()
}
