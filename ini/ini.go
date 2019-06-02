package ini

import (
	"errors"
	"strings"

	"github.com/gandrille/go-commons/filesystem"
)

// IMPORTANT! READ ME FIRST!
// All the functions in this file are designed for providing nice status messages,
// NOT for efficiency optimization.
// Do NOT use this functions if you need performance.

// GetValue gets a value from an ini file.
func GetValue(file, section, key string, failIfFileNotExists bool, failIfKeyNotExists bool) (string, error) {

	// Check is file exists
	if exists, err := filesystem.RegularFileExists(file); err != nil {
		return "", errors.New("Error while checking if the file " + file + " exists: " + err.Error())
	} else if !exists && failIfFileNotExists {
		return "", errors.New("The file " + file + " does NOT exist")
	}

	// Get file content
	content, err := filesystem.ReadFileAsStringOrEmptyIfNotExists(file)
	if err != nil {
		return "", errors.New("Error while reading file " + file + " content: " + err.Error())
	}

	// Find key
	lines := strings.Split(content, "\n")
	if idx, keyVal := findKey(lines, section, key); idx != -1 {
		return keyVal, nil
	}

	return notFound(file, section, key, failIfKeyNotExists)
}

// SetValue sets a value in an ini file.
// with enclosewithspaces == true, the line is written with space around the '=' sign
// returns true if the file has been created or modified
func SetValue(file, section, key, value string, failIfFileNotExists, enclosewithspaces bool) (bool, error) {

	// Compute line
	var newline string
	if enclosewithspaces {
		newline = key + " = " + value
	} else {
		newline = key + "=" + value
	}

	// Check is file exists
	if exists, err := filesystem.RegularFileExists(file); err != nil {
		return false, errors.New("Error while checking if the file " + file + " exists: " + err.Error())
	} else if !exists {
		if failIfFileNotExists {
			return false, errors.New("The file " + file + " does NOT exist")
		} else {
			fileContent := "[" + section + "]\n" + newline + "\n"
			res := filesystem.WriteStringFile(file, fileContent, true)
			if res.IsSuccess() {
				return true, nil
			} else {
				return false, errors.New(res.Message())
			}
		}
	}

	// Get file content
	content, err := filesystem.ReadFileAsStringOrEmptyIfNotExists(file)
	if err != nil {
		return false, errors.New("Error while reading file " + file + " content: " + err.Error())
	}

	// Find key
	lines := strings.Split(content, "\n")
	idx, curVal := findKey(lines, section, key)

	// updates content
	if idx == -1 {
		secIdx := getLineNumberOfSection(lines, section)
		if secIdx == -1 {
			lines = append(lines, "["+section+"]", newline)
		} else {
			lines = append(lines[:secIdx], append([]string{newline}, lines[secIdx:]...)...)
		}
	} else if idx == len(lines) {
		lines = append(lines, newline)
	} else {
		if curVal == value {
			// No need to update the file
			return false, nil
		} else {
			lines[idx] = newline
		}
	}

	// write content back
	if res := filesystem.WriteStringFile(file, strings.Join(lines, "\n"), true); res.IsFailure() {
		return false, errors.New("Error while writing file " + file + " with updated content: " + res.Message())
	} else {
		return true, nil
	}
}

// RemoveValue removes a value in an ini file.
// returns true if the key was found and removed
func RemoveValue(file, section, key string) (bool, error) {

	// Check is file exists
	if exists, err := filesystem.RegularFileExists(file); err != nil {
		return false, errors.New("Error while checking if the file " + file + " exists: " + err.Error())
	} else if !exists {
		return false, errors.New("The file " + file + " does NOT exist")
	}

	// Get file content
	content, err := filesystem.ReadFileAsStringOrEmptyIfNotExists(file)
	if err != nil {
		return false, errors.New("Error while reading file " + file + " content: " + err.Error())
	}

	// Find key
	lines := strings.Split(content, "\n")
	if idx, _ := findKey(lines, section, key); idx == -1 {
		return false, nil
	} else {
		lines = append(lines[:idx], lines[idx+1:]...)
	}

	// write content back
	if res := filesystem.WriteStringFile(file, strings.Join(lines, "\n"), true); res.IsFailure() {
		return false, errors.New("Error while writing file " + file + " with updated content: " + res.Message())
	} else {
		return true, nil
	}
}

func findKey(lines []string, sectionName, keyName string) (int, string) {
	idx := getLineNumberOfSection(lines, sectionName)

	if idx == -1 {
		return -1, ""
	}

	for i := idx; i < len(lines); i++ {
		line := lines[i]

		if isSec, _ := isSectionLine(line); isSec {
			return -1, ""
		}

		if isVal, curkey, curval := isKeyValLine(line); isVal {
			if keyName == curkey {
				return i, curval
			}
		}
	}

	return -1, ""
}

func getLineNumberOfSection(lines []string, sectionName string) int {
	sectionName = strings.Trim(sectionName, " ")
	if len(sectionName) == 0 {
		return 0
	}

	for i, line := range lines {
		if isSec, secName := isSectionLine(line); isSec {
			if sectionName == secName {
				return i + 1
			}
		}
	}

	return -1
}

func isSectionLine(line string) (bool, string) {
	line = strings.Trim(line, " ")
	if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
		sectionName := strings.Trim(line[1:len(line)-1], " ")
		return true, sectionName
	}
	return false, ""
}

func isKeyValLine(line string) (bool, string, string) {
	idx := strings.Index(line, "=")
	if idx != -1 {
		curKey := strings.Trim(line[:idx], " ")
		curVal := strings.Trim(line[idx+1:], " ")
		if curKey != "" {
			return true, curKey, curVal
		}
	}
	return false, "", ""
}

func notFound(file, section, key string, failIfKeyNotExists bool) (string, error) {
	if failIfKeyNotExists {
		return "", errors.New("Key '" + key + "' not found in section '" + section + "' of file " + file)
	} else {
		return "", nil
	}
}
