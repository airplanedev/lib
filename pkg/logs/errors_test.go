package logs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScanForErrorNodeESM(tt *testing.T) {
	for _, test := range []struct {
		name   string
		log    string
		module string
		ok     bool
	}{
		{
			name:   "match: node 18",
			log:    `Error [ERR_REQUIRE_ESM]: require() of ES Module /airplane/node_modules/execa/index.js from /airplane/.airplane/dist/shim.js not supported.`,
			module: "execa",
			ok:     true,
		},
		{
			name:   "match: node 16",
			log:    `Error [ERR_REQUIRE_ESM]: Must use import to load ES Module: /airplane/node_modules/execa/index.js require() of ES modules is not supported.`,
			module: "execa",
			ok:     true,
		},
		{
			name:   "match: node 14",
			log:    `Error [ERR_REQUIRE_ESM]: Must use import to load ES Module: /airplane/node_modules/execa/index.js`,
			module: "execa",
			ok:     true,
		},
		{
			name:   "match: different module",
			log:    `Error [ERR_REQUIRE_ESM]: Must use import to load ES Module: /test/foo/bar/node_modules/node-fetch/dist/index.js`,
			module: "node-fetch",
			ok:     true,
		},
		{
			name:   "match: unknown format",
			log:    `Error [ERR_REQUIRE_ESM]: some output format we don't support yet`,
			module: "",   // No module found...
			ok:     true, // ...but we still identify it was an ESM issue
		},
		{
			name: "no match",
			log:  `Error [ERR_LOOKS_LIKE_PYTHON]`,
			ok:   false,
		},
	} {
		tt.Run(test.name, func(t *testing.T) {
			module, ok := ScanForErrorNodeESM(test.log)
			require.Equal(t, test.ok, ok)
			if test.ok {
				require.Equal(t, test.module, module)
			}
		})
	}
}
