package shell

import (
	"os"
	"testing"

	"github.com/airplanedev/lib/pkg/runtime"
	"github.com/stretchr/testify/require"
)

func TestShellRuntime(t *testing.T) {
	require := require.New(t)

	r := Runtime{}
	code, fileMode, err := r.Generate(&runtime.Task{
		URL: "https://app.airplane.dev/t/shell_simple",
	})
	require.NoError(err)
	require.Equal(os.FileMode(0744), fileMode)
	require.Equal(`#!/bin/bash
# Linked to https://app.airplane.dev/t/shell_simple [do not edit this line]

# Params are in environment variables as PARAM_{SLUG}, e.g. PARAM_USER_ID
echo "Hello World!"
echo "Printing env for debugging purposes:"
env
`, string(code))
}
