package filesystem

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/gandrille/go-commons/result"
)

// FloderContent returns the list of elements (files, folders,...) inside a subtree.
func FloderContent(folderPath string) ([]string, error) {
	folderPath = strings.Replace(folderPath, "~", HomeDir(), 1)
	var fileList []string

	err := filepath.Walk(folderPath, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return err
	})

	return fileList, err
}

// FloderFiles returns the list of regular files inside a subtree.
func FloderFiles(folderPath string) ([]string, error) {
	folderPath = strings.Replace(folderPath, "~", HomeDir(), 1)
	var fileList []string

	err := filepath.Walk(folderPath, func(path string, f os.FileInfo, err error) error {
		if f.Mode().IsRegular() {
			fileList = append(fileList, path)
		}
		return err
	})

	return fileList, err
}

// FolderExists checks if a folder exists.
// Return true if the folder exists, false otherwise.
// error != nil if a file exists and is not a folder (ie a regular file or link).
// TODO : return value is a bit overcomplicated
func FolderExists(folderPath string) (bool, error) {
	folderName := strings.Replace(folderPath, HomeDir(), "~", 1)
	folderPath = strings.Replace(folderPath, "~", HomeDir(), 1)

	if stat, err := os.Stat(folderPath); err == nil {
		if stat.IsDir() {
			return true, nil
		} else {
			return false, errors.New(folderName + " is not a folder")
		}
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

// CreateFolderIfNeeded creates a folder if it does NOT exists.
// Returns a string which describes what has been done, or an error message.
func CreateFolderIfNeeded(folderPath string) result.Result {
	folderPath = strings.Replace(folderPath, "~", HomeDir(), 1)

	if exists, err := FolderExists(folderPath); err != nil {
		return result.New(false, "Don't know if folder "+folderPath+" exists")
	} else if exists {
		return result.New(true, "Folder "+folderPath+" already exists")
	}

	// create folder
	if err := os.MkdirAll(folderPath, 0775); err != nil {
		return result.New(false, "Error while creating "+folderPath+": "+err.Error())
	}

	return result.New(true, "Folder "+folderPath+" created")
}
