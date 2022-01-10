package build

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/airplanedev/lib/pkg/utils/fsx"
	"github.com/pkg/errors"
)

// node creates a dockerfile for Node (typescript/javascript).
func node(root string, options KindOptions) (string, error) {
	var err error

	// For backwards compatibility, continue to build old Node tasks
	// in the same way. Tasks built with the latest CLI will set
	// shim=true which enables the new code path.
	if shim, ok := options["shim"].(string); !ok || shim != "true" {
		return nodeLegacyBuilder(root, options)
	}

	// Assert that the entrypoint file exists:
	entrypoint, _ := options["entrypoint"].(string)
	if entrypoint == "" {
		return "", errors.New("expected an entrypoint")
	}
	if err := fsx.AssertExistsAll(filepath.Join(root, entrypoint)); err != nil {
		return "", err
	}

	workdir, _ := options["workdir"].(string)
	pathPackageJSON := filepath.Join(root, "package.json")
	hasPackageJSON := fsx.AssertExistsAll(pathPackageJSON) == nil
	pathYarnLock := filepath.Join(root, "yarn.lock")
	var usesWorkspaces bool

	if hasPackageJSON {
		// Check to see if the package.json uses yarn/npm workspaces.
		// If the package.json has a "workspaces" key, it uses workspaces!
		// We want to know this because if we are in a workspace, our install
		// has to honor all of the package.json in the workspace.
		buf, err := os.ReadFile(pathPackageJSON)
		if err != nil {
			return "", errors.Wrapf(err, "node: reading %s", pathPackageJSON)
		}

		var pkg pkgJSONWorkspaces
		if err := json.Unmarshal(buf, &pkg); err != nil {
			return "", fmt.Errorf("node: parsing %s - %w", pathPackageJSON, err)
		}
		usesWorkspaces = len(pkg.workspaces) > 0
	}

	cfg := struct {
		Workdir               string
		Base                  string
		HasPackageJSON        bool
		UsesWorkspaces        bool
		IsYarn                bool
		InlineShim            string
		InlineShimPackageJSON string
		NodeVersion           string
		ExternalFlags         string
	}{
		Workdir:        workdir,
		HasPackageJSON: hasPackageJSON,
		IsYarn:         fsx.AssertExistsAll(pathYarnLock) == nil,
		UsesWorkspaces: usesWorkspaces,
		// esbuild is relatively generous in the node versions it supports:
		// https://esbuild.github.io/api/#target
		NodeVersion: GetNodeVersion(options),
	}

	// Workaround to get esbuild to not bundle dependencies.
	// See build.ExternalPackages for details.
	if cfg.HasPackageJSON {
		deps, err := ExternalPackages(pathPackageJSON)
		if err != nil {
			return "", err
		}
		var flags []string
		for _, dep := range deps {
			flags = append(flags, fmt.Sprintf("--external:%s", dep))
		}
		cfg.ExternalFlags = strings.Join(flags, " ")
	}

	if !strings.HasPrefix(cfg.Workdir, "/") {
		cfg.Workdir = "/" + cfg.Workdir
	}

	cfg.Base, err = getBaseNodeImage(cfg.NodeVersion)
	if err != nil {
		return "", err
	}

	pjson, err := GenShimPackageJSON()
	if err != nil {
		return "", err
	}
	cfg.InlineShimPackageJSON = inlineString(string(pjson))

	shim, err := NodeShim(entrypoint)
	if err != nil {
		return "", err
	}
	cfg.InlineShim = inlineString(shim)

	// The following Dockerfile can build both JS and TS tasks. In general, we're
	// aiming for recent EC202x support and for support for import/export syntax.
	// The former is easier, since recent versions of Node have excellent coverage
	// of the ECMAScript spec. The latter could be achieved through ECMAScript
	// modules (ESM), but those are not well-supported within the Node community.
	// Basic functionality of ESM is also still in the experimental stage, such as
	// module resolution for relative paths (f.e. ./main.js vs. ./main). Therefore,
	// we have to fallback to a separate build step to offer import/export support.
	// We have a few options -- f.e. babel, tsc, or swc -- but we go with esbuild
	// since it is native to Go.
	//
	// Down the road, we may want to give customers more control over this build process
	// in which case we could introduce an extra step for performing build commands.
	return applyTemplate(heredoc.Doc(`
		FROM {{.Base}}
		WORKDIR /airplane{{.Workdir}}
		# Support setting BUILD_NPM_RC or BUILD_NPM_TOKEN to configure private registry auth
		ARG BUILD_NPM_RC
		ARG BUILD_NPM_TOKEN
		RUN [ -z "${BUILD_NPM_RC}" ] || echo "${BUILD_NPM_RC}" > .npmrc
		RUN [ -z "${BUILD_NPM_TOKEN}" ] || echo "//registry.npmjs.org/:_authToken=${BUILD_NPM_TOKEN}" > .npmrc
		# qemu (on m1 at least) segfaults while looking up a UID/GID for running
		# postinstall scripts. We run as root with --unsafe-perm instead, skipping
		# that lookup. Possibly could fix by building for linux/arm on m1 instead
		# of always building for linux/amd64.
		RUN npm install -g typescript@4.2 && \
			npm install -g esbuild@0.12 --unsafe-perm
		RUN mkdir -p /airplane/.airplane && \
			cd /airplane/.airplane && \
			{{.InlineShimPackageJSON}} > package.json && \
			npm install

		{{if .HasPackageJSON}}
		COPY package*.json yarn.* /airplane/
		{{else}}
		RUN echo '{}' > /airplane/package.json
		{{end}}

		{{if .UsesWorkspaces}}
		COPY . /airplane
		{{end}}

		{{if .IsYarn}}
		RUN yarn --non-interactive
		{{else}}
		RUN npm install
		{{end}}

		{{if not .UsesWorkspaces}}
		COPY . /airplane
		{{end}}
		
		RUN {{.InlineShim}} > /airplane/.airplane/shim.js && \
			esbuild /airplane/.airplane/shim.js \
				--bundle \
				--external:airplane \
				--platform=node {{.ExternalFlags}} \
				--target=node{{.NodeVersion}} \
				--outfile=/airplane/.airplane/dist/shim.js
		ENTRYPOINT ["node", "/airplane/.airplane/dist/shim.js"]
	`), cfg)
}

func GenShimPackageJSON() ([]byte, error) {
	b, err := json.Marshal(struct {
		Dependencies map[string]string `json:"dependencies"`
	}{
		Dependencies: map[string]string{
			"airplane":    "~0.1.2",
			"@types/node": "^16",
		},
	})
	return b, errors.Wrap(err, "generating shim dependencies")
}

func GetNodeVersion(opts KindOptions) string {
	defaultVersion := "16"
	if opts == nil || opts["nodeVersion"] == nil {
		return defaultVersion
	}
	nv, ok := opts["nodeVersion"].(string)
	if !ok {
		return defaultVersion
	}

	return nv
}

//go:embed node-shim.js
var nodeShim string

func NodeShim(entrypoint string) (string, error) {
	// Remove the `.ts` suffix if one exists, since tsc doesn't accept
	// import paths with `.ts` endings. `.js` endings are fine.
	entrypoint = strings.TrimSuffix(entrypoint, ".ts")
	// The shim is stored under the .airplane directory.
	entrypoint = filepath.Join("../", entrypoint)
	// Escape for embedding into a string
	entrypoint = backslashEscape(entrypoint, `"`)

	shim, err := applyTemplate(nodeShim, struct {
		Entrypoint string
	}{
		Entrypoint: entrypoint,
	})
	if err != nil {
		return "", errors.Wrap(err, "templating shim")
	}

	return shim, nil
}

// nodeLegacyBuilder creates a dockerfile for Node (typescript/javascript).
//
// TODO(amir): possibly just run `npm start` instead of exposing lots
// of options to users?
func nodeLegacyBuilder(root string, options KindOptions) (string, error) {
	entrypoint, _ := options["entrypoint"].(string)
	main := filepath.Join(root, entrypoint)
	deps := filepath.Join(root, "package.json")
	yarnlock := filepath.Join(root, "yarn.lock")
	pkglock := filepath.Join(root, "package-lock.json")
	lang, _ := options["language"].(string)
	// `workdir` is fixed usually - `buildWorkdir` is a subdirectory of `workdir` if there's
	// `buildCommand` and is ultimately where `entrypoint` is run from.
	buildCommand, _ := options["buildCommand"].(string)
	buildDir, _ := options["buildDir"].(string)
	workdir := "/airplane"
	buildWorkdir := "/airplane"
	cmds := []string{}

	// Make sure that entrypoint and `package.json` exist.
	if err := fsx.AssertExistsAll(main, deps); err != nil {
		return "", err
	}

	// Determine the install command to use.
	if err := fsx.AssertExistsAll(pkglock); err == nil {
		cmds = append(cmds, `npm install package-lock.json`)
	} else if err := fsx.AssertExistsAll(yarnlock); err == nil {
		cmds = append(cmds, `yarn install`)
	}

	// Language specific.
	switch lang {
	case "typescript":
		if buildDir == "" {
			buildDir = ".airplane"
		}
		cmds = append(cmds, `npm install -g typescript@4.1`)
		cmds = append(cmds, `[ -f tsconfig.json ] || echo '{"include": ["*", "**/*"], "exclude": ["node_modules"]}' >tsconfig.json`)
		cmds = append(cmds, fmt.Sprintf(`rm -rf %s && tsc --outDir %s --rootDir .`, buildDir, buildDir))
		if buildCommand != "" {
			// It's not totally expected, but if you do set buildCommand we'll run it after tsc
			cmds = append(cmds, buildCommand)
		}
		buildWorkdir = path.Join(workdir, buildDir)
		// If entrypoint ends in .ts, replace it with .js
		entrypoint = strings.TrimSuffix(entrypoint, ".ts") + ".js"
	case "javascript":
		if buildCommand != "" {
			cmds = append(cmds, buildCommand)
		}
		if buildDir != "" {
			buildWorkdir = path.Join(workdir, buildDir)
		}
	default:
		return "", errors.Errorf("build: unknown language %q, expected \"javascript\" or \"typescript\"", lang)
	}
	entrypoint = path.Join(buildWorkdir, entrypoint)

	baseImage, err := getBaseNodeImage(GetNodeVersion(options))
	if err != nil {
		return "", err
	}

	return applyTemplate(heredoc.Doc(`
		FROM {{ .Base }}
		WORKDIR {{ .Workdir }}
		# Support setting BUILD_NPM_RC or BUILD_NPM_TOKEN to configure private registry auth
		ARG BUILD_NPM_RC
		ARG BUILD_NPM_TOKEN
		RUN [ -z "${BUILD_NPM_RC}" ] || echo "${BUILD_NPM_RC}" > .npmrc
		RUN [ -z "${BUILD_NPM_TOKEN}" ] || echo "//registry.npmjs.org/:_authToken=${BUILD_NPM_TOKEN}" > .npmrc
		COPY . {{ .Workdir }}
		{{ range .Commands }}
		RUN {{ . }}
		{{ end }}
		WORKDIR {{ .BuildWorkdir }}
		ENTRYPOINT ["node", "{{ .Main }}"]
	`), struct {
		Base         string
		Workdir      string
		BuildWorkdir string
		Commands     []string
		Main         string
	}{
		Base:         baseImage,
		Workdir:      workdir,
		BuildWorkdir: buildWorkdir,
		Commands:     cmds,
		Main:         entrypoint,
	})
}

func getBaseNodeImage(version string) (string, error) {
	if version == "" {
		version = "16"
	}
	v, err := GetVersion(NameNode, version)
	if err != nil {
		return "", err
	}
	base := v.String()
	if base == "" {
		// Assume the version is already a more-specific version - default to just returning it back
		base = "node:" + version + "-buster"
	}

	return base, nil
}

type pkgJSONWorkspaces struct {
	workspaces []string
}

func (p *pkgJSONWorkspaces) UnmarshalJSON(data []byte) error {
	// Workspaces might be an array of strings...
	var workspacesArray struct {
		Workspaces []string `json:"workspaces"`
	}
	if err := json.Unmarshal(data, &workspacesArray); err == nil {
		p.workspaces = workspacesArray.Workspaces
		return nil
	}

	// Or it might be an object with an array of strings.
	var workspacesObject struct {
		Workspaces struct {
			Packages []string `json:"packages"`
		} `json:"workspaces"`
	}
	if err := json.Unmarshal(data, &workspacesObject); err != nil {
		return err
	}
	p.workspaces = workspacesObject.Workspaces.Packages
	return nil

}
