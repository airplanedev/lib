package discover

import (
	"context"
	"path/filepath"

	"github.com/airplanedev/lib/pkg/api"
	"github.com/airplanedev/lib/pkg/deploy/taskdir"
	"github.com/airplanedev/lib/pkg/deploy/taskdir/definitions"
	"github.com/airplanedev/lib/pkg/runtime"
	"github.com/airplanedev/lib/pkg/utils/logger"
)

type DefnDiscoverer struct {
	Client             api.IAPIClient
	AssumeYes          bool
	AssumeNo           bool
	Logger             logger.Logger
	MissingTaskHandler func(context.Context, definitions.DefinitionInterface) (*api.Task, error)
}

var _ TaskDiscoverer = &DefnDiscoverer{}

func (dd *DefnDiscoverer) IsAirplaneTask(ctx context.Context, file string) (slug string, err error) {
	if !definitions.IsTaskDef(file) {
		return "", nil
	}

	def, err := getDef(file)
	if err != nil {
		return "", err
	}

	return def.GetSlug(), nil
}

func (dd *DefnDiscoverer) GetTaskConfig(ctx context.Context, task api.Task, file string) (TaskConfig, error) {
	dir, err := taskdir.Open(file, true)
	if err != nil {
		return TaskConfig{}, err
	}
	defer dir.Close()

	def, err := dir.ReadDefinition_0_3()
	if err != nil {
		return TaskConfig{}, err
	}

	tc := TaskConfig{
		Task: task,
		Def:  &def,
	}

	entrypoint, err := def.Entrypoint()
	if err == definitions.ErrNoEntrypoint {
		// nothing
	} else if err != nil {
		return TaskConfig{}, err
	} else {
		defnDir := filepath.Dir(dir.DefinitionPath())
		absEntrypoint, err := filepath.Abs(filepath.Join(defnDir, entrypoint))
		if err != nil {
			return TaskConfig{}, err
		}
		tc.TaskEntrypoint = absEntrypoint

		r, err := runtime.Lookup(entrypoint, task.Kind)
		if err != nil {
			return TaskConfig{}, err
		}

		taskroot, err := r.Root(absEntrypoint)
		if err != nil {
			return TaskConfig{}, err
		}
		tc.TaskRoot = taskroot

		wd, err := r.Workdir(absEntrypoint)
		if err != nil {
			return TaskConfig{}, err
		}
		def.SetWorkdir(taskroot, wd)

		// Entrypoint for builder needs to be relative to taskroot, not definition directory.
		if defnDir != taskroot {
			ep, err := filepath.Rel(taskroot, absEntrypoint)
			if err != nil {
				return TaskConfig{}, err
			}
			def.SetEntrypoint(ep)
		}
	}

	return tc, nil
}

func (dd *DefnDiscoverer) TaskConfigSource() TaskConfigSource {
	return TaskConfigSourceDefn
}

func (dd *DefnDiscoverer) HandleMissingTask(ctx context.Context, file string) (*api.Task, error) {
	if dd.MissingTaskHandler != nil {
		def, err := getDef(file)
		if err != nil {
			return nil, err
		}
		return dd.MissingTaskHandler(ctx, def)
	}
	return nil, nil
}

func getDef(file string) (definitions.DefinitionInterface, error) {
	dir, err := taskdir.Open(file, true)
	if err != nil {
		return &definitions.Definition_0_3{}, err
	}
	defer dir.Close()

	def, err := dir.ReadDefinition_0_3()
	if err != nil {
		return &definitions.Definition_0_3{}, err
	}

	return &def, nil
}
