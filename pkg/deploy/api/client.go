package api

import (
	"context"
	"time"

	"github.com/airplanedev/lib/pkg/build"
)

type APIClient interface {
	GetTask(ctx context.Context, slug string) (res Task, err error)
	ListResources(ctx context.Context) (res ListResourcesResponse, err error)
	CreateBuildUpload(ctx context.Context, req CreateBuildUploadRequest) (res CreateBuildUploadResponse, err error)
}

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

type CreateBuildUploadRequest struct {
	SizeBytes int `json:"sizeBytes"`
}

type CreateBuildUploadResponse struct {
	Upload       Upload `json:"upload"`
	WriteOnlyURL string `json:"writeOnlyURL"`
}

type Upload struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// UpdateTaskRequest updates a task.
type UpdateTaskRequest struct {
	Slug                       string            `json:"slug"`
	Name                       string            `json:"name"`
	Description                string            `json:"description"`
	Image                      *string           `json:"image"`
	Command                    []string          `json:"command"`
	Arguments                  []string          `json:"arguments"`
	Parameters                 Parameters        `json:"parameters"`
	Constraints                RunConstraints    `json:"constraints"`
	Env                        TaskEnv           `json:"env"`
	ResourceRequests           map[string]string `json:"resourceRequests"`
	Resources                  map[string]string `json:"resources"`
	Kind                       build.TaskKind    `json:"kind"`
	KindOptions                build.KindOptions `json:"kindOptions"`
	Repo                       string            `json:"repo"`
	RequireExplicitPermissions bool              `json:"requireExplicitPermissions"`
	Permissions                Permissions       `json:"permissions"`
	// TODO(amir): friendly type here (120s, 5m ...)
	Timeout int     `json:"timeout"`
	BuildID *string `json:"buildID"`

	InterpolationMode string `json:"interpolationMode" yaml:"-"`
}

type ListResourcesResponse struct {
	Resources []Resource `json:"resources"`
}

type Resource struct {
	ID         string                 `json:"id" db:"id"`
	TeamID     string                 `json:"teamID" db:"team_id"`
	Name       string                 `json:"name" db:"name"`
	Kind       ResourceKind           `json:"kind" db:"kind"`
	KindConfig map[string]interface{} `json:"kindConfig" db:"kind_config"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	CreatedBy string    `json:"createdBy" db:"created_by"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
	UpdatedBy string    `json:"updatedBy" db:"updated_by"`

	IsPrivate bool `json:"isPrivate" db:"is_private"`
}

type ResourceKind string

type Permissions []Permission

type Permission struct {
	Action     Action  `json:"action"`
	SubUserID  *string `json:"subUserID"`
	SubGroupID *string `json:"subGroupID"`
}

type Action string

type ResourceRequests map[string]string

type Resources map[string]string

type TaskEnv map[string]EnvVarValue

type EnvVarValue struct {
	Value  *string `json:"value" yaml:"value,omitempty"`
	Config *string `json:"config" yaml:"config,omitempty"`
}

// Parameters represents a slice of task parameters.
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

// Value represents a value.
type Value interface{}

// Type enumerates parameter types.
type Type string

// Component enumerates components.
type Component string

// All Component types.
const (
	ComponentNone      Component = ""
	ComponentEditorSQL Component = "editor-sql"
	ComponentTextarea  Component = "textarea"
)

// RunConstraints represents run constraints.
type RunConstraints struct {
	Labels []AgentLabel `json:"labels" yaml:"labels"`
}

// AgentLabel represents an agent label.
type AgentLabel struct {
	Key   string `json:"key" yaml:"key"`
	Value string `json:"value" yaml:"value"`
}
