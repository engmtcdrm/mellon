package app

import (
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/engmtcdrm/go-ansi"
)

// Application constants
const (
	// Name of the application
	Name = "minno"

	// DotName is the name of the application with a dot prefix
	DotName = "." + Name

	// LongDesc provides a detailed description of the application
	LongDesc = "A lightweight CLI tool for securing and obtaining secrets"

	// ShortDesc provides a brief description of the application
	ShortDesc = "A lightweight CLI tool for securing and obtaining secrets"

	// RepoUrl is the URL of the application's repository
	RepoUrl = "https://github.com/engmtcdrm/minno"
)

var (
	// Version of the application
	Version = "dev"
)

// ThemeMinno returns a new theme based on the Minno color scheme.
func ThemeMinno() *huh.Theme {
	t := huh.ThemeBase()

	var (
		black       = lipgloss.Color(strconv.Itoa(ansi.ANSIBlack))
		green       = lipgloss.Color(strconv.Itoa(ansi.ANSIGreen))
		yellow      = lipgloss.Color(strconv.Itoa(ansi.ANSIYellow))
		magenta     = lipgloss.Color(strconv.Itoa(ansi.ANSIMagenta))
		cyan        = lipgloss.Color(strconv.Itoa(ansi.ANSICyan))
		white       = lipgloss.Color(strconv.Itoa(ansi.ANSIWhite))
		brightBlack = lipgloss.Color(strconv.Itoa(ansi.ANSIBrightBlack))
		red         = lipgloss.Color(strconv.Itoa(ansi.ANSIBrightRed))
	)

	t.Focused.Base = t.Focused.Base.BorderForeground(yellow)
	t.Focused.Title = t.Focused.Title.Foreground(cyan)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(cyan)
	t.Focused.Directory = t.Focused.Directory.Foreground(cyan)
	t.Focused.Description = t.Focused.Description.Foreground(brightBlack)
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

	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(green)
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(brightBlack)
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(yellow)

	t.Blurred = t.Focused
	t.Blurred.Base = t.Blurred.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.NoteTitle = t.Blurred.NoteTitle.Foreground(brightBlack)
	t.Blurred.Title = t.Blurred.NoteTitle.Foreground(brightBlack)

	t.Blurred.TextInput.Prompt = t.Blurred.TextInput.Prompt.Foreground(brightBlack)
	t.Blurred.TextInput.Text = t.Blurred.TextInput.Text.Foreground(white)

	t.Blurred.NextIndicator = lipgloss.NewStyle()
	t.Blurred.PrevIndicator = lipgloss.NewStyle()

	return t
}
