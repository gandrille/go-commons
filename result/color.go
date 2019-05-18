package result

import (
	"fmt"

	"github.com/fatih/color"
)

var green = color.New(color.FgGreen).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()
var cyan = color.New(color.FgCyan).SprintFunc()

// PrintOK prints a success message.
// A newline is appended.
func PrintOK(message string) {
	fmt.Printf("%s %s\n", green("[OK]"), message)
}

// PrintError prints an error message.
// A newline is appended.
func PrintError(message string) {
	fmt.Printf("%s %s\n", red("[ERROR]"), message)
}

// PrintRed prints a message in red.
// A newline is appended.
func PrintRed(message string) {
	fmt.Printf("%s\n", red(message))
}

// PrintInfo prints an info message.
// A newline is appended.
func PrintInfo(message string) {
	color.Cyan(message)
}

// Describe prints the name in color, and the shortDesc using normal color.
// A newline is appended.
func Describe(name, shortDesc string) {
	if shortDesc == "" {
		fmt.Printf("%s\n", cyan(name))
	} else {
		fmt.Printf("%s %s\n", cyan(name), shortDesc)
	}
}
