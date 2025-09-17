package header

import (
	"fmt"
	"strings"

	pp "github.com/engmtcdrm/go-prettyprint"
	"github.com/engmtcdrm/mellon/app"
)

// PrintHeader prints the header of the application
func PrintHeader() {
	fmt.Println(pp.Magenta("           _ _"))
	fmt.Println(pp.Magenta(" _____ ___| | |___ ___"))
	fmt.Println(pp.Red("|     | -_| | | . |   |"))
	fmt.Println(pp.Red("|_|_|_|___|_|_|___|_|_| ") + pp.Green(app.Version))
	fmt.Println(app.LongDesc)
	fmt.Println(pp.Cyan(app.RepoUrl))
	fmt.Println(strings.Repeat("-", max(len(app.LongDesc), len(app.RepoUrl))))
	fmt.Println()
}
