// fsx includes extensions to the stdlib fs package.
package fsx

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Exists returns true if the given path exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// FilesExist ensures that all paths exists or returns an error.
func AssertExistsAll(paths ...string) error {
	for _, p := range paths {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			return fmt.Errorf("could not find file %s", path.Base(p))
		}
	}
	return nil
}

// FilesExist ensures that any of the paths exists or returns an error.
func AssertExistsAny(paths ...string) error {
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return nil
		} else if os.IsNotExist(err) {
			continue
		} else {
			return err
		}
	}
	return fmt.Errorf("could not find any files %s", paths)
}

// Find attempts to find the path of the given filename.
//
// The method recursively visits parent dirs until the given
// filename is found, If the file is not found the method
// returns false.
//
// Continues recursively until the root directory is reached.
func Find(dir, filename string) (string, bool) {
	return FindUntil(dir, "", filename)
}

// Find attempts to find the path of the given filename.
//
// The method recursively visits parent dirs until the given
// filename is found, If the file is not found the method
// returns false.
//
// Continues until the `end` directory is reached (inclusively).
// If `end` is an empty string, continues until the root directory.
func FindUntil(start, end, filename string) (string, bool) {
	dst := filepath.Join(start, filename)

	if !Exists(dst) {
		next := filepath.Dir(start)
		if next == start || next == "." || (end != "" && strings.HasPrefix(end, next)) || next == string(filepath.Separator) {
			return "", false
		}
		return FindUntil(next, end, filename)
	}

	return start, true
}
