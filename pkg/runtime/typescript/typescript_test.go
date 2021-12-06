package typescript

import (
	"context"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/airplanedev/lib/pkg/build"
	"github.com/airplanedev/lib/pkg/examples"
	"github.com/airplanedev/lib/pkg/runtime"
	"github.com/airplanedev/lib/pkg/runtime/runtimetest"
	"github.com/airplanedev/lib/pkg/utils/fsx"
	"github.com/stretchr/testify/require"
)

func TestDev(tt *testing.T) {
	ctx := context.Background()

	tests := []runtimetest.Test{
		{
			Kind: build.TaskKindNode,
			Opts: runtime.PrepareRunOptions{Path: "typescript/simple/main.ts"},
		},
		{
			Kind: build.TaskKindNode,
			Opts: runtime.PrepareRunOptions{Path: "typescript/npm/main.ts"},
		},
		{
			Kind: build.TaskKindNode,
			Opts: runtime.PrepareRunOptions{Path: "typescript/yarn/main.ts"},
		},
		{
			Kind: build.TaskKindNode,
			Opts: runtime.PrepareRunOptions{Path: "typescript/imports/task/main.ts"},
		},
		{
			Kind:         build.TaskKindNode,
			Opts:         runtime.PrepareRunOptions{Path: "typescript/noparams/main.ts"},
			SearchString: "success",
		},
		{
			Kind: build.TaskKindNode,
			Opts: runtime.PrepareRunOptions{Path: "typescript/esnext/main.ts"},
		},
		{
			Kind: build.TaskKindNode,
			Opts: runtime.PrepareRunOptions{
				Path: "typescript/esnext/main.ts",
				// This tests that we can compile to an older version of Node syntax.
				// However, since this runs against whatever version of Node is installed
				// (f.e. in CI in GHA, or on your laptop, both of which should have the
				// latest version of Node) it doesn't actually test that the compiled
				// code _runs_ in an older version of Node. We have more control over
				// this in our Docker-based builder tests, so we test that there.
				KindOptions: build.KindOptions{"nodeVersion": "12"},
			},
		},
		{
			Kind: build.TaskKindNode,
			Opts: runtime.PrepareRunOptions{
				Path: "typescript/esnext/main.ts",
				// See nodeVersion comment above.
				KindOptions: build.KindOptions{"nodeVersion": "14"},
			},
		},
		{
			Kind: build.TaskKindNode,
			Opts: runtime.PrepareRunOptions{Path: "typescript/esm/main.ts"},
		},
		{
			Kind: build.TaskKindNode,
			Opts: runtime.PrepareRunOptions{Path: "typescript/aliases/main.ts"},
		},
		// TODO: debug why yarn workspaces aren't working. Seems like we would need to compile
		// pkg1 before compiling pkg2. Once we do that, add an npm workspaces variant along with
		// JS variants.
		// {
		// 	Kind: build.TaskKindNode,
		// 	Opts: runtime.PrepareRunOptions{Path: "typescript/yarnworkspaces/pkg2/src/index.ts"},
		// },
	}

	// For the dev workflow, we expect users to run `npm install` themselves before
	// running the dev command. Therefore, perform an `npm install` on each example:
	for _, test := range tests {
		p := examples.Path(tt, test.Opts.Path)

		// Check if this example uses npm or yarn:
		r, err := runtime.Lookup(p, test.Kind)
		require.NoError(tt, err)
		root, err := r.Root(p)
		require.NoError(tt, err)
		var cmd *exec.Cmd
		if fsx.Exists(filepath.Join(root, "yarn.lock")) {
			cmd = exec.CommandContext(ctx, "yarn", "install", "--frozen-lockfile")
		} else {
			cmd = exec.CommandContext(ctx, "npm", "install", "--no-save")
		}

		// Install dependencies:
		cmd.Dir = root
		out, err := cmd.CombinedOutput()
		require.NoError(tt, err, "Failed to run %q for %q:\n%s", cmd.String(), test.Opts.Path, string(out))
	}

	runtimetest.Run(tt, ctx, tests)
}
