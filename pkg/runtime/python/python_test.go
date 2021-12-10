package python

import (
	"context"
	"testing"

	"github.com/airplanedev/lib/pkg/common/logger"
	"github.com/stretchr/testify/require"
)

func TestCheckPythonInstalled(t *testing.T) {
	require := require.New(t)

	// Assumes python3 is installed in test environment...
	err := checkPythonInstalled(context.Background(), &logger.NoopLogger{})
	require.NoError(err)
}
