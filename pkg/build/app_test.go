package build

import (
	"context"
	"testing"
)

func TestAppBuilder(t *testing.T) {
	ctx := context.Background()

	tests := []Test{
		{
			Root:        "app/simple",
			Kind:        "app",
			SkipCommand: true,
		},
	}

	RunTests(t, ctx, tests)
}
