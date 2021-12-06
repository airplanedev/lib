package python

import (
	"context"
	"testing"

	"github.com/airplanedev/lib/pkg/runtime/runtimetest"
	"github.com/stretchr/testify/require"
)

func TestCheckPythonInstalled(t *testing.T) {
	require := require.New(t)

	// Assumes python3 is installed in test environment...
	err := checkPythonInstalled(context.Background(), &runtimetest.NoopLogger{})
	require.NoError(err)
}
