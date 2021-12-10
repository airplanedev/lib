package discover

import (
	"context"
	"path/filepath"

	"github.com/airplanedev/lib/pkg/common/api"
	"github.com/airplanedev/lib/pkg/common/logger"
	"github.com/airplanedev/lib/pkg/deploy/taskdir/definitions"
	"github.com/airplanedev/lib/pkg/runtime"
	"github.com/pkg/errors"
)

type ScriptDiscoverer struct {
	Logger logger.Logger
}

var _ TaskDiscoverer = &ScriptDiscoverer{}

func (sd *ScriptDiscoverer) IsAirplaneTask(ctx context.Context, file string) (slug string, err error) {
	slug, _ = runtime.Slug(file)
	return
}

func (sd *ScriptDiscoverer) GetTaskConfig(ctx context.Context, task api.Task, file string) (TaskConfig, error) {
	r, err := runtime.Lookup(file, task.Kind)
	if err != nil {
		return TaskConfig{}, errors.Wrapf(err, "cannot determine how to deploy %q - check your CLI is up to date", file)
	}

	def, err := definitions.NewDefinitionFromTask(task)
	if err != nil {
		return TaskConfig{}, err
	}

	absFile, err := filepath.Abs(file)
	if err != nil {
		return TaskConfig{}, err
	}

	taskroot, err := r.Root(absFile)
	if err != nil {
		return TaskConfig{}, err
	}
	if err := def.SetEntrypoint(sd.Logger, taskroot, absFile); err != nil {
		return TaskConfig{}, err
	}

	wd, err := r.Workdir(absFile)
	if err != nil {
		return TaskConfig{}, err
	}
	def.SetWorkdir(taskroot, wd)

	kind, kindOptions, err := def.GetKindAndOptions()
	if err != nil {
		return TaskConfig{}, err
	}

	return TaskConfig{
		TaskRoot:         taskroot,
		WorkingDirectory: wd,
		TaskFilePath:     absFile,
		Def:              &def,
		Kind:             kind,
		KindOptions:      kindOptions,
		Task:             task,
	}, nil
}

func (sd *ScriptDiscoverer) TaskConfigSource() TaskConfigSource {
	return TaskConfigSourceScript
}

func (sd *ScriptDiscoverer) HandleMissingTask(ctx context.Context, file string) (api.Task, error) {
	return api.Task{}, nil
}
