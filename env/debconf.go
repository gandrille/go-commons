package env

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/gandrille/go-commons/filesystem"
	"github.com/gandrille/go-commons/result"

	"github.com/gandrille/go-commons/misc"
)

// IMPORTANT! READ ME FIRST!
// All the functions in this file are designed for providing nice status messages,
// NOT for efficiency optimization.
// Do NOT use this functions if you need performance.

// debconf is a way to configure Debian packages.
// WARNING : it has no link with dconf
const debconfShowExe = "/usr/bin/debconf-show"
const debconfUpdateExe = "/usr/bin/debconf-set-selections"

// DebConfKey type
type DebConfKey struct {
	packageName  string
	keyName      string
	alreadyAsked bool
	value        string
}

// isNull checks if the key is a null object
func (key *DebConfKey) isNull() bool {
	return *key == NullDebConfKey()
}

// NullDebConfKey gets a null object
func NullDebConfKey() DebConfKey {
	return DebConfKey{"", "", false, ""}
}

// ReadDebconfKey retreive a key inside a package
// If the key is not found, the null key is returned and NO error is reported
func ReadDebconfKey(packageName, keyName string) (DebConfKey, error) {
	keys, err := ReadDebconfKeys(packageName)
	if err != nil {
		return NullDebConfKey(), err
	}

	for _, key := range keys {
		if keyName == key.keyName {
			return key, nil
		}
	}

	return NullDebConfKey(), nil
}

// ReadDebconfKeys retreive the list of debconf values relative to a package
func ReadDebconfKeys(packageName string) ([]DebConfKey, error) {

	// Check if executable exists
	if exists, err := filesystem.RegularFileExists(debconfShowExe); err != nil || exists == false {
		return nil, errors.New("File " + debconfShowExe + " does NOT exist")
	}

	// Retreive lines
	out, err := exec.Command(debconfShowExe, packageName).Output()

	// Error
	if err != nil {
		return nil, err
	}

	// No error
	keys := []DebConfKey{}
	for _, line := range strings.Split(string(out), "\n") {

		alreadyAsked := strings.HasPrefix(line, "*")
		if alreadyAsked {
			line = line[1:]
		}

		for strings.HasPrefix(line, " ") {
			line = line[1:]
		}

		// line is well formed
		sep := strings.Index(line, ":")
		if sep != -1 {
			keyName := line[:sep]
			keyValue := strings.Trim(line[sep+1:], " ")
			key := DebConfKey{packageName, keyName, alreadyAsked, keyValue}
			keys = append(keys, key)
		}
	}
	return keys, nil
}

// WriteDebconfKey writes a debconf key
func WriteDebconfKey(packageName, keyName, typeName, value string) result.Result {

	// Read
	key, err := ReadDebconfKey(packageName, keyName)

	if err == nil && !key.isNull() && value == key.value {
		return result.NewUnchanged("Key '" + keyName + "' already has value " + value)
	}

	// Check if executable exists
	if exists, err := filesystem.RegularFileExists(debconfUpdateExe); err != nil || exists == false {
		return result.NewError("File " + debconfUpdateExe + " does NOT exist")
	}

	// Write
	str := packageName + " " + keyName + " " + typeName + " " + value
	cmd := exec.Command(debconfUpdateExe, "-v")
	res := misc.RunCmdStdIn("Debconf "+str, str+"\n", cmd)
	if res.IsSuccess() {
		return result.NewUpdated("Key '" + keyName + "' updated with value " + value)
	}

	return result.NewError("Error while writing '" + keyName + "' " + err.Error())
}
