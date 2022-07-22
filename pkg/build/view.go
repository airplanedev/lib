package build

import (
	_ "embed"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/airplanedev/lib/pkg/utils/fsx"
	"github.com/pkg/errors"
)

// view creates a dockerfile for an view.
func view(root string, options KindOptions) (string, error) {
	// Assert that the entrypoint file exists:
	entrypoint, _ := options["entrypoint"].(string)
	if entrypoint == "" {
		return "", errors.New("expected an entrypoint")
	}
	if err := fsx.AssertExistsAll(filepath.Join(root, entrypoint)); err != nil {
		return "", err
	}

	// Assert that API host is set.
	apiHost, _ := options["apiHost"].(string)
	if apiHost == "" {
		return "", errors.New("expected an api host")
	}
	if !strings.HasPrefix(apiHost, "https://") {
		apiHost = "https://" + apiHost
	}

	// TODO: possibly support multiple build tools.
	base, err := getBaseNodeImage("")
	if err != nil {
		return "", err
	}

	mainTsxStr, err := MainTsxString(entrypoint)
	if err != nil {
		return "", err
	}
	indexHtmlStr, err := IndexHtmlString()
	if err != nil {
		return "", err
	}
	viteConfigStr, err := ViteConfigString()
	if err != nil {
		return "", err
	}

	packageJSONPath := filepath.Join(root, "package.json")
	var packageJSON interface{}
	if fsx.Exists(packageJSONPath) {
		packageJSONFile, err := ioutil.ReadFile(packageJSONPath)
		if err != nil {
			return "", errors.Wrap(err, "reading package JSON")
		}
		if err := json.Unmarshal([]byte(packageJSONFile), &packageJSON); err != nil {
			return "", errors.Wrap(err, "parsing package JSON")
		}
	}

	packagesToCheck := []string{"vite", "@vitejs/plugin-react", "react", "react-dom", "@airplane/views"}
	packagesToAdd := []string{}
	deps, ok := packageJSON.(map[string]interface{})["dependencies"].(map[string]interface{})
	if ok {
		for _, pkg := range packagesToCheck {
			if _, ok := deps[pkg]; !ok {
				packagesToAdd = append(packagesToAdd, pkg)
			}
		}
	}

	cfg := struct {
		Base             string
		InstallCommand   string
		OutDir           string
		InlineMainTsx    string
		InlineIndexHtml  string
		InlineViteConfig string
		APIHost          string
		PackagesToAdd    string
	}{
		Base:             base,
		InstallCommand:   "yarn install --non-interactive --frozen-lockfile",
		OutDir:           "dist",
		InlineMainTsx:    inlineString(mainTsxStr),
		InlineIndexHtml:  inlineString(indexHtmlStr),
		InlineViteConfig: inlineString(viteConfigStr),
		APIHost:          apiHost,
		PackagesToAdd:    strings.Join(packagesToAdd, " "),
	}

	return applyTemplate(heredoc.Doc(`
		FROM {{.Base}} as builder
		WORKDIR /airplane
		RUN apt update && apt install -y jq

		COPY package*.json yarn.* /airplane/
		RUN [ -f package.json ] || { echo "{}" > package.json; }
		{{if .PackagesToAdd }}
			RUN yarn add {{.PackagesToAdd}}
		{{end}}
		RUN {{.InstallCommand}}

		RUN mkdir /airplane/src/
		RUN {{.InlineIndexHtml}} > /airplane/index.html
		RUN {{.InlineMainTsx}} > /airplane/main.tsx
		RUN {{.InlineViteConfig}} > /airplane/vite.config.ts
		ENV AIRPLANE_API_HOST={{.APIHost}}

		COPY . /airplane/src/
		RUN /airplane/node_modules/.bin/vite build --outDir {{.OutDir}}

		# Docker's minimal image - we just need an empty place to copy the build artifacts.
		FROM scratch
		COPY --from=builder /airplane/{{.OutDir}}/ .
	`), cfg)
}

//go:embed views/vite.config.ts
var viteConfigTemplateStr string

func ViteConfigString() (string, error) {
	return viteConfigTemplateStr, nil
}

//go:embed views/index.html
var indexHtmlTemplateStr string

func IndexHtmlString() (string, error) {
	return indexHtmlTemplateStr, nil
}

//go:embed views/main.tsx
var mainTsxTemplateStr string

func MainTsxString(entrypoint string) (string, error) {
	if strings.HasSuffix(entrypoint, ".tsx") {
		entrypoint = entrypoint[:len(entrypoint)-4]
	}
	return applyTemplate(mainTsxTemplateStr, struct {
		Entrypoint string
	}{
		Entrypoint: entrypoint,
	})
}
