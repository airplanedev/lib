package definitions

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/airplanedev/lib/pkg/api"
	"github.com/airplanedev/lib/pkg/build"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type Definition_0_2 struct {
	Slug             string               `yaml:"slug"`
	Name             string               `yaml:"name"`
	Description      string               `yaml:"description,omitempty"`
	Arguments        []string             `yaml:"arguments,omitempty"`
	Parameters       api.Parameters       `yaml:"parameters,omitempty"`
	Constraints      api.RunConstraints   `yaml:"constraints,omitempty"`
	Env              api.TaskEnv          `yaml:"env,omitempty"`
	ResourceRequests api.ResourceRequests `yaml:"resourceRequests,omitempty"`
	Resources        api.Resources        `yaml:"resources,omitempty"`
	Repo             string               `yaml:"repo,omitempty"`
	Timeout          int                  `yaml:"timeout,omitempty"`

	Deno       *DenoDefinition       `yaml:"deno,omitempty"`
	Image      *ImageDefinition      `yaml:"image,omitempty"`
	Dockerfile *DockerfileDefinition `yaml:"dockerfile,omitempty"`
	Go         *GoDefinition         `yaml:"go,omitempty"`
	Node       *NodeDefinition       `yaml:"node,omitempty"`
	Python     *PythonDefinition     `yaml:"python,omitempty"`
	Shell      *ShellDefinition      `yaml:"shell,omitempty"`

	SQL  *SQLDefinition  `yaml:"sql,omitempty"`
	REST *RESTDefinition `yaml:"rest,omitempty"`

	// Root is a directory path relative to the parent directory of this
	// task definition which defines what directory should be included
	// in the task's Docker image.
	//
	// If not set, defaults to "." (in other words, the parent directory of this task definition).
	//
	// This field is ignored when using the "image" builder.
	Root string `yaml:"root,omitempty"`

	buildConfig build.BuildConfig
}

type ImageDefinition struct {
	Image   string   `yaml:"image,omitempty"`
	Command []string `yaml:"command,omitempty"`
}

type DenoDefinition struct {
	Entrypoint string `yaml:"entrypoint" mapstructure:"entrypoint"`
}

type DockerfileDefinition struct {
	Dockerfile string `yaml:"dockerfile" mapstructure:"dockerfile"`
}

type GoDefinition struct {
	Entrypoint string `yaml:"entrypoint" mapstructure:"entrypoint"`
}

type NodeDefinition struct {
	Entrypoint  string `yaml:"entrypoint" mapstructure:"entrypoint"`
	Language    string `yaml:"language" mapstructure:"language"`
	NodeVersion string `yaml:"nodeVersion" mapstructure:"nodeVersion"`
}

type PythonDefinition struct {
	Entrypoint string `yaml:"entrypoint" mapstructure:"entrypoint"`
}

type ShellDefinition struct {
	Entrypoint string `yaml:"entrypoint" mapstructure:"entrypoint"`
}

type SQLDefinition struct {
	Query string `yaml:"query" mapstructure:"query"`
}

type RESTDefinition struct {
	Headers            map[string]string `yaml:"headers,omitempty" mapstructure:"headers"`
	Method             string            `yaml:"method" mapstructure:"method"`
	Path               string            `yaml:"path" mapstructure:"path"`
	URLParams          map[string]string `yaml:"urlParams,omitempty" mapstructure:"urlParams"`
	Body               string            `yaml:"body,omitempty" mapstructure:"body,omitempty"`
	JSONBody           interface{}       `yaml:"jsonBody,omitempty" mapstructure:"jsonBody,omitempty"`
	FormURLEncodedBody map[string]string `yaml:"formUrlEncodedBody,omitempty" mapstructure:"formUrlEncodedBody,omitempty"`
	FormDataBody       map[string]string `yaml:"formDataBody,omitempty" mapstructure:"formDataBody,omitempty"`
}

func NewDefinitionFromTask_0_2(task api.Task) (Definition_0_2, error) {
	def := Definition_0_2{
		Slug:             task.Slug,
		Name:             task.Name,
		Description:      task.Description,
		Arguments:        task.Arguments,
		Parameters:       task.Parameters,
		Constraints:      task.Constraints,
		Env:              task.Env,
		ResourceRequests: task.ResourceRequests,
		Repo:             task.Repo,
		Timeout:          task.Timeout,
	}

	var taskDef interface{}
	if task.Kind == build.TaskKindDeno {
		def.Deno = &DenoDefinition{}
		taskDef = &def.Deno

	} else if task.Kind == build.TaskKindDockerfile {
		def.Dockerfile = &DockerfileDefinition{}
		taskDef = &def.Dockerfile

	} else if task.Kind == build.TaskKindGo {
		def.Go = &GoDefinition{}
		taskDef = &def.Go

	} else if task.Kind == build.TaskKindNode {
		def.Node = &NodeDefinition{}
		taskDef = &def.Node

	} else if task.Kind == build.TaskKindPython {
		def.Python = &PythonDefinition{}
		taskDef = &def.Python

	} else if task.Kind == build.TaskKindImage {
		def.Image = &ImageDefinition{
			Command: task.Command,
		}
		if task.Image != nil {
			def.Image.Image = *task.Image
		}

	} else if task.Kind == build.TaskKindShell {
		def.Shell = &ShellDefinition{}
		taskDef = &def.Shell

	} else if task.Kind == build.TaskKindSQL {
		def.SQL = &SQLDefinition{}
		taskDef = &def.SQL

	} else if task.Kind == build.TaskKindREST {
		def.REST = &RESTDefinition{}
		taskDef = &def.REST

	} else {
		return Definition_0_2{}, errors.Errorf("unknown kind specified: %s", task.Kind)
	}

	if taskDef != nil && task.KindOptions != nil {
		if err := mapstructure.Decode(task.KindOptions, taskDef); err != nil {
			return Definition_0_2{}, errors.Wrap(err, "decoding options")
		}
	}

	return def, nil
}

func (def *Definition_0_2) GetKindAndOptions() (build.TaskKind, build.KindOptions, error) {
	options := build.KindOptions{}
	if def.Deno != nil {
		if err := mapstructure.Decode(def.Deno, &options); err != nil {
			return "", build.KindOptions{}, errors.Wrap(err, "decoding Deno definition")
		}
		return build.TaskKindDeno, options, nil
	} else if def.Dockerfile != nil {
		if err := mapstructure.Decode(def.Dockerfile, &options); err != nil {
			return "", build.KindOptions{}, errors.Wrap(err, "decoding Dockerfile definition")
		}
		return build.TaskKindDockerfile, options, nil
	} else if def.Image != nil {
		return build.TaskKindImage, build.KindOptions{}, nil
	} else if def.Go != nil {
		if err := mapstructure.Decode(def.Go, &options); err != nil {
			return "", build.KindOptions{}, errors.Wrap(err, "decoding Go definition")
		}
		return build.TaskKindGo, options, nil
	} else if def.Node != nil {
		if err := mapstructure.Decode(def.Node, &options); err != nil {
			return "", build.KindOptions{}, errors.Wrap(err, "decoding Node definition")
		}
		return build.TaskKindNode, options, nil
	} else if def.Python != nil {
		if err := mapstructure.Decode(def.Python, &options); err != nil {
			return "", build.KindOptions{}, errors.Wrap(err, "decoding Python definition")
		}
		return build.TaskKindPython, options, nil
	} else if def.Shell != nil {
		if err := mapstructure.Decode(def.Shell, &options); err != nil {
			return "", build.KindOptions{}, errors.Wrap(err, "decoding Shell definition")
		}
		return build.TaskKindShell, options, nil
	} else if def.SQL != nil {
		if err := mapstructure.Decode(def.SQL, &options); err != nil {
			return "", build.KindOptions{}, errors.Wrap(err, "decoding SQL definition")
		}
		return build.TaskKindSQL, options, nil
	} else if def.REST != nil {
		if err := mapstructure.Decode(def.REST, &options); err != nil {
			return "", build.KindOptions{}, errors.Wrap(err, "decoding REST definition")
		}

		// API expects a single body field to be a string. For convenience, we allow the YAML definition to be a
		// structured object when the jsonBody is actually valid JSON. In that case, if it's not a string, we
		// JSON-serialize it into a string.
		// API also expects a bodyType key.
		if options["jsonBody"] != nil {
			if _, ok := options["jsonBody"].(string); !ok && options["jsonBody"] != nil {
				jsonBody, err := json.Marshal(options["jsonBody"])
				if err != nil {
					return "", build.KindOptions{}, errors.Wrap(err, "marshalling JSON body")
				}
				options["body"] = string(jsonBody)
			} else {
				options["body"] = options["jsonBody"]
			}
			options["bodyType"] = "json"
			delete(options, "jsonBody")

		} else if options["formUrlEncodedBody"] != nil {
			options["formData"] = options["formUrlEncodedBody"]
			options["bodyType"] = "x-www-form-urlencoded"
			delete(options, "formUrlEncodedBody")

		} else if options["formDataBody"] != nil {
			options["formData"] = options["formDataBody"]
			options["bodyType"] = "form-data"
			delete(options, "formDataBody")

		} else {
			options["bodyType"] = "raw"
		}

		return build.TaskKindREST, options, nil
	}

	return "", build.KindOptions{}, errors.New("No kind specified")
}

func (def *Definition_0_2) SetWorkdir(taskroot, workdir string) error {
	// TODO: currently only a concept on Node - should be generalized to all builders.
	if def.Node == nil {
		return nil
	}

	def.SetBuildConfig("workdir", strings.TrimPrefix(workdir, taskroot))
	return nil
}

func (def Definition_0_2) Validate() (Definition_0_2, error) {
	if def.Slug == "" {
		return def, errors.New("Expected a task slug")
	}

	defs := []string{}
	if def.Deno != nil {
		defs = append(defs, "deno")
	}
	if def.Dockerfile != nil {
		defs = append(defs, "dockerfile")
	}
	if def.Image != nil {
		defs = append(defs, "image")
	}
	if def.Go != nil {
		defs = append(defs, "go")
	}
	if def.Node != nil {
		defs = append(defs, "node")
	}
	if def.Python != nil {
		defs = append(defs, "python")
	}
	if def.SQL != nil {
		defs = append(defs, "sql")
	}
	if def.REST != nil {
		defs = append(defs, "rest")
	}

	if len(defs) == 0 {
		return def, errors.New("No task type defined")
	}
	if len(defs) > 1 {
		return def, errors.Errorf("Too many task types defined: only one of (%s) expected", strings.Join(defs, ", "))
	}

	// TODO: validate the rest of the fields!

	return def, nil
}

// Upgrades this task definition for JST interpolation.
// Assumes only usage of expressions is {{JSON}}.
func (def *Definition_0_2) UpgradeJST() error {
	def.Arguments = upgradeArguments(def.Arguments)
	return nil
}

func (def *Definition_0_2) GetEnv() (api.TaskEnv, error) {
	return def.Env, nil
}

func (def *Definition_0_2) GetSlug() string {
	return def.Slug
}

func (def *Definition_0_2) GetUpdateTaskRequest(ctx context.Context, client api.IAPIClient, currentTask *api.Task) (api.UpdateTaskRequest, error) {
	kind, options, err := def.GetKindAndOptions()
	if err != nil {
		return api.UpdateTaskRequest{}, err
	}

	utr := api.UpdateTaskRequest{
		Slug:             def.Slug,
		Name:             def.Name,
		Description:      def.Description,
		Command:          []string{},
		Arguments:        def.Arguments,
		Parameters:       def.Parameters,
		Constraints:      def.Constraints,
		Env:              def.Env,
		ResourceRequests: def.ResourceRequests,
		Resources:        def.Resources,
		Kind:             kind,
		KindOptions:      options,
		Repo:             def.Repo,
		Timeout:          def.Timeout,
	}
	if currentTask != nil {
		utr.Permissions = currentTask.Permissions
		utr.RequireExplicitPermissions = currentTask.RequireExplicitPermissions
		utr.ExecuteRules = currentTask.ExecuteRules
	}
	return utr, nil
}

func (def *Definition_0_2) GetBuildConfig() (build.BuildConfig, error) {
	config := build.BuildConfig{}

	_, options, err := def.GetKindAndOptions()
	if err != nil {
		return nil, err
	}
	for key, val := range options {
		config[key] = val
	}

	for key, val := range def.buildConfig {
		if val == nil { // Nil masks out the value.
			delete(config, key)
		} else {
			config[key] = val
		}
	}

	return config, nil
}

func (def *Definition_0_2) SetBuildConfig(key string, value interface{}) {
	if def.buildConfig == nil {
		def.buildConfig = map[string]interface{}{}
	}
	def.buildConfig[key] = value
}

func (def *Definition_0_2) Entrypoint() (string, error) {
	switch kind, _, _ := def.GetKindAndOptions(); kind {
	case build.TaskKindDeno:
		return def.Deno.Entrypoint, nil
	case build.TaskKindGo:
		return def.Go.Entrypoint, nil
	case build.TaskKindNode:
		return def.Node.Entrypoint, nil
	case build.TaskKindPython:
		return def.Python.Entrypoint, nil
	case build.TaskKindShell:
		return def.Shell.Entrypoint, nil
	case build.TaskKindImage, build.TaskKindDockerfile, build.TaskKindSQL, build.TaskKindREST:
		return "", ErrNoEntrypoint
	default:
		return "", errors.Errorf("unexpected kind %q", kind)
	}
}

func (def *Definition_0_2) Write(path string) error {
	return errors.New("not implemented")
}

func (d Definition_0_2) upgrade() (DefinitionInterface, error) {
	return nil, errors.New("not implemented")
}
