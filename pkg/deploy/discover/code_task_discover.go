package discover

import (
	"context"
	_ "embed"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/airplanedev/lib/pkg/api"
	"github.com/airplanedev/lib/pkg/build"
	"github.com/airplanedev/lib/pkg/deploy/taskdir/definitions"
	"github.com/airplanedev/lib/pkg/runtime"
	"github.com/airplanedev/lib/pkg/utils/logger"
	"github.com/airplanedev/lib/pkg/utils/pathcase"
	"github.com/pkg/errors"
)

//go:embed parser/node/parser.ts
var parserScript []byte

type CodeTaskDiscoverer struct {
	Client  api.IAPIClient
	Logger  logger.Logger
	EnvSlug string
}

var _ TaskDiscoverer = &CodeTaskDiscoverer{}

func (c *CodeTaskDiscoverer) IsAirplaneTask(ctx context.Context, file string) (string, error) {
	// Run thing that extracts slug from something in the file.
	return "", nil
}

func (c *CodeTaskDiscoverer) GetTaskConfigs(ctx context.Context, file string) ([]TaskConfig, error) {
	// Only analyze typescript files that end with .tasks.ts
	if !strings.HasSuffix(file, ".tasks.ts") {
		return nil, nil
	}

	// Create a temp file
	tempFile, err := os.CreateTemp("", "airplane.parser.node.*.ts")
	if err != nil {
		return nil, nil
	}
	_, err = tempFile.Write(parserScript)
	if err != nil {
		return nil, err
	}

	// Run parser on the file
	out, err := exec.Command("npx", "-p", "typescript", "-p", "@types/node", "-p", "ts-node",
		"ts-node", tempFile.Name(), file).Output()
	if err != nil {
		return nil, err
	}

	var parsedTasks []map[string]interface{}
	if err := json.Unmarshal(out, &parsedTasks); err != nil {
		return nil, err
	}

	if len(parsedTasks) == 0 {
		// Unable to find any Airplane tasks in the file
		return nil, nil
	}

	// Entrypoint needs to be relative to the taskroot.
	r, err := runtime.Lookup(file, build.TaskKindNode)
	if err != nil {
		return nil, err
	}

	absFile, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}

	taskroot, err := r.Root(absFile)
	if err != nil {
		return nil, err
	}

	absEntrypoint, err := pathcase.ActualFilename(absFile)
	if err != nil {
		return nil, err
	}
	ep, err := filepath.Rel(taskroot, absEntrypoint)
	if err != nil {
		return nil, err
	}

	wd, err := r.Workdir(absFile)
	if err != nil {
		return nil, err
	}

	var taskConfigs []TaskConfig
	for _, parsedTask := range parsedTasks {
		// Construct task config
		def := definitions.Definition_0_3{
			Name: parsedTask["name"].(string),
			Slug: parsedTask["slug"].(string),
			Node: &definitions.NodeDefinition_0_3{
				NodeVersion: "18",
			},
		}

		task, err := c.Client.GetTask(ctx, api.GetTaskRequest{
			Slug:    def.Slug,
			EnvSlug: c.EnvSlug,
		})
		if err != nil {
			var merr *api.TaskMissingError
			if !errors.As(err, &merr) {
				return nil, errors.Wrap(err, "unable to get task")
			}

			c.Logger.Warning(`Task with slug %s does not exist, skipping deployment.`, def.Slug)
			return nil, nil
		}
		if task.IsArchived {
			c.Logger.Warning(`Task with slug %s is archived, skipping deployment.`, def.Slug)
			return nil, nil
		}

		parameters := parsedTask["parameters"].(map[string]interface{})
		var params []definitions.ParameterDefinition_0_3
		for taskParamSlug, taskParam := range parameters {
			constructedParam := taskParam.(map[string]interface{})
			params = append(params, definitions.ParameterDefinition_0_3{
				Name: constructedParam["name"].(string),
				Slug: taskParamSlug,
				Type: constructedParam["kind"].(string),
			})
		}
		def.Parameters = params

		if err := def.SetAbsoluteEntrypoint(absFile); err != nil {
			return nil, err
		}

		def.SetBuildConfig("entrypoint", ep)
		def.SetBuildConfig("entrypointFunc", parsedTask["entrypointFunc"].(string))
		def.SetBuildConfig("isCodeDefinedTask", true)

		var paramSlugs []string
		for _, param := range def.Parameters {
			// Preserve order of parameters
			paramSlugs = append(paramSlugs, param.Slug)
		}
		def.SetBuildConfig("paramSlugs", paramSlugs)

		if err := def.SetWorkdir(taskroot, wd); err != nil {
			return nil, err
		}

		taskConfigs = append(taskConfigs, TaskConfig{
			TaskID:         task.ID,
			TaskRoot:       taskroot,
			TaskEntrypoint: absFile,
			Def:            &def,
			Source:         ConfigSourceCode,
		})
	}

	return taskConfigs, nil
}

func (c *CodeTaskDiscoverer) ConfigSource() ConfigSource {
	return ConfigSourceCode
}
