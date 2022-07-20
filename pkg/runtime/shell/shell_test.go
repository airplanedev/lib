package shell

import (
	"os"
	"testing"

	"github.com/airplanedev/lib/pkg/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateShell(t *testing.T) {
	testCases := []struct {
		desc            string
		generateComment bool
		code            string
	}{
		{
			desc: "generates code",
			code: `#!/bin/bash
# Params are in environment variables as PARAM_{SLUG}, e.g. PARAM_USER_ID
echo "Hello World!"
echo "Printing env for debugging purposes:"
env
`,
		},
		{
			desc:            "generates code with comment",
			generateComment: true,
			code: `#!/bin/bash
# Linked to https://app.airplane.dev/t/shell_simple [do not edit this line]

# Params are in environment variables as PARAM_{SLUG}, e.g. PARAM_USER_ID
echo "Hello World!"
echo "Printing env for debugging purposes:"
env
`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)

			r := Runtime{}
			code, fileMode, err := r.Generate(&runtime.Task{
				URL: "https://app.airplane.dev/t/shell_simple",
			}, runtime.GenerateOpts{GenerateComment: tC.generateComment})
			require.NoError(err)
			assert.Equal(os.FileMode(0744), fileMode)
			assert.Equal(tC.code, string(code))
		})
	}
}
