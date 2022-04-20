package build

import (
	"context"
	"testing"
)

func TestPythonBuilder(t *testing.T) {
	ctx := context.Background()

	tests := []Test{
		{
			Root: "python/simple",
			Kind: TaskKindPython,
			Options: KindOptions{
				"shim":       "true",
				"entrypoint": "main.py",
			},
		},
		{
			Root: "python/requirements",
			Kind: TaskKindPython,
			Options: KindOptions{
				"shim":       "true",
				"entrypoint": "main.py",
			},
			SearchString: "[1]",
		},
		{
			Root: "python/requirementswithbuildargs",
			Kind: TaskKindPython,
			Options: KindOptions{
				"shim":       "true",
				"entrypoint": "main.py",
			},
			BuildArgs: map[string]string{
				"VER": "3.1.0",
			},
			SearchString: "[1]",
		},
	}

	RunTests(t, ctx, tests)
}
