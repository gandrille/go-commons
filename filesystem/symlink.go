package filesystem

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/gandrille/go-commons/result"
)

// IMPORTANT! READ ME FIRST!
// All the functions in this file are designed for providing nice status messages,
// NOT for efficiency optimization.
// Do NOT use this functions if you need performance.

// UpdateOrCreateSymlink
func UpdateOrCreateSymlink(existing, linkname string) result.Result {

	// Check if existing (source) exists
	if exists, err := Exists(existing); !exists && err != nil {
		return result.NewError("Error while checking if " + existing + " exists: " + err.Error())
	} else if !exists {
		return result.NewError("Symbolic link destination " + existing + " does NOT exist")
	}

	// checking if symlink exists
	if exists, err := SymlinkExists(linkname); err != nil {
		return result.NewError("Error while retreiving link info: " + err.Error())
	} else if exists {
		if target, err := os.Readlink(linkname); err != nil {
			return result.NewError("Error while reading symlink destination: " + err.Error())
		} else {
			expected := filepath.Clean(strings.Replace(existing, "~", HomeDir(), 1) + "/")
			actual := filepath.Clean(strings.Replace(target, "~", HomeDir(), 1) + "/")

			// Nothing to update
			if actual == expected {
				return result.NewUnchanged("Symbolic link " + linkname + " is already pointing to " + existing)
			}

			// Unlink
			if err := os.Remove(linkname); err != nil {
				return result.NewError("Error while removing symbolic link " + actual + ": " + err.Error())
			}
		}
	}

	// Create link
	if err := os.Symlink(existing, linkname); err != nil {
		return result.NewError("Error while creating symbolic link " + linkname + ": " + err.Error())
	}

	return result.NewUpdated("symbolic link " + linkname + " is now pointing to " + existing)
}

// SymlinkExists checks if a symbolic link exists.
// Returns true if the file exists and is a symbolic link, false otherwise.
// error != nil if the file exists and is not a symbolic link (ie a directory).
// TODO : return value is a bit overcomplicated
func SymlinkExists(filePath string) (bool, error) {
	filePath = strings.Replace(filePath, "~", HomeDir(), 1)

	if stat, err := os.Lstat(filePath); err == nil {
		if stat.Mode()&os.ModeSymlink != 0 {
			return true, nil
		} else {
			return false, errors.New(filePath + " is not a symbolic link")
		}
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}
