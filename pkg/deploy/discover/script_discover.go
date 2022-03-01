package discover

import (
	"context"
	"path/filepath"

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

func (sd *ScriptDiscoverer) GetTaskConfig(ctx context.Context, file string) (*TaskConfig, error) {
	slug, _ = runtime.Slug(file)

	// TODO: get task by slug, handle missing task

	def, err := definitions.NewDefinitionFromTask(task)
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
	if err := def.SetEntrypoint(taskroot, absFile); err != nil {
		return nil, err
	}

	wd, err := r.Workdir(absFile)
	if err != nil {
		return nil, err
	}
	def.SetWorkdir(taskroot, wd)

	return &TaskConfig{
		// TaskID: "TODO",
		TaskRoot:       taskroot,
		TaskEntrypoint: absFile,
		Def:            &def,
		From:           sd.TaskConfigSource(),
	}, nil
}

func (sd *ScriptDiscoverer) TaskConfigSource() TaskConfigSource {
	return TaskConfigSourceScript
}
