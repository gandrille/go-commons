package filesystem

import (
	"os"
	"strings"

	"github.com/gandrille/go-commons/result"
)

// IMPORTANT! READ ME FIRST!
// All the functions in this file are designed for providing nice status messages,
// NOT for efficiency optimization.
// Do NOT use this functions if you need performance.

// CreateOrAppendIfNotInFile checks if a content is already present inside a file.
// The content MUST be a FULL line (strict line equality) or a set of full lines with a '\n' between them.
// If the content is already present in the file, the file remains unchanged and the function returns false.
// Otherwise, the content is appended at the end of the file, and the function returns true.
func CreateOrAppendIfNotInFile(filePath, content string) (bool, error) {
	filePath = strings.Replace(filePath, "~", HomeDir(), 1)

	// Checks if the file exists
	exists, err := RegularFileExists(filePath)
	if err != nil {
		return false, err
	}

	// The file exists
	if exists {
		contains, err1 := StringFileContains(filePath, content)
		if err1 != nil {
			return false, err1
		}
		if contains {
			return false, nil
		}

		err2 := createOrAppendInFile(filePath, content)
		if err2 != nil {
			return false, err2
		}
		return true, nil
	}

	if err := createOrAppendInFile(filePath, content); err != nil {
		return false, err
	}
	return true, nil
}

func createOrAppendInFile(filePath, fileText string) error {
	return writeStringInFile(filePath, fileText, os.O_CREATE|os.O_APPEND|os.O_WRONLY)
}

// CopyFileWithUpdate copies srcFile to dstFile replacing all the lines starting with startwith by replacement.
func CopyFileWithUpdate(srcFile, dstFile, startwith, replacement string, appendIfNoMatch bool) result.Result {

	originalSrcContentStr, err1 := ReadFileAsString(srcFile)
	if err1 != nil {
		return result.NewError(err1.Error())
	}

	originalDstContentStr, err2 := ReadFileAsStringOrEmptyIfNotExists(dstFile)
	if err2 != nil {
		return result.NewError(err2.Error())
	}

	finalContent := updateLine(originalSrcContentStr, startwith, replacement, appendIfNoMatch)

	// Dst file has expected content
	if originalDstContentStr == finalContent {
		if originalSrcContentStr == originalDstContentStr {
			return result.NewUnchanged("Content of " + dstFile + " is the same as " + srcFile + " and we have no line strating with " + startwith)
		} else {
			return result.NewUnchanged("Content of " + dstFile + " is already the same as " + srcFile + " with lines strating with " + startwith + " updated")
		}
	}

	// update needed
	if res := WriteStringFile(dstFile, finalContent, true); !res.IsSuccess() {
		return res
	}

	return result.NewUpdated("Content of " + dstFile + " written with content from " + srcFile + " with lines strating with " + startwith + " updated")
}

// UpdateLineInFile replaces all the lines of filepath which are starting with startwith by replacement.
// if appendIfNoMatch == true, appends the replacement at the end of the file if no replacement have been made before.
func UpdateLineInFile(filePath, startwith, replacement string, appendIfNoMatch bool) result.Result {

	originalContent, err1 := ReadFileAsString(filePath)
	if err1 != nil {
		return result.NewError(err1.Error())
	}

	finalContent := updateLine(originalContent, startwith, replacement, appendIfNoMatch)

	// Nothing to do
	if originalContent == finalContent {
		return result.NewUnchanged(filePath + " lines strating with " + startwith + " already updated")
	}

	// update needed
	if res := WriteStringFile(filePath, finalContent, true); !res.IsSuccess() {
		return res
	}

	return result.NewUpdated(filePath + " lines strating with " + startwith + " updated")
}

// RemoveLineInFile removes all the lines of filepath which are starting with startwith.
// if fullline == true, an exact match (instead of start with) is required to remove the line.
func RemoveLineInFile(filePath, startwith string, fullline bool) result.Result {

	originalContent, err1 := ReadFileAsString(filePath)
	if err1 != nil {
		return result.NewError(err1.Error())
	}

	finalContent := removeLine(originalContent, startwith, fullline)

	// Nothing to do
	if originalContent == finalContent {
		if fullline {
			return result.NewUnchanged(filePath + " does NOT have lines equal to " + startwith)
		} else {
			return result.NewUnchanged(filePath + " does NOT have lines starting with " + startwith)
		}
	}

	// update needed
	if res := WriteStringFile(filePath, finalContent, true); !res.IsSuccess() {
		return res
	}

	return result.NewUpdated(filePath + " lines strating with " + startwith + " updated")
}

func updateLine(content, startwith, replacement string, appendIfNoMatch bool) string {

	contentSlice := strings.Split(content, "\n")

	found := false
	for i, line := range contentSlice {
		if strings.HasPrefix(line, startwith) {
			contentSlice[i] = replacement
			found = true
		}
	}
	if appendIfNoMatch && !found {
		contentSlice = append(contentSlice, replacement)
	}

	return strings.Join(contentSlice, "\n")
}

func removeLine(content, startwith string, fullline bool) string {

	contentSlice := strings.Split(content, "\n")
	finalContent := []string{}

	for _, line := range contentSlice {
		if line != startwith && (!strings.HasPrefix(line, startwith) || fullline) {
			finalContent = append(finalContent, line)
		}
	}

	return strings.Join(finalContent, "\n")
}
