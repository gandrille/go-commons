package zipfile

import (
	"archive/zip"
	"bytes"
	"errors"
	"strings"
)

// ====================================
// Ease the usage of SMALL zip files
// Why only SMALL ones?
// Because the full content of the zip
// file is loaded into memory...
// ====================================

// ZipFile is the in memory representation of an existing zip file.
type ZipFile struct {
	Files []ZipElement
	Err   error
}

// Open loads a zip file into memory for easy usage.
// It does NOT requires closing it.
func Open(filePath string) ZipFile {
	var files []ZipElement

	// Open zip file
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return ZipFile{files, errors.New(filePath + " can't open file")}
	}
	defer r.Close()

	// Read files embedded into the zip file
	for _, f := range r.File {
		name := f.Name
		rc, err := f.Open()
		if err != nil {
			files = append(files, ZipElement{name, []byte{}, false})
		} else {
			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(rc)
			if err != nil {
				files = append(files, ZipElement{name, []byte{}, false})
			} else {
				files = append(files, ZipElement{name, buf.Bytes(), true})
			}
			rc.Close()
		}
	}

	// Build invalid files list
	var invalid []string
	for _, file := range files {
		if !file.IsValid() {
			invalid = append(invalid, file.name)
		}
	}

	// Return
	if len(invalid) == 0 {
		return ZipFile{files, nil}
	} else {
		return ZipFile{files, errors.New(filePath + " can't read embedded files: " + strings.Join(invalid, ","))}
	}
}

// IsValid checks if a zip file has been loaded into memory.
// It does NOT implies that all elements are in a valid state.
func (zip ZipFile) IsValid() bool {
	return zip.Err == nil
}

// HasFile checks if a zip element is part of a zipfile.
func (zip ZipFile) HasFile(name string) bool {
	return zip.GetFile(name) != nil
}

// GetFile returns a zip element matching a given name, or nil if NOT found.
func (zip ZipFile) GetFile(name string) *ZipElement {
	for _, file := range zip.Files {
		if file.name == name {
			return &file
		}
	}
	return nil
}

// FilesStartingWith returns a slice of zip element with their names starting with a given prefix.
func (zip ZipFile) FilesStartingWith(prefix string) []ZipElement {
	var list []ZipElement
	for _, file := range zip.Files {
		if strings.HasPrefix(file.name, prefix) {
			list = append(list, file)
		}
	}
	return list
}
