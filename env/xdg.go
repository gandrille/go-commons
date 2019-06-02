package env

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gandrille/go-commons/filesystem"
)

// IMPORTANT! READ ME FIRST!
// All the functions in this file are designed for providing nice status messages,
// NOT for efficiency optimization.
// Do NOT use this functions if you need performance.

const xdgRead = "/usr/bin/xdg-user-dir"
const xdgUpdate = "/usr/bin/xdg-user-dirs-update"

// ReadXdgKey
func ReadXdgDir(key string) (string, error) {

	// Check if executable exists
	if exists, err := filesystem.RegularFileExists(xdgRead); err != nil || exists == false {
		return "", errors.New("File " + xdgRead + " does NOT exist")
	}

	dir, err := exec.Command(xdgRead, key).Output()
	if err != nil {
		return "", errors.New("Can't retreive xdg dir " + key + ": " + err.Error())
	}
	path := filepath.Clean(strings.Trim(string(dir), "\n") + "/")

	return path, nil
}

// UpdateXdgDir an XDK key referencing a directory
// Returns an error in case of failure.
// Return true if the new value is different than the previous one.
func UpdateXdgDir(key, value string) (bool, error) {

	hostName := Hostname()
	homeDir := filesystem.HomeDir()

	// Read
	src, errReadXdgKey := ReadXdgDir(key)
	if errReadXdgKey != nil {
		return false, errReadXdgKey
	}

	// Normalize dest
	dst := filepath.Clean(strings.Replace(strings.Replace(value, "$HOME", homeDir, -1), "$HOSTNAME", hostName, -1) + "/")

	// Is update needed ?
	if src == dst {
		return false, nil
	}

	// Check if executable exists
	if exists, err := filesystem.RegularFileExists(xdgUpdate); err != nil || exists == false {
		return false, errors.New("File " + xdgUpdate + " does NOT exist")
	}

	// Write new value
	if err := exec.Command(xdgUpdate, "--set", key, dst).Run(); err != nil {
		return false, errors.New("Can't write xdg dir " + key + ": " + err.Error())
	}

	// Remove src directory if needed
	if str, err := filesystem.IsEmptyFolder(src); err != nil {
		fmt.Println("[WARNING] Error while checking if " + src + " is empty: " + err.Error())
	} else if str == "EMPTY" {
		if err := os.Remove(src); err != nil {
			fmt.Println("[WARNING] empty folder " + src + " was NOT removed: " + err.Error())
		} else {
			fmt.Println("[INFO] empty folder " + src + " removed")
		}
	}

	// Create destination if needed
	if result := filesystem.CreateFolderIfNeeded(dst); !result.IsSuccess() {
		fmt.Println("[WARNING] Error creating " + dst + " " + result.Message())
	}

	return true, nil
}
