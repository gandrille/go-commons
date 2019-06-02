package zipfile

import (
	"errors"

	"github.com/gandrille/go-commons/filesystem"
	"github.com/gandrille/go-commons/result"
)

// IMPORTANT! READ ME FIRST!
// All the functions in this file are designed for providing nice status messages,
// NOT for efficiency optimization.
// Do NOT use this functions if you need performance.

// ZipElement is the in memory representation of a file inside a zip file.
type ZipElement struct {
	name    string
	content []byte
	isValid bool
}

// Name gets the name (zip relative path) of this element.
func (element ZipElement) Name() string {
	return element.name
}

// IsValid checks if a ZipElement has been loaded corectly in memory,
// and therefore if we can read its content.
func (element ZipElement) IsValid() bool {
	return element.isValid
}

// Write dumps the content of the element of a zip file to the filesystem.
func (element ZipElement) Write(filePath string) result.Result {
	if !element.isValid {
		return result.NewError(element.name + " has not been loaded corectly")
	}

	return filesystem.WriteBinaryFile(filePath, element.content, true)
}

// BytesContent gets the contant as a slice of bytes
func (element ZipElement) BytesContent() ([]byte, error) {
	if element.isValid {
		return element.content, nil
	}
	return nil, errors.New(element.name + " is invalid")
}

// StringContent gets the contant as a string
func (element ZipElement) StringContent() (string, error) {
	if element.isValid {
		return string(element.content), nil
	}
	return "", errors.New(element.name + " is invalid")
}
