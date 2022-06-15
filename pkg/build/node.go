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

type templateParams struct {
	Workdir                            string
	Base                               string
	HasPackageJSON                     bool
	UsesWorkspaces                     bool
	InlineShim                         string
	InlineShimPackageJSON              string
	InlineWorkflowBundlerScript        string
	InlineWorkflowInterceptorsScript   string
	InlineWorkflowShimScript           string
	InlineWorkflowShimActivitiesScript string
	IsWorkflow                         bool
	HasCustomActivities                bool
	NodeVersion                        string
	ExternalFlags                      string
	InstallCommand                     string
	PostInstallCommand                 string
	Args                               string
}

// node creates a dockerfile for Node (typescript/javascript).
func node(root string, options KindOptions, buildArgs []string) (string, error) {
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
	rootPackageJSON := filepath.Join(root, "package.json")
	hasPackageJSON := fsx.AssertExistsAll(rootPackageJSON) == nil
	pathYarnLock := filepath.Join(root, "yarn.lock")
	pathPackageLock := filepath.Join(root, "package-lock.json")
	hasPackageLock := fsx.AssertExistsAll(pathPackageLock) == nil
	isYarn := fsx.AssertExistsAll(pathYarnLock) == nil

	hasCustomActivities := fsx.AssertExistsAny(
		filepath.Join(root, "activities.ts"),
		filepath.Join(root, "activities.js"),
	) == nil

	runtime, ok := options["runtime"]
	var isWorkflow bool

	if ok {
		// Depending on how the options were serialized, the runtime can be
		// either a string or TaskRuntime; handle both.
		switch v := runtime.(type) {
		case string:
			isWorkflow = v == string(TaskRuntimeWorkflow)
		case TaskRuntime:
			isWorkflow = v == TaskRuntimeWorkflow
		default:
		}
	}

	var pkg pkgJSON
	if hasPackageJSON {
		// Check to see if the package.json uses yarn/npm workspaces.
		// If the package.json has a "workspaces" key, it uses workspaces!
		// We want to know this because if we are in a workspace, our install
		// has to honor all of the package.json in the workspace.
		buf, err := os.ReadFile(rootPackageJSON)
		if err != nil {
			return "", errors.Wrapf(err, "node: reading %s", rootPackageJSON)
		}

		if err := json.Unmarshal(buf, &pkg); err != nil {
			return "", fmt.Errorf("node: parsing %s - %w", rootPackageJSON, err)
		}
	}

	for i, a := range buildArgs {
		buildArgs[i] = fmt.Sprintf("ARG %s", a)
	}
	argsCommand := strings.Join(buildArgs, "\n")

	cfg := templateParams{
		Workdir:        workdir,
		HasPackageJSON: hasPackageJSON,
		UsesWorkspaces: len(pkg.Workspaces.workspaces) > 0,
		// esbuild is relatively generous in the node versions it supports:
		// https://esbuild.github.io/api/#target
		NodeVersion:         GetNodeVersion(options),
		PostInstallCommand:  pkg.Settings.PostInstallCommand,
		Args:                argsCommand,
		IsWorkflow:          isWorkflow,
		HasCustomActivities: hasCustomActivities,
	}

	if cfg.HasPackageJSON {
		// Workaround to get esbuild to not bundle dependencies.
		// See build.ExternalPackages for details.
		deps, err := ExternalPackages(rootPackageJSON)
		if err != nil {
			return "", err
		}
		var flags []string
		for _, dep := range deps {
			flags = append(flags, fmt.Sprintf("--external:%s", dep))
		}
		if isWorkflow {
			// Even if these are imported, we need to mark the root packages
			// as external for esbuild to work properly. Esbuild doesn't
			// care about repeats, so no need to dedupe.
			flags = append(flags, "--external:@temporalio", "--external:@swc")
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

	pjson, err := GenShimPackageJSON(rootPackageJSON, isWorkflow)
	if err != nil {
		return "", err
	}
	cfg.InlineShimPackageJSON = inlineString(string(pjson))

	if isWorkflow {
		cfg.InlineShim = inlineString(workerAndActivityShim)
		cfg.InlineWorkflowBundlerScript = inlineString(workflowBundlerScript)
		cfg.InlineWorkflowInterceptorsScript = inlineString(workflowInterceptorsScript)

		workflowShim, err := TemplateEntrypoint(workflowShimScript, entrypoint)
		if err != nil {
			return "", err
		}
		cfg.InlineWorkflowShimScript = inlineString(workflowShim)
		cfg.InlineWorkflowShimActivitiesScript = inlineString(workflowShimActivitiesScript)
	} else {
		shim, err := TemplatedNodeShim(entrypoint)
		if err != nil {
			return "", err
		}
		cfg.InlineShim = inlineString(shim)
	}

	installCommand := "npm install --production"
	if pkg.Settings.InstallCommand != "" {
		installCommand = pkg.Settings.InstallCommand
	} else if isYarn {
		installCommand = "yarn install --non-interactive --production --frozen-lockfile"
	} else if hasPackageLock {
		// Use npm ci if possible, since it's faster and behaves better:
		// https://docs.npmjs.com/cli/v8/commands/npm-ci
		installCommand = "npm ci --production"
	}
	cfg.InstallCommand = strings.ReplaceAll(installCommand, "\n", "\\n")

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
		ENV NODE_ENV=production
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
		RUN npm install -g esbuild@0.12 --unsafe-perm
		
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

		{{.Args}}

		RUN {{.InstallCommand}}

		{{if not .UsesWorkspaces}}
		COPY . /airplane
		{{end}}

		{{if .PostInstallCommand}}
		RUN {{.PostInstallCommand}}
		{{end}}

		{{if .IsWorkflow}}
		{{if not .HasCustomActivities}}
		RUN touch /airplane/activities.js
		{{end}}
		RUN {{.InlineWorkflowShimScript}} >> /airplane/.airplane/workflow-shim.js
		RUN {{.InlineWorkflowShimActivitiesScript}} >> /airplane/.airplane/workflow-shim-activities.js
		RUN {{.InlineWorkflowInterceptorsScript}} >> /airplane/.airplane/workflow-interceptors.js
		RUN {{.InlineWorkflowBundlerScript}} >> /airplane/.airplane/workflow-bundler.js
		RUN node /airplane/.airplane/workflow-bundler.js
		{{end}}

		RUN {{.InlineShim}} > /airplane/.airplane/shim.js && \
			esbuild /airplane/.airplane/shim.js \
				--bundle \
				--platform=node {{.ExternalFlags}} \
				--target=node{{.NodeVersion}} \
				--outfile=/airplane/.airplane/dist/shim.js

		ENTRYPOINT ["node", "/airplane/.airplane/dist/shim.js"]
	`), cfg)
}

func GenShimPackageJSON(pathPackageJSON string, isWorkflow bool) ([]byte, error) {
	deps, err := ListDependencies(pathPackageJSON)
	if err != nil {
		return nil, err
	}

	pjson := struct {
		Dependencies map[string]string `json:"dependencies"`
		Type         string            `json:"type,omitempty"`
	}{
		Dependencies: map[string]string{
			"airplane": "~0.1.2",
		},
	}

	if isWorkflow {
		// airplane>=0.2.0-6 already includes Temporal as a dependency, and so we don't include it here.
		pjson.Dependencies["airplane"] = "0.2.0-6"
	}

	// Allow users to override any shim dependencies. Given shim code is bundled
	// with user code, we cannot use separate versions of these dependencies so
	// default to whichever version the user requests.
	for _, dep := range deps {
		delete(pjson.Dependencies, dep)
	}

	b, err := json.Marshal(pjson)
	return b, errors.Wrap(err, "marshalling shim dependencies")
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

//go:embed workflow/worker-and-activity-shim.js
var workerAndActivityShim string

//go:embed workflow/workflow-bundler.js
var workflowBundlerScript string

//go:embed workflow/workflow-interceptors.js
var workflowInterceptorsScript string

//go:embed workflow/workflow-shim.js
var workflowShimScript string

//go:embed workflow/workflow-shim-activities.js
var workflowShimActivitiesScript string

func TemplatedNodeShim(entrypoint string) (string, error) {
	return TemplateEntrypoint(nodeShim, entrypoint)
}

func TemplateEntrypoint(script string, entrypoint string) (string, error) {
	// Remove the `.ts` suffix if one exists, since tsc doesn't accept
	// import paths with `.ts` endings. `.js` endings are fine.
	entrypoint = strings.TrimSuffix(entrypoint, ".ts")
	// The shim is stored under the .airplane directory.
	entrypoint = filepath.Join("../", entrypoint)
	// Escape for embedding into a string
	entrypoint = backslashEscape(entrypoint, `"`)

	shim, err := applyTemplate(script, struct {
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

// Settings represent Airplane specific settings.
type Settings struct {
	Root               string `json:"root"`
	InstallCommand     string `json:"install"`
	PostInstallCommand string `json:"postinstall"`
}

type pkgJSON struct {
	Settings   Settings          `json:"airplane"`
	Workspaces pkgJSONWorkspaces `json:"workspaces"`
}

type pkgJSONWorkspaces struct {
	workspaces []string
}

func (p *pkgJSONWorkspaces) UnmarshalJSON(data []byte) error {
	// Workspaces might be an array of strings...
	var workspaces []string
	if err := json.Unmarshal(data, &workspaces); err == nil {
		p.workspaces = workspaces
		return nil
	}

	// Or it might be an object with an array of strings.
	var workspacesObject struct {
		Packages []string `json:"packages"`
	}
	if err := json.Unmarshal(data, &workspacesObject); err != nil {
		return err
	}
	p.workspaces = workspacesObject.Packages
	return nil

}
