// runtimetest is a testing library for constructing runtime
// tests using example tasks from the `pkg/examples` directory.
package runtimetest

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/airplanedev/lib/pkg/build"
	"github.com/airplanedev/lib/pkg/common/logger"
	"github.com/airplanedev/lib/pkg/examples"
	"github.com/airplanedev/lib/pkg/runtime"
	"github.com/otiai10/copy"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/require"
)

type Test struct {
	Kind build.TaskKind
	Opts runtime.PrepareRunOptions
	// SearchString is a string to look for in the example's output
	// to validate that the task completed successfully. If not set,
	// defaults to a random value which is passed into the example
	// via the `id` parameter.
	SearchString string
}

func Run(tt *testing.T, ctx context.Context, tests []Test) {
	for _, test := range tests {
		test := test // get a local loop reference
		opts := test.Opts

		r, err := runtime.Lookup(opts.Path, test.Kind)
		require.NoError(tt, err)

		tt.Run(toName(tt, r, opts), func(t *testing.T) {
			t.Parallel()
			require := require.New(t)

			// Generate a random ID that we can look for in the output to make
			// sure the task ran correctly.
			if opts.ParamValues == nil {
				opts.ParamValues = runtime.Values{}
			}
			if test.SearchString == "" {
				test.SearchString = ksuid.New().String()
				opts.ParamValues["id"] = test.SearchString
			}

			cmds, closer, err := r.PrepareRun(ctx, &logger.NoopLogger{}, runtime.PrepareRunOptions{
				Path:        copyExample(t, r, examples.Path(t, opts.Path)),
				ParamValues: opts.ParamValues,
				KindOptions: opts.KindOptions,
			})
			require.NoError(err)
			defer func() {
				require.NoError(closer.Close())
			}()

			// Execute the run and look for our search string to indicate that
			// it built and ran correctly.
			cmd := exec.CommandContext(ctx, cmds[0], cmds[1:]...)
			out, err := cmd.CombinedOutput()
			require.NoError(err, "unable to run dev command:\n%s", string(out))
			require.True(strings.Contains(string(out), test.SearchString), "unable to find %q in output:\n%s", test.SearchString, string(out))
		})
	}
}

// toName is a helper to extract the test name from an example path.
func toName(t *testing.T, r runtime.Interface, opts runtime.PrepareRunOptions) string {
	root, err := r.Root(examples.Path(t, opts.Path))
	require.NoError(t, err, "unable to lookup root for %q", opts.Path)
	return filepath.Base(root)
}

// copyExample is a helper to copy the contents of an example into a
// temporary directory where we can safely perform side effects (like
// generating a .airplane) without colliding with other parallel tests
// (f.e. testing different kindOptions on the same example.
func copyExample(t *testing.T, r runtime.Interface, path string) string {
	require := require.New(t)

	// Create the temporary directory.
	tmpdir, err := os.MkdirTemp("", "runtimes-*")
	require.NoError(err)
	t.Cleanup(func() {
		require.NoError(os.RemoveAll(tmpdir))
	})

	// Copy the example's root directory into the temporary directory.
	root, err := r.Root(path)
	require.NoError(err, "unable to lookup root")
	require.NoError(copy.Copy(root, tmpdir))

	// Recompute the example path to be relative to the temporary directory.
	relpath, err := filepath.Rel(root, path)
	require.NoError(err)
	return filepath.Join(tmpdir, relpath)
}
