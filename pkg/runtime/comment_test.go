package runtime

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSlug(tt *testing.T) {
	for _, test := range []struct {
		name string
		in   string
		slug string
	}{
		{
			name: "empty file",
			in:   ``,
		},
		{
			name: "missing comment",
			in:   `import airplane from 'airplane'`,
		},
		{
			name: "unrelated comment",
			in: `// Airplane (https://airplane.dev) is great!
console.log('ship it')`,
		},
		{
			name: "extracts slug correctly",
			in: `// Linked to https://app.airplane.dev/t/myslug [do not edit this line]
console.log('ship it')`,
			slug: "myslug",
		},
		{
			name: "extracts slug correctly in staging",
			in: `// Linked to https://web.airstage.app/t/myslug [do not edit this line]
console.log('ship it')`,
			slug: "myslug",
		},
		{
			name: "extracts slug correctly in dev",
			in: `// Linked to https://app.airplane.so:5000/t/myslug [do not edit this line]
console.log('ship it')`,
			slug: "myslug",
		},
	} {
		tt.Run(test.name, func(t *testing.T) {
			buf := strings.NewReader(test.in)
			require.Equal(t, test.slug, slugFromReader(buf))
		})
	}
}
