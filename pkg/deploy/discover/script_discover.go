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
	"github.com/pkg/errors"
)

type ScriptDiscoverer struct {
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
	if err := def.SetEntrypoint(taskroot, absFile); err != nil {
		return TaskConfig{}, err
	}

	wd, err := r.Workdir(absFile)
	if err != nil {
		return TaskConfig{}, err
	}
	if err := def.SetWorkdir(taskroot, wd); err != nil {
		return TaskConfig{}, err
	}

	return TaskConfig{
		TaskRoot:       taskroot,
		TaskEntrypoint: absFile,
		Def:            &def,
		Task:           task,
	}, nil
}

func (sd *ScriptDiscoverer) TaskConfigSource() TaskConfigSource {
	return TaskConfigSourceScript
}

func (sd *ScriptDiscoverer) HandleMissingTask(ctx context.Context, file string) (*api.Task, error) {
	return nil, nil
}
