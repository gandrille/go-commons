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

const gioExe = "/usr/bin/gio"

// ReadMimeDefautApp Returns the default application value for a mime type, or an error
func ReadMimeDefautApp(typeName string) (string, error) {

	// Check if executable exists
	if exists, err := filesystem.RegularFileExists(gioExe); err != nil || exists == false {
		return "", errors.New("File " + gioExe + " does NOT exist")
	}

	// App reading
	out, err := exec.Command(gioExe, "mime", typeName).Output()
	if err != nil {
		return "", errors.New("Can't read default application for mime type '" + typeName + "'")
	}
	firstLine := strings.Split(string(out), "\n")[0]

	sep := strings.Index(firstLine, ":")
	if sep == -1 {
		return "", nil
	}

	value := strings.Trim(firstLine[sep+1:], " ")
	return value, nil
}

// WriteMimeDefautApp
func WriteMimeDefautApp(typeName, newValue string) result.Result {

	// Read old value
	oldValue, err1 := ReadMimeDefautApp(typeName)
	exists := (err1 == nil && oldValue != "")

	// No update needed
	if exists && oldValue == newValue {
		return result.NewUnchanged("Mime application associated with mime type " + typeName + " already has value " + newValue)
	}

	// Write new value
	_, err2 := exec.Command(gioExe, "mime", typeName, newValue).Output()
	if err2 != nil {
		return result.NewError("Can't update application associated with mime type " + typeName)
	}

	// Let's compute the final success message
	if exists {
		return result.NewUpdated("Mime application associated with mime type " + typeName + " updated with value " + newValue)
	} else {
		return result.NewCreated("Mime application associated with mime type " + typeName + " created with value " + newValue)
	}
}
