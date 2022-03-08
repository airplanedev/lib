package build

import (
	"context"
	"testing"
)

func TestNodeBuilder(t *testing.T) {
	ctx := context.Background()

	tests := []Test{
		// {
		// 	Root: "javascript/simple",
		// 	Kind: TaskKindNode,
		// 	Options: KindOptions{
		// 		"shim":       "true",
		// 		"entrypoint": "main.js",
		// 	},
		// },
		// {
		// 	Root: "typescript/simple",
		// 	Kind: TaskKindNode,
		// 	Options: KindOptions{
		// 		"shim":       "true",
		// 		"entrypoint": "main.ts",
		// 	},
		// },
		// {
		// 	Root: "typescript/npm",
		// 	Kind: TaskKindNode,
		// 	Options: KindOptions{
		// 		"shim":       "true",
		// 		"entrypoint": "main.ts",
		// 	},
		// },
		// {
		// 	Root: "typescript/yarn",
		// 	Kind: TaskKindNode,
		// 	Options: KindOptions{
		// 		"shim":       "true",
		// 		"entrypoint": "main.ts",
		// 	},
		// },
		// {
		// 	Root: "typescript/imports",
		// 	Kind: TaskKindNode,
		// 	Options: KindOptions{
		// 		"shim":       "true",
		// 		"entrypoint": "task/main.ts",
		// 	},
		// },
		// {
		// 	Root: "typescript/noparams",
		// 	Kind: TaskKindNode,
		// 	Options: KindOptions{
		// 		"shim":       "true",
		// 		"entrypoint": "main.ts",
		// 	},
		// 	// Since this example does not take parameters, override the default SearchString.
		// 	SearchString: "success",
		// },
		// {
		// 	Root: "typescript/esnext",
		// 	Kind: TaskKindNode,
		// 	Options: KindOptions{
		// 		"shim":       "true",
		// 		"entrypoint": "main.ts",
		// 	},
		// },
		// {
		// 	Root: "typescript/esnext",
		// 	Kind: TaskKindNode,
		// 	Options: KindOptions{
		// 		"shim":        "true",
		// 		"entrypoint":  "main.ts",
		// 		"nodeVersion": "12",
		// 	},
		// },
		// {
		// 	Root: "typescript/esnext",
		// 	Kind: TaskKindNode,
		// 	Options: KindOptions{
		// 		"shim":        "true",
		// 		"entrypoint":  "main.ts",
		// 		"nodeVersion": "14",
		// 	},
		// },
		// {
		// 	Root: "typescript/esm",
		// 	Kind: TaskKindNode,
		// 	Options: KindOptions{
		// 		"shim":       "true",
		// 		"entrypoint": "main.ts",
		// 	},
		// },
		{
			Root: "typescript/aliases",
			Kind: TaskKindNode,
			Options: KindOptions{
				"shim":       "true",
				"entrypoint": "main.ts",
			},
		},
		// {
		// 	Root: "typescript/externals",
		// 	Kind: TaskKindNode,
		// 	Options: KindOptions{
		// 		"shim":       "true",
		// 		"entrypoint": "main.ts",
		// 	},
		// },
		// {
		// 	Root: "typescript/yarnworkspaces",
		// 	Kind: TaskKindNode,
		// 	Options: KindOptions{
		// 		"shim":       "true",
		// 		"entrypoint": "pkg2/src/index.ts",
		// 	},
		// },
		// {
		// 	Root: "typescript/yarnworkspacesobject",
		// 	Kind: TaskKindNode,
		// 	Options: KindOptions{
		// 		"shim":       "true",
		// 		"entrypoint": "pkg2/src/index.ts",
		// 	},
		// },
		// {
		// 	Root: "typescript/nodeworkspaces",
		// 	Kind: TaskKindNode,
		// 	Options: KindOptions{
		// 		"shim":       "true",
		// 		"entrypoint": "pkg2/src/index.ts",
		// 	},
		// },
		// {
		// 	Root: "typescript/nopackagejson",
		// 	Kind: TaskKindNode,
		// 	Options: KindOptions{
		// 		"shim":       "true",
		// 		"entrypoint": "main.ts",
		// 	},
		// },
		// {
		// 	Root: "typescript/custominstall",
		// 	Kind: TaskKindNode,
		// 	Options: KindOptions{
		// 		"shim":       "true",
		// 		"entrypoint": "main.ts",
		// 	},
		// },
		// {
		// 	Root: "typescript/custompostinstall",
		// 	Kind: TaskKindNode,
		// 	Options: KindOptions{
		// 		"shim":       "true",
		// 		"entrypoint": "main.ts",
		// 	},
		// },
		// {
		// 	Root: "typescript/prisma",
		// 	Kind: TaskKindNode,
		// 	Options: KindOptions{
		// 		"shim":       "true",
		// 		"entrypoint": "main.ts",
		// 	},
		// },
	}

	RunTests(t, ctx, tests)
}
