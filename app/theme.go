package app

import (
	"github.com/engmtcdrm/go-pardon"
	pp "github.com/engmtcdrm/go-prettyprint"
)

func init() {
	minnoTheme()
}

func minnoTheme() {
	pardon.SetDefaultIconFunc(func(icon string) string { return pp.Cyan(icon) })
	pardon.SetDefaultAnswerFunc(func(answer string) string { return pp.Yellow(answer) })
	pardon.SetDefaultCursorFunc(func(cursor string) string { return pp.Red(cursor) })
	pardon.SetDefaultSelectFunc(func(s string) string { return pp.Magenta(s) })
	pardon.Icons.QuestionMark = "? "
}
