package log

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	SprintfRed    = color.New(color.FgRed).SprintfFunc()
	SprintfGreen  = color.New(color.FgGreen).SprintfFunc()
	SprintfYellow = color.New(color.FgYellow).SprintfFunc()
	SprintfBlue   = color.New(color.FgBlue).SprintfFunc()
	SprintfGray   = color.New(color.FgWhite).SprintfFunc()
)

var indent = "  "

func Info(msg string) {
	fmt.Fprintf(color.Output, "%s %s %s\n", indent, SprintfBlue("•"), msg)
}

func Infof(msg string, v ...interface{}) {
	fmt.Fprintf(color.Output, "%s %s %s", indent, SprintfBlue("•"), fmt.Sprintf(msg, v...))
}

func Success(msg string) {
	fmt.Fprintf(color.Output, "%s%s %s", indent, SprintfGreen("✔"), msg)
}

func Successf(msg string, v ...interface{}) {
	fmt.Fprintf(color.Output, "%s%s %s", indent, SprintfGreen("✔"), fmt.Sprintf(msg, v...))
}

func Plain(msg string) {
	fmt.Printf("%s%s", indent, msg)
}

func Plainf(msg string, v ...interface{}) {
	fmt.Fprintf(color.Output, "%s%s", indent, fmt.Sprintf(msg, v...))
}

func Warnf(msg string, v ...interface{}) {
	fmt.Fprintf(color.Output, "%s%s %s", indent, SprintfRed("•"), fmt.Sprintf(msg, v...))
}

func Error(msg string) {
	fmt.Fprintf(color.Output, "%s%s %s\n", indent, SprintfRed("⨯"), msg)
}

func Printf(msg string, v ...interface{}) {
	fmt.Fprintf(color.Output, "%s%s %s", indent, SprintfGray("•"), fmt.Sprintf(msg, v...))
}
