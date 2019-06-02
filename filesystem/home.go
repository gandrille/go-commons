package filesystem

import (
	"log"
	"os"
)

// IMPORTANT! READ ME FIRST!
// All the functions in this file are designed for providing nice status messages,
// NOT for efficiency optimization.
// Do NOT use this functions if you need performance.

// HomeDir gets the value of HOME environment variable.
// Exit(1) on failure.
func HomeDir() string {
	value := os.Getenv("HOME")
	if value == "" {
		log.Fatal("FATAL ERROR: Can't access env variable 'HOME'")
	}
	return value
}
