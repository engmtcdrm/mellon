package prettyprint

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/engmtcdrm/minno/utils/ansi"
)

var (
	IconComplete = Green("[✓] ")
	IconAlert    = Yellow("[!] ")
	IconFailed   = Red("[✗] ")
	IconInfo     = Cyan("[i] ")
)

// Black formats using the default formats for its operands and returns the
// resulting string with a black foreground. Spaces are added between
// operands when neither is a string.
func Black(a ...any) string {
	return ansi.Black + fmt.Sprint(a...) + ansi.Reset
}

// Blackf formats according to a format specifier and returns the resulting
// string with a black foreground.
func Blackf(format string, a ...any) string {
	return ansi.Black + fmt.Sprintf(format, a...) + ansi.Reset
}

// Red formats using the default formats for its operands and returns the
// resulting string with a red foreground. Spaces are added between
// operands when neither is a string.
func Red(a ...any) string {
	return ansi.Red + fmt.Sprint(a...) + ansi.Reset
}

// Redf formats according to a format specifier and returns the resulting
// string with a red foreground.
func Redf(format string, a ...any) string {
	return ansi.Red + fmt.Sprintf(format, a...) + ansi.Reset
}

// Green formats using the default formats for its operands and returns the
// resulting string with a green foreground. Spaces are added between
// operands when neither is a string.
func Green(a ...any) string {
	return ansi.Green + fmt.Sprint(a...) + ansi.Reset
}

// Greenf formats according to a format specifier and returns the resulting
// string with a green foreground.
func Greenf(format string, a ...any) string {
	return ansi.Green + fmt.Sprintf(format, a...) + ansi.Reset
}

// Yellow formats using the default formats for its operands and returns the
// resulting string with a yellow foreground. Spaces are added between
// operands when neither is a string.
func Yellow(a ...any) string {
	return ansi.Yellow + fmt.Sprint(a...) + ansi.Reset
}

// Yellowf formats according to a format specifier and returns the resulting
// string with a yellow foreground.
func Yellowf(format string, a ...any) string {
	return ansi.Yellow + fmt.Sprintf(format, a...) + ansi.Reset
}

// Blue formats using the default formats for its operands and returns the
// resulting string with a blue foreground. Spaces are added between
// operands when neither is a string.
func Blue(a ...any) string {
	return ansi.Blue + fmt.Sprint(a...) + ansi.Reset
}

// Bluef formats according to a format specifier and returns the resulting
// string with a blue foreground.
func Bluef(format string, a ...any) string {
	return ansi.Blue + fmt.Sprintf(format, a...) + ansi.Reset
}

// Magenta formats using the default formats for its operands and returns the
// resulting string with a magenta foreground. Spaces are added between
// operands when neither is a string.
func Magenta(a ...any) string {
	return ansi.Magenta + fmt.Sprint(a...) + ansi.Reset
}

// Magentaf formats according to a format specifier and returns the resulting
// string with a magenta foreground.
func Magentaf(format string, a ...any) string {
	return ansi.Magenta + fmt.Sprintf(format, a...) + ansi.Reset
}

// Cyan formats using the default formats for its operands and returns the
// resulting string with a cyan foreground. Spaces are added between
// operands when neither is a string.
func Cyan(a ...any) string {
	return ansi.Cyan + fmt.Sprint(a...) + ansi.Reset
}

// Cyanf formats according to a format specifier and returns the resulting
// string with a cyan foreground.
func Cyanf(format string, a ...any) string {
	return ansi.Cyan + fmt.Sprintf(format, a...) + ansi.Reset
}

// White formats using the default formats for its operands and returns the
// resulting string with a white foreground. Spaces are added between
// operands when neither is a string.
func White(a ...any) string {
	return ansi.White + fmt.Sprint(a...) + ansi.Reset
}

// Whitef formats according to a format specifier and returns the resulting
// string with a white foreground.
func Whitef(format string, a ...any) string {
	return ansi.White + fmt.Sprintf(format, a...) + ansi.Reset
}

func Complete(msg string) string {
	return IconComplete + msg
}

func Alert(msg string) string {
	return IconAlert + msg
}

func Fail(msg string) string {
	return IconFailed + msg
}

func Info(msg string) string {
	return IconInfo + msg
}

func Var(variable string, value string) string {
	return Info(Cyan(variable) + " is set to " + Green(value))
}

func ThemeMinno() *huh.Theme {
	t := huh.ThemeBase()

	var (
		black       = lipgloss.Color(strconv.Itoa(ansi.ANSIBlack))
		green       = lipgloss.Color(strconv.Itoa(ansi.ANSIGreen))
		yellow      = lipgloss.Color(strconv.Itoa(ansi.ANSIYellow))
		magenta     = lipgloss.Color(strconv.Itoa(ansi.ANSIMagenta))
		cyan        = lipgloss.Color(strconv.Itoa(ansi.ANSICyan))
		white       = lipgloss.Color(strconv.Itoa(ansi.ANSIWhite))
		blackBright = lipgloss.Color(strconv.Itoa(ansi.ANSIBrightBlack))
		red         = lipgloss.Color(strconv.Itoa(ansi.ANSIBrightRed))
	)

	t.Focused.Base = t.Focused.Base.BorderForeground(yellow)
	t.Focused.Title = t.Focused.Title.Foreground(cyan)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(cyan)
	t.Focused.Directory = t.Focused.Directory.Foreground(cyan)
	t.Focused.Description = t.Focused.Description.Foreground(blackBright)
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(red)
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(red)
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(yellow)
	t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(yellow)
	t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(yellow)
	t.Focused.Option = t.Focused.Option.Foreground(white)
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(yellow)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(green)
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(white)
	t.Focused.SelectedPrefix = t.Focused.SelectedPrefix.Foreground(green).SetString("✓ ")
	t.Focused.UnselectedPrefix = t.Focused.UnselectedPrefix.Foreground(white).SetString("• ")
	t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(white).Background(magenta)
	t.Focused.BlurredButton = t.Focused.BlurredButton.Foreground(white).Background(black)

	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(red)
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(blackBright)
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(yellow)

	t.Blurred = t.Focused
	t.Blurred.Base = t.Blurred.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.NoteTitle = t.Blurred.NoteTitle.Foreground(blackBright)
	t.Blurred.Title = t.Blurred.NoteTitle.Foreground(blackBright)

	t.Blurred.TextInput.Prompt = t.Blurred.TextInput.Prompt.Foreground(blackBright)
	t.Blurred.TextInput.Text = t.Blurred.TextInput.Text.Foreground(white)

	t.Blurred.NextIndicator = lipgloss.NewStyle()
	t.Blurred.PrevIndicator = lipgloss.NewStyle()

	return t
}
