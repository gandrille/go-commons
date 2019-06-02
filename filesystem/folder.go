package filesystem

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gandrille/go-commons/result"
)

// IMPORTANT! READ ME FIRST!
// All the functions in this file are designed for providing nice status messages,
// NOT for efficiency optimization.
// Do NOT use this functions if you need performance.

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

// EmptyFolder checks if a folder is empty.
// Returns a string explaining the folder status
// NOT_EXIST
// NOT_FOLDER
// NOT_EMPTY
// EMPTY
// or an error if diagnosis failed
func IsEmptyFolder(folderPath string) (string, error) {
	folderPath = strings.Replace(folderPath, "~", HomeDir(), 1)

	if stat, err := os.Stat(folderPath); err == nil {
		if stat.IsDir() {
			f, err := os.Open(folderPath)
			if err != nil {
				return "", errors.New("Error while reading " + folderPath + " content: " + err.Error())
			}
			defer f.Close()

			if _, err := f.Readdirnames(1); err == io.EOF {
				return "EMPTY", nil
			} else if err != nil {
				return "", errors.New("Error while reading " + folderPath + " file names: " + err.Error())
			} else {
				return "NOT_EMPTY", nil
			}
		} else {
			return "NOT_FOLDER", nil
		}
	} else if os.IsNotExist(err) {
		return "NOT_EXIST", nil
	} else {
		return "", errors.New("Error while reading " + folderPath + ": " + err.Error())
	}
}

// CreateFolderIfNeeded creates a folder if it does NOT exists.
// Returns a string which describes what has been done, or an error message.
func CreateFolderIfNeeded(folderPath string) result.Result {
	folderPath = strings.Replace(folderPath, "~", HomeDir(), 1)

	if exists, err := FolderExists(folderPath); err != nil {
		return result.NewError("Don't know if folder " + folderPath + " exists")
	} else if exists {
		return result.NewUnchanged("Folder " + folderPath + " already exists")
	}

	// create folder
	if err := os.MkdirAll(folderPath, 0775); err != nil {
		return result.NewError("Error while creating " + folderPath + ": " + err.Error())
	}

	return result.NewCreated("Folder " + folderPath + " created")
}
