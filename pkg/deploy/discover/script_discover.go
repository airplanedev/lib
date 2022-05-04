package discover

import (
	"context"
	"path/filepath"

	"github.com/airplanedev/lib/pkg/api"
	"github.com/airplanedev/lib/pkg/deploy/taskdir/definitions"
	"github.com/airplanedev/lib/pkg/runtime"
	_ "github.com/airplanedev/lib/pkg/runtime/javascript"
	_ "github.com/airplanedev/lib/pkg/runtime/python"
	_ "github.com/airplanedev/lib/pkg/runtime/shell"
	_ "github.com/airplanedev/lib/pkg/runtime/sql"
	_ "github.com/airplanedev/lib/pkg/runtime/typescript"
	"github.com/airplanedev/lib/pkg/utils/logger"
	"github.com/airplanedev/lib/pkg/utils/pathcase"
	"github.com/pkg/errors"
)

type ScriptDiscoverer struct {
	Client  api.IAPIClient
	Logger  logger.Logger
	EnvSlug string
}

var _ TaskDiscoverer = &ScriptDiscoverer{}

func (sd *ScriptDiscoverer) IsAirplaneTask(ctx context.Context, file string) (string, error) {
	slug := runtime.Slug(file)
	return slug, nil
}

func (sd *ScriptDiscoverer) GetTaskConfig(ctx context.Context, file string) (*TaskConfig, error) {
	slug := runtime.Slug(file)
	if slug == "" {
		return nil, nil
	}

	task, err := sd.Client.GetTask(ctx, api.GetTaskRequest{
		Slug:    slug,
		EnvSlug: sd.EnvSlug,
	})
	if err != nil {
		var merr *api.TaskMissingError
		if !errors.As(err, &merr) {
			return nil, errors.Wrap(err, "unable to get task")
		}

		sd.Logger.Warning(`Task with slug %s does not exist, skipping deployment.`, slug)
		return nil, nil
	}
	if task.IsArchived {
		sd.Logger.Warning(`Task with slug %s is archived, skipping deployment.`, slug)
		return nil, nil
	}

	def, err := definitions.NewDefinitionFromTask(ctx, sd.Client, task)
	if err != nil {
		return nil, err
	}

	r, err := runtime.Lookup(file, task.Kind)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot determine how to deploy %q - check your CLI is up to date", file)
	}

	absFile, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}

	taskroot, err := r.Root(absFile)
	if err != nil {
		return nil, err
	}

	// Entrypoint needs to be relative to the taskroot.
	absEntrypoint, err := pathcase.ActualFilename(absFile)
	if err != nil {
		return nil, err
	}
	ep, err := filepath.Rel(taskroot, absEntrypoint)
	if err != nil {
		return nil, err
	}
	def.SetBuildConfig("entrypoint", ep)

	wd, err := r.Workdir(absFile)
	if err != nil {
		return nil, err
	}
	if err := def.SetWorkdir(taskroot, wd); err != nil {
		return nil, err
	}

	return &TaskConfig{
		TaskID:         task.ID,
		TaskRoot:       taskroot,
		TaskEntrypoint: absFile,
		Def:            def,
		Source:         sd.TaskConfigSource(),
	}, nil
}

func (sd *ScriptDiscoverer) TaskConfigSource() TaskConfigSource {
	return TaskConfigSourceScript
}
