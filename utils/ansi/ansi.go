package ansi

import (
	"regexp"
	"strconv"
)

var (
	// Move the cursor to the top left corner of screen
	CursorTopLeft = CursorPosition(1, 1)
	// Move the cursor to beginning of the row (line).
	CursorLineBegin = CursorHorizontalAbsolute(1)
	// Move cursor down one row (line).
	CursorNextLine = CursorNextLineN(1)
	// Move cursor up one row (line).
	CursorPreviousLine = CursorPreviousLineN(1)

	// Scroll the screen up one row (line).
	ScrollNext = ScrollUpN(1)
	// Scroll the screen down one row (line).
	ScrollPrev = ScrollDownN(1)
)

// ANSI color codes.
const (
	ANSIBlack = iota
	ANSIRed
	ANSIGreen
	ANSIYellow
	ANSIBlue
	ANSIMagenta
	ANSICyan
	ANSIWhite
	ANSIBrightBlack
	ANSIBrightRed
	ANSIBrightGreen
	ANSIBrightYellow
	ANSIBrightBlue
	ANSIBrightMagenta
	ANSIBrightCyan
	ANSIBrightWhite
)

const (
	csi = "\x1b["
	// Clears from the cursor to the end of the screen.
	ClearFromCursorToEndScreen = csi + "0J"
	// Clears from the cursor to the beginning of the screen.
	ClearFromCursorToBeginScreen = csi + "1J"
	// Clear the entire screen.
	ClearScreen = csi + "2J"
	// Clear from the cursor to the end of the row (line).
	ClearToEnd = csi + "K"
	// Clear from the cursor to the beginning of the row (line).
	ClearToBegin = csi + "1K"
	// Clear the entire row (line).
	ClearLine = csi + "2K"
	// Clear the entire row (line) and move the cursor to the start of the row (line).
	ClearLineReset = csi + "2k\r"

	// Save the cursor position
	SaveCursorPos = csi + "s"
	// Restore the cursor position
	RestoreCursorPos = csi + "u"

	// Hide the cursor
	HideCursor = csi + "?25l"

	// Show the cursor
	ShowCursor = csi + "?25h"

	// Text Formatting

	Reset      = csi + "0m"
	Bold       = csi + "1m"
	Dim        = csi + "2m"
	Italic     = csi + "3m"
	Underline  = csi + "4m"
	SlowBlink  = csi + "5m"
	RapidBlink = csi + "6m"
	Reverse    = csi + "7m"
	Hidden     = csi + "8m"
	Strike     = csi + "9m"

	// Reset Formatting

	DoubleUnderline = csi + "21m"
	ResetIntensity  = csi + "22m"
	ResetItalic     = csi + "23m"
	ResetUnderline  = csi + "24m"
	ResetBlink      = csi + "25m"
	ResetReverse    = csi + "27m"
	ResetHidden     = csi + "28m"
	ResetStrike     = csi + "29m"

	// Foreground Color

	Black   = csi + "30m"
	Red     = csi + "31m"
	Green   = csi + "32m"
	Yellow  = csi + "33m"
	Blue    = csi + "34m"
	Magenta = csi + "35m"
	Cyan    = csi + "36m"
	White   = csi + "37m"

	// Background Color

	BlackBg   = csi + "40m"
	RedBg     = csi + "41m"
	GreenBg   = csi + "42m"
	YellowBg  = csi + "43m"
	BlueBg    = csi + "44m"
	MagentaBg = csi + "45m"
	CyanBg    = csi + "46m"
	WhiteBg   = csi + "47m"
	DefaultBg = csi + "49m"

	// Intense Foreground Color

	IntenseBlack   = csi + "90m"
	IntenseRed     = csi + "91m"
	IntenseGreen   = csi + "92m"
	IntenseYellow  = csi + "93m"
	IntenseBlue    = csi + "94m"
	IntenseMagenta = csi + "95m"
	IntenseCyan    = csi + "96m"
	IntenseWhite   = csi + "97m"

	// Intense Background Color

	IntenseBlackBg   = csi + "100m"
	IntenseRedBg     = csi + "101m"
	IntenseGreenBg   = csi + "102m"
	IntenseYellowBg  = csi + "103m"
	IntenseBlueBg    = csi + "104m"
	IntenseMagentaBg = csi + "105m"
	IntenseCyanBg    = csi + "106m"
	IntenseWhiteBg   = csi + "107m"
)

// 8-bit foreground color
// color must be between 0 and 255 otherwise it will return an empty string
func Foreground8Bit(color int) string {
	if color < 0 || color > 255 {
		return ""
	}

	return csi + "38;5;" + strconv.Itoa(color) + "m"
}

// 8-bit background color
// color must be between 0 and 255 othrewise it will return an empty string
func Background8Bit(color int) string {
	if color < 0 || color > 255 {
		return ""
	}

	return csi + "48;5;" + strconv.Itoa(color) + "m"
}

// 24-bit foreground color
// r, g, b must be between 0 and 255 otherwise it will return an empty string
func Foreground24Bit(r, g, b int) string {
	if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
		return ""
	}

	return csi + "38;2;" + strconv.Itoa(r) + ";" + strconv.Itoa(g) + ";" + strconv.Itoa(b) + "m"
}

// 24-bit background color
// r, g, b must be between 0 and 255 otherwise it will return an empty string
func Background24Bit(r, g, b int) string {
	if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
		return ""
	}

	return csi + "48;2;" + strconv.Itoa(r) + ";" + strconv.Itoa(g) + ";" + strconv.Itoa(b) + "m"
}

// CursorUp moves the cursor up n rows (lines).
func CursorUp(n int) string {
	return csi + strconv.Itoa(n) + "A"
}

// CursorDown moves the cursor down n rows (lines).
func CursorDown(n int) string {
	return csi + strconv.Itoa(n) + "B"
}

// CursorForward moves the cursor forward n columns.
func CursorForward(n int) string {
	return csi + strconv.Itoa(n) + "C"
}

// CursorBackward moves the cursor backward n columns.
func CursorBackward(n int) string {
	return csi + strconv.Itoa(n) + "D"
}

// CursorNextLine moves the cursor down n rows (lines).
func CursorNextLineN(n int) string {
	return csi + strconv.Itoa(n) + "E"
}

// CursorPreviousLine moves the cursor up n rows (lines).
func CursorPreviousLineN(n int) string {
	return csi + strconv.Itoa(n) + "F"
}

// CursorHorizontalAbsolute moves the cursor to the nth column.
func CursorHorizontalAbsolute(n int) string {
	return csi + strconv.Itoa(n) + "G"
}

// CursorPosition moves the cursor to the specified position of row (line) and column.
func CursorPosition(row, column int) string {
	return csi + strconv.Itoa(row) + ";" + strconv.Itoa(column) + "H"
}

// ScrollUp scrolls the screen up n rows (lines).
func ScrollUpN(n int) string {
	return csi + strconv.Itoa(n) + "S"
}

// ScrollDownN scrolls the screen down n rows (lines).
func ScrollDownN(n int) string {
	return csi + strconv.Itoa(n) + "T"
}

// StripCodes removes all ANSI escape codes from the input string.
func StripCodes(input string) string {
	re := regexp.MustCompile(`\x1b\[[;?0-9]*[a-zA-Z]`)
	return re.ReplaceAllString(input, "")
}
