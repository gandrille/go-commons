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
// WARNING : it has no link with debconf
// https://developer.gnome.org/dconf/unstable/dconf-overview.html

// ===============================
// ===============================
// == TODO migrate to gsettings ==
// ===============================
// ===============================

// For modifying the dconf backend storage itself, use the dconf tool; but gsettings should be used by preference.
// https://developer.gnome.org/GSettings/
const dconfExe = "/usr/bin/dconf"

// ReadDconfKey reads a dconf key.
// Returns the key value
func ReadDconfKey(key string) (string, error) {

	// Check if executable exists
	if exists, err := filesystem.RegularFileExists(dconfExe); err != nil || exists == false {
		return "", errors.New("File " + dconfExe + " does NOT exist")
	}

	// Key reading
	out, err := exec.Command(dconfExe, "read", key).Output()
	if err != nil {
		return "", errors.New("Can't read key " + key + " using dconf")
	}
	value := strings.TrimSuffix(string(out), "\n")

	return value, nil
}

// WriteDconfKey writes a dconf key.
func WriteDconfKey(key, newValue string) result.Result {

	// Read old value
	oldValue, err := ReadDconfKey(key)
	if err != nil {
		return result.NewError(err.Error())
	}

	// Update needed: write new value
	if oldValue != newValue {
		if err := exec.Command(dconfExe, "write", key, newValue).Run(); err != nil {
			return result.NewError("Can't write key '" + key + "' with dconf")
		}
	}

	// At this point, we have managed to read and write the key
	switch {
	case oldValue == "":
		return result.NewCreated("Dconf key '" + key + "' initialized with " + newValue)
	case oldValue == newValue:
		return result.NewUnchanged("Dconf key '" + key + "' already has value " + newValue)
	default:
		return result.NewUpdated("Dconf key '" + key + "' updated from " + oldValue + " to " + newValue)
	}
}
