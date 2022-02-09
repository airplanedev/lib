package discover

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"sort"

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

	// EnvSlug is the slug of the environment to look for discovered tasks in.
	//
	// If a task is discovered, but doesn't exist in this environment, then the task
	// is treated as missing.
	EnvSlug string
}

// DiscoverTasks recursively discovers Airplane tasks. Only one task config per slug is returned.
// If there are multiple tasks discovered with the same slug, the order of the discoverers takes
// precedence; if a single discoverer discovers multiple tasks with the same slug, the first task
// discovered takes precedence. Task configs are returned in alphabetical order of their slugs.
func (d *Discoverer) DiscoverTasks(ctx context.Context, paths ...string) ([]TaskConfig, error) {
	taskConfigsBySlug := map[string][]TaskConfig{}
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
			for _, tc := range nestedTaskConfigs {
				slug := tc.Task.Slug
				if _, ok := taskConfigsBySlug[slug]; !ok {
					taskConfigsBySlug[slug] = []TaskConfig{}
				}
				taskConfigsBySlug[slug] = append(taskConfigsBySlug[slug], tc)
			}
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
			if _, ok := taskConfigsBySlug[slug]; !ok {
				taskConfigsBySlug[slug] = []TaskConfig{}
			}
			taskConfigsBySlug[slug] = append(taskConfigsBySlug[slug], taskConfig)
		}
	}

	if len(taskConfigsBySlug) == 0 {
		return nil, nil
	}

	// Sort by slugs, so we have a deterministic order.
	slugs := make([]string, 0, len(taskConfigsBySlug))
	for slug := range taskConfigsBySlug {
		slugs = append(slugs, slug)
	}
	sort.Strings(slugs)

	taskConfigs := make([]TaskConfig, len(slugs))
	for i, slug := range slugs {
		tcs := taskConfigsBySlug[slug]
		if len(tcs) == 1 {
			taskConfigs[i] = tcs[0]
			continue
		}
		found := false
		for _, td := range d.TaskDiscoverers {
			if found {
				break
			}
			for _, tc := range tcs {
				if td.TaskConfigSource() == tc.From {
					taskConfigs[i] = tc
					found = true
					break
				}
			}
		}
	}
	return taskConfigs, nil
}
