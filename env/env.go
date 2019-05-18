package env

import (
	"log"
	"os"
	"os/user"
)

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
