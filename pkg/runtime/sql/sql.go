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
var code = []byte(`-- Add your SQL queries here.
-- See SQL documentation: https://docs.airplane.dev/creating-tasks/sql
SELECT 1;
`)

// Runtime implementation.
type Runtime struct{}

// PrepareRun implementation.
func (r Runtime) PrepareRun(ctx context.Context, logger logger.Logger, opts runtime.PrepareRunOptions) (rexprs []string, rcloser io.Closer, rerr error) {
	return nil, nil, runtime.ErrNotImplemented
}

// Generate implementation.
func (r Runtime) Generate(t *runtime.Task, opts runtime.GenerateOpts) ([]byte, os.FileMode, error) {
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

// SupportsLocalExecution implementation.
func (r Runtime) SupportsLocalExecution() bool {
	return false
}
