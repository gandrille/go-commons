package env

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/gandrille/go-commons/filesystem"
	"github.com/gandrille/go-commons/result"
)

// IMPORTANT! READ ME FIRST!
// All the functions in this file are designed for providing nice status messages,
// NOT for efficiency optimization.
// Do NOT use this functions if you need performance.

// dconf is Gnome configuration system
// https://developer.gnome.org/dconf/unstable/dconf-overview.html
const dconfExe = "/usr/bin/dconf"

// ReadDconfKey reads a dconf key.
// Returns the key value
func ReadDconfKey(key string) (string, error) {

	// Check if executable exists
	if exists, err := filesystem.FileExists(dconfExe); err != nil || exists == false {
		return "", errors.New("File " + dconfExe + " does NOT exist")
	}

	// Key reading
	out, err := exec.Command(dconfExe, "read", key).Output()
	if err != nil {
		return "", errors.New("Can't read key " + key + " with dconf")
	}
	value := strings.TrimSuffix(string(out), "\n")

	return value, nil
}

// WriteDconfKey writes a dconf key.
func WriteDconfKey(key, newValue string) result.Result {

	// Read old value
	oldValue, err := ReadDconfKey(key)
	if err != nil {
		result.Failure(err.Error())
	}

	// Update needed: write new value
	if oldValue != newValue {
		if err := exec.Command(dconfExe, "write", key, newValue).Run(); err != nil {
			result.Failure("Can't write key " + key + " with dconf")
		}
	}

	// At this point, we have managed to read and write the key
	// Let's compute the final success message
	msg := ""
	switch {
	case oldValue == "":
		msg = "Dconf key " + key + " initialized with " + newValue
	case oldValue == newValue:
		msg = "Dconf key " + key + " already has value " + newValue
	default:
		msg = "Dconf key " + key + " updated from " + oldValue + " to " + newValue
	}

	return result.Success(msg)
}
