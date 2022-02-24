// Package runtime generates code to match a runtime.
//
// The runtime package is capable of writing airplane specific
// comments that are used to link a task file to a remote task.
//
// All runtimes are also capable of generating initial code to
// match the task, including the parameters.
package runtime

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/airplanedev/lib/pkg/build"
	"github.com/airplanedev/lib/pkg/utils/logger"
	"github.com/pkg/errors"
)

var (
	// ErrMissing is returned when a resource was not found.
	//
	// It can be checked via `errors.Is(err, ErrMissing)`.
	ErrMissing = errors.New("runtime: resource is missing")

	// ErrNotImplemented is returned when a runtime does not
	// support preparing a run.
	//
	// It can be checked via `errors.Is(err, ErrNotImplemented)`.
	ErrNotImplemented = errors.New("runtime: not implemented")
)

// Interface repersents a runtime.
type Interface interface {
	// Generate accepts a task and generates code to match the task.
	//
	// os.FileMode is used for the permissions of the generated file. Files will typically use 0644
	// but might use 0744 for executable scripts (e.g. shell scripts).
	Generate(task *Task) ([]byte, os.FileMode, error)

	// Workdir attempts to detect the root of the given task path.
	//
	// Unlike root it decides the dockerfile's `workdir` directive
	// this might be different than root because it decides where
	// the build commands are run.
	Workdir(path string) (dir string, err error)

	// Root attempts to detect the root of the given task path.
	//
	// Typically runtimes will look for a specific file such as
	// `package.json` or `requirements.txt`, they'll use `fs.Find()`.
	Root(path string) (dir string, err error)

	// Kind returns a task kind that matches the runtime.
	//
	// Generate and other methods should not be called
	// for a task that doesn't match the returned kind.
	Kind() build.TaskKind

	// FormatComment formats a string into a comment using
	// the relevant comment characters for this runtime.
	FormatComment(s string) string

	// PrepareRun should prepare a local run of a task.
	//
	// It must create a temporary directory, install any dependencies
	// and prepare the script to be run.
	//
	// On success the method returns a slice that represents an `cmd.Exec`
	// options which contains the command to be run and its arguments. It
	// should also return an io.Closer that will be called when the local
	// run has completed. If not needed, a nil closer can be returned. If
	// an error is encountered during PrepareRun, the runtime should perform
	// its own cleanup.
	//
	// If running the script locally is not supported the method returns
	// an `ErrNotImplemented`.
	PrepareRun(ctx context.Context, logger logger.Logger, opts PrepareRunOptions) (rexprs []string, closer io.Closer, err error)

	// SupportsLocalExecution returns true if local execution is supported.
	// This is expected to match whether PrepareRun returns `ErrNotImplemented`.
	SupportsLocalExecution() bool
}

type PrepareRunOptions struct {
	// Path is the file path leading to the task's entrypoint.
	//
	// It should be an absolute path.
	Path string

	// ParamValues specifies the user-provided parameter values to
	// execute this run with.
	ParamValues Values

	// KindOptions specifies any runtime-specific task configuration.
	KindOptions build.KindOptions
}

// Runtimes is a collection of registered runtimes.
//
// The key is the file extension used for the runtime.
var runtimes = make(map[string]Interface)

// Register registers the given ext with r.
func Register(ext string, r Interface) {
	if _, ok := runtimes[ext]; ok {
		panic(fmt.Sprintf("runtime: %s already registered", ext))
	}
	runtimes[ext] = r
}

// Lookup returns a runtime by kind and path.
// If an extension match is found, use that runtime. Otherwise rely on the task kind.
func Lookup(path string, kind build.TaskKind) (Interface, error) {
	ext := filepath.Ext(path)
	if runtime, ok := runtimes[ext]; ok {
		return runtime, nil
	}

	// There was no exact match on the extension. Fallback to checking if there
	// is exactly one match on the task kind, which can occur for task kinds that
	// support arbitrary extensions (f.e. shell).
	possible := []Interface{}
	for _, runtime := range runtimes {
		if runtime.Kind() == kind {
			possible = append(possible, runtime)
		}
	}
	if len(possible) > 1 {
		return nil, errors.Errorf("found %d runtimes for task type, expecting 1", len(possible))
	}
	if len(possible) == 0 {
		return nil, errors.New("did not find any runtimes for task type")
	}
	return possible[0], nil
}

// SuggestExt returns the default extension for a given TaskKind, if any.
func SuggestExt(kind build.TaskKind) string {
	for ext, runtime := range runtimes {
		if runtime.Kind() == kind {
			return ext
		}
	}
	return ""
}

func SuggestKind(ext string) (build.TaskKind, error) {
	if runtime, ok := runtimes[ext]; ok {
		return runtime.Kind(), nil
	}
	return "", errors.New("No kind to suggest")
}

func CloseFunc(f func() error) io.Closer {
	return closer{f: f}
}

type closer struct {
	f func() error
}

func (c closer) Close() error {
	return c.f()
}
