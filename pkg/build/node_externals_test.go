package build

import (
	"testing"

	"github.com/airplanedev/lib/pkg/build/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExternalPackages(t *testing.T) {
	testCases := []struct {
		desc             string
		packageJSON      string
		externalPackages []string
	}{
		{
			desc:        "no packages",
			packageJSON: fixtures.Path(t, "node_externals/empty/package.json"),
		},
		{
			desc:             "marks external dependencies, dev dependencies and optional dependencies",
			packageJSON:      fixtures.Path(t, "node_externals/dependencies/package.json"),
			externalPackages: []string{"react", "@types/react", "react-table"},
		},
		{
			desc:             "does not mark esm modules as external",
			packageJSON:      fixtures.Path(t, "node_externals/esm/package.json"),
			externalPackages: []string{"react"},
		},
		{
			desc:             "marks external all packages in yarn workspace",
			packageJSON:      fixtures.Path(t, "node_externals/yarnworkspace/package.json"),
			externalPackages: []string{"react", "@types/react", "react-table"},
		},
		{
			desc:             "does not mark local yarn workspace import as external",
			packageJSON:      fixtures.Path(t, "node_externals/yarnworkspace_importlocal/package.json"),
			externalPackages: []string{"react", "@types/react", "react-table"},
		},
		{
			desc:             "marks local yarn workspace import as external if mismatched version",
			packageJSON:      fixtures.Path(t, "node_externals/yarnworkspace_importlocalmismatched/package.json"),
			externalPackages: []string{"react", "@types/react", "react-table", "lib"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)
			externalPackages, err := ExternalPackages(tC.packageJSON)
			require.NoError(err)

			assert.ElementsMatch(tC.externalPackages, externalPackages)
		})
	}
}
