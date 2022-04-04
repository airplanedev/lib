package build

import (
	"context"
	"testing"
)

func TestShellBuilder(t *testing.T) {
	ctx := context.Background()

	tests := []Test{
		{
			Root: "shell/simple",
			Kind: TaskKindShell,
			Options: KindOptions{
				"shim":       "true",
				"entrypoint": "main.sh",
			},
		},
		{
			Root: "shell/ubuntu",
			Kind: TaskKindShell,
			Options: KindOptions{
				"shim":       "true",
				"entrypoint": "main.sh",
			},
			SearchString: "bar",
		},
		{
			Root: "shell/zcli",
			Kind: TaskKindShell,
			Options: KindOptions{
				"shim":       "true",
				"entrypoint": "main.sh",
			},
			SearchString: "bar",
		},
		{
			Root: "shell/diff-workdir",
			Kind: TaskKindShell,
			Options: KindOptions{
				"shim":       "true",
				"entrypoint": "main.sh",
			},
			SearchString: "bar",
		},
		{
			Root: "shell/ubuntu-no-newline",
			Kind: TaskKindShell,
			Options: KindOptions{
				"shim":       "true",
				"entrypoint": "main.sh",
			},
		},
	}

	RunTests(t, ctx, tests)
}
