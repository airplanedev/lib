package build

import (
	_ "embed"
	"os"
	"path/filepath"

	"github.com/MakeNowJust/heredoc"
	"github.com/pkg/errors"
)

// app creates a dockerfile for an app.
func app(root string) (string, error) {
	// TODO: create vite.config.ts if it does not exist.
	// TODO: possibly support multiple build tools.
	_, err := os.Stat(filepath.Join(root, "vite.config.ts"))
	if err != nil {
		if os.IsNotExist(err) {
			return "", errors.New("only vite is supported. Root directory must have a vite.config.ts")
		}
		return "", err
	}

	cfg := struct {
		Base           string
		InstallCommand string
		OutDir         string
	}{
		Base:           "node@sha256:d388b6e0648a0f56765260daa20ee4533225845ee1e4d48886c3d0f7aaaa6384",
		InstallCommand: "yarn install --non-interactive --frozen-lockfile",
		OutDir:         "dist",
	}

	return applyTemplate(heredoc.Doc(`
		FROM {{.Base}} as builder
		WORKDIR /airplane

		COPY package*.json yarn.* /airplane/
		RUN {{.InstallCommand}}

		COPY . /airplane/
		RUN yarn build --outDir {{.OutDir}}

		FROM scratch
		COPY --from=builder /airplane/{{.OutDir}}/ .
	`), cfg)
}
