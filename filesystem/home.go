package filesystem

import (
	"log"
	"os"
)

// HomeDir gets the value of HOME environment variable.
// Exit(1) on failure.
func HomeDir() string {
	value := os.Getenv("HOME")
	if value == "" {
		log.Fatal("FATAL ERROR: Can't access env variable 'HOME'")
	}
	return value
}
