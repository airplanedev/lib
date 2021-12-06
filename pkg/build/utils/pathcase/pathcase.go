package pathcase

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// ActualFilename attempts to return the filename but with correct casing, based on the current
// filesystem. This is needed because case-insensitive FS (e.g. APFS on macOS) will happily read
// `myentrypoint.js` even if the file is named `myEntrypoint.js`, but when we tar and ship this
// to an ext4 (Linux) FS, we error when trying to open `myentrypoint.js` (it's case-sensitive on
// ext4 and `myentrypoint.js` does not exist.)
// Note this only fixes the file part of the path - directories in the path are returned as is,
// even if their casing is incorrect. Further work would be needed to support recursively fixing
// all parent directories.
func ActualFilename(filename string) (string, error) {
	dir := filepath.Dir(filename)
	files, readErr := os.ReadDir(dir)
	if readErr != nil {
		// Going to iterate through what files were found anyways
	}

	for _, file := range files {
		if strings.EqualFold(file.Name(), filepath.Base(filename)) {
			return filepath.Join(dir, file.Name()), nil
		}
	}

	// If we haven't found the file, return the error (if any) from reading the directory.
	return "", errors.Wrap(readErr, "listing directory")
}
