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
	if exists, err := filesystem.RegularFileExists(xfconfExe); err != nil || exists == false {
		return "", errors.New("File " + xfconfExe + " does NOT exist")
	}

	// Property reading
	out, err := exec.Command(xfconfExe, "--channel", channel, "--property", property).Output()
	if err != nil {
		return "", errors.New("Can't read property " + property + " using xfconf")
	}
	value := strings.TrimSuffix(string(out), "\n")

	return value, nil
}

// SetXfconfProperty sets an Xconf property.
func SetXfconfProperty(channel, property, newValue string) result.Result {
	return CreateOrSetXfconfProperty(channel, property, "", newValue)
}

// CreateOrSetXfconfProperty sets an Xconf property.
func CreateOrSetXfconfProperty(channel, property, propType, newValue string) result.Result {

	// Read old value
	oldValue, _ := ReadXfconfProperty(channel, property)

	// Update needed: write new value (in case of reading error, oldValue is empty)
	if oldValue != newValue {
		params := []string{"--channel", channel, "--property", property, "--create", "--set", newValue}
		if propType != "" {
			params = append(params, "--type", propType)
		}
		if err := exec.Command(xfconfExe, params...).Run(); err != nil {
			return result.NewError("Can't write property " + property + " on channel " + channel + " with xconf. Reason : " + err.Error())
		}
	}

	// At this point, we have managed to read and write property
	// Let's compute the final success message
	switch {
	case oldValue == "":
		return result.NewCreated("Xconf property " + property + " on channel " + channel + " initialized with " + newValue)
	case oldValue == newValue:
		return result.NewUnchanged("Xconf property " + property + " on channel " + channel + " already has value " + newValue)
	default:
		return result.NewUpdated("Xconf property " + property + " on channel " + channel + " updated from " + oldValue + " to " + newValue)
	}
}

// ListXfconfProperties gives the list of all properties from a given channel.
func ListXfconfProperties(channel string) ([]string, error) {

	// Check if executable exists
	if exists, err := filesystem.RegularFileExists(xfconfExe); err != nil || exists == false {
		return nil, errors.New("File " + xfconfExe + " does NOT exist")
	}

	out, err := exec.Command(xfconfExe, "--channel", channel, "--list").Output()
	if err != nil {
		return []string{}, err
	}
	return strings.Split(string(out), "\n"), nil
}
