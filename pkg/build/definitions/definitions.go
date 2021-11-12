package definitions

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/airplanedev/lib/pkg/build/logger"
	"github.com/airplanedev/lib/pkg/build/utils/pathcase"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// Definition represents a YAML-based task definition that can be used to create
// or update Airplane tasks.
//
// Note this is the subset of fields that can be represented with a revision,
// and therefore isolated to a specific environment.
type Definition Definition_0_2

func NewDefinitionFromTask(task Task) (Definition, error) {
	def := Definition{
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
	if task.Kind == TaskKindDeno {
		def.Deno = &DenoDefinition{}
		taskDef = &def.Deno

	} else if task.Kind == TaskKindDockerfile {
		def.Dockerfile = &DockerfileDefinition{}
		taskDef = &def.Dockerfile

	} else if task.Kind == TaskKindGo {
		def.Go = &GoDefinition{}
		taskDef = &def.Go

	} else if task.Kind == TaskKindNode {
		def.Node = &NodeDefinition{}
		taskDef = &def.Node

	} else if task.Kind == TaskKindPython {
		def.Python = &PythonDefinition{}
		taskDef = &def.Python

	} else if task.Kind == TaskKindImage {
		def.Image = &ImageDefinition{
			Command: task.Command,
		}
		if task.Image != nil {
			def.Image.Image = *task.Image
		}

	} else if task.Kind == TaskKindShell {
		def.Shell = &ShellDefinition{}
		taskDef = &def.Shell

	} else if task.Kind == TaskKindSQL {
		def.SQL = &SQLDefinition{}
		taskDef = &def.SQL

	} else if task.Kind == TaskKindREST {
		def.REST = &RESTDefinition{}
		taskDef = &def.REST

	} else {
		return Definition{}, errors.Errorf("unknown kind specified: %s", task.Kind)
	}

	if taskDef != nil && task.KindOptions != nil {
		if err := mapstructure.Decode(task.KindOptions, taskDef); err != nil {
			return Definition{}, errors.Wrap(err, "decoding options")
		}
	}

	return def, nil
}

func (def Definition) GetKindAndOptions() (TaskKind, KindOptions, error) {
	options := KindOptions{}
	if def.Deno != nil {
		if err := mapstructure.Decode(def.Deno, &options); err != nil {
			return "", KindOptions{}, errors.Wrap(err, "decoding Deno definition")
		}
		return TaskKindDeno, options, nil
	} else if def.Dockerfile != nil {
		if err := mapstructure.Decode(def.Dockerfile, &options); err != nil {
			return "", KindOptions{}, errors.Wrap(err, "decoding Dockerfile definition")
		}
		return TaskKindDockerfile, options, nil
	} else if def.Image != nil {
		return TaskKindImage, KindOptions{}, nil
	} else if def.Go != nil {
		if err := mapstructure.Decode(def.Go, &options); err != nil {
			return "", KindOptions{}, errors.Wrap(err, "decoding Go definition")
		}
		return TaskKindGo, options, nil
	} else if def.Node != nil {
		if err := mapstructure.Decode(def.Node, &options); err != nil {
			return "", KindOptions{}, errors.Wrap(err, "decoding Node definition")
		}
		return TaskKindNode, options, nil
	} else if def.Python != nil {
		if err := mapstructure.Decode(def.Python, &options); err != nil {
			return "", KindOptions{}, errors.Wrap(err, "decoding Python definition")
		}
		return TaskKindPython, options, nil
	} else if def.Shell != nil {
		if err := mapstructure.Decode(def.Shell, &options); err != nil {
			return "", KindOptions{}, errors.Wrap(err, "decoding Shell definition")
		}
		return TaskKindShell, options, nil
	} else if def.SQL != nil {
		if err := mapstructure.Decode(def.SQL, &options); err != nil {
			return "", KindOptions{}, errors.Wrap(err, "decoding SQL definition")
		}
		return TaskKindSQL, options, nil
	} else if def.REST != nil {
		if err := mapstructure.Decode(def.REST, &options); err != nil {
			return "", KindOptions{}, errors.Wrap(err, "decoding REST definition")
		}

		// API expects a single body field to be a string. For convenience, we allow the YAML definition to be a
		// structured object when the jsonBody is actually valid JSON. In that case, if it's not a string, we
		// JSON-serialize it into a string.
		// API also expects a bodyType key.
		if options["jsonBody"] != nil {
			if _, ok := options["jsonBody"].(string); !ok && options["jsonBody"] != nil {
				jsonBody, err := json.Marshal(options["jsonBody"])
				if err != nil {
					return "", KindOptions{}, errors.Wrap(err, "marshalling JSON body")
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

		return TaskKindREST, options, nil
	}

	return "", KindOptions{}, errors.New("No kind specified")
}

// SetEntrypoint computes and normalizes the entrypoint based on the task root and absolute
// path to the entrypoint.
func (def *Definition) SetEntrypoint(taskroot, absEntrypoint string) error {
	var err error
	// Fix casing on entrypoint for case-insensitive filesystems.
	origEntrypoint := absEntrypoint
	absEntrypoint, err = pathcase.ActualFilename(absEntrypoint)
	if err != nil {
		return err
	}
	if absEntrypoint != origEntrypoint {
		logger.Warning(
			"Using %q instead of %q - different casing may not work on all operating systems",
			filepath.Base(absEntrypoint), filepath.Base(origEntrypoint),
		)
	}

	ep, err := filepath.Rel(taskroot, absEntrypoint)
	if err != nil {
		return err
	}

	switch kind, _, _ := def.GetKindAndOptions(); kind {
	case TaskKindNode:
		def.Node.Entrypoint = ep
	case TaskKindPython:
		def.Python.Entrypoint = ep
	case TaskKindShell:
		def.Shell.Entrypoint = ep
	default:
		return errors.Errorf("unexpected kind %q", kind)
	}
	return nil
}

func (def *Definition) SetWorkdir(taskroot, workdir string) {
	// TODO: currently only a concept on Node - should be generalized to all builders.
	if def.Node != nil {
		def.Node.Workdir = strings.TrimPrefix(workdir, taskroot)
	}
}

func (def Definition) Validate() (Definition, error) {
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

var jsonRegex = regexp.MustCompile(`{{ *JSON *}}`)

// Upgrades this task definition for JST interpolation.
// Assumes only usage of expressions is {{JSON}}.
func (def *Definition) UpgradeJST() {
	args := []string{}
	for _, arg := range def.Arguments {
		jstArg := jsonRegex.ReplaceAllString(arg, "{{JSON.stringify(params)}}")
		args = append(args, jstArg)
	}
	def.Arguments = args
}

func UnmarshalDefinition(buf []byte, defPath string) (Definition, error) {
	// Validate definition against our Definition struct
	if err := validateYAML(buf, Definition{}); err != nil {
		// Try older definitions?
		if def, oerr := tryOlderDefinitions(buf); oerr == nil {
			return def, nil
		}

		// Print any "expected" validation errors
		switch err := errors.Cause(err).(type) {
		case ErrInvalidYAML:
			return Definition{}, newErrReadDefinition(fmt.Sprintf("Error reading %s, invalid YAML:\n  %s", defPath, err))
		case ErrSchemaValidation:
			errorMsgs := []string{}
			for _, verr := range err.Errors {
				errorMsgs = append(errorMsgs, fmt.Sprintf("%s: %s", verr.Field(), verr.Description()))
			}
			return Definition{}, newErrReadDefinition(fmt.Sprintf("Error reading %s", defPath), errorMsgs...)
		default:
			return Definition{}, errors.Wrapf(err, "reading %s", defPath)
		}
	}

	var def Definition
	if err := yaml.Unmarshal(buf, &def); err != nil {
		return Definition{}, errors.Wrap(err, "unmarshalling task definition")
	}

	return def, nil
}

func tryOlderDefinitions(buf []byte) (Definition, error) {
	var err error
	if err = validateYAML(buf, Definition_0_1{}); err == nil {
		var def Definition_0_1
		if e := yaml.Unmarshal(buf, &def); e != nil {
			return Definition{}, err
		}
		return def.upgrade()
	}
	return Definition{}, err
}
