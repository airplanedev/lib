package definitions

import (
	"context"
	_ "embed"
	"encoding/json"
	"os"
	"path"

	"github.com/airplanedev/lib/pkg/api"
	"github.com/airplanedev/lib/pkg/build"
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

	Permissions *PermissionDefinition_0_3 `json:"permissions,omitempty"`
	Constraints *api.RunConstraints       `json:"constraints,omitempty"`
	Timeout     int                       `json:"timeout,omitempty"`
}

type taskKind_0_3 interface {
	fillInUpdateTaskRequest(context.Context, api.IAPIClient, *api.UpdateTaskRequest) error
	upgradeJST() error
	getKindOptions() (build.KindOptions, error)
	getEntrypoint() (string, error)
	getRoot() (string, error)
	getEnv() (api.TaskEnv, error)
}

var _ taskKind_0_3 = &ImageDefinition_0_3{}

type ImageDefinition_0_3 struct {
	Image   string      `json:"image"`
	Command []string    `json:"command"`
	Root    string      `json:"root,omitempty"`
	Env     api.TaskEnv `json:"env,omitempty"`
}

func (d *ImageDefinition_0_3) fillInUpdateTaskRequest(ctx context.Context, client api.IAPIClient, req *api.UpdateTaskRequest) error {
	req.Image = &d.Image
	req.Command = d.Command
	return nil
}

func (d *ImageDefinition_0_3) upgradeJST() error {
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
	return d.Env, nil
}

var _ taskKind_0_3 = &DenoDefinition_0_3{}

type DenoDefinition_0_3 struct {
	Entrypoint string      `json:"entrypoint"`
	Arguments  []string    `json:"arguments,omitempty"`
	Root       string      `json:"root,omitempty"`
	Env        api.TaskEnv `json:"env,omitempty"`
}

func (d *DenoDefinition_0_3) fillInUpdateTaskRequest(ctx context.Context, client api.IAPIClient, req *api.UpdateTaskRequest) error {
	req.Arguments = d.Arguments
	return nil
}

func (d *DenoDefinition_0_3) upgradeJST() error {
	d.Arguments = upgradeArguments(d.Arguments)
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
	return d.Env, nil
}

var _ taskKind_0_3 = &DockerfileDefinition_0_3{}

type DockerfileDefinition_0_3 struct {
	Dockerfile string      `json:"dockerfile"`
	Root       string      `json:"root,omitempty"`
	Env        api.TaskEnv `json:"env,omitempty"`
}

func (d *DockerfileDefinition_0_3) fillInUpdateTaskRequest(ctx context.Context, client api.IAPIClient, req *api.UpdateTaskRequest) error {
	return nil
}

func (d *DockerfileDefinition_0_3) upgradeJST() error {
	return nil
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
	return d.Env, nil
}

var _ taskKind_0_3 = &GoDefinition_0_3{}

type GoDefinition_0_3 struct {
	Entrypoint string      `json:"entrypoint"`
	Arguments  []string    `json:"arguments,omitempty"`
	Root       string      `json:"root,omitempty"`
	Env        api.TaskEnv `json:"env,omitempty"`
}

func (d *GoDefinition_0_3) fillInUpdateTaskRequest(ctx context.Context, client api.IAPIClient, req *api.UpdateTaskRequest) error {
	req.Arguments = d.Arguments
	return nil
}

func (d *GoDefinition_0_3) upgradeJST() error {
	d.Arguments = upgradeArguments(d.Arguments)
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
	return d.Env, nil
}

var _ taskKind_0_3 = &NodeDefinition_0_3{}

type NodeDefinition_0_3 struct {
	Entrypoint  string      `json:"entrypoint"`
	NodeVersion string      `json:"nodeVersion"`
	Arguments   []string    `json:"arguments,omitempty"`
	Root        string      `json:"root,omitempty"`
	Env         api.TaskEnv `json:"env,omitempty"`
}

func (d *NodeDefinition_0_3) fillInUpdateTaskRequest(ctx context.Context, client api.IAPIClient, req *api.UpdateTaskRequest) error {
	req.Arguments = d.Arguments
	return nil
}

func (d *NodeDefinition_0_3) upgradeJST() error {
	d.Arguments = upgradeArguments(d.Arguments)
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
	return d.Env, nil
}

var _ taskKind_0_3 = &PythonDefinition_0_3{}

type PythonDefinition_0_3 struct {
	Entrypoint string      `json:"entrypoint"`
	Arguments  []string    `json:"arguments,omitempty"`
	Root       string      `json:"root,omitempty"`
	Env        api.TaskEnv `json:"env,omitempty"`
}

func (d *PythonDefinition_0_3) fillInUpdateTaskRequest(ctx context.Context, client api.IAPIClient, req *api.UpdateTaskRequest) error {
	req.Arguments = d.Arguments
	return nil
}

func (d *PythonDefinition_0_3) upgradeJST() error {
	d.Arguments = upgradeArguments(d.Arguments)
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
	return d.Env, nil
}

var _ taskKind_0_3 = &ShellDefinition_0_3{}

type ShellDefinition_0_3 struct {
	Entrypoint string      `json:"entrypoint"`
	Arguments  []string    `json:"arguments,omitempty"`
	Root       string      `json:"root,omitempty"`
	Env        api.TaskEnv `json:"env,omitempty"`
}

func (d *ShellDefinition_0_3) fillInUpdateTaskRequest(ctx context.Context, client api.IAPIClient, req *api.UpdateTaskRequest) error {
	req.Arguments = d.Arguments
	return nil
}

func (d *ShellDefinition_0_3) upgradeJST() error {
	d.Arguments = upgradeArguments(d.Arguments)
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
	return d.Env, nil
}

var _ taskKind_0_3 = &SQLDefinition_0_3{}

type SQLDefinition_0_3 struct {
	Resource   string                 `json:"resource"`
	Entrypoint string                 `json:"entrypoint"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
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

func (d *SQLDefinition_0_3) upgradeJST() error {
	return nil
}

func (d *SQLDefinition_0_3) getKindOptions() (build.KindOptions, error) {
	queryBytes, err := os.ReadFile(d.Entrypoint)
	if err != nil {
		return nil, errors.Wrapf(err, "reading SQL entrypoint %s", d.Entrypoint)
	}
	return build.KindOptions{
		"query":     string(queryBytes),
		"queryArgs": d.Parameters,
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
	URLParams map[string]string      `json:"urlParams,omitempty"`
	Headers   map[string]string      `json:"headers,omitempty"`
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

func (d *RESTDefinition_0_3) upgradeJST() error {
	return nil
}

func (d *RESTDefinition_0_3) getKindOptions() (build.KindOptions, error) {
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
}

type OptionDefinition_0_3 struct {
	Label string `json:"label"`
	Value string `json:"value"`
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

type PermissionDefinition_0_3 struct {
	Viewers    []string `json:"viewers,omitempty"`
	Requesters []string `json:"requesters,omitempty"`
	Executers  []string `json:"executers,omitempty"`
	Admins     []string `json:"admins,omitempty"`
}

func (d PermissionDefinition_0_3) isEmpty() bool {
	return len(d.Viewers) == 0 && len(d.Requesters) == 0 && len(d.Executers) == 0 && len(d.Admins) == 0
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
			Entrypoint: entrypoint,
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

	if err := d.addPermissionsToUpdateTaskRequest(ctx, client, &req); err != nil {
		return api.UpdateTaskRequest{}, err
	}

	if d.Constraints != nil {
		req.Constraints = *d.Constraints
	}

	if err := d.addKindSpecificsToUpdateTaskRequest(ctx, client, &req); err != nil {
		return api.UpdateTaskRequest{}, err
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

func (d Definition_0_3) addPermissionsToUpdateTaskRequest(ctx context.Context, client api.IAPIClient, req *api.UpdateTaskRequest) error {
	if d.Permissions != nil && !d.Permissions.isEmpty() {
		req.RequireExplicitPermissions = true
		// TODO: convert permissions.
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
	taskKind, err := d.taskKind()
	if err != nil {
		return err
	}
	return taskKind.upgradeJST()
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

func getResourcesByName(ctx context.Context, client api.IAPIClient) (map[string]api.Resource, error) {
	// Remap resources from ref -> name to ref -> id.
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
