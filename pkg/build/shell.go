package build

import (
	_ "embed"
	"io/ioutil"
	"path/filepath"

	"github.com/MakeNowJust/heredoc"
	"github.com/airplanedev/lib/pkg/utils/fsx"
	"github.com/pkg/errors"
)

func shell(root string, options KindOptions) (string, error) {
	// Assert that the entrypoint file exists:
	entrypoint, _ := options["entrypoint"].(string)
	if entrypoint == "" {
		return "", errors.New("entrypoint is unexpectedly missing")
	}
	if err := fsx.AssertExistsAll(filepath.Join(root, entrypoint)); err != nil {
		return "", err
	}

	// Build off of the dockerfile if provided:
	var dockerfileTemplate string
	var workDir string
	if dockerfilePath := FindDockerfile(root); dockerfilePath != "" {
		contents, err := ioutil.ReadFile(dockerfilePath)
		if err != nil {
			return "", errors.Wrap(err, "opening dockerfile")
		}
		dockerfileTemplate = string(contents)
		workDir = "."
	} else {
		dockerfileTemplate = heredoc.Doc(`
			FROM ubuntu:21.04
			# Install some common libraries
			RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
				&& apt-get -y install --no-install-recommends \
					apt-utils \
					openssh-client \
					gnupg2 \
					iproute2 \
					procps \
					lsof \
					htop \
					net-tools \
					curl \
					wget \
					ca-certificates \
					unzip \
					zip \
					nano \
					vim-tiny \
					less \
					jq \
					lsb-release \
					apt-transport-https \
					dialog \
					zlib1g \
					locales \
					strace \
				&& apt-get autoremove -y && apt-get clean -y && rm -rf /var/lib/apt/lists/*
		`)
		workDir = "/airplane"
	}

	// Extend template with our own logic - set up a WORKDIR and shim.
	dockerfileTemplate = dockerfileTemplate + heredoc.Doc(`
		WORKDIR {{.Workdir}}
		RUN mkdir -p .airplane && {{.InlineShim}} > .airplane/shim.sh
		
		COPY . .
		RUN chmod +x {{.Entrypoint}}
		
		ENTRYPOINT ["bash", ".airplane/shim.sh", "./{{.Entrypoint}}"]
	`)
	return applyTemplate(dockerfileTemplate, struct {
		InlineShim string
		Entrypoint string
		Workdir    string
	}{
		InlineShim: inlineString(ShellShim()),
		Entrypoint: backslashEscape(entrypoint, `"`),
		Workdir:    workDir,
	})
}

//go:embed shell-shim.sh
var shellShim string

func ShellShim() string {
	return shellShim
}

func DockerfilePaths() []string {
	return []string{
		"Dockerfile.airplane",
		"Dockerfile",
	}
}

// FindDockerfile looks for variants of supported Dockerfile locations and returns non-empty path
// to the file, if found.
func FindDockerfile(root string) string {
	for _, filePath := range DockerfilePaths() {
		dockerfilePath := filepath.Join(root, filePath)
		if fsx.Exists(dockerfilePath) {
			return dockerfilePath
		}
	}
	return ""
}
