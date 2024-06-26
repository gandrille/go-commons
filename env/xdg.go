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

const xdgExec = "/usr/bin/xdg-settings"
const xdgRead = "/usr/bin/xdg-user-dir"
const xdgUpdate = "/usr/bin/xdg-user-dirs-update"

// ReadXdgSettings Reads an XDG settings
func ReadXdgSettings(key string) (string, error) {

	// Check if executable exists
	if exists, err := filesystem.RegularFileExists(xdgExec); err != nil || !exists {
		return "", errors.New("File " + xdgExec + " does NOT exist")
	}

	value, err := exec.Command(xdgExec, "get", key).Output()
	if err != nil {
		return "", errors.New("Can't retreive xdg-settings for " + key + ": " + err.Error())
	}
	result := strings.Trim(string(value), "\n")

	return result, nil
}

// WriteXdgSettings Updates an XDG property
func WriteXdgSettings(key, value string) (bool, error) {

	// Read
	src, errReadXdgKey := ReadXdgSettings(key)
	if errReadXdgKey != nil {
		return false, errReadXdgKey
	}

	// Is update needed ?
	if src == value {
		return false, nil
	}

	// Write new value
	if err := exec.Command(xdgExec, "set", key, value).Run(); err != nil {
		return false, errors.New("Can't write xdg value for " + key + ": " + err.Error())
	}

	return true, nil
}

// ReadXdgDir Reads a XDG folder key
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

// UpdateXdgDir Updates the location of an XDG directory
// It updates the ~/.config/user-dirs.dirs file using XDG executables
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
			fmt.Println("[UPDATED] empty folder " + src + " removed")
		}
	} else {
		fmt.Println("[WARNING] folder " + src + " was NOT removed because it is not empty")
	}

	// Create destination if needed
	if result := filesystem.CreateFolderIfNeeded(dst); !result.IsSuccess() {
		fmt.Println("[WARNING] Error creating " + dst + " " + result.Message())
	}

	return true, nil
}
