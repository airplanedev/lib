package fixtures

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

// Path converts a path (relative to the examples/ folder) into an absolute path.
func Path(t *testing.T, relpath string) string {
	// Get the current filename so that we can import test data relative to the
	// current file (not the process's working directory).
	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok, "unable to access runtime information")
	dir := filepath.Dir(filename)
	return filepath.Join(dir, relpath)
}
