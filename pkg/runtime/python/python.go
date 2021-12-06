package python

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/airplanedev/lib/pkg/build"

	// "github.com/airplanedev/lib/pkg/logger"
	"github.com/airplanedev/lib/pkg/runtime"
	"github.com/airplanedev/lib/pkg/utils/fsx"
	"github.com/pkg/errors"
)

// Init register the runtime.
func init() {
	runtime.Register(".py", Runtime{})
}

// Code template.
var code = template.Must(template.New("py").Parse(`{{with .Comment -}}
{{.}}

{{end -}}
# Put the main logic of the task in the main function.
def main(params):
    print("parameters:", params)

    # You can return data to show outputs to users.
    # Outputs documentation: https://docs.airplane.dev/tasks/outputs
    return [
        {"element": "hydrogen", "weight": 1.008},
        {"element": "helium", "weight": 4.0026}
    ]
`))

// Data represents the data template.
type data struct {
	Comment string
}

// Runtime implementation.
type Runtime struct{}

// PrepareRun implementation.
func (r Runtime) PrepareRun(ctx context.Context, logger runtime.Logger, opts runtime.PrepareRunOptions) (rexprs []string, rcloser io.Closer, rerr error) {
	if err := checkPythonInstalled(ctx, logger); err != nil {
		return nil, nil, err
	}

	root, err := r.Root(opts.Path)
	if err != nil {
		return nil, nil, err
	}

	tmpdir := filepath.Join(root, ".airplane")
	if err := os.Mkdir(tmpdir, os.ModeDir|0777); err != nil && !os.IsExist(err) {
		return nil, nil, errors.Wrap(err, "creating .airplane directory")
	}
	closer := runtime.CloseFunc(func() error {
		logger.Debug("Cleaning up temporary directory...")
		return errors.Wrap(os.RemoveAll(tmpdir), "unable to remove temporary directory")
	})
	defer func() {
		// If we encountered an error before returning, then we're responsible
		// for performing our own cleanup.
		if rerr != nil {
			closer.Close()
		}
	}()

	entrypoint, err := filepath.Rel(root, opts.Path)
	if err != nil {
		return nil, nil, errors.Wrap(err, "entrypoint is not within the task root")
	}
	shim, err := build.PythonShim(root, entrypoint)
	if err != nil {
		return nil, nil, err
	}

	if err := os.WriteFile(filepath.Join(tmpdir, "shim.py"), []byte(shim), 0644); err != nil {
		return nil, nil, errors.Wrap(err, "writing shim file")
	}

	pv, err := json.Marshal(opts.ParamValues)
	if err != nil {
		return nil, nil, errors.Wrap(err, "serializing param values")
	}

	return []string{"python3", filepath.Join(tmpdir, "shim.py"), string(pv)}, closer, nil
}

// Checks for python3 binary, as per PEP 0394:
// https://www.python.org/dev/peps/pep-0394/#recommendation
func checkPythonInstalled(ctx context.Context, logger runtime.Logger) error {
	cmd := exec.CommandContext(ctx, "python3", "--version")
	logger.Debug("Running %s", strings.Join(cmd.Args, " "))
	err := cmd.Run()
	if err != nil {
		return errors.New(heredoc.Doc(`
		It looks like the python3 command is not installed.

		Ensure Python 3 is installed and the python3 command exists: https://www.python.org/downloads
	`))
	}
	return nil
}

// Generate implementation.
func (r Runtime) Generate(t *runtime.Task) ([]byte, fs.FileMode, error) {
	d := data{}
	if t != nil {
		d.Comment = runtime.Comment(r, *t)
	}

	var buf bytes.Buffer
	if err := code.Execute(&buf, d); err != nil {
		return nil, 0, fmt.Errorf("python: template execute - %w", err)
	}

	return buf.Bytes(), 0644, nil
}

// Workdir implementation.
func (r Runtime) Workdir(path string) (string, error) {
	return r.Root(path)
}

// Root implementation.
func (r Runtime) Root(path string) (string, error) {
	root, ok := fsx.Find(path, "requirements.txt")
	if !ok {
		return filepath.Dir(path), nil
	}
	return root, nil
}

// Kind implementation.
func (r Runtime) Kind() build.TaskKind {
	return build.TaskKindPython
}

// FormatComment implementation.
func (r Runtime) FormatComment(s string) string {
	var lines []string

	for _, line := range strings.Split(s, "\n") {
		lines = append(lines, "# "+line)
	}

	return strings.Join(lines, "\n")
}
