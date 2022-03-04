package javascript

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/airplanedev/lib/pkg/build"

	"github.com/airplanedev/lib/pkg/runtime"
	"github.com/airplanedev/lib/pkg/utils/fsx"
	"github.com/airplanedev/lib/pkg/utils/logger"
	"github.com/blang/semver/v4"
	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/pkg/errors"
)

// Init register the runtime.
func init() {
	runtime.Register(".js", Runtime{})
}

// Code template.
var code = template.Must(template.New("js").Parse(`{{with .Comment -}}
{{.}}

{{end -}}
// Put the main logic of the task in this function.
export default async function(params) {
  console.log('parameters:', params);

  // You can return data to show outputs to users.
  // Outputs documentation: https://docs.airplane.dev/tasks/outputs
  return [
    {element: 'hydrogen', weight: 1.008},
    {element: 'helium', weight: 4.0026},
  ];
}
`))

// Data represents the data template.
type data struct {
	Comment string
}

// Runtime implementaton.
type Runtime struct{}

// Generate implementation.
func (r Runtime) Generate(t *runtime.Task) ([]byte, fs.FileMode, error) {
	d := data{}
	if t != nil {
		d.Comment = runtime.Comment(r, t.URL)
	}

	var buf bytes.Buffer
	if err := code.Execute(&buf, d); err != nil {
		return nil, 0, fmt.Errorf("javascript: template execute - %w", err)
	}

	return buf.Bytes(), 0644, nil
}

// Workdir picks the working directory for commands to be executed from.
//
// For JS, that is the nearest parent directory containing a `package.json`.
func (r Runtime) Workdir(path string) (string, error) {
	if p, ok := fsx.Find(path, "package.json"); ok {
		return p, nil
	}

	// Otherwise default to immediate directory of path
	return filepath.Dir(path), nil
}

// Root picks which directory to use as the root of a task's code.
// All code in that directory will be available at runtime.
//
// For JS, this is usually just the workdir. However, this can be overridden
// with the `airplane.root` property in the `package.json`.
func (r Runtime) Root(path string) (string, error) {
	// By default, the root is the workdir.
	root, err := r.Workdir(path)
	if err != nil {
		return "", err
	}

	// Unless the root is overridden with an `airplane.root` field
	// in a `package.json`.
	pkgjson := filepath.Join(root, "package.json")
	buf, err := os.ReadFile(pkgjson)
	if err != nil {
		// No package.json, use workdir as root.
		if os.IsNotExist(err) {
			return root, nil
		}
		return "", errors.Wrapf(err, "javascript: reading %s", pkgjson)
	}

	var pkg struct {
		Settings build.Settings `json:"airplane"`
	}

	if err := json.Unmarshal(buf, &pkg); err != nil {
		return "", fmt.Errorf("javascript: reading %s - %w", root, err)
	}

	if pkgjsonRoot := pkg.Settings.Root; pkgjsonRoot != "" {
		return filepath.Join(root, pkgjsonRoot), nil
	}

	return root, nil
}

// Kind implementation.
func (r Runtime) Kind() build.TaskKind {
	return build.TaskKindNode
}

func (r Runtime) FormatComment(s string) string {
	lines := []string{}
	for _, line := range strings.Split(s, "\n") {
		lines = append(lines, "// "+line)
	}
	return strings.Join(lines, "\n")
}

func (r Runtime) PrepareRun(ctx context.Context, logger logger.Logger, opts runtime.PrepareRunOptions) (rexprs []string, rcloser io.Closer, rerr error) {
	checkNodeVersion(ctx, logger, opts.KindOptions)

	root, err := r.Root(opts.Path)
	if err != nil {
		return nil, nil, err
	}
	workdir, err := r.Workdir(opts.Path)
	if err != nil {
		return nil, nil, err
	}

	tmpdir := filepath.Join(root, ".airplane")
	if err := os.Mkdir(tmpdir, os.ModeDir|0777); err != nil && !os.IsExist(err) {
		return nil, nil, errors.Wrap(err, "creating .airplane directory")
	}
	closer := runtime.CloseFunc(func() error {
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
	shim, err := build.NodeShim(entrypoint)
	if err != nil {
		return nil, nil, err
	}

	if err := os.WriteFile(filepath.Join(tmpdir, "shim.js"), []byte(shim), 0644); err != nil {
		return nil, nil, errors.Wrap(err, "writing shim file")
	}

	// Install the dependencies we need for our shim file:
	pjson, err := build.GenShimPackageJSON()
	if err != nil {
		return nil, nil, err
	}
	if err := os.WriteFile(filepath.Join(tmpdir, "package.json"), pjson, 0644); err != nil {
		return nil, nil, errors.Wrap(err, "writing shim package.json")
	}
	cmd := exec.CommandContext(ctx, "npm", "install")
	cmd.Dir = filepath.Join(root, ".airplane")
	logger.Debug("Running %s (in %s)", strings.Join(cmd.Args, " "), root)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log(strings.TrimSpace(string(out)))
		return nil, nil, errors.New("failed to install shim deps")
	}

	if err := os.RemoveAll(filepath.Join(tmpdir, "dist")); err != nil {
		return nil, nil, errors.Wrap(err, "cleaning dist folder")
	}

	// Confirm we have a `package.json`, otherwise we might install shim dependencies
	// in the wrong folder.
	pathPkgJSON := filepath.Join(workdir, "package.json")
	hasPkgJSON := fsx.AssertExistsAll(pathPkgJSON) == nil
	if !hasPkgJSON {
		return nil, nil, errors.New("a package.json is missing")
	}
	// Workaround to get esbuild to not bundle dependencies.
	// See build.ExternalPackages for details.
	externalDeps, err := build.ExternalPackages(pathPkgJSON)
	if err != nil {
		return nil, nil, err
	}

	start := time.Now()
	res := esbuild.Build(esbuild.BuildOptions{
		Bundle: true,

		EntryPoints: []string{filepath.Join(tmpdir, "shim.js")},
		Outfile:     filepath.Join(tmpdir, "dist/shim.js"),
		Write:       true,

		External: externalDeps,
		Platform: esbuild.PlatformNode,
		Engines: []esbuild.Engine{
			// esbuild is relatively generous in the node versions it supports:
			// https://esbuild.github.io/api/#target
			{Name: esbuild.EngineNode, Version: build.GetNodeVersion(opts.KindOptions)},
		},
	})
	for _, w := range res.Warnings {
		logger.Debug("esbuild(warn): %v", w)
	}
	for _, e := range res.Errors {
		logger.Warning("esbuild(error): %v", e)
	}
	logger.Debug("Compiled JS in %s", time.Since(start).String())

	pv, err := json.Marshal(opts.ParamValues)
	if err != nil {
		return nil, nil, errors.Wrap(err, "serializing param values")
	}

	return []string{"node", res.OutputFiles[0].Path, string(pv)}, closer, nil
}

// checkNodeVersion compares the major version of the currently installed
// node binary with that of the configured task and logs a warning if they
// do not match.
func checkNodeVersion(ctx context.Context, logger logger.Logger, opts build.KindOptions) {
	nodeVersion, ok := opts["nodeVersion"].(string)
	if !ok {
		return
	}

	v, err := semver.ParseTolerant(nodeVersion)
	if err != nil {
		logger.Debug("Unable to parse node version (%s): ignoring", nodeVersion)
		return
	}

	cmd := exec.CommandContext(ctx, "node", "--version")
	logger.Debug("Running %s", strings.Join(cmd.Args, " "))
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Debug("failed to check node version: is node installed?")
		return
	}

	logger.Debug("node version: %s", strings.TrimSpace(string(out)))
	if !strings.HasPrefix(string(out), fmt.Sprintf("v%d", v.Major)) {
		logger.Warning("Your local version of Node (%s) does not match the version your task is configured to run against (v%s).", strings.TrimSpace(string(out)), v)
	}

	cmd = exec.CommandContext(ctx, "npx", "--version")
	logger.Debug("Running %s", strings.Join(cmd.Args, " "))
	out, err = cmd.CombinedOutput()
	if err != nil {
		logger.Debug("failed to check npx version: are you running a recent enough version of node?")
		return
	}

	logger.Debug("npx version: %s", strings.TrimSpace(string(out)))
}
