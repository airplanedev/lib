package runtime

import "github.com/airplanedev/lib/pkg/build"

// Task represents a task.
type Task struct {
	URL                        string            `json:"-" yaml:"-"`
	ID                         string            `json:"taskID" yaml:"id"`
	Name                       string            `json:"name" yaml:"name"`
	Slug                       string            `json:"slug" yaml:"slug"`
	Description                string            `json:"description" yaml:"description"`
	Image                      *string           `json:"image" yaml:"image"`
	Command                    []string          `json:"command" yaml:"command"`
	Arguments                  []string          `json:"arguments" yaml:"arguments"`
	Parameters                 Parameters        `json:"parameters" yaml:"parameters"`
	Constraints                RunConstraints    `json:"constraints" yaml:"constraints"`
	Env                        TaskEnv           `json:"env" yaml:"env"`
	ResourceRequests           ResourceRequests  `json:"resourceRequests" yaml:"resourceRequests"`
	Resources                  Resources         `json:"resources" yaml:"resources"`
	Kind                       build.TaskKind    `json:"kind" yaml:"kind"`
	KindOptions                build.KindOptions `json:"kindOptions" yaml:"kindOptions"`
	Repo                       string            `json:"repo" yaml:"repo"`
	RequireExplicitPermissions bool              `json:"requireExplicitPermissions" yaml:"-"`
	Permissions                Permissions       `json:"permissions" yaml:"-"`
	Timeout                    int               `json:"timeout" yaml:"timeout"`
	InterpolationMode          string            `json:"interpolationMode" yaml:"-"`
}

type ResourceRequests map[string]string

type TaskEnv map[string]EnvVarValue

type Resources map[string]string

type EnvVarValue struct {
	Value  *string `json:"value" yaml:"value,omitempty"`
	Config *string `json:"config" yaml:"config,omitempty"`
}

// RunConstraints represents run constraints.
type RunConstraints struct {
	Labels []AgentLabel `json:"labels" yaml:"labels"`
}

// AgentLabel represents an agent label.
type AgentLabel struct {
	Key   string `json:"key" yaml:"key"`
	Value string `json:"value" yaml:"value"`
}

type Permissions []Permission

type Permission struct {
	Action     Action  `json:"action"`
	SubUserID  *string `json:"subUserID"`
	SubGroupID *string `json:"subGroupID"`
}

// Type enumerates parameter types.
type Type string

type Action string

// All Parameter types.
const (
	TypeString    Type = "string"
	TypeBoolean   Type = "boolean"
	TypeUpload    Type = "upload"
	TypeInteger   Type = "integer"
	TypeFloat     Type = "float"
	TypeDate      Type = "date"
	TypeDatetime  Type = "datetime"
	TypeConfigVar Type = "configvar"
)

type Parameters []Parameter

// Parameter represents a task parameter.
type Parameter struct {
	Name        string      `json:"name" yaml:"name"`
	Slug        string      `json:"slug" yaml:"slug"`
	Type        Type        `json:"type" yaml:"type"`
	Desc        string      `json:"desc" yaml:"desc,omitempty"`
	Component   Component   `json:"component" yaml:"component,omitempty"`
	Default     Value       `json:"default" yaml:"default,omitempty"`
	Constraints Constraints `json:"constraints" yaml:"constraints,omitempty"`
}

// Component enumerates components.
type Component string

// All Component types.
const (
	ComponentNone      Component = ""
	ComponentEditorSQL Component = "editor-sql"
	ComponentTextarea  Component = "textarea"
)

// Value represents a value.
type Value interface{}

// Constraints represent constraints.
type Constraints struct {
	Optional bool               `json:"optional" yaml:"optional,omitempty"`
	Regex    string             `json:"regex" yaml:"regex,omitempty"`
	Options  []ConstraintOption `json:"options,omitempty" yaml:"options,omitempty"`
}

type ConstraintOption struct {
	Label string `json:"label"`
	Value Value  `json:"value"`
}

// Values represent parameters values.
//
// An alias is used because we want the type
// to be `map[string]interface{}` and not a custom one.
//
// They're keyed by the parameter "slug".
type Values = map[string]interface{}
