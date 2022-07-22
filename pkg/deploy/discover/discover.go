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

type ConfigSource string

const (
	ConfigSourceScript ConfigSource = "script"
	ConfigSourceDefn   ConfigSource = "defn"
	ConfigSourceCode   ConfigSource = "code"
)

type TaskConfig struct {
	TaskID         string
	TaskRoot       string
	TaskEntrypoint string
	Def            definitions.DefinitionInterface
	Source         ConfigSource
}

func (c TaskConfig) GetSource() ConfigSource {
	return c.Source
}

type ViewConfig struct {
	ID     string
	Root   string
	Def    definitions.ViewDefinition
	Source ConfigSource
}

func (c ViewConfig) GetSource() ConfigSource {
	return c.Source
}

type ConfigDiscoverer interface {
	ConfigSource() ConfigSource
}

type TaskDiscoverer interface {
	// IsAirplaneTask inspects a file and if that file represents an Airplane task, it returns
	// that task's slug. If that file is not a task, it will return an empty string.
	IsAirplaneTask(ctx context.Context, file string) (string, error)
	// GetTaskConfig converts an Airplane task file into a fully-qualified task definition.
	// If the task should not be discovered as an Airplane task, a nil task config is returned.
	GetTaskConfigs(ctx context.Context, file string) ([]TaskConfig, error)
	// ConfigSource returns a unique identifier of this Discoverer.
	ConfigSource() ConfigSource
}

type ViewDiscoverer interface {
	// GetTaskConfig converts an Airplane task file into a fully-qualified task definition.
	// If the task should not be discovered as an Airplane task, a nil task config is returned.
	GetViewConfig(ctx context.Context, file string) (*ViewConfig, error)
	// ConfigSource returns a unique identifier of this Discoverer.
	ConfigSource() ConfigSource
}

type Discoverer struct {
	TaskDiscoverers []TaskDiscoverer
	ViewDiscoverers []ViewDiscoverer
	Client          api.IAPIClient
	Logger          logger.Logger

	// EnvSlug is the slug of the environment to look for discovered tasks in.
	//
	// If a task is discovered, but doesn't exist in this environment, then the task
	// is treated as missing.
	EnvSlug string
}

// Discover recursively discovers Airplane tasks & views. Only one config per slug is returned.
// If there are multiple configs discovered with the same slug, the order of the discoverers takes
// precedence; if a single discoverer discovers multiple configs with the same slug, the first config
// discovered takes precedence. Configs are returned in alphabetical order of their slugs.
func (d *Discoverer) Discover(ctx context.Context, paths ...string) ([]TaskConfig, []ViewConfig, error) {
	taskConfigsBySlug := map[string][]TaskConfig{}
	viewConfigsBySlug := map[string][]ViewConfig{}
	for _, p := range paths {
		if ignoredDirectories[p] {
			continue
		}
		fileInfo, err := os.Stat(p)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "determining if %s is file or directory", p)
		}

		if fileInfo.IsDir() {
			// We found a directory. Recursively explore all of the files and directories in it.
			nestedFiles, err := ioutil.ReadDir(p)
			if err != nil {
				return nil, nil, errors.Wrapf(err, "reading directory %s", p)
			}
			var nestedPaths []string
			for _, nestedFile := range nestedFiles {
				nestedPaths = append(nestedPaths, path.Join(p, nestedFile.Name()))
			}
			nestedTaskConfigs, nestedViewConfigs, err := d.Discover(ctx, nestedPaths...)
			if err != nil {
				return nil, nil, err
			}
			for _, tc := range nestedTaskConfigs {
				slug := tc.Def.GetSlug()
				if _, ok := taskConfigsBySlug[slug]; !ok {
					taskConfigsBySlug[slug] = []TaskConfig{}
				}
				taskConfigsBySlug[slug] = append(taskConfigsBySlug[slug], tc)
			}
			for _, vc := range nestedViewConfigs {
				slug := vc.Def.Slug
				if _, ok := viewConfigsBySlug[slug]; !ok {
					viewConfigsBySlug[slug] = []ViewConfig{}
				}
				viewConfigsBySlug[slug] = append(viewConfigsBySlug[slug], vc)
			}
		} else {
			// We found a file.
			for _, td := range d.TaskDiscoverers {
				taskConfigs, err := td.GetTaskConfigs(ctx, p)
				if err != nil {
					return nil, nil, err
				}
				if len(taskConfigs) == 0 {
					// This file is not an Airplane task.
					continue
				}

				for _, taskConfig := range taskConfigs {
					slug := taskConfig.Def.GetSlug()
					if _, ok := taskConfigsBySlug[slug]; !ok {
						taskConfigsBySlug[slug] = []TaskConfig{}
					}
					taskConfigsBySlug[slug] = append(taskConfigsBySlug[slug], taskConfig)
				}
			}
			for _, vd := range d.ViewDiscoverers {
				viewConfig, err := vd.GetViewConfig(ctx, p)
				if err != nil {
					return nil, nil, err
				}
				if viewConfig == nil {
					// This file is not an Airplane view.
					continue
				}
				slug := viewConfig.Def.Slug
				if _, ok := viewConfigsBySlug[slug]; !ok {
					viewConfigsBySlug[slug] = []ViewConfig{}
				}
				viewConfigsBySlug[slug] = append(viewConfigsBySlug[slug], *viewConfig)
			}
		}
	}

	return deduplicateConfigs(taskConfigsBySlug, d.TaskDiscoverers), deduplicateConfigs(viewConfigsBySlug, d.ViewDiscoverers), nil
}

// deduplicateConfigs returns a list of configs unique by slug, sorted by slug
// from a map of slug -> [task config, ...]. Configs are chosen based on order of Discoverers & order of discovery.
func deduplicateConfigs[C interface{ GetSource() ConfigSource }, D ConfigDiscoverer](taskConfigsBySlug map[string][]C, configDiscoverers []D) []C {
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

	taskConfigs := make([]C, len(slugs))
	for i, slug := range slugs {
		tcs := taskConfigsBySlug[slug]

		// Short-circuit if there's only one task config in the list.
		if len(tcs) == 1 {
			taskConfigs[i] = tcs[0]
			continue
		}

		// Otherwise, loop through the Discoverers. Take the first task config that matches the
		// discoverer in this order.
		found := false
		for _, td := range configDiscoverers {
			for _, tc := range tcs {
				if td.ConfigSource() == tc.GetSource() {
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
