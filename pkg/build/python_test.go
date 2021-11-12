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
	}

	RunTests(t, ctx, tests)
}
