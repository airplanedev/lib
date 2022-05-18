package build

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var esmModules = []string{
	"node-fetch",
	// airplane>=0.2.0 depends on node-fetch
	"airplane",
}

// ExternalPackages reads package.json and returns all dependencies and dev dependencies.
// This is used as a bit of a workaround for esbuild - we're using esbuild to transform code but
// don't actually want it to bundle. We hit issues when it tries to bundle optional packages
// (and the packages aren't installed) - for example, pg optionally depends on pg-native, and
// using just pg causes esbuild to bundle pg which bundles pg-native, which errors.
// TODO: replace this with a cleaner esbuild plugin that can mark dependencies as external:
// https://github.com/evanw/esbuild/issues/619#issuecomment-751995294
func ExternalPackages(rootPackageJSON string) ([]string, error) {
	usesWorkspaces, err := hasWorkspaces(rootPackageJSON)
	if err != nil {
		return nil, err
	}
	pathPackageJSONs := []string{rootPackageJSON}
	if usesWorkspaces {
		nestedPackageJSONs, err := findNestedPackageJSONs(filepath.Dir(rootPackageJSON))
		if err != nil {
			return nil, err
		}
		for _, j := range nestedPackageJSONs {
			if j != rootPackageJSON {
				pathPackageJSONs = append(pathPackageJSONs, j)
			}
		}
	}

	var deps []string
	for _, pathPackageJSON := range pathPackageJSONs {
		if pathPackageJSON == "" {
			continue
		}

		if usesWorkspaces {
			// If we are in a npm/yarn workspace, we want to bundle all packages in the same
			// workspaces so they are run through esbuild.
			yarnWorkspacePackages, err := getYarnWorkspacePackages(pathPackageJSON)
			if err != nil {
				return nil, err
			}
			for _, p := range yarnWorkspacePackages {
				esmModules = append(esmModules, p)
			}
		}

		allDeps, err := ListDependencies(pathPackageJSON)
		if err != nil {
			return nil, err
		}
		for _, dep := range allDeps {
			// Mark all dependencies as external, except for known ESM-only deps. These deps
			// need to be bundled by esbuild so that esbuild can convert them to CommonJS.
			// As long as these modules don't happen to pull in any optional modules, we should be OK.
			// This is a bandaid until we figure out how to handle ESM without bundling.
			if !contains(esmModules, dep) {
				deps = append(deps, dep)
			}
		}
	}

	return deps, nil
}

// ListDependencies lists all dependencies (including dev and optional) in a `package.json` file.
func ListDependencies(pathPackageJSON string) ([]string, error) {
	var deps []string

	f, err := os.Open(pathPackageJSON)
	if err != nil {
		// There is no package.json (or we can't open it). Treat as having no dependencies.
		return []string{}, nil
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.Wrap(err, "reading package.json")
	}
	var d struct {
		Dependencies         map[string]string `json:"dependencies"`
		DevDependencies      map[string]string `json:"devDependencies"`
		OptionalDependencies map[string]string `json:"optionalDependencies"`
	}
	if err := json.Unmarshal(b, &d); err != nil {
		return nil, errors.Wrap(err, "unmarshaling package.json")
	}

	for k := range d.Dependencies {
		deps = append(deps, k)
	}
	for k := range d.DevDependencies {
		deps = append(deps, k)
	}
	for k := range d.OptionalDependencies {
		deps = append(deps, k)
	}
	return deps, nil
}

// contains returns true if `list` includes `needle`.
func contains(list []string, needle string) bool {
	for _, elem := range list {
		if elem == needle {
			return true
		}
	}
	return false
}

func findNestedPackageJSONs(dir string) ([]string, error) {
	var pathPackageJSONs []string
	err := filepath.WalkDir(dir,
		func(path string, di fs.DirEntry, err error) error {
			if !strings.Contains(path, "node_modules") && di.Name() == "package.json" {
				pathPackageJSONs = append(pathPackageJSONs, path)
			}
			return nil
		})
	return pathPackageJSONs, err
}

func hasWorkspaces(pathPackageJSON string) (bool, error) {
	if _, err := os.Stat(pathPackageJSON); errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	var pkg pkgJSON
	buf, err := os.ReadFile(pathPackageJSON)
	if err != nil {
		return false, errors.Wrapf(err, "node: reading %s", pathPackageJSON)
	}

	if err := json.Unmarshal(buf, &pkg); err != nil {
		return false, fmt.Errorf("node: parsing %s - %w", pathPackageJSON, err)
	}
	return len(pkg.Workspaces.workspaces) > 0, nil
}

func getYarnWorkspacePackages(pathPackageJSON string) ([]string, error) {
	cmd := exec.Command("yarn", "workspaces", "info")
	cmd.Dir = filepath.Dir(pathPackageJSON)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.Wrap(err, "reading yarn/npm workspaces. Do you have yarn installed?")
	}
	// out will be something like:
	// yarn workspaces v1.22.17
	// {
	//   "pkg1": {
	//     "location": "pkg1",
	//     "workspaceDependencies": [],
	//     "mismatchedWorkspaceDependencies": []
	//   },
	//   "pkg2": {
	//     "location": "pkg2",
	//     "workspaceDependencies": [
	//       "pkg1"
	//     ],
	//     "mismatchedWorkspaceDependencies": []
	//   }
	// }
	// Done in 0.02s.
	//
	// We want to grab the keys of the JSON object.
	r := regexp.MustCompile(`{[\S\s]+}`)
	yarnWorkspaceJSON := r.FindString(string(out))
	if yarnWorkspaceJSON == "" {
		return nil, errors.New("empty yarn workspace info")
	}
	keys, err := getJSONKeys(yarnWorkspaceJSON)
	if err != nil {
		return nil, errors.Wrap(err, "output of `yarn workspace info` is not JSON")
	}
	return keys, nil
}

func getJSONKeys(jsonString string) ([]string, error) {
	c := make(map[string]json.RawMessage)
	if err := json.Unmarshal([]byte(jsonString), &c); err != nil {
		return nil, err
	}

	keys := make([]string, len(c))

	i := 0
	for key := range c {
		keys[i] = key
		i++
	}

	return keys, nil
}
