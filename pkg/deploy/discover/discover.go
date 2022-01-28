package discover

import (
	"context"
	"io/ioutil"
	"os"
	"path"

	"github.com/airplanedev/lib/pkg/api"
	"github.com/airplanedev/lib/pkg/deploy/taskdir/definitions"
	"github.com/airplanedev/lib/pkg/utils/logger"
	"github.com/pkg/errors"
)

var ignoredDirectories = map[string]bool{
	"node_modules": true,
	"__pycache__":  true,
	".git":         true,
}

type TaskConfigSource string

const (
	TaskConfigSourceScript TaskConfigSource = "script"
	TaskConfigSourceDefn   TaskConfigSource = "defn"
)

type TaskConfig struct {
	TaskRoot         string
	WorkingDirectory string
	TaskEntrypoint   string
	Task             api.Task
	Def              definitions.DefinitionInterface
	From             TaskConfigSource
}

type TaskDiscoverer interface {
	IsAirplaneTask(ctx context.Context, file string) (slug string, err error)
	GetTaskConfig(ctx context.Context, task api.Task, file string) (TaskConfig, error)
	TaskConfigSource() TaskConfigSource
	HandleMissingTask(ctx context.Context, file string) (*api.Task, error)
}

type Discoverer struct {
	TaskDiscoverers []TaskDiscoverer
	Client          api.IAPIClient
	Logger          logger.Logger
	EnvSlug         string
}

// DiscoverTasks recursively discovers Airplane tasks.
func (d *Discoverer) DiscoverTasks(ctx context.Context, paths ...string) ([]TaskConfig, error) {
	var taskConfigs []TaskConfig
	for _, p := range paths {
		if ignoredDirectories[p] {
			continue
		}
		fileInfo, err := os.Stat(p)
		if err != nil {
			return nil, errors.Wrapf(err, "determining if %s is file or directory", p)
		}

		if fileInfo.IsDir() {
			// We found a directory. Recursively explore all of the files and directories in it.
			nestedFiles, err := ioutil.ReadDir(p)
			if err != nil {
				return nil, errors.Wrapf(err, "reading directory %s", p)
			}
			var nestedPaths []string
			for _, nestedFile := range nestedFiles {
				nestedPaths = append(nestedPaths, path.Join(p, nestedFile.Name()))
			}
			nestedTaskConfigs, err := d.DiscoverTasks(ctx, nestedPaths...)
			if err != nil {
				return nil, err
			}
			taskConfigs = append(taskConfigs, nestedTaskConfigs...)
			continue
		}
		// We found a file.
		for _, td := range d.TaskDiscoverers {
			slug, err := td.IsAirplaneTask(ctx, p)
			if err != nil {
				return nil, err
			}
			if slug == "" {
				// The file is not an Airplane task.
				continue
			}
			task, err := d.Client.GetTask(ctx, api.GetTaskRequest{
				Slug:    slug,
				EnvSlug: d.EnvSlug,
			})
			if err != nil {
				var missingErr *api.TaskMissingError
				if errors.As(err, &missingErr) {
					taskPtr, err := td.HandleMissingTask(ctx, p)
					if err != nil {
						return nil, err
					} else if taskPtr == nil {
						d.Logger.Warning(`Task with slug %s does not exist, skipping deploy.`, slug)
						continue
					}
					task = *taskPtr
				} else {
					return nil, err
				}
			}
			taskConfig, err := td.GetTaskConfig(ctx, task, p)
			if err != nil {
				return nil, err
			}
			taskConfig.From = td.TaskConfigSource()
			taskConfigs = append(taskConfigs, taskConfig)
		}
	}

	return taskConfigs, nil
}
