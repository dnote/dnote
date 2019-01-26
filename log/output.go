package log

import (
	"fmt"
)

// PrintContent prints the note content with an appropriate format.
func PrintContent(content string) {
	fmt.Printf("\n-----------------------content-----------------------\n")
	fmt.Printf("%s", content)
	fmt.Printf("\n-----------------------------------------------------\n")
}
