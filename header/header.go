package header

import (
	"fmt"
	"strings"

	"github.com/engmtcdrm/minno/app"
	pp "github.com/engmtcdrm/minno/utils/prettyprint"
)

func PrintBanner() {
	fmt.Println(pp.Magenta("       _"))
	fmt.Println(pp.Magenta(" _____|_|___ ___ ___"))
	fmt.Println(pp.Red("|     | |   |   | . |"))
	fmt.Println(pp.Red("|_|_|_|_|_|_|_|_|___| ") + pp.Greenf("v%s", app.Version))
	fmt.Println(app.LongDesc)
	fmt.Println(pp.Cyan(app.RepoUrl))
	fmt.Println(strings.Repeat("-", max(len(app.LongDesc), len(app.RepoUrl))))
	fmt.Println()
}
