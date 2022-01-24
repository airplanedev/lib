package sql

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/airplanedev/lib/pkg/build"
	"github.com/airplanedev/lib/pkg/runtime"
	"github.com/airplanedev/lib/pkg/utils/logger"
)

// Init register the runtime.
func init() {
	runtime.Register(".sql", Runtime{})
}

// Code.
var code = []byte(`-- Add your SQL code here:
SELECT 1;

-- Define SQL parameters in your definition file:
-- SELECT * from users where user_id = :user_id;
`)

// Runtime implementation.
type Runtime struct{}

// PrepareRun implementation.
func (r Runtime) PrepareRun(ctx context.Context, logger logger.Logger, opts runtime.PrepareRunOptions) (rexprs []string, rcloser io.Closer, rerr error) {
	return nil, nil, runtime.ErrNotImplemented
}

// Generate implementation.
func (r Runtime) Generate(t *runtime.Task) ([]byte, os.FileMode, error) {
	return code, 0644, nil
}

// Workdir implementation.
func (r Runtime) Workdir(path string) (string, error) {
	return r.Root(path)
}

// Root implementation.
func (r Runtime) Root(path string) (string, error) {
	return filepath.Dir(path), nil
}

// Kind implementation.
func (r Runtime) Kind() build.TaskKind {
	return build.TaskKindSQL
}

// FormatComment implementation.
func (r Runtime) FormatComment(s string) string {
	var lines []string

	for _, line := range strings.Split(s, "\n") {
		lines = append(lines, "-- "+line)
	}

	return strings.Join(lines, "\n")
}
