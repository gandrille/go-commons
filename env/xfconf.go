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

// Xfconf is Xfce configuration system
// https://docs.xfce.org/xfce/xfconf/start
const xfconfExe = "/usr/bin/xfconf-query"

// ReadXfconfProperty reads an Xfconf property.
// Returns the key value, or an error.
func ReadXfconfProperty(channel, property string) (string, error) {

	// Check if executable exists
	if exists, err := filesystem.FileExists(xfconfExe); err != nil || exists == false {
		return "", errors.New("File " + xfconfExe + " does NOT exist")
	}

	// Property reading
	out, err := exec.Command(xfconfExe, "--channel", channel, "--property", property).Output()
	if err != nil {
		return "", errors.New("Can't read property " + property + " with xfconf")
	}
	value := strings.TrimSuffix(string(out), "\n")

	return value, nil
}

// SetXfconfProperty sets an Xconf property.
func SetXfconfProperty(channel, property, newValue string) result.Result {

	// Read old value
	oldValue, err := ReadXfconfProperty(channel, property)
	if err != nil {
		return result.Failure(err.Error())
	}

	// Update needed: write new value
	if oldValue != newValue {
		if err := exec.Command(xfconfExe, "--channel", channel, "--property", property, "--set", newValue).Run(); err != nil {
			return result.Failure("Can't write property " + property + " on channel " + channel + " with xconf. Reason : " + err.Error())
		}
	}

	// At this point, we have managed to read and write property
	// Let's compute the final success message
	msg := ""
	switch {
	case oldValue == "":
		msg = "Xconf property " + property + " on channel " + channel + " initialized with " + newValue
	case oldValue == newValue:
		msg = "Xconf property " + property + " on channel " + channel + " already has value " + newValue
	default:
		msg = "Xconf property " + property + " on channel " + channel + " updated from " + oldValue + " to " + newValue
	}

	return result.Success(msg)
}

// ListXfconfProperties gives the list of all properties from a given channel.
func ListXfconfProperties(channel string) ([]string, error) {
	out, err := exec.Command(xfconfExe, "--channel", channel, "--list").Output()
	if err != nil {
		return []string{}, err
	}
	return strings.Split(string(out), "\n"), nil
}
