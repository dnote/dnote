package log

import (
	"fmt"
)

var (
	ColorRed    = 31
	ColorGreen  = 32
	ColorYellow = 33
	ColorBlue   = 34
	ColorGray   = 37
)

var indent = "  "

func Info(msg string) {
	fmt.Printf("%s\033[%dm%s\033[0m %s\n", indent, ColorBlue, "•", msg)
}

func Infof(msg string, v ...interface{}) {
	fmt.Printf("%s\033[%dm%s\033[0m %s", indent, ColorBlue, "•", fmt.Sprintf(msg, v...))
}

func Warnf(msg string, v ...interface{}) {
	fmt.Printf("%s\033[%dm%s\033[0m %s", indent, ColorRed, "•", fmt.Sprintf(msg, v...))
}

func Error(msg string) {
	fmt.Printf("%s\033[%dm%s\033[0m %s\n", indent, ColorRed, "⨯", msg)
}

func Printf(msg string, v ...interface{}) {
	fmt.Printf("%s\033[%dm%s\033[0m %s", indent, ColorGray, "•", fmt.Sprintf(msg, v...))
}

func WithPrefixf(prefixColor int, prefix, msg string, v ...interface{}) {
	fmt.Printf("  \033[%dm%s\033[0m %s\n", prefixColor, prefix, fmt.Sprintf(msg, v...))
}
