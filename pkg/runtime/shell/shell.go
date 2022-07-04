package shell

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/airplanedev/lib/pkg/build"
	"github.com/airplanedev/lib/pkg/runtime"
	"github.com/airplanedev/lib/pkg/utils"
	"github.com/airplanedev/lib/pkg/utils/fsx"
	"github.com/airplanedev/lib/pkg/utils/handlebars"
	"github.com/airplanedev/lib/pkg/utils/logger"
	"github.com/pkg/errors"
)

// Init register the runtime.
func init() {
	runtime.Register(".sh", Runtime{})
}

// Code template.
var code = template.Must(template.New("sh").Parse(`#!/bin/bash
{{with .Comment -}}
{{.}}

{{end -}}
# Params are in environment variables as PARAM_{SLUG}, e.g. PARAM_USER_ID
echo "Hello World!"
echo "Printing env for debugging purposes:"
env
`))

// Data represents the data template.
type data struct {
	Comment string
}

// Runtime implementation.
type Runtime struct{}

// PrepareRun implementation.
func (r Runtime) PrepareRun(ctx context.Context, logger logger.Logger, opts runtime.PrepareRunOptions) (rexprs []string, rcloser io.Closer, rerr error) {
	if err := checkAndPromptFileExecutable(opts.Path); err != nil {
		return nil, nil, err
	}

	root, err := r.Root(opts.Path)
	if err != nil {
		return nil, nil, err
	}

	if dockerfilePath := build.FindDockerfile(root); dockerfilePath != "" {
		logger.Warning("Found Dockerfile at %s.", dockerfilePath)
		logger.Warning("`airplane dev` does not currently support running inside a Docker image.")
		logger.Warning("The script will run inside your local machine environment.")
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

	shim := build.ShellShim()
	if err := os.WriteFile(filepath.Join(tmpdir, "shim.sh"), []byte(shim), 0644); err != nil {
		return nil, nil, errors.Wrap(err, "writing shim file")
	}

	entrypoint, err := filepath.Rel(root, opts.Path)
	if err != nil {
		return nil, nil, errors.Wrap(err, "entrypoint is not within the task root")
	}

	cmd := []string{
		"bash", filepath.Join(tmpdir, "shim.sh"),
		filepath.Join(root, entrypoint),
	}
	// TODO: this is a rough approximation of how interpolateParameters works in prod
	for slug := range opts.ParamValues {
		tmpl := fmt.Sprintf("%s={{%s}}", slug, slug)
		val, err := handlebars.Render(tmpl, opts.ParamValues)
		if err != nil {
			return nil, nil, errors.Wrap(err, "rendering shell command")
		}
		cmd = append(cmd, val)
	}
	return cmd, closer, nil
}

// Generate implementation.
func (r Runtime) Generate(t *runtime.Task) ([]byte, os.FileMode, error) {
	d := data{}
	if t != nil && t.URL != "" {
		d.Comment = runtime.Comment(r, t.URL)
	}

	var buf bytes.Buffer
	if err := code.Execute(&buf, d); err != nil {
		return nil, 0, fmt.Errorf("shell: template execute - %w", err)
	}

	// 0744 has +x to execute shell scripts.
	return buf.Bytes(), 0744, nil
}

// Workdir implementation.
func (r Runtime) Workdir(path string) (string, error) {
	return r.Root(path)
}

// Root implementation.
func (r Runtime) Root(path string) (string, error) {
	for _, filePath := range build.DockerfilePaths() {
		if root, ok := fsx.Find(path, filePath); ok {
			return root, nil
		}
	}
	return filepath.Dir(path), nil
}

// Kind implementation.
func (r Runtime) Kind() build.TaskKind {
	return build.TaskKindShell
}

// FormatComment implementation.
func (r Runtime) FormatComment(s string) string {
	var lines []string

	for _, line := range strings.Split(s, "\n") {
		lines = append(lines, "# "+line)
	}

	return strings.Join(lines, "\n")
}

// SupportsLocalExecution implementation.
func (r Runtime) SupportsLocalExecution() bool {
	return true
}

// checkAndPromptFileExecutable checks that a file is executable. If it isn't, it prompts the user to make it
// executable. Returns an error if the file is not executable.
func checkAndPromptFileExecutable(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return errors.Wrap(err, "describing file")
	}
	if isExecOwner(info.Mode()) {
		return nil
	}
	var allow bool
	readablePath := path
	if utils.CanPrompt() {
		wd, err := os.Getwd()
		if err == nil {
			relPath, err := filepath.Rel(wd, path)
			if err == nil {
				readablePath = relPath
			}
		}
		allow, _ = utils.Confirm(fmt.Sprintf("File %s is not executable. Do you want to make your file executable?", readablePath))
	}
	if !allow {
		return errors.Errorf("File %s is not executable. Run `chmod +x %s` to make it executable.", readablePath, readablePath)
	}
	if err = os.Chmod(path, info.Mode()|0100); err != nil {
		return errors.Wrap(err, "making file executable")
	}
	return nil
}

// isExecOwner returns whether a file is executable by its owner.
func isExecOwner(mode os.FileMode) bool {
	return mode&0100 != 0
}
