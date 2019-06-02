package env

import (
	"log"
	"os"
	"os/user"
)

// IMPORTANT! READ ME FIRST!
// All the functions in this file are designed for providing nice status messages,
// NOT for efficiency optimization.
// Do NOT use this functions if you need performance.

// Hostname gets the value of HOSTNAME environment variable.
// os.Exit(1) on failure.
func Hostname() string {
	value, err := os.Hostname()
	if err != nil {
		log.Fatal("FATAL ERROR: Can't access env variable 'HOSTNAME'. Reason: " + err.Error() + "\n")
	}
	if value == "" {
		log.Fatal("FATAL ERROR: Can't get env variable 'HOSTNAME'\n")
	}
	return value
}

// Username gets the login name.
// os.Exit(1) on failure.
func Username() string {
	user, err := user.Current()
	if err != nil {
		log.Fatal("FATAL ERROR: Can't get username. Reason: " + err.Error() + "\n")
	}
	if user == nil || user.Username == "" {
		log.Fatal("FATAL ERROR: Can't get username\n")
	}
	return user.Username
}
