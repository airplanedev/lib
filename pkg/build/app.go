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

	base, err := getBaseNodeImage("")
	if err != nil {
		return "", err
	}

	cfg := struct {
		Base           string
		InstallCommand string
		OutDir         string
	}{
		Base:           base,
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

		# Docker's minimal image - we just need an empty place to copy the build artifacts.
		FROM scratch
		COPY --from=builder /airplane/{{.OutDir}}/ .
	`), cfg)
}