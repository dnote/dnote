package log

import (
	"fmt"
	"github.com/dnote/color"
	"os"
)

var (
	// ColorRed is a red foreground color
	ColorRed = color.New(color.FgRed)
	// ColorGreen is a green foreground color
	ColorGreen = color.New(color.FgGreen)
	// ColorYellow is a yellow foreground color
	ColorYellow = color.New(color.FgYellow)
	// ColorBlue is a blue foreground color
	ColorBlue = color.New(color.FgBlue)
	// ColorGray is a gray foreground color
	ColorGray = color.New(color.FgHiBlack)
)

var indent = "  "

// Info prints information
func Info(msg string) {
	fmt.Fprintf(color.Output, "%s%s %s\n", indent, ColorBlue.Sprint("•"), msg)
}

// Infof prints information with optional format verbs
func Infof(msg string, v ...interface{}) {
	fmt.Fprintf(color.Output, "%s%s %s", indent, ColorBlue.Sprint("•"), fmt.Sprintf(msg, v...))
}

// Success prints a success message
func Success(msg string) {
	fmt.Fprintf(color.Output, "%s%s %s", indent, ColorGreen.Sprint("✔"), msg)
}

// Successf prints a success message with optional format verbs
func Successf(msg string, v ...interface{}) {
	fmt.Fprintf(color.Output, "%s%s %s", indent, ColorGreen.Sprint("✔"), fmt.Sprintf(msg, v...))
}

// Plain prints a plain message without any prefix symbol
func Plain(msg string) {
	fmt.Printf("%s%s", indent, msg)
}

// Plainf prints a plain message without any prefix symbol. It takes optional format verbs.
func Plainf(msg string, v ...interface{}) {
	fmt.Printf("%s%s", indent, fmt.Sprintf(msg, v...))
}

// Warnf prints a warning message with optional format verbs
func Warnf(msg string, v ...interface{}) {
	fmt.Fprintf(color.Output, "%s%s %s", indent, ColorRed.Sprint("•"), fmt.Sprintf(msg, v...))
}

// Error prints an error message
func Error(msg string) {
	fmt.Fprintf(color.Output, "%s%s %s\n", indent, ColorRed.Sprint("⨯"), msg)
}

// Errorf prints an error message with optional format verbs
func Errorf(msg string, v ...interface{}) {
	fmt.Fprintf(color.Output, "%s%s %s", indent, ColorRed.Sprintf("⨯"), fmt.Sprintf(msg, v...))
}

// Printf prints an normal message
func Printf(msg string, v ...interface{}) {
	fmt.Fprintf(color.Output, "%s%s %s", indent, ColorGray.Sprint("•"), fmt.Sprintf(msg, v...))
}

// Askf prints an question with optional format verbs. The leading symbol differs in color depending
// on whether the input is masked.
func Askf(msg string, masked bool, v ...interface{}) {
	symbolChar := "[?]"

	var symbol string
	if masked {
		symbol = ColorGray.Sprintf(symbolChar)
	} else {
		symbol = ColorGreen.Sprintf(symbolChar)
	}

	fmt.Fprintf(color.Output, "%s%s %s: ", indent, symbol, fmt.Sprintf(msg, v...))
}

// Debug prints to the console if DNOTE_DEBUG is set
func Debug(msg string, v ...interface{}) {
	if os.Getenv("DNOTE_DEBUG") == "1" {
		fmt.Fprintf(color.Output, "%s %s", ColorGray.Sprint("DEBUG:"), fmt.Sprintf(msg, v...))
	}
}
