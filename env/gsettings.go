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

// Gsettings is Gnome configuration system
// https://developer.gnome.org/GSettings/
const gsettingsExe = "/usr/bin/gsettings"

// ReadGsettingsKey reads a gsettings property.
// Returns the key value, or an error if the key was not found.
func ReadGsettingsKey(schema, key string) (string, error) {

	// Check if executable exists
	if exists, err := filesystem.RegularFileExists(gsettingsExe); err != nil || exists == false {
		return "", errors.New("File " + gsettingsExe + " does NOT exist")
	}

	// Key reading
	out, err := exec.Command(gsettingsExe, "get", schema, key).Output()
	if err != nil {
		return "", errors.New("Can't read key '" + key + "' in schema '" + schema + "' using gsettings")
	}
	value := strings.TrimSuffix(string(out), "\n")

	return value, nil
}

// WriteGsettingsKey writes a gsettings key.
func WriteGsettingsKey(schema, key, newValue string) result.Result {

	// Read old value
	oldValue, err := ReadGsettingsKey(schema, key)
	exists := (err == nil)

	// No update needed
	if exists && oldValue == newValue {
		return result.NewUnchanged("Gsettings key " + key + " already has value " + newValue)
	}

	// Write new value
	if oldValue != newValue {
		if err := exec.Command(gsettingsExe, "set", schema, key, newValue).Run(); err != nil {
			result.NewError("Can't write key '" + key + "' in schema '" + schema + "' using gsettings")
		}
	}

	// Let's compute the final success message
	if exists {
		return result.NewUpdated("gsettings key " + key + " updated from " + oldValue + " to " + newValue)
	} else {
		return result.NewCreated("gsettings key " + key + " initialized with " + newValue)
	}
}
