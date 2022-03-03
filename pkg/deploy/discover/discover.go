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
	TaskID         string
	TaskRoot       string
	TaskEntrypoint string
	Def            definitions.DefinitionInterface
	Source         TaskConfigSource
}

type TaskDiscoverer interface {
	// IsAirplaneTask inspects a file and if that file represents an Airplane task, it returns
	// that task's slug. If that file is not a task, it will return an empty string.
	IsAirplaneTask(ctx context.Context, file string) (string, error)
	// GetTaskConfig converts an Airplane task file into a fully-qualified task definition.
	// If the task should not be discovered as an Airplane task, a nil task config is returned.
	GetTaskConfig(ctx context.Context, file string) (*TaskConfig, error)
	// TaskConfigSource returns a unique identifier of this TaskDiscoverer.
	TaskConfigSource() TaskConfigSource
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
				slug := tc.Def.GetSlug()
				if _, ok := taskConfigsBySlug[slug]; !ok {
					taskConfigsBySlug[slug] = []TaskConfig{}
				}
				taskConfigsBySlug[slug] = append(taskConfigsBySlug[slug], tc)
			}
		} else {
			// We found a file.
			for _, td := range d.TaskDiscoverers {
				taskConfig, err := td.GetTaskConfig(ctx, p)
				if err != nil {
					return nil, err
				}
				if taskConfig == nil {
					// This file is not an Airplane task.
					continue
				}
				slug := taskConfig.Def.GetSlug()
				if _, ok := taskConfigsBySlug[slug]; !ok {
					taskConfigsBySlug[slug] = []TaskConfig{}
				}
				taskConfigsBySlug[slug] = append(taskConfigsBySlug[slug], *taskConfig)
			}
		}
	}

	return d.deduplicateTaskConfigs(taskConfigsBySlug), nil
}

// Given a map of slug -> [task config, ...], returns a list of task configs unique by slug, sorted
// by slug. Task configs are chosen based on order of TaskDiscoverers & order of discovery.
func (d Discoverer) deduplicateTaskConfigs(taskConfigsBySlug map[string][]TaskConfig) []TaskConfig {
	// Short-circuit if we have no task configs.
	if len(taskConfigsBySlug) == 0 {
		return nil
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

		// Short-circuit if there's only one task config in the list.
		if len(tcs) == 1 {
			taskConfigs[i] = tcs[0]
			continue
		}

		// Otherwise, loop through the TaskDiscoverers. Take the first task config that matches the
		// discoverer in this order.
		found := false
		for _, td := range d.TaskDiscoverers {
			for _, tc := range tcs {
				if td.TaskConfigSource() == tc.Source {
					taskConfigs[i] = tc
					found = true
					break
				}
			}
			if found {
				break
			}
		}
	}
	return taskConfigs
}
