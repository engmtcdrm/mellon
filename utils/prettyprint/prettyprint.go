package prettyprint

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
)

func IconAlert(icon string) string {
	if icon == "" {
		return color.YellowString("[!] ")
	}

	return color.YellowString(icon)
}

func IconComplete(icon string) string {
	if icon == "" {
		return color.GreenString("[\u2713] ")
	}

	return color.GreenString(icon)
}

func IconInfo(icon string) string {
	if icon == "" {
		return color.CyanString("[i] ")
	}

	return color.CyanString(icon)
}

func IconFailed(icon string) string {
	if icon == "" {
		return color.RedString("[\u2717] ")
	}

	return color.RedString(icon)
}

func Complete(msg string, icon ...string) string {
	_icon := ""

	if len(icon) > 0 {
		_icon = icon[0]
	}

	return IconComplete(_icon) + msg
}

func Alert(msg string, icon ...string) string {
	_icon := ""

	if len(icon) > 0 {
		_icon = icon[0]
	}

	return IconAlert(_icon) + msg
}

func Fail(msg string, icon ...string) string {
	_icon := ""

	if len(icon) > 0 {
		_icon = icon[0]
	}
	return IconFailed(_icon) + msg
}

func Info(msg string, icon ...string) string {
	_icon := ""

	if len(icon) > 0 {
		_icon = icon[0]
	}

	return IconInfo(_icon) + msg
}

func Var(variable string, value string) string {
	return Info(color.CyanString(variable) + " is set to " + color.GreenString(value))
}

func ThemeMinno() *huh.Theme {
	t := huh.ThemeBase()

	t.Focused.Base = lipgloss.NewStyle().PaddingLeft(1).BorderStyle(lipgloss.ThickBorder()).BorderLeft(true)

	// GBNT
	// var (
	// 	background = lipgloss.AdaptiveColor{Dark: "#314056"}
	// 	selection  = lipgloss.AdaptiveColor{Dark: "#435775"}
	// 	foreground = lipgloss.AdaptiveColor{Dark: "#f8f8f2"}
	// 	comment    = lipgloss.AdaptiveColor{Dark: "#747572"}
	// 	green      = lipgloss.AdaptiveColor{Dark: "#B5D2AD"}
	// 	purple     = lipgloss.AdaptiveColor{Dark: "#196589"}
	// 	red        = lipgloss.AdaptiveColor{Dark: "#C5283D"}
	// 	// red        = lipgloss.AdaptiveColor{Dark: "#F51000"}
	// 	yellow = lipgloss.AdaptiveColor{Dark: "#00B0A5"}
	// )

	// DWARF
	// var (
	// 	background = lipgloss.AdaptiveColor{Dark: "#414342"}
	// 	selection  = lipgloss.AdaptiveColor{Dark: "#5A5E5C"}
	// 	foreground = lipgloss.AdaptiveColor{Dark: "#E7E4DB"}
	// 	comment    = lipgloss.AdaptiveColor{Dark: "#968D8E"}
	// 	green      = lipgloss.AdaptiveColor{Dark: "#7BF9FF"}
	// 	purple     = lipgloss.AdaptiveColor{Dark: "#9573FF"}
	// 	red        = lipgloss.AdaptiveColor{Dark: "#BF4A40"}
	// 	yellow     = lipgloss.AdaptiveColor{Dark: "#EE9739"}
	// )

	// BASIC
	var (
		background = lipgloss.AdaptiveColor{Dark: "#414342"}
		selection  = lipgloss.AdaptiveColor{Dark: "#5A5E5C"}
		foreground = lipgloss.AdaptiveColor{Dark: "#E7E4DB"}
		comment    = lipgloss.AdaptiveColor{Dark: "#968D8E"}
		green      = lipgloss.AdaptiveColor{Dark: "#7BF9FF"}
		purple     = lipgloss.AdaptiveColor{Dark: "#9573FF"}
		red        = lipgloss.AdaptiveColor{Dark: "#BF4A40"}
		// yellow     = lipgloss.AdaptiveColor{Dark: "#EE9739"}
	)

	yellow2 := lipgloss.Color("0") // gray, unset
	yellow2 = lipgloss.Color("1")  // red
	yellow2 = lipgloss.Color("2")  // green
	yellow2 = lipgloss.Color("3")  // yellow
	yellow2 = lipgloss.Color("4")  // blue
	yellow2 = lipgloss.Color("5")  // magenta
	yellow2 = lipgloss.Color("6")  // cyan
	yellow2 = lipgloss.Color("7")  // white
	yellow2 = lipgloss.Color("8")  // gray?
	yellow2 = lipgloss.Color("9")  // red highlight
	// yellow2 = lipgloss.Color("10") // green highlight
	// yellow2 = lipgloss.Color("11") // yellow highlight
	// yellow2 = lipgloss.Color("12") // blue highlight
	yellow2 = lipgloss.Color("13") // magenta highlight
	// yellow2 = lipgloss.Color("14") // cyan highlight
	// yellow2 = lipgloss.Color("15") // white highlight

	t.Focused.Base = t.Focused.Base.BorderForeground(selection)
	t.Focused.Title = t.Focused.Title.Foreground(purple)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(purple)
	t.Focused.Description = t.Focused.Description.Foreground(comment)
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(red)
	t.Focused.Directory = t.Focused.Directory.Foreground(purple)
	t.Focused.File = t.Focused.File.Foreground(foreground)
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(red)
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(yellow2)
	t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(yellow2)
	t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(yellow2)
	t.Focused.Option = t.Focused.Option.Foreground(foreground)
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(yellow2)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(green)
	t.Focused.SelectedPrefix = t.Focused.SelectedPrefix.Foreground(green)
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(foreground)
	t.Focused.UnselectedPrefix = t.Focused.UnselectedPrefix.Foreground(comment)
	t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(yellow2).Background(purple).Bold(true)
	t.Focused.BlurredButton = t.Focused.BlurredButton.Foreground(foreground).Background(background)

	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(yellow2)
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(comment)
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(yellow2)

	t.Blurred = t.Focused
	t.Blurred.Base = t.Blurred.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.NextIndicator = lipgloss.NewStyle()
	t.Blurred.PrevIndicator = lipgloss.NewStyle()

	return t
}
