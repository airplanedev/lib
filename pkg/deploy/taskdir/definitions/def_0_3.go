package definitions

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"text/template"

	"github.com/airplanedev/lib/pkg/api"
	"github.com/airplanedev/lib/pkg/build"
	"github.com/airplanedev/lib/pkg/utils/pointers"
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
	Resources   map[string]string         `json:"resources,omitempty"`

	Image  *ImageDefinition_0_3  `json:"docker,omitempty"`
	Node   *NodeDefinition_0_3   `json:"node,omitempty"`
	Python *PythonDefinition_0_3 `json:"python,omitempty"`
	Shell  *ShellDefinition_0_3  `json:"shell,omitempty"`

	SQL  *SQLDefinition_0_3  `json:"sql,omitempty"`
	REST *RESTDefinition_0_3 `json:"rest,omitempty"`

	Constraints        map[string]string        `json:"constraints,omitempty"`
	RequireRequests    bool                     `json:"requireRequests,omitempty"`
	AllowSelfApprovals DefaultTrueDefinition    `json:"allowSelfApprovals,omitempty"`
	Timeout            DefaultTimeoutDefinition `json:"timeout,omitempty"`
	Runtime            build.TaskRuntime        `json:"runtime,omitempty"`

	Schedules map[string]ScheduleDefinition_0_3 `json:"schedules,omitempty"`

	buildConfig  build.BuildConfig
	defnFilePath string
}

var _ DefinitionInterface = &Definition_0_3{}

type taskKind_0_3 interface {
	fillInUpdateTaskRequest(context.Context, api.IAPIClient, *api.UpdateTaskRequest) error
	hydrateFromTask(context.Context, api.IAPIClient, *api.Task) error
	setEntrypoint(string) error
	setAbsoluteEntrypoint(string) error
	getAbsoluteEntrypoint() (string, error)
	getKindOptions() (build.KindOptions, error)
	getEntrypoint() (string, error)
	getEnv() (api.TaskEnv, error)
	getConfigAttachments() []api.ConfigAttachment
}

var _ taskKind_0_3 = &ImageDefinition_0_3{}

type ImageDefinition_0_3 struct {
	Image      string      `json:"image"`
	Entrypoint string      `json:"entrypoint,omitempty"`
	Command    string      `json:"command"`
	EnvVars    api.TaskEnv `json:"envVars,omitempty"`
}

func (d *ImageDefinition_0_3) fillInUpdateTaskRequest(ctx context.Context, client api.IAPIClient, req *api.UpdateTaskRequest) error {
	if d.Image != "" {
		req.Image = &d.Image
	}
	if args, err := shlex.Split(d.Command); err != nil {
		return err
	} else {
		req.Arguments = args
	}
	if cmd, err := shlex.Split(d.Entrypoint); err != nil {
		return err
	} else {
		req.Command = cmd
	}
	req.Env = d.EnvVars
	return nil
}

func (d *ImageDefinition_0_3) hydrateFromTask(ctx context.Context, client api.IAPIClient, t *api.Task) error {
	if t.Image != nil {
		d.Image = *t.Image
	}
	d.Command = shellescape.QuoteCommand(t.Arguments)
	d.Entrypoint = shellescape.QuoteCommand(t.Command)
	d.EnvVars = t.Env
	return nil
}

func (d *ImageDefinition_0_3) setEntrypoint(entrypoint string) error {
	return ErrNoEntrypoint
}

func (d *ImageDefinition_0_3) setAbsoluteEntrypoint(entrypoint string) error {
	return ErrNoEntrypoint
}

func (d *ImageDefinition_0_3) getAbsoluteEntrypoint() (string, error) {
	return "", ErrNoEntrypoint
}

func (d *ImageDefinition_0_3) getKindOptions() (build.KindOptions, error) {
	return nil, nil
}

func (d *ImageDefinition_0_3) getEntrypoint() (string, error) {
	return "", ErrNoEntrypoint
}

func (d *ImageDefinition_0_3) getEnv() (api.TaskEnv, error) {
	return d.EnvVars, nil
}

func (d *ImageDefinition_0_3) getConfigAttachments() []api.ConfigAttachment {
	return []api.ConfigAttachment{}
}

var _ taskKind_0_3 = &NodeDefinition_0_3{}

type NodeDefinition_0_3 struct {
	Entrypoint  string      `json:"entrypoint"`
	NodeVersion string      `json:"nodeVersion"`
	EnvVars     api.TaskEnv `json:"envVars,omitempty"`

	absoluteEntrypoint string `json:"-"`
}

func (d *NodeDefinition_0_3) fillInUpdateTaskRequest(ctx context.Context, client api.IAPIClient, req *api.UpdateTaskRequest) error {
	req.Env = d.EnvVars
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
	d.EnvVars = t.Env
	return nil
}

func (d *NodeDefinition_0_3) setEntrypoint(entrypoint string) error {
	d.Entrypoint = entrypoint
	return nil
}

func (d *NodeDefinition_0_3) setAbsoluteEntrypoint(entrypoint string) error {
	d.absoluteEntrypoint = entrypoint
	return nil
}

func (d *NodeDefinition_0_3) getAbsoluteEntrypoint() (string, error) {
	if d.absoluteEntrypoint == "" {
		return "", ErrNoAbsoluteEntrypoint
	}
	return d.absoluteEntrypoint, nil
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

func (d *NodeDefinition_0_3) getEnv() (api.TaskEnv, error) {
	return d.EnvVars, nil
}

func (d *NodeDefinition_0_3) getConfigAttachments() []api.ConfigAttachment {
	return []api.ConfigAttachment{}
}

var _ taskKind_0_3 = &PythonDefinition_0_3{}

type PythonDefinition_0_3 struct {
	Entrypoint string      `json:"entrypoint"`
	EnvVars    api.TaskEnv `json:"envVars,omitempty"`

	absoluteEntrypoint string `json:"-"`
}

func (d *PythonDefinition_0_3) fillInUpdateTaskRequest(ctx context.Context, client api.IAPIClient, req *api.UpdateTaskRequest) error {
	req.Env = d.EnvVars
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
	d.EnvVars = t.Env
	return nil
}

func (d *PythonDefinition_0_3) setEntrypoint(entrypoint string) error {
	d.Entrypoint = entrypoint
	return nil
}

func (d *PythonDefinition_0_3) setAbsoluteEntrypoint(entrypoint string) error {
	d.absoluteEntrypoint = entrypoint
	return nil
}

func (d *PythonDefinition_0_3) getAbsoluteEntrypoint() (string, error) {
	if d.absoluteEntrypoint == "" {
		return "", ErrNoAbsoluteEntrypoint
	}
	return d.absoluteEntrypoint, nil
}

func (d *PythonDefinition_0_3) getKindOptions() (build.KindOptions, error) {
	return build.KindOptions{
		"entrypoint": d.Entrypoint,
	}, nil
}

func (d *PythonDefinition_0_3) getEntrypoint() (string, error) {
	return d.Entrypoint, nil
}

func (d *PythonDefinition_0_3) getEnv() (api.TaskEnv, error) {
	return d.EnvVars, nil
}

func (d *PythonDefinition_0_3) getConfigAttachments() []api.ConfigAttachment {
	return []api.ConfigAttachment{}
}

var _ taskKind_0_3 = &ShellDefinition_0_3{}

type ShellDefinition_0_3 struct {
	Entrypoint string      `json:"entrypoint"`
	EnvVars    api.TaskEnv `json:"envVars,omitempty"`

	absoluteEntrypoint string `json:"-"`
}

func (d *ShellDefinition_0_3) fillInUpdateTaskRequest(ctx context.Context, client api.IAPIClient, req *api.UpdateTaskRequest) error {
	req.Env = d.EnvVars
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
	d.EnvVars = t.Env
	return nil
}

func (d *ShellDefinition_0_3) setEntrypoint(entrypoint string) error {
	d.Entrypoint = entrypoint
	return nil
}

func (d *ShellDefinition_0_3) setAbsoluteEntrypoint(entrypoint string) error {
	d.absoluteEntrypoint = entrypoint
	return nil
}

func (d *ShellDefinition_0_3) getAbsoluteEntrypoint() (string, error) {
	if d.absoluteEntrypoint == "" {
		return "", ErrNoAbsoluteEntrypoint
	}
	return d.absoluteEntrypoint, nil
}

func (d *ShellDefinition_0_3) getKindOptions() (build.KindOptions, error) {
	return build.KindOptions{
		"entrypoint": d.Entrypoint,
	}, nil
}

func (d *ShellDefinition_0_3) getEntrypoint() (string, error) {
	return d.Entrypoint, nil
}

func (d *ShellDefinition_0_3) getEnv() (api.TaskEnv, error) {
	return d.EnvVars, nil
}

func (d *ShellDefinition_0_3) getConfigAttachments() []api.ConfigAttachment {
	return []api.ConfigAttachment{}
}

var _ taskKind_0_3 = &SQLDefinition_0_3{}

type SQLDefinition_0_3 struct {
	Resource        string                 `json:"resource"`
	Entrypoint      string                 `json:"entrypoint"`
	QueryArgs       map[string]interface{} `json:"queryArgs,omitempty"`
	TransactionMode SQLTransactionMode     `json:"transactionMode,omitempty"`
	Configs         []string               `json:"configs,omitempty"`

	// Contents of Entrypoint, cached
	entrypointContents string `json:"-"`
	absoluteEntrypoint string `json:"-"`
}

type SQLTransactionMode string

var _ yaml.IsZeroer = SQLTransactionMode("")

func (tm SQLTransactionMode) IsZero() bool {
	return tm == "auto" || tm == ""
}

func (tm SQLTransactionMode) Value() string {
	if tm == "" {
		return "auto"
	}
	return string(tm)
}

func (d *SQLDefinition_0_3) GetQuery() (string, error) {
	if d.entrypointContents == "" {
		if d.absoluteEntrypoint == "" {
			return "", ErrNoAbsoluteEntrypoint
		}
		queryBytes, err := os.ReadFile(d.absoluteEntrypoint)
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
		req.Resources["db"] = res.ID
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
			d.QueryArgs = mv
		} else {
			return errors.Errorf("expected map queryArgs, got %T instead", v)
		}
	}
	if v, ok := t.KindOptions["transactionMode"]; ok {
		if sv, ok := v.(string); ok {
			d.TransactionMode = SQLTransactionMode(sv)
		} else {
			return errors.Errorf("expected string transactionMode, got %T instead", v)
		}
	}

	d.Configs = make([]string, len(t.Configs))
	for idx, config := range t.Configs {
		d.Configs[idx] = config.NameTag
	}

	return nil
}

func (d *SQLDefinition_0_3) setEntrypoint(entrypoint string) error {
	d.Entrypoint = entrypoint
	return nil
}

func (d *SQLDefinition_0_3) setAbsoluteEntrypoint(entrypoint string) error {
	d.absoluteEntrypoint = entrypoint
	return nil
}

func (d *SQLDefinition_0_3) getAbsoluteEntrypoint() (string, error) {
	if d.absoluteEntrypoint == "" {
		return "", ErrNoAbsoluteEntrypoint
	}
	return d.absoluteEntrypoint, nil
}

func (d *SQLDefinition_0_3) getKindOptions() (build.KindOptions, error) {
	query, err := d.GetQuery()
	if err != nil {
		return nil, err
	}
	if d.QueryArgs == nil {
		d.QueryArgs = map[string]interface{}{}
	}
	return build.KindOptions{
		"entrypoint":      d.Entrypoint,
		"query":           query,
		"queryArgs":       d.QueryArgs,
		"transactionMode": d.TransactionMode.Value(),
	}, nil
}

func (d *SQLDefinition_0_3) getEntrypoint() (string, error) {
	return d.Entrypoint, nil
}

func (d *SQLDefinition_0_3) getEnv() (api.TaskEnv, error) {
	return nil, nil
}

func (d *SQLDefinition_0_3) getConfigAttachments() []api.ConfigAttachment {
	configAttachments := make([]api.ConfigAttachment, len(d.Configs))
	for i, configName := range d.Configs {
		configAttachments[i] = api.ConfigAttachment{NameTag: configName}
	}

	return configAttachments
}

var _ taskKind_0_3 = &RESTDefinition_0_3{}

type RESTDefinition_0_3 struct {
	Resource  string                 `json:"resource"`
	Method    string                 `json:"method"`
	Path      string                 `json:"path"`
	URLParams map[string]interface{} `json:"urlParams,omitempty"`
	Headers   map[string]interface{} `json:"headers,omitempty"`
	BodyType  string                 `json:"bodyType"`
	Body      interface{}            `json:"body,omitempty"`
	FormData  map[string]interface{} `json:"formData,omitempty"`
	Configs   []string               `json:"configs,omitempty"`
}

func (d *RESTDefinition_0_3) fillInUpdateTaskRequest(ctx context.Context, client api.IAPIClient, req *api.UpdateTaskRequest) error {
	resourcesByName, err := getResourcesByName(ctx, client)
	if err != nil {
		return err
	}
	if res, ok := resourcesByName[d.Resource]; ok {
		req.Resources["rest"] = res.ID
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
		d.Body = v
	}
	if v, ok := t.KindOptions["formData"]; ok {
		if mv, ok := v.(map[string]interface{}); ok {
			d.FormData = mv
		} else {
			return errors.Errorf("expected map formData, got %T instead", v)
		}
	}

	d.Configs = make([]string, len(t.Configs))
	for idx, config := range t.Configs {
		d.Configs[idx] = config.NameTag
	}

	return nil
}

func (d *RESTDefinition_0_3) setEntrypoint(entrypoint string) error {
	return ErrNoEntrypoint
}

func (d *RESTDefinition_0_3) setAbsoluteEntrypoint(entrypoint string) error {
	return ErrNoEntrypoint
}

func (d *RESTDefinition_0_3) getAbsoluteEntrypoint() (string, error) {
	return "", ErrNoEntrypoint
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

func (d *RESTDefinition_0_3) getEnv() (api.TaskEnv, error) {
	return nil, nil
}

func (d *RESTDefinition_0_3) getConfigAttachments() []api.ConfigAttachment {
	configAttachments := make([]api.ConfigAttachment, len(d.Configs))
	for i, configName := range d.Configs {
		configAttachments[i] = api.ConfigAttachment{NameTag: configName}
	}

	return configAttachments
}

type ParameterDefinition_0_3 struct {
	Name        string                 `json:"name"`
	Slug        string                 `json:"slug"`
	Type        string                 `json:"type"`
	Description string                 `json:"description,omitempty"`
	Default     interface{}            `json:"default,omitempty"`
	Required    DefaultTrueDefinition  `json:"required,omitempty"`
	Options     []OptionDefinition_0_3 `json:"options,omitempty"`
	Regex       string                 `json:"regex,omitempty"`
}

type OptionDefinition_0_3 struct {
	Label  string      `json:"label"`
	Value  interface{} `json:"value,omitempty"`
	Config *string     `json:"config,omitempty"`
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

type ScheduleDefinition_0_3 struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	CronExpr    string                 `json:"cron"`
	ParamValues map[string]interface{} `json:"paramValues,omitempty"`
}

//go:embed schema_0_3.json
var schemaStr string

func NewDefinition_0_3(name string, slug string, kind build.TaskKind, entrypoint string) (Definition_0_3, error) {
	def := Definition_0_3{
		Name: name,
		Slug: slug,
	}

	switch kind {
	case build.TaskKindImage:
		def.Image = &ImageDefinition_0_3{
			Image:   "alpine:3",
			Command: `echo "hello world"`,
		}
	case build.TaskKindNode:
		def.Node = &NodeDefinition_0_3{
			Entrypoint:  entrypoint,
			NodeVersion: "16",
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
		def.REST = &RESTDefinition_0_3{
			Method:   "POST",
			Path:     "/",
			BodyType: "json",
			Body:     "{}",
		}
	default:
		return Definition_0_3{}, errors.Errorf("unknown kind: %s", kind)
	}

	return def, nil
}

func (d Definition_0_3) Marshal(format DefFormat) ([]byte, error) {
	switch format {
	case DefFormatYAML:
		// Use the JSON marshaler so we use MarshalJSON methods.
		buf, err := yaml.MarshalWithOptions(d,
			yaml.UseJSONMarshaler(),
			yaml.UseLiteralStyleIfMultiline(true))
		if err != nil {
			return nil, err
		}
		return buf, nil

	case DefFormatJSON:
		// Use the YAML marshaler so we can take advantage of the yaml.IsZeroer check on omitempty.
		// But make it use the JSON marshaler so we use MarshalJSON methods.
		buf, err := yaml.MarshalWithOptions(d,
			yaml.UseJSONMarshaler(),
			yaml.JSON())
		if err != nil {
			return nil, err
		}
		// `yaml.Marshal` doesn't allow configuring JSON indentation, so do it after the fact.
		var out bytes.Buffer
		if err := json.Indent(&out, buf, "", "\t"); err != nil {
			return nil, err
		}
		return out.Bytes(), nil

	default:
		return nil, errors.Errorf("unknown format: %s", format)
	}
}

// GenerateCommentedFile generates a commented YAML file under certain circumstances. If the format
// requested isn't YAML, or if the definition has other things filled in, this method defaults to
// calling Marshal(format).
func (d Definition_0_3) GenerateCommentedFile(format DefFormat) ([]byte, error) {
	// If it's not YAML, or you have other things defined on your task def, bail.
	if format != DefFormatYAML ||
		d.Description != "" ||
		len(d.Parameters) > 0 ||
		len(d.Resources) > 0 ||
		len(d.Constraints) > 0 ||
		d.RequireRequests ||
		!d.AllowSelfApprovals.IsZero() ||
		!d.Timeout.IsZero() {
		return d.Marshal(format)
	}

	kind, err := d.Kind()
	if err != nil {
		return nil, err
	}

	taskDefinition := new(bytes.Buffer)
	var paramsExtraInfo string
	switch kind {
	case build.TaskKindImage:
		if d.Image.Entrypoint != "" || len(d.Image.EnvVars) > 0 {
			return d.Marshal(format)
		}
		tmpl, err := template.New("image").Parse(imageTemplate)
		if err != nil {
			return nil, errors.Wrap(err, "parsing image template")
		}
		if err := tmpl.Execute(taskDefinition, d.Image); err != nil {
			return nil, errors.Wrap(err, "executing image template")
		}
		paramsExtraInfo = imageParamsExtraDescription
	case build.TaskKindNode:
		if len(d.Node.EnvVars) > 0 {
			return d.Marshal(format)
		}
		tmpl, err := template.New("node").Parse(nodeTemplate)
		if err != nil {
			return nil, errors.Wrap(err, "parsing node template")
		}
		if err := tmpl.Execute(taskDefinition, d.Node); err != nil {
			return nil, errors.Wrap(err, "executing node template")
		}
	case build.TaskKindPython:
		if len(d.Python.EnvVars) > 0 {
			return d.Marshal(format)
		}
		tmpl, err := template.New("python").Parse(pythonTemplate)
		if err != nil {
			return nil, errors.Wrap(err, "parsing python template")
		}
		if err := tmpl.Execute(taskDefinition, d.Python); err != nil {
			return nil, errors.Wrap(err, "executing python template")
		}
	case build.TaskKindShell:
		if len(d.Shell.EnvVars) > 0 {
			return d.Marshal(format)
		}
		tmpl, err := template.New("shell").Parse(shellTemplate)
		if err != nil {
			return nil, errors.Wrap(err, "parsing shell template")
		}
		if err := tmpl.Execute(taskDefinition, d.Shell); err != nil {
			return nil, errors.Wrap(err, "executing shell template")
		}
		paramsExtraInfo = shellParamsExtraDescription
	case build.TaskKindSQL:
		if d.SQL.Resource != "" || len(d.SQL.QueryArgs) > 0 {
			return d.Marshal(format)
		}
		tmpl, err := template.New("sql").Parse(sqlTemplate)
		if err != nil {
			return nil, errors.Wrap(err, "parsing SQL template")
		}
		if err := tmpl.Execute(taskDefinition, d.SQL); err != nil {
			return nil, errors.Wrap(err, "executing sql template")
		}
	case build.TaskKindREST:
		if d.REST.Resource != "" ||
			len(d.REST.URLParams) > 0 ||
			len(d.REST.Headers) > 0 ||
			len(d.REST.FormData) > 0 {
			return d.Marshal(format)
		}
		tmpl, err := template.New("rest").Parse(restTemplate)
		if err != nil {
			return nil, errors.Wrap(err, "parsing REST template")
		}
		if err := tmpl.Execute(taskDefinition, d.REST); err != nil {
			return nil, errors.Wrap(err, "executing rest template")
		}
	default:
		return d.Marshal(format)
	}

	tmpl, err := template.New("definition").Parse(definitionTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "parsing definition template")
	}
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, map[string]interface{}{
		"slug":                   d.Slug,
		"name":                   d.Name,
		"taskDefinition":         taskDefinition.String(),
		"paramsExtraDescription": paramsExtraInfo,
	}); err != nil {
		return nil, errors.Wrap(err, "executing definition template")
	}
	return buf.Bytes(), nil
}

func (d *Definition_0_3) Unmarshal(format DefFormat, buf []byte) error {
	var err error
	switch format {
	case DefFormatYAML:
		buf, err = yaml.YAMLToJSON(buf)
		if err != nil {
			return err
		}
	case DefFormatJSON:
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

func (d *Definition_0_3) SetAbsoluteEntrypoint(entrypoint string) error {
	taskKind, err := d.taskKind()
	if err != nil {
		return err
	}

	return taskKind.setAbsoluteEntrypoint(entrypoint)
}

func (d *Definition_0_3) GetAbsoluteEntrypoint() (string, error) {
	taskKind, err := d.taskKind()
	if err != nil {
		return "", err
	}

	return taskKind.getAbsoluteEntrypoint()
}

func (d Definition_0_3) Kind() (build.TaskKind, error) {
	if d.Image != nil {
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
	if d.Image != nil {
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

func (d Definition_0_3) GetUpdateTaskRequest(ctx context.Context, client api.IAPIClient) (api.UpdateTaskRequest, error) {
	req := api.UpdateTaskRequest{
		Slug:        d.Slug,
		Name:        d.Name,
		Description: d.Description,
		Timeout:     d.Timeout.Value(),
		Runtime:     d.Runtime,
		ExecuteRules: api.UpdateExecuteRulesRequest{
			RequireRequests: &d.RequireRequests,
		},
		Resources: make(map[string]string),
	}

	if err := d.addParametersToUpdateTaskRequest(ctx, &req); err != nil {
		return api.UpdateTaskRequest{}, err
	}

	if err := d.addResourcesToUpdateTaskRequest(ctx, client, &req); err != nil {
		return api.UpdateTaskRequest{}, err
	}

	if len(d.Constraints) > 0 {
		labels := []api.AgentLabel{}
		for key, val := range d.Constraints {
			labels = append(labels, api.AgentLabel{
				Key:   key,
				Value: val,
			})
		}
		req.Constraints = api.RunConstraints{
			Labels: labels,
		}
	}

	req.ExecuteRules.DisallowSelfApprove = pointers.Bool(!d.AllowSelfApprovals.Value())

	if err := d.addKindSpecificsToUpdateTaskRequest(ctx, client, &req); err != nil {
		return api.UpdateTaskRequest{}, err
	}

	return req, nil
}

func (d Definition_0_3) addParametersToUpdateTaskRequest(ctx context.Context, req *api.UpdateTaskRequest) error {
	req.Parameters = make([]api.Parameter, len(d.Parameters))
	for i, pd := range d.Parameters {
		param := api.Parameter{
			Name: pd.Name,
			Slug: pd.Slug,
			Desc: pd.Description,
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

		if pd.Default != nil {
			if pd.Type == "configvar" {
				switch reflect.ValueOf(pd.Default).Kind() {
				case reflect.Map:
					m, ok := pd.Default.(map[string]interface{})
					if !ok {
						return errors.Errorf("expected map but got %T", pd.Default)
					}
					if configName, ok := m["config"]; !ok {
						return errors.Errorf("missing config property from configvar type: %v", pd.Default)
					} else {
						param.Default = map[string]interface{}{
							"__airplaneType": "configvar",
							"name":           configName,
						}
					}
				case reflect.String:
					param.Default = map[string]interface{}{
						"__airplaneType": "configvar",
						"name":           pd.Default,
					}
				default:
					return errors.Errorf("unsupported type for default value: %T", pd.Default)
				}
			} else {
				switch reflect.ValueOf(pd.Default).Kind() {
				case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
					reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
					reflect.Float32, reflect.Float64:
					param.Default = pd.Default
				default:
					return errors.Errorf("unsupported type for default value: %T", pd.Default)
				}
			}
		}

		if !pd.Required.Value() {
			param.Constraints.Optional = true
		}

		param.Constraints.Regex = pd.Regex

		if len(pd.Options) > 0 {
			param.Constraints.Options = make([]api.ConstraintOption, len(pd.Options))
			for j, od := range pd.Options {
				param.Constraints.Options[j].Label = od.Label
				if od.Config != nil {
					param.Constraints.Options[j].Value = map[string]interface{}{
						"__airplaneType": "configvar",
						"name":           *od.Config,
					}
				} else if pd.Type == "configvar" {
					param.Constraints.Options[j].Value = map[string]interface{}{
						"__airplaneType": "configvar",
						"name":           od.Value,
					}
				} else {
					param.Constraints.Options[j].Value = od.Value
				}
			}
		}

		req.Parameters[i] = param
	}
	return nil
}

func (d Definition_0_3) addResourcesToUpdateTaskRequest(ctx context.Context, client api.IAPIClient, req *api.UpdateTaskRequest) error {
	if len(d.Resources) == 0 {
		return nil
	}

	resourcesBySlug, err := getResourcesBySlug(ctx, client)
	if err != nil {
		return errors.Wrap(err, "fetching resources")
	}

	for alias, slug := range d.Resources {
		r, ok := resourcesBySlug[slug]
		if ok {
			req.Resources[alias] = r.ID
		} else {
			return errors.Errorf("unknown resource: %s", slug)
		}
	}

	return nil
}

func (d Definition_0_3) addKindSpecificsToUpdateTaskRequest(ctx context.Context, client api.IAPIClient, req *api.UpdateTaskRequest) error {
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

	configAttachments, err := d.GetConfigAttachments()
	if err != nil {
		return err
	}
	req.Configs = &configAttachments

	taskKind, err := d.taskKind()
	if err != nil {
		return err
	}
	if err := taskKind.fillInUpdateTaskRequest(ctx, client, req); err != nil {
		return err
	}
	return nil
}

func (d Definition_0_3) Entrypoint() (string, error) {
	taskKind, err := d.taskKind()
	if err != nil {
		return "", err
	}
	return taskKind.getEntrypoint()
}

func (d Definition_0_3) GetDefnFilePath() string {
	return d.defnFilePath
}

func (d *Definition_0_3) SetDefnFilePath(filePath string) {
	d.defnFilePath = filePath
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

func (d *Definition_0_3) GetConfigAttachments() ([]api.ConfigAttachment, error) {
	taskKind, err := d.taskKind()
	if err != nil {
		return nil, err
	}
	return taskKind.getConfigAttachments(), nil
}

func (d *Definition_0_3) GetSlug() string {
	return d.Slug
}

func (d *Definition_0_3) GetName() string {
	return d.Name
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

func (d *Definition_0_3) GetSchedules() map[string]api.Schedule {
	if len(d.Schedules) == 0 {
		return nil
	}

	schedules := make(map[string]api.Schedule)
	for slug, def := range d.Schedules {
		schedules[slug] = api.Schedule{
			Name:        def.Name,
			Description: def.Description,
			CronExpr:    def.CronExpr,
			ParamValues: def.ParamValues,
		}
	}
	return schedules
}

func NewDefinitionFromTask_0_3(ctx context.Context, client api.IAPIClient, t api.Task) (Definition_0_3, error) {
	d := Definition_0_3{
		Name:            t.Name,
		Slug:            t.Slug,
		Description:     t.Description,
		RequireRequests: t.ExecuteRules.RequireRequests,
		Runtime:         t.Runtime,
	}

	if err := d.convertParametersFromTask(ctx, client, &t); err != nil {
		return Definition_0_3{}, err
	}

	if err := d.convertResourcesFromTask(ctx, client, &t); err != nil {
		return Definition_0_3{}, err
	}

	if err := d.convertTaskKindFromTask(ctx, client, &t); err != nil {
		return Definition_0_3{}, err
	}

	if !t.Constraints.IsEmpty() {
		d.Constraints = map[string]string{}
		for _, label := range t.Constraints.Labels {
			d.Constraints[label.Key] = label.Value
		}
	}

	d.AllowSelfApprovals.value = pointers.Bool(!t.ExecuteRules.DisallowSelfApprove)
	d.Timeout.value = t.Timeout

	schedules := make(map[string]ScheduleDefinition_0_3)
	for _, trigger := range t.Triggers {
		if trigger.Kind != api.TriggerKindSchedule || trigger.Slug == nil {
			// Trigger is not a schedule deployed via code
			continue
		}
		if trigger.ArchivedAt != nil || trigger.DisabledAt != nil {
			// Trigger is archived or disabled, so don't add to task defn file
			continue
		}

		schedules[*trigger.Slug] = ScheduleDefinition_0_3{
			Name:        trigger.Name,
			Description: trigger.Description,
			CronExpr:    trigger.KindConfig.Schedule.CronExpr.String(),
			ParamValues: trigger.KindConfig.Schedule.ParamValues,
		}
	}
	if len(schedules) > 0 {
		d.Schedules = schedules
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

		if param.Default != nil {
			if param.Type == "configvar" {
				switch k := reflect.ValueOf(param.Default).Kind(); k {
				case reflect.Map:
					configName, err := extractConfigVarValue(param.Default)
					if err != nil {
						return errors.Wrap(err, "unhandled default configvar")
					}
					p.Default = configName
				default:
					return errors.Errorf("unsupported type for default value: %T", param.Default)
				}
			} else {
				switch k := reflect.ValueOf(param.Default).Kind(); k {
				case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
					reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
					reflect.Float32, reflect.Float64:
					p.Default = param.Default
				default:
					return errors.Errorf("unsupported type for default value: %T", param.Default)
				}
			}
		}

		p.Required.value = pointers.Bool(!param.Constraints.Optional)

		p.Regex = param.Constraints.Regex

		if len(param.Constraints.Options) > 0 {
			p.Options = make([]OptionDefinition_0_3, len(param.Constraints.Options))
			for j, opt := range param.Constraints.Options {
				if param.Type == "configvar" {
					switch k := reflect.ValueOf(opt.Value).Kind(); k {
					case reflect.Map:
						configName, err := extractConfigVarValue(opt.Value)
						if err != nil {
							return errors.Wrap(err, "unhandled option")
						}
						p.Options[j] = OptionDefinition_0_3{
							Label: opt.Label,
							Value: configName,
						}
					default:
						return errors.Errorf("unhandled option type: %s", k)
					}
				} else {
					switch k := reflect.ValueOf(opt.Value).Kind(); k {
					case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
						reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
						reflect.Float32, reflect.Float64:
						p.Options[j] = OptionDefinition_0_3{
							Label: opt.Label,
							Value: opt.Value,
						}
					default:
						return errors.Errorf("unhandled option type: %s", k)
					}
				}
			}
		}

		d.Parameters[idx] = p
	}
	return nil
}

func (d *Definition_0_3) convertResourcesFromTask(ctx context.Context, client api.IAPIClient, t *api.Task) error {
	if len(t.Resources) == 0 {
		return nil
	}

	resourcesByID, err := getResourcesByID(ctx, client)
	if err != nil {
		return errors.Wrap(err, "fetching resources")
	}

	d.Resources = make(map[string]string)
	for alias, id := range t.Resources {
		// Ignore SQL/REST resources; they get routed elsewhere.
		if (t.Kind == build.TaskKindSQL && alias == "db") ||
			(t.Kind == build.TaskKindREST && alias == "rest") {
			continue
		}
		r, ok := resourcesByID[id]
		if ok {
			d.Resources[alias] = r.Slug
		}
	}

	return nil
}

func extractConfigVarValue(v interface{}) (string, error) {
	m, ok := v.(map[string]interface{})
	if !ok {
		return "", errors.Errorf("expected map but got %T", v)
	}
	if airplaneType, ok := m["__airplaneType"]; !ok || airplaneType != "configvar" {
		return "", errors.Errorf("expected airplaneType=configvar but got %v", airplaneType)
	}
	if configName, ok := m["name"]; !ok {
		return "", errors.Errorf("missing name property from configvar type: %v", v)
	} else if configNameStr, ok := configName.(string); !ok {
		return "", errors.Errorf("expected name to be string but got %T", configName)
	} else {
		return configNameStr, nil
	}
}

func (d *Definition_0_3) convertTaskKindFromTask(ctx context.Context, client api.IAPIClient, t *api.Task) error {
	switch t.Kind {
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

	// Pass runtime through to builder
	config["runtime"] = d.Runtime

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

func getResourcesBySlug(ctx context.Context, client api.IAPIClient) (map[string]api.Resource, error) {
	resp, err := client.ListResources(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "fetching resources")
	}
	resourcesBySlug := map[string]api.Resource{}
	for _, resource := range resp.Resources {
		resourcesBySlug[resource.Slug] = resource
	}
	return resourcesBySlug, nil
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
