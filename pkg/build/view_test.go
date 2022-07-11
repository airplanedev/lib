package build

import (
	"context"
	"testing"
)

func TestViewBuilder(t *testing.T) {
	ctx := context.Background()

	tests := []Test{
		{
			Root:    "view/simple",
			Kind:    "view",
			SkipRun: true,
		},
	}

	RunTests(t, ctx, tests)
}
