package build

import (
	"context"
	"testing"
)

func TestViewBuilder(t *testing.T) {
	ctx := context.Background()

	tests := []Test{
		{
			Root: "view/simple",
			Kind: "view",
			Options: KindOptions{
				"entrypoint": "src/App.tsx",
				"apiHost":    "https://api:5000",
			},
			SkipRun: true,
		},
	}

	RunTests(t, ctx, tests)
}
