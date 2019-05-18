package filesystem

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/gandrille/go-commons/result"
)

// IMPORTANT! READ ME FIRST!
// All the functions in this file are designed for providing nice status messages,
// NOT for efficiency optimization.
// Do NOT use this functions if you need performance.

// ReadFileAsBinary gets the content of a binary file.
func ReadFileAsBinary(filePath string) ([]byte, error) {
	filePath = strings.Replace(filePath, "~", HomeDir(), 1)

	if exists, err := FileExists(filePath); err != nil || exists == false {
		return []byte{}, errors.New("File " + filePath + " does NOT exist")
	}

	byteArray, err := ioutil.ReadFile(filePath)
	if err != nil {
		return []byte{}, errors.New("Can't read file " + filePath + ": " + err.Error())
	}
	return byteArray, nil
}

// ReadFileAsString gets the content of a text file.
func ReadFileAsString(filePath string) (string, error) {
	byteArray, err := ReadFileAsBinary(filePath)
	if err != nil {
		return "", err
	}
	return string(byteArray), nil
}

// FileExists checks if a regular file exists.
// Returns true if the file exists and is a regular file, false otherwise.
// error != nil if the file exists and is not a regular file (ie a directory).
// TODO : return value is a bit overcomplicated
func FileExists(filePath string) (bool, error) {
	filePath = strings.Replace(filePath, "~", HomeDir(), 1)

	if stat, err := os.Stat(filePath); err == nil {
		if stat.Mode().IsRegular() {
			return true, nil
		} else {
			return false, errors.New(filePath + " is not a regular file")
		}
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

// StringFileContains checks if a content is already present inside a file.
// The content MUST be a FULL line (with strict line equality) or a set of full lines with a '\n' between them.
func StringFileContains(filePath string, content string) (bool, error) {
	filePath = strings.Replace(filePath, "~", HomeDir(), 1)

	fileContent, err := ReadFileAsString(filePath)
	if err != nil {
		return false, err
	}
	fileLines := strings.Split(fileContent, "\n")

	strLines := strings.Split(content, "\n")

	idx := 0
	for _, line := range fileLines {
		if line == strLines[idx] {
			idx++
			if idx == len(strLines) {
				return true, nil
			}
		} else {
			idx = 0
		}
	}

	return false, nil
}

// ReadFileAsStringOrEmptyIfNotExists gets the content of a text file.
// If the file does NOT exists, the method does NOT send an error but returns an empty string.
func ReadFileAsStringOrEmptyIfNotExists(filePath string) (string, error) {
	filePath = strings.Replace(filePath, "~", HomeDir(), 1)

	// Checks if the file exists
	if exists, err := FileExists(filePath); err != nil {
		return "", err
	} else if !exists {
		return "", nil
	}

	// If file exists
	if byteArray, err := ioutil.ReadFile(filePath); err != nil {
		return "", err
	} else {
		return string(byteArray), nil
	}
}

// WriteStringFile creates a file and writes the content of a string into it.
// if overwrite  == true, replaces the file content if the file exists.
func WriteStringFile(filePath, newContent string, overwrite bool) result.Result {
	fileName := strings.Replace(filePath, HomeDir(), "~", 1)
	filePath = strings.Replace(filePath, "~", HomeDir(), 1)

	// Check if file exists
	exists, errExists := FileExists(filePath)
	if errExists != nil {
		return result.Failure(errExists.Error())
	}

	// The file does NOT exist
	if !exists {
		if err := replaceStringFileContent(filePath, newContent); err != nil {
			return result.Failure(fileName + " writing error: " + err.Error())
		} else {
			return result.Success(fileName + " created")
		}
	}

	// The file exists
	curContent, err := ReadFileAsString(filePath)
	if err != nil {
		return result.Failure(fileName + " already exists but we can't read its content: " + err.Error())
	}
	if curContent == newContent {
		return result.Success(fileName + " already has expected content")
	}
	if overwrite {
		if err := replaceStringFileContent(filePath, newContent); err != nil {
			return result.Failure("Can't update " + fileName + ": " + err.Error())
		}
		return result.Success(fileName + " updated")
	}
	return result.Success(fileName + " user defined content left unchanged")
}

// WriteBinaryFile creates a file and writes the content of a byte slice into it.
// if overwrite  == true, replaces the file content if the file exists.
func WriteBinaryFile(filePath string, newContent []byte, writeIfFileExists bool) result.Result {
	fileName := strings.Replace(filePath, HomeDir(), "~", 1)
	filePath = strings.Replace(filePath, "~", HomeDir(), 1)

	// Check if file exists
	exists, errExists := FileExists(filePath)
	if errExists != nil {
		return result.Failure(fileName + " exists but is not a regular file")
	}

	// The file does NOT exist
	if !exists {
		if err := replaceBinaryFileContent(filePath, newContent); err != nil {
			return result.Failure(fileName + " writing error: " + err.Error())
		} else {
			return result.Success(fileName + " created")
		}
	}

	// The file exists
	curContent, err := ReadFileAsBinary(filePath)
	if err != nil {
		return result.Failure(fileName + " already exists but we can't read its content: " + err.Error())
	}
	if bytes.Equal(curContent, newContent) {
		return result.Success(fileName + " already has expected content")
	}
	if writeIfFileExists {
		if err := replaceBinaryFileContent(filePath, newContent); err != nil {
			return result.Failure("Can't update " + fileName + ": " + err.Error())
		}
		return result.Success(fileName + " updated")
	}
	return result.Success(fileName + " user defined content left unchanged")
}

func replaceStringFileContent(filePath, newContent string) error {
	return writeStringInFile(filePath, newContent, os.O_CREATE|os.O_TRUNC|os.O_WRONLY)
}

func replaceBinaryFileContent(filePath string, newContent []byte) error {
	return writeBinaryInFile(filePath, newContent, os.O_CREATE|os.O_TRUNC|os.O_WRONLY)
}

func writeStringInFile(filePath, fileText string, flag int) error {
	return writeBinaryInFile(filePath, []byte(fileText), flag)
}

func writeBinaryInFile(filePath string, content []byte, flag int) error {
	filePath = strings.Replace(filePath, "~", HomeDir(), 1)

	if err := os.MkdirAll(path.Dir(filePath), 0775); err != nil {
		return err
	}

	f, err := os.OpenFile(filePath, flag, 0644)
	if err != nil {
		return err
	}

	if _, err := f.Write(content); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}
