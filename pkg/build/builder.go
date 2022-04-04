package build

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/airplanedev/lib/pkg/build/ignore"
	"github.com/airplanedev/lib/pkg/utils/bufiox"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	dockerJSONMessage "github.com/docker/docker/pkg/jsonmessage"
	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"
)

// RegistryAuth represents the registry auth.
type RegistryAuth struct {
	Token string
	Repo  string
}

// Response represents a build response.
type Response struct {
	ImageURL string
	// Optional, only if applicable
	BuildID string
}

// Host returns the registry hostname.
func (r RegistryAuth) host() string {
	return strings.SplitN(r.Repo, "/", 2)[0]
}

// LocalConfig configures a (local) builder.
type LocalConfig struct {
	// Root is the root directory.
	//
	// It must be an absolute path to the project directory.
	Root string

	// Builder is the builder name to use.
	//
	// There are various built-in builders, along with the dockerfile
	// builder and image builder.
	//
	// If empty, it assumes the "image" builder.
	Builder string

	// Options are the build arguments to use.
	//
	// When nil, it uses an empty map of options.
	Options KindOptions

	// Auth represents the registry auth to use.
	//
	// If nil, Push will produce an error.
	Auth *RegistryAuth

	// BuildEnv is a map of build-time environment variables to use.
	BuildEnv map[string]string
}

type DockerfileConfig struct {
	Builder string
	Root    string
	Options KindOptions
}

// Builder implements an image builder.
type Builder struct {
	root     string
	name     string
	options  KindOptions
	auth     *RegistryAuth
	buildEnv map[string]string
	client   *client.Client
}

// New returns a new local builder with c.
func New(c LocalConfig) (*Builder, error) {
	if !filepath.IsAbs(c.Root) {
		return nil, fmt.Errorf("build: expected an absolute root path, got %q", c.Root)
	}

	if c.Builder == "" {
		c.Builder = string(NameImage)
	}

	if c.Options == nil {
		c.Options = KindOptions{}
	}

	client, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}

	return &Builder{
		root:     c.Root,
		name:     c.Builder,
		options:  c.Options,
		auth:     c.Auth,
		buildEnv: c.BuildEnv,
		client:   client,
	}, nil
}

func (b *Builder) Close() error {
	return b.client.Close()
}

// Build runs the docker build.
//
// Depending on the configured `Config.Builder` the method verifies that
// the directory can be built and all necessary files exist.
//
// The method creates a Dockerfile depending on the configured builder
// and adds it to the tree, it passes the tree as the build context
// and initializes the build.
func (b *Builder) Build(ctx context.Context, taskID, version string) (*Response, error) {
	name := "task-" + SanitizeTaskID(taskID)
	uri := name + ":" + version
	if b.auth != nil {
		uri = b.auth.Repo + "/" + uri
	}

	patterns, err := ignore.DockerignorePatterns(b.root)
	if err != nil {
		return nil, err
	}
	tree, err := NewTree(TreeOptions{
		ExcludePatterns: patterns,
	})
	if err != nil {
		return nil, errors.Wrap(err, "new tree")
	}
	defer tree.Close()

	dockerfile, err := BuildDockerfile(DockerfileConfig{
		Builder: b.name,
		Root:    b.root,
		Options: b.options,
	})
	if err != nil {
		return nil, errors.Wrap(err, "creating dockerfile")
	}

	dockerfilePath := ".airplane/Dockerfile"
	if err := tree.MkdirAll(filepath.Dir(dockerfilePath)); err != nil {
		return nil, err
	}
	if err := tree.Write(dockerfilePath, strings.NewReader(dockerfile)); err != nil {
		return nil, err
	}

	if err := tree.Copy(b.root); err != nil {
		return nil, err
	}

	bc, err := tree.Archive()
	if err != nil {
		return nil, err
	}
	defer bc.Close()

	buildArgs := make(map[string]*string)
	for k, v := range b.buildEnv {
		value := v
		buildArgs[k] = &value
	}

	opts := types.ImageBuildOptions{
		Dockerfile:  dockerfilePath,
		Tags:        []string{uri},
		BuildArgs:   buildArgs,
		Platform:    "linux/amd64",
		AuthConfigs: b.authconfigs(),
	}

	resp, err := b.client.ImageBuild(ctx, bc, opts)
	if err != nil {
		return nil, errors.Wrap(err, "image build")
	}
	defer resp.Body.Close()

	scanner := bufiox.NewScanner(resp.Body)
	for scanner.Scan() {
		var event *dockerJSONMessage.JSONMessage
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			return nil, errors.Wrap(err, "unmarshalling docker build event")
		}

		if err := event.Display(os.Stderr, isatty.IsTerminal(os.Stderr.Fd())); err != nil {
			return nil, errors.Wrap(err, "docker build")
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, errors.Wrap(err, "scanning")
	}

	return &Response{
		ImageURL: uri,
	}, nil
}

// Push pushes the given image.
func (b *Builder) Push(ctx context.Context, uri string) error {
	if b.auth == nil {
		return errors.New("push requires registry auth")
	}

	authjson, err := json.Marshal(b.registryAuth())
	if err != nil {
		return err
	}

	resp, err := b.client.ImagePush(ctx, uri, types.ImagePushOptions{
		RegistryAuth: base64.URLEncoding.EncodeToString(authjson),
	})
	if err != nil {
		return err
	}
	defer resp.Close()

	scanner := bufiox.NewScanner(resp)
	for scanner.Scan() {
		var event *dockerJSONMessage.JSONMessage
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			return errors.Wrap(err, "unmarshalling docker build event")
		}

		if err := event.Display(os.Stderr, isatty.IsTerminal(os.Stderr.Fd())); err != nil {
			return errors.Wrap(err, "docker push")
		}
	}
	if err := scanner.Err(); err != nil {
		return errors.Wrap(err, "scanning")
	}

	return nil
}

// RegistryAuth returns the registry auth.
func (b *Builder) registryAuth() types.AuthConfig {
	return types.AuthConfig{
		Username: "oauth2accesstoken",
		Password: b.auth.Token,
	}
}

// Authconfigs returns the authconfigs to use.
func (b *Builder) authconfigs() map[string]types.AuthConfig {
	if b.auth == nil {
		return map[string]types.AuthConfig{}
	}

	return map[string]types.AuthConfig{
		b.auth.host(): b.registryAuth(),
	}
}

// SanitizeTaskID sanitizes the given task ID.
//
// Names may only contain lowercase letters, numbers, and
// hyphens, and must begin with a letter and end with a letter or number.
//
// We are planning to tweak our team/task ID generation to fit this:
// https://linear.app/airplane/issue/AIR-355/restrict-task-id-charset
//
// The following string manipulations won't matter for non-ksuid
// IDs (the current scheme).
func SanitizeTaskID(s string) string {
	s = strings.ToLower(s)
	if unicode.IsDigit(rune(s[len(s)-1])) {
		s = s[:len(s)-1] + "a"
	}
	return s
}

// TODO: this can merge with TaskKind
type Name string

const (
	NameGo         Name = "go"
	NameDeno       Name = "deno"
	NameImage      Name = "image"
	NamePython     Name = "python"
	NameNode       Name = "node"
	NameDockerfile Name = "dockerfile"
	NameShell      Name = "shell"

	NameSQL  Name = "sql"
	NameREST Name = "rest"
)

func NeedsBuilding(kind TaskKind) (bool, error) {
	switch Name(kind) {
	case NameGo, NameDeno, NamePython, NameNode, NameDockerfile, NameShell:
		return true, nil
	case NameImage, NameSQL, NameREST:
		return false, nil
	default:
		return false, errors.Errorf("NeedsBuilding got unexpected kind %s", kind)
	}
}

func BuildDockerfile(c DockerfileConfig) (string, error) {
	switch Name(c.Builder) {
	case NameGo:
		return golang(c.Root, c.Options)
	case NameDeno:
		return deno(c.Root, c.Options)
	case NamePython:
		return python(c.Root, c.Options)
	case NameNode:
		return node(c.Root, c.Options)
	case NameDockerfile:
		return dockerfile(c.Root, c.Options)
	case NameShell:
		return shell(c.Root, c.Options)
	default:
		return "", errors.Errorf("build: unknown builder type %q", c.Builder)
	}
}

func applyTemplate(t string, data interface{}) (string, error) {
	tmpl, err := template.New("airplane").Parse(t)
	if err != nil {
		return "", errors.Wrap(err, "parsing template")
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", errors.Wrap(err, "executing template")
	}

	return buf.String(), nil
}

func inlineString(s string) string {
	// To inline a multi-line string into a Dockerfile, insert `\n\` characters:
	s = strings.Join(strings.Split(s, "\n"), "\\n\\\n")
	// Since the string is wrapped in single-quotes, escape any single-quotes
	// inside of the target string.
	s = strings.ReplaceAll(s, "'", `'"'"'`)
	return "printf '" + s + "'"
}

// backslashEscape escapes s by replacing `\` with `\\` and all runes in chars with `\{rune}`.
// Typically should backslashEscape(s, `"`) to escape backslashes and double quotes.
func backslashEscape(s string, chars string) string {
	// Always escape backslash
	s = strings.ReplaceAll(s, `\`, `\\`)
	for _, char := range chars {
		s = strings.ReplaceAll(s, string(char), `\`+string(char))
	}
	return s
}
