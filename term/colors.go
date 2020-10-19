package term

import (
	"fmt"
)

type colorValue int

// Color values
const (
	CBlack        colorValue = 30
	CRed          colorValue = 31
	CGreen        colorValue = 32
	CYellow       colorValue = 33
	CBlue         colorValue = 34
	CMagenta      colorValue = 35
	CCyan         colorValue = 36
	CLightGray    colorValue = 37
	CDarkGray     colorValue = 90
	CLightRed     colorValue = 91
	CLightGreen   colorValue = 92
	CLightYellow  colorValue = 93
	CLightBlue    colorValue = 94
	CLightMagenta colorValue = 95
	CLightCyan    colorValue = 96
	CWhite        colorValue = 97
)

// Colored wraps msg into escape codes coloring the text into given colorValue
func Colored(msg string, c colorValue, bold bool) string {
	bstr := ""
	if bold {
		bstr = ";1"
	}
	return fmt.Sprintf("\033[%d%sm%s\033[0m", c, bstr, msg)
}

// Blue returns a blue version of msg
func Blue(msg string) string {
	return Colored(msg, CLightBlue, false)
}

// Red returns a red version of msg
func Red(msg string) string {
	return Colored(msg, CLightRed, false)
}

// Green returns a green version of msg
func Green(msg string) string {
	return Colored(msg, CLightGreen, false)
}

// Yellow returns a yellow version of msg
func Yellow(msg string) string {
	return Colored(msg, CLightYellow, false)
}

// Cyan returns a cyan version of msg
func Cyan(msg string) string {
	return Colored(msg, CLightCyan, false)
}

// Errorf is a red-colored version of fmt.Printf
func Errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Print(Red(msg))
}

// Successf is a green-colored version of fmt.Printf
func Successf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Print(Green(msg))
}

// Warnf is a yellow-colored version of fmt.Printf
func Warnf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Print(Yellow(msg))
}
