package definitions

import (
	"context"
	_ "embed"
	"encoding/json"
	"os"
	"path"
	"strings"

	"github.com/airplanedev/lib/pkg/api"
	"github.com/airplanedev/lib/pkg/build"
	"github.com/alessio/shellescape"
	"github.com/flynn/go-shlex"
	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
	"github.com/xeipuuv/gojsonschema"
)

type Definition_0_3 struct {
	Name        string                    `json:"name"`
	Slug        string                    `json:"slug"`
	Description string                    `json:"description,omitempty"`
	Parameters  []ParameterDefinition_0_3 `json:"parameters,omitempty"`

	Deno       *DenoDefinition_0_3       `json:"deno,omitempty"`
	Dockerfile *DockerfileDefinition_0_3 `json:"dockerfile,omitempty"`
	Go         *GoDefinition_0_3         `json:"go,omitempty"`
	Image      *ImageDefinition_0_3      `json:"image,omitempty"`
	Node       *NodeDefinition_0_3       `json:"node,omitempty"`
	Python     *PythonDefinition_0_3     `json:"python,omitempty"`
	Shell      *ShellDefinition_0_3      `json:"shell,omitempty"`

	SQL  *SQLDefinition_0_3  `json:"sql,omitempty"`
	REST *RESTDefinition_0_3 `json:"rest,omitempty"`

	Constraints        *api.RunConstraints `json:"constraints,omitempty"`
	RequireRequests    bool                `json:"requireRequests,omitempty"`
	AllowSelfApprovals *bool               `json:"allowSelfApprovals,omitempty"`
	Timeout            int                 `json:"timeout,omitempty"`

	buildConfig build.BuildConfig
}

type taskKind_0_3 interface {
	fillInUpdateTaskRequest(context.Context, api.IAPIClient, *api.UpdateTaskRequest) error
	hydrateFromTask(context.Context, api.IAPIClient, *api.Task) error
	setEntrypoint(string) error
	getKindOptions() (build.KindOptions, error)
	getEntrypoint() (string, error)
	getRoot() (string, error)
	getEnv() (api.TaskEnv, error)
}

var _ taskKind_0_3 = &ImageDefinition_0_3{}

type ImageDefinition_0_3 struct {
	Image      string      `json:"image"`
	Entrypoint string      `json:"entrypoint,omitempty"`
	Command    []string    `json:"command"`
	Root       string      `json:"root,omitempty"`
	EnvVars    api.TaskEnv `json:"envVars,omitempty"`
}

func (d *ImageDefinition_0_3) fillInUpdateTaskRequest(ctx context.Context, client api.IAPIClient, req *api.UpdateTaskRequest) error {
	if d.Image != "" {
		req.Image = &d.Image
	}
	req.Arguments = d.Command
	if cmd, err := shlex.Split(d.Entrypoint); err != nil {
		return err
	} else {
		req.Command = cmd
	}
	return nil
}

func (d *ImageDefinition_0_3) hydrateFromTask(ctx context.Context, client api.IAPIClient, t *api.Task) error {
	if t.Image != nil {
		d.Image = *t.Image
	}
	d.Command = t.Arguments
	d.Entrypoint = shellescape.QuoteCommand(t.Command)
	return nil
}

func (d *ImageDefinition_0_3) setEntrypoint(entrypoint string) error {
	d.Entrypoint = entrypoint
	return nil
}

func (d *ImageDefinition_0_3) getKindOptions() (build.KindOptions, error) {
	return nil, nil
}

func (d *ImageDefinition_0_3) getEntrypoint() (string, error) {
	return "", ErrNoEntrypoint
}

func (d *ImageDefinition_0_3) getRoot() (string, error) {
	return d.Root, nil
}

func (d *ImageDefinition_0_3) getEnv() (api.TaskEnv, error) {
	return d.EnvVars, nil
}

var _ taskKind_0_3 = &DenoDefinition_0_3{}

type DenoDefinition_0_3 struct {
	Entrypoint string      `json:"entrypoint"`
	Root       string      `json:"root,omitempty"`
	EnvVars    api.TaskEnv `json:"envVars,omitempty"`
}

func (d *DenoDefinition_0_3) fillInUpdateTaskRequest(ctx context.Context, client api.IAPIClient, req *api.UpdateTaskRequest) error {
	return nil
}

func (d *DenoDefinition_0_3) hydrateFromTask(ctx context.Context, client api.IAPIClient, t *api.Task) error {
	if v, ok := t.KindOptions["entrypoint"]; ok {
		if sv, ok := v.(string); ok {
			d.Entrypoint = sv
		} else {
			return errors.Errorf("expected string entrypoint, got %T instead", v)
		}
	}
	return nil
}

func (d *DenoDefinition_0_3) setEntrypoint(entrypoint string) error {
	d.Entrypoint = entrypoint
	return nil
}

func (d *DenoDefinition_0_3) getKindOptions() (build.KindOptions, error) {
	return build.KindOptions{
		"entrypoint": d.Entrypoint,
	}, nil
}

func (d *DenoDefinition_0_3) getEntrypoint() (string, error) {
	return d.Entrypoint, nil
}

func (d *DenoDefinition_0_3) getRoot() (string, error) {
	return d.Root, nil
}

func (d *DenoDefinition_0_3) getEnv() (api.TaskEnv, error) {
	return d.EnvVars, nil
}

var _ taskKind_0_3 = &DockerfileDefinition_0_3{}

type DockerfileDefinition_0_3 struct {
	Dockerfile string      `json:"dockerfile"`
	Root       string      `json:"root,omitempty"`
	EnvVars    api.TaskEnv `json:"envVars,omitempty"`
}

func (d *DockerfileDefinition_0_3) fillInUpdateTaskRequest(ctx context.Context, client api.IAPIClient, req *api.UpdateTaskRequest) error {
	return nil
}

func (d *DockerfileDefinition_0_3) hydrateFromTask(ctx context.Context, client api.IAPIClient, t *api.Task) error {
	if v, ok := t.KindOptions["dockerfile"]; ok {
		if sv, ok := v.(string); ok {
			d.Dockerfile = sv
		} else {
			return errors.Errorf("expected string dockerfile, got %T instead", v)
		}
	}
	return nil
}

func (d *DockerfileDefinition_0_3) setEntrypoint(entrypoint string) error {
	return ErrNoEntrypoint
}

func (d *DockerfileDefinition_0_3) getKindOptions() (build.KindOptions, error) {
	return build.KindOptions{
		"dockerfile": d.Dockerfile,
	}, nil
}

func (d *DockerfileDefinition_0_3) getEntrypoint() (string, error) {
	return "", ErrNoEntrypoint
}

func (d *DockerfileDefinition_0_3) getRoot() (string, error) {
	return d.Root, nil
}

func (d *DockerfileDefinition_0_3) getEnv() (api.TaskEnv, error) {
	return d.EnvVars, nil
}

var _ taskKind_0_3 = &GoDefinition_0_3{}

type GoDefinition_0_3 struct {
	Entrypoint string      `json:"entrypoint"`
	Root       string      `json:"root,omitempty"`
	EnvVars    api.TaskEnv `json:"envVars,omitempty"`
}

func (d *GoDefinition_0_3) fillInUpdateTaskRequest(ctx context.Context, client api.IAPIClient, req *api.UpdateTaskRequest) error {
	return nil
}

func (d *GoDefinition_0_3) hydrateFromTask(ctx context.Context, client api.IAPIClient, t *api.Task) error {
	if v, ok := t.KindOptions["entrypoint"]; ok {
		if sv, ok := v.(string); ok {
			d.Entrypoint = sv
		} else {
			return errors.Errorf("expected string entrypoint, got %T instead", v)
		}
	}
	return nil
}

func (d *GoDefinition_0_3) setEntrypoint(entrypoint string) error {
	d.Entrypoint = entrypoint
	return nil
}

func (d *GoDefinition_0_3) getKindOptions() (build.KindOptions, error) {
	return build.KindOptions{
		"entrypoint": d.Entrypoint,
	}, nil
}

func (d *GoDefinition_0_3) getEntrypoint() (string, error) {
	return d.Entrypoint, nil
}

func (d *GoDefinition_0_3) getRoot() (string, error) {
	return d.Root, nil
}

func (d *GoDefinition_0_3) getEnv() (api.TaskEnv, error) {
	return d.EnvVars, nil
}

var _ taskKind_0_3 = &NodeDefinition_0_3{}

type NodeDefinition_0_3 struct {
	Entrypoint  string      `json:"entrypoint"`
	NodeVersion string      `json:"nodeVersion"`
	Root        string      `json:"root,omitempty"`
	EnvVars     api.TaskEnv `json:"envVars,omitempty"`
}

func (d *NodeDefinition_0_3) fillInUpdateTaskRequest(ctx context.Context, client api.IAPIClient, req *api.UpdateTaskRequest) error {
	return nil
}

func (d *NodeDefinition_0_3) hydrateFromTask(ctx context.Context, client api.IAPIClient, t *api.Task) error {
	if v, ok := t.KindOptions["entrypoint"]; ok {
		if sv, ok := v.(string); ok {
			d.Entrypoint = sv
		} else {
			return errors.Errorf("expected string entrypoint, got %T instead", v)
		}
	}
	if v, ok := t.KindOptions["nodeVersion"]; ok {
		if sv, ok := v.(string); ok {
			d.NodeVersion = sv
		} else {
			return errors.Errorf("expected string nodeVersion, got %T instead", v)
		}
	}
	return nil
}

func (d *NodeDefinition_0_3) setEntrypoint(entrypoint string) error {
	d.Entrypoint = entrypoint
	return nil
}

func (d *NodeDefinition_0_3) getKindOptions() (build.KindOptions, error) {
	return build.KindOptions{
		"entrypoint":  d.Entrypoint,
		"nodeVersion": d.NodeVersion,
	}, nil
}

func (d *NodeDefinition_0_3) getEntrypoint() (string, error) {
	return d.Entrypoint, nil
}

func (d *NodeDefinition_0_3) getRoot() (string, error) {
	return d.Root, nil
}

func (d *NodeDefinition_0_3) getEnv() (api.TaskEnv, error) {
	return d.EnvVars, nil
}

var _ taskKind_0_3 = &PythonDefinition_0_3{}

type PythonDefinition_0_3 struct {
	Entrypoint string      `json:"entrypoint"`
	Root       string      `json:"root,omitempty"`
	EnvVars    api.TaskEnv `json:"envVars,omitempty"`
}

func (d *PythonDefinition_0_3) fillInUpdateTaskRequest(ctx context.Context, client api.IAPIClient, req *api.UpdateTaskRequest) error {
	return nil
}

func (d *PythonDefinition_0_3) hydrateFromTask(ctx context.Context, client api.IAPIClient, t *api.Task) error {
	if v, ok := t.KindOptions["entrypoint"]; ok {
		if sv, ok := v.(string); ok {
			d.Entrypoint = sv
		} else {
			return errors.Errorf("expected string entrypoint, got %T instead", v)
		}
	}
	return nil
}

func (d *PythonDefinition_0_3) setEntrypoint(entrypoint string) error {
	d.Entrypoint = entrypoint
	return nil
}

func (d *PythonDefinition_0_3) getKindOptions() (build.KindOptions, error) {
	return build.KindOptions{
		"entrypoint": d.Entrypoint,
	}, nil
}

func (d *PythonDefinition_0_3) getEntrypoint() (string, error) {
	return d.Entrypoint, nil
}

func (d *PythonDefinition_0_3) getRoot() (string, error) {
	return d.Root, nil
}

func (d *PythonDefinition_0_3) getEnv() (api.TaskEnv, error) {
	return d.EnvVars, nil
}

var _ taskKind_0_3 = &ShellDefinition_0_3{}

type ShellDefinition_0_3 struct {
	Entrypoint string      `json:"entrypoint"`
	Root       string      `json:"root,omitempty"`
	EnvVars    api.TaskEnv `json:"envVars,omitempty"`
}

func (d *ShellDefinition_0_3) fillInUpdateTaskRequest(ctx context.Context, client api.IAPIClient, req *api.UpdateTaskRequest) error {
	return nil
}

func (d *ShellDefinition_0_3) hydrateFromTask(ctx context.Context, client api.IAPIClient, t *api.Task) error {
	if v, ok := t.KindOptions["entrypoint"]; ok {
		if sv, ok := v.(string); ok {
			d.Entrypoint = sv
		} else {
			return errors.Errorf("expected string entrypoint, got %T instead", v)
		}
	}
	return nil
}

func (d *ShellDefinition_0_3) setEntrypoint(entrypoint string) error {
	d.Entrypoint = entrypoint
	return nil
}

func (d *ShellDefinition_0_3) getKindOptions() (build.KindOptions, error) {
	return build.KindOptions{
		"entrypoint": d.Entrypoint,
	}, nil
}

func (d *ShellDefinition_0_3) getEntrypoint() (string, error) {
	return d.Entrypoint, nil
}

func (d *ShellDefinition_0_3) getRoot() (string, error) {
	return d.Root, nil
}

func (d *ShellDefinition_0_3) getEnv() (api.TaskEnv, error) {
	return d.EnvVars, nil
}

var _ taskKind_0_3 = &SQLDefinition_0_3{}

type SQLDefinition_0_3 struct {
	Resource   string                 `json:"resource"`
	Entrypoint string                 `json:"entrypoint"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`

	// Contents of Entrypoint, cached
	entrypointContents string `json:"-"`
}

func (d *SQLDefinition_0_3) GetQuery() (string, error) {
	if d.entrypointContents == "" {
		queryBytes, err := os.ReadFile(d.Entrypoint)
		if err != nil {
			return "", errors.Wrapf(err, "reading SQL entrypoint %s", d.Entrypoint)
		}
		d.entrypointContents = string(queryBytes)
	}
	return d.entrypointContents, nil
}

func (d *SQLDefinition_0_3) fillInUpdateTaskRequest(ctx context.Context, client api.IAPIClient, req *api.UpdateTaskRequest) error {
	resourcesByName, err := getResourcesByName(ctx, client)
	if err != nil {
		return err
	}
	if res, ok := resourcesByName[d.Resource]; ok {
		req.Resources = map[string]string{
			"db": res.ID,
		}
	} else {
		return errors.Errorf("unknown resource: %s", d.Resource)
	}
	return nil
}

func (d *SQLDefinition_0_3) hydrateFromTask(ctx context.Context, client api.IAPIClient, t *api.Task) error {
	if resID, ok := t.Resources["db"]; ok {
		resourcesByID, err := getResourcesByID(ctx, client)
		if err != nil {
			return err
		}
		if res, ok := resourcesByID[resID]; ok {
			d.Resource = res.Name
		}
	}
	if v, ok := t.KindOptions["entrypoint"]; ok {
		if sv, ok := v.(string); ok {
			d.Entrypoint = sv
		} else {
			return errors.Errorf("expected string entrypoint, got %T instead", v)
		}
	}
	if v, ok := t.KindOptions["query"]; ok {
		if sv, ok := v.(string); ok {
			d.entrypointContents = sv
		} else {
			return errors.Errorf("expected string query, got %T instead", v)
		}
	}
	if v, ok := t.KindOptions["queryArgs"]; ok {
		if mv, ok := v.(map[string]interface{}); ok {
			d.Parameters = mv
		} else {
			return errors.Errorf("expected map queryArgs, got %T instead", v)
		}
	}
	return nil
}

func (d *SQLDefinition_0_3) setEntrypoint(entrypoint string) error {
	d.Entrypoint = entrypoint
	return nil
}

func (d *SQLDefinition_0_3) getKindOptions() (build.KindOptions, error) {
	query, err := d.GetQuery()
	if err != nil {
		return nil, err
	}
	if d.Parameters == nil {
		d.Parameters = map[string]interface{}{}
	}
	return build.KindOptions{
		"entrypoint": d.Entrypoint,
		"query":      query,
		"queryArgs":  d.Parameters,
	}, nil
}

func (d *SQLDefinition_0_3) getEntrypoint() (string, error) {
	return d.Entrypoint, nil
}

func (d *SQLDefinition_0_3) getRoot() (string, error) {
	return "", nil
}

func (d *SQLDefinition_0_3) getEnv() (api.TaskEnv, error) {
	return nil, nil
}

var _ taskKind_0_3 = &RESTDefinition_0_3{}

type RESTDefinition_0_3 struct {
	Resource  string                 `json:"resource"`
	Method    string                 `json:"method"`
	Path      string                 `json:"path"`
	URLParams map[string]interface{} `json:"urlParams,omitempty"`
	Headers   map[string]interface{} `json:"headers,omitempty"`
	BodyType  string                 `json:"bodyType"`
	Body      string                 `json:"body,omitempty"`
	FormData  map[string]interface{} `json:"formData,omitempty"`
}

func (d *RESTDefinition_0_3) fillInUpdateTaskRequest(ctx context.Context, client api.IAPIClient, req *api.UpdateTaskRequest) error {
	resourcesByName, err := getResourcesByName(ctx, client)
	if err != nil {
		return err
	}
	if res, ok := resourcesByName[d.Resource]; ok {
		req.Resources = map[string]string{
			"rest": res.ID,
		}
	} else {
		return errors.Errorf("unknown resource: %s", d.Resource)
	}
	return nil
}

func (d *RESTDefinition_0_3) hydrateFromTask(ctx context.Context, client api.IAPIClient, t *api.Task) error {
	if resID, ok := t.Resources["rest"]; ok {
		resourcesByID, err := getResourcesByID(ctx, client)
		if err != nil {
			return err
		}
		if res, ok := resourcesByID[resID]; ok {
			d.Resource = res.Name
		}
	}
	if v, ok := t.KindOptions["method"]; ok {
		if sv, ok := v.(string); ok {
			d.Method = sv
		} else {
			return errors.Errorf("expected string method, got %T instead", v)
		}
	}
	if v, ok := t.KindOptions["path"]; ok {
		if sv, ok := v.(string); ok {
			d.Path = sv
		} else {
			return errors.Errorf("expected string path, got %T instead", v)
		}
	}
	if v, ok := t.KindOptions["urlParams"]; ok {
		if mv, ok := v.(map[string]interface{}); ok {
			d.URLParams = mv
		} else {
			return errors.Errorf("expected map urlParams, got %T instead", v)
		}
	}
	if v, ok := t.KindOptions["headers"]; ok {
		if mv, ok := v.(map[string]interface{}); ok {
			d.Headers = mv
		} else {
			return errors.Errorf("expected map headers, got %T instead", v)
		}
	}
	if v, ok := t.KindOptions["bodyType"]; ok {
		if sv, ok := v.(string); ok {
			d.BodyType = sv
		} else {
			return errors.Errorf("expected string bodyType, got %T instead", v)
		}
	}
	if v, ok := t.KindOptions["body"]; ok {
		if sv, ok := v.(string); ok {
			d.Body = sv
		} else {
			return errors.Errorf("expected string body, got %T instead", v)
		}
	}
	if v, ok := t.KindOptions["formData"]; ok {
		if mv, ok := v.(map[string]interface{}); ok {
			d.FormData = mv
		} else {
			return errors.Errorf("expected map formData, got %T instead", v)
		}
	}
	return nil
}

func (d *RESTDefinition_0_3) setEntrypoint(entrypoint string) error {
	return ErrNoEntrypoint
}

func (d *RESTDefinition_0_3) getKindOptions() (build.KindOptions, error) {
	if d.URLParams == nil {
		d.URLParams = map[string]interface{}{}
	}
	if d.Headers == nil {
		d.Headers = map[string]interface{}{}
	}
	if d.FormData == nil {
		d.FormData = map[string]interface{}{}
	}
	return build.KindOptions{
		"method":    d.Method,
		"path":      d.Path,
		"urlParams": d.URLParams,
		"headers":   d.Headers,
		"bodyType":  d.BodyType,
		"body":      d.Body,
		"formData":  d.FormData,
	}, nil
}

func (d *RESTDefinition_0_3) getEntrypoint() (string, error) {
	return "", ErrNoEntrypoint
}

func (d *RESTDefinition_0_3) getRoot() (string, error) {
	return "", nil
}

func (d *RESTDefinition_0_3) getEnv() (api.TaskEnv, error) {
	return nil, nil
}

type ParameterDefinition_0_3 struct {
	Name        string                 `json:"name"`
	Slug        string                 `json:"slug"`
	Type        string                 `json:"type"`
	Description string                 `json:"description,omitempty"`
	Default     interface{}            `json:"default,omitempty"`
	Required    *bool                  `json:"required,omitempty"`
	Options     []OptionDefinition_0_3 `json:"options,omitempty"`
	Regex       string                 `json:"regex,omitempty"`
}

type OptionDefinition_0_3 struct {
	Label string      `json:"label"`
	Value interface{} `json:"value"`
}

var _ json.Unmarshaler = &OptionDefinition_0_3{}

func (o *OptionDefinition_0_3) UnmarshalJSON(b []byte) error {
	// If it's just a string, dump it in the value field.
	var value string
	if err := json.Unmarshal(b, &value); err == nil {
		o.Value = value
		return nil
	}

	// Otherwise, perform a normal unmarshal operation.
	// Note we need a new type, otherwise we recursively call this
	// method and end up stack overflowing.
	type option OptionDefinition_0_3
	var opt option
	if err := json.Unmarshal(b, &opt); err != nil {
		return err
	}
	*o = OptionDefinition_0_3(opt)

	return nil
}

//go:embed schema_0_3.json
var schemaStr string

func NewDefinition_0_3(name string, slug string, kind build.TaskKind, entrypoint string) (Definition_0_3, error) {
	def := Definition_0_3{
		Name: name,
		Slug: slug,
	}

	switch kind {
	case build.TaskKindDeno:
		def.Deno = &DenoDefinition_0_3{
			Entrypoint: entrypoint,
		}
	case build.TaskKindDockerfile:
		def.Dockerfile = &DockerfileDefinition_0_3{
			Dockerfile: entrypoint,
		}
	case build.TaskKindGo:
		def.Go = &GoDefinition_0_3{
			Entrypoint: entrypoint,
		}
	case build.TaskKindImage:
		def.Image = &ImageDefinition_0_3{}
	case build.TaskKindNode:
		def.Node = &NodeDefinition_0_3{
			Entrypoint:  entrypoint,
			NodeVersion: "14",
		}
	case build.TaskKindPython:
		def.Python = &PythonDefinition_0_3{
			Entrypoint: entrypoint,
		}
	case build.TaskKindShell:
		def.Shell = &ShellDefinition_0_3{
			Entrypoint: entrypoint,
		}
	case build.TaskKindSQL:
		def.SQL = &SQLDefinition_0_3{
			Entrypoint: entrypoint,
		}
	case build.TaskKindREST:
		def.REST = &RESTDefinition_0_3{}
	default:
		return Definition_0_3{}, errors.Errorf("unknown kind: %s", kind)
	}

	return def, nil
}

func (d Definition_0_3) Marshal(format TaskDefFormat) ([]byte, error) {
	buf, err := json.MarshalIndent(d, "", "\t")
	if err != nil {
		return nil, err
	}

	switch format {
	case TaskDefFormatYAML:
		buf, err = yaml.JSONToYAML(buf)
		if err != nil {
			return nil, err
		}
	case TaskDefFormatJSON:
		// nothing
	default:
		return nil, errors.Errorf("unknown format: %s", format)
	}

	return buf, nil
}

func (d *Definition_0_3) Unmarshal(format TaskDefFormat, buf []byte) error {
	var err error
	switch format {
	case TaskDefFormatYAML:
		buf, err = yaml.YAMLToJSON(buf)
		if err != nil {
			return err
		}
	case TaskDefFormatJSON:
		// nothing
	default:
		return errors.Errorf("unknown format: %s", format)
	}

	schemaLoader := gojsonschema.NewStringLoader(schemaStr)
	docLoader := gojsonschema.NewBytesLoader(buf)

	result, err := gojsonschema.Validate(schemaLoader, docLoader)
	if err != nil {
		return errors.Wrap(err, "validating schema")
	}

	if !result.Valid() {
		return errors.WithStack(ErrSchemaValidation{Errors: result.Errors()})
	}

	if err = json.Unmarshal(buf, &d); err != nil {
		return err
	}
	return nil
}

func (d Definition_0_3) Kind() (build.TaskKind, error) {
	if d.Deno != nil {
		return build.TaskKindDeno, nil
	} else if d.Dockerfile != nil {
		return build.TaskKindDockerfile, nil
	} else if d.Go != nil {
		return build.TaskKindGo, nil
	} else if d.Image != nil {
		return build.TaskKindImage, nil
	} else if d.Node != nil {
		return build.TaskKindNode, nil
	} else if d.Python != nil {
		return build.TaskKindPython, nil
	} else if d.Shell != nil {
		return build.TaskKindShell, nil
	} else if d.SQL != nil {
		return build.TaskKindSQL, nil
	} else if d.REST != nil {
		return build.TaskKindREST, nil
	} else {
		return "", errors.New("incomplete task definition")
	}
}

func (d Definition_0_3) taskKind() (taskKind_0_3, error) {
	if d.Deno != nil {
		return d.Deno, nil
	} else if d.Dockerfile != nil {
		return d.Dockerfile, nil
	} else if d.Go != nil {
		return d.Go, nil
	} else if d.Image != nil {
		return d.Image, nil
	} else if d.Node != nil {
		return d.Node, nil
	} else if d.Python != nil {
		return d.Python, nil
	} else if d.Shell != nil {
		return d.Shell, nil
	} else if d.SQL != nil {
		return d.SQL, nil
	} else if d.REST != nil {
		return d.REST, nil
	} else {
		return nil, errors.New("incomplete task definition")
	}
}

func (d Definition_0_3) GetUpdateTaskRequest(ctx context.Context, client api.IAPIClient, currentTask *api.Task) (api.UpdateTaskRequest, error) {
	req := api.UpdateTaskRequest{
		Slug:        d.Slug,
		Name:        d.Name,
		Description: d.Description,
		Timeout:     d.Timeout,
	}

	if err := d.addParametersToUpdateTaskRequest(ctx, &req); err != nil {
		return api.UpdateTaskRequest{}, err
	}

	if d.Constraints != nil && !d.Constraints.IsEmpty() {
		req.Constraints = *d.Constraints
	}

	if d.RequireRequests {
		req.ExecuteRules.RequireRequests = true
	}
	if d.AllowSelfApprovals != nil && !*d.AllowSelfApprovals {
		req.ExecuteRules.DisallowSelfApprove = true
	}

	if err := d.addKindSpecificsToUpdateTaskRequest(ctx, client, &req); err != nil {
		return api.UpdateTaskRequest{}, err
	}

	if currentTask != nil {
		req.RequireExplicitPermissions = currentTask.RequireExplicitPermissions
		req.Permissions = currentTask.Permissions
	}

	return req, nil
}

func (d Definition_0_3) addParametersToUpdateTaskRequest(ctx context.Context, req *api.UpdateTaskRequest) error {
	req.Parameters = make([]api.Parameter, len(d.Parameters))
	for i, pd := range d.Parameters {
		param := api.Parameter{
			Name:    pd.Name,
			Slug:    pd.Slug,
			Desc:    pd.Description,
			Default: pd.Default,
		}

		switch pd.Type {
		case "shorttext":
			param.Type = "string"
		case "longtext":
			param.Type = "string"
			param.Component = api.ComponentTextarea
		case "sql":
			param.Type = "string"
			param.Component = api.ComponentEditorSQL
		case "boolean", "upload", "integer", "float", "date", "datetime", "configvar":
			param.Type = api.Type(pd.Type)
		default:
			return errors.Errorf("unknown parameter type: %s", pd.Type)
		}

		if pd.Required != nil && !*pd.Required {
			param.Constraints.Optional = true
		}

		param.Constraints.Regex = pd.Regex

		if len(pd.Options) > 0 {
			param.Constraints.Options = make([]api.ConstraintOption, len(pd.Options))
			for j, od := range pd.Options {
				param.Constraints.Options[j].Label = od.Label
				param.Constraints.Options[j].Value = od.Value
			}
		}

		req.Parameters[i] = param
	}
	return nil
}

func (d Definition_0_3) addKindSpecificsToUpdateTaskRequest(ctx context.Context, client api.IAPIClient, req *api.UpdateTaskRequest) error {
	resourcesByName := map[string]api.Resource{}
	if d.SQL != nil || d.REST != nil {
		// Remap resources from ref -> name to ref -> id.
		resp, err := client.ListResources(ctx)
		if err != nil {
			return errors.Wrap(err, "fetching resources")
		}
		for _, resource := range resp.Resources {
			resourcesByName[resource.Name] = resource
		}
	}

	kind, options, err := d.GetKindAndOptions()
	if err != nil {
		return err
	}
	req.Kind = kind
	req.KindOptions = options

	env, err := d.GetEnv()
	if err != nil {
		return err
	}
	req.Env = env

	taskKind, err := d.taskKind()
	if err != nil {
		return err
	}
	if err := taskKind.fillInUpdateTaskRequest(ctx, client, req); err != nil {
		return err
	}
	return nil
}

func (d Definition_0_3) Root(dir string) (string, error) {
	taskKind, err := d.taskKind()
	if err != nil {
		return "", err
	}
	root, err := taskKind.getRoot()
	if err != nil {
		return "", err
	}
	return path.Join(dir, root), nil
}

var ErrNoEntrypoint = errors.New("No entrypoint")

func (d Definition_0_3) Entrypoint() (string, error) {
	taskKind, err := d.taskKind()
	if err != nil {
		return "", err
	}
	return taskKind.getEntrypoint()
}

func (d *Definition_0_3) UpgradeJST() error {
	return nil
}

func (d *Definition_0_3) GetKindAndOptions() (build.TaskKind, build.KindOptions, error) {
	kind, err := d.Kind()
	if err != nil {
		return "", nil, err
	}

	taskKind, err := d.taskKind()
	if err != nil {
		return "", nil, err
	}

	options, err := taskKind.getKindOptions()
	if err != nil {
		return "", nil, err
	}

	return kind, options, nil
}

func (d *Definition_0_3) GetEnv() (api.TaskEnv, error) {
	taskKind, err := d.taskKind()
	if err != nil {
		return nil, err
	}
	return taskKind.getEnv()
}

func (d *Definition_0_3) GetSlug() string {
	return d.Slug
}

func (d *Definition_0_3) SetEntrypoint(entrypoint string) error {
	taskKind, err := d.taskKind()
	if err != nil {
		return err
	}

	return taskKind.setEntrypoint(entrypoint)
}

func (d *Definition_0_3) SetWorkdir(taskroot, workdir string) error {
	// TODO: currently only a concept on Node - should be generalized to all builders.
	if d.Node == nil {
		return nil
	}

	d.SetBuildConfig("workdir", strings.TrimPrefix(workdir, taskroot))

	return nil
}

func NewDefinitionFromTask_0_3(ctx context.Context, client api.IAPIClient, t api.Task) (Definition_0_3, error) {
	d := Definition_0_3{
		Name:            t.Name,
		Slug:            t.Slug,
		Description:     t.Description,
		RequireRequests: t.ExecuteRules.RequireRequests,
		Timeout:         t.Timeout,
	}

	if err := d.convertParametersFromTask(ctx, client, &t); err != nil {
		return Definition_0_3{}, err
	}

	if err := d.convertTaskKindFromTask(ctx, client, &t); err != nil {
		return Definition_0_3{}, err
	}

	if !t.Constraints.IsEmpty() {
		d.Constraints = &t.Constraints
	}

	if t.ExecuteRules.DisallowSelfApprove {
		v := false
		d.AllowSelfApprovals = &v
	}

	return d, nil
}

func (d *Definition_0_3) convertParametersFromTask(ctx context.Context, client api.IAPIClient, t *api.Task) error {
	if len(t.Parameters) == 0 {
		return nil
	}
	d.Parameters = make([]ParameterDefinition_0_3, len(t.Parameters))
	for idx, param := range t.Parameters {
		p := ParameterDefinition_0_3{
			Name:        param.Name,
			Slug:        param.Slug,
			Description: param.Desc,
			Default:     param.Default,
		}

		switch param.Type {
		case "string":
			switch param.Component {
			case api.ComponentTextarea:
				p.Type = "longtext"
			case api.ComponentEditorSQL:
				p.Type = "sql"
			case api.ComponentNone:
				p.Type = "shorttext"
			default:
				return errors.Errorf("unexpected component for type=string: %s", param.Component)
			}
		case "boolean", "upload", "integer", "float", "date", "datetime", "configvar":
			p.Type = string(param.Type)
		default:
			return errors.Errorf("unknown parameter type: %s", param.Type)
		}

		if param.Constraints.Optional {
			v := false
			p.Required = &v
		}

		p.Regex = param.Constraints.Regex

		if len(param.Constraints.Options) > 0 {
			p.Options = make([]OptionDefinition_0_3, len(param.Constraints.Options))
			for j, opt := range param.Constraints.Options {
				p.Options[j] = OptionDefinition_0_3{
					Label: opt.Label,
					Value: opt.Value,
				}
			}
		}

		d.Parameters[idx] = p
	}
	return nil
}

func (d *Definition_0_3) convertTaskKindFromTask(ctx context.Context, client api.IAPIClient, t *api.Task) error {
	switch t.Kind {
	case build.TaskKindDeno:
		d.Deno = &DenoDefinition_0_3{}
		return d.Deno.hydrateFromTask(ctx, client, t)
	case build.TaskKindDockerfile:
		d.Dockerfile = &DockerfileDefinition_0_3{}
		return d.Dockerfile.hydrateFromTask(ctx, client, t)
	case build.TaskKindGo:
		d.Go = &GoDefinition_0_3{}
		return d.Go.hydrateFromTask(ctx, client, t)
	case build.TaskKindImage:
		d.Image = &ImageDefinition_0_3{}
		return d.Image.hydrateFromTask(ctx, client, t)
	case build.TaskKindNode:
		d.Node = &NodeDefinition_0_3{}
		return d.Node.hydrateFromTask(ctx, client, t)
	case build.TaskKindPython:
		d.Python = &PythonDefinition_0_3{}
		return d.Python.hydrateFromTask(ctx, client, t)
	case build.TaskKindShell:
		d.Shell = &ShellDefinition_0_3{}
		return d.Shell.hydrateFromTask(ctx, client, t)
	case build.TaskKindSQL:
		d.SQL = &SQLDefinition_0_3{}
		return d.SQL.hydrateFromTask(ctx, client, t)
	case build.TaskKindREST:
		d.REST = &RESTDefinition_0_3{}
		return d.REST.hydrateFromTask(ctx, client, t)
	default:
		return errors.Errorf("unknown task kind: %s", t.Kind)
	}
}

func (d *Definition_0_3) GetBuildConfig() (build.BuildConfig, error) {
	config := build.BuildConfig{}

	_, options, err := d.GetKindAndOptions()
	if err != nil {
		return nil, err
	}
	for key, val := range options {
		config[key] = val
	}

	for key, val := range d.buildConfig {
		if val == nil { // Nil masks out the value.
			delete(config, key)
		} else {
			config[key] = val
		}
	}

	return config, nil
}

func (d *Definition_0_3) SetBuildConfig(key string, value interface{}) {
	if d.buildConfig == nil {
		d.buildConfig = map[string]interface{}{}
	}
	d.buildConfig[key] = value
}

func getResourcesByName(ctx context.Context, client api.IAPIClient) (map[string]api.Resource, error) {
	resp, err := client.ListResources(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "fetching resources")
	}
	resourcesByName := map[string]api.Resource{}
	for _, resource := range resp.Resources {
		resourcesByName[resource.Name] = resource
	}
	return resourcesByName, nil
}

func getResourcesByID(ctx context.Context, client api.IAPIClient) (map[string]api.Resource, error) {
	resp, err := client.ListResources(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "fetching resources")
	}
	resourcesByID := map[string]api.Resource{}
	for _, resource := range resp.Resources {
		resourcesByID[resource.ID] = resource
	}
	return resourcesByID, nil
}
