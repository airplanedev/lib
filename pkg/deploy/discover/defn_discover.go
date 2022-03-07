package discover

import (
	"context"
	"path/filepath"

	"github.com/airplanedev/lib/pkg/api"
	"github.com/airplanedev/lib/pkg/deploy/taskdir"
	"github.com/airplanedev/lib/pkg/deploy/taskdir/definitions"
	"github.com/airplanedev/lib/pkg/runtime"
	"github.com/airplanedev/lib/pkg/utils/logger"
	"github.com/pkg/errors"
)

type DefnDiscoverer struct {
	Client    api.IAPIClient
	AssumeYes bool
	AssumeNo  bool
	Logger    logger.Logger

	// MissingTaskHandler is called from `GetTaskConfig` if a task ID cannot be found for a definition
	// file. The handler should either create the task and return the created task's TaskMetadata, or
	// it should return `nil` to signal that the definition should be ignored. If not set, these
	// definitions are ignored.
	MissingTaskHandler func(context.Context, definitions.DefinitionInterface) (*api.TaskMetadata, error)
}

var _ TaskDiscoverer = &DefnDiscoverer{}

func (dd *DefnDiscoverer) IsAirplaneTask(ctx context.Context, file string) (string, error) {
	if !definitions.IsTaskDef(file) {
		return "", nil
	}

	dir, err := taskdir.Open(file)
	if err != nil {
		return "", err
	}
	defer dir.Close()

	def, err := dir.ReadDefinition()
	if err != nil {
		return "", err
	}

	return def.GetSlug(), nil
}

func (dd *DefnDiscoverer) GetTaskConfig(ctx context.Context, file string) (*TaskConfig, error) {
	if !definitions.IsTaskDef(file) {
		return nil, nil
	}

	dir, err := taskdir.Open(file)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	def, err := dir.ReadDefinition()
	if err != nil {
		return nil, err
	}

	tc := TaskConfig{
		Def:    def,
		Source: dd.TaskConfigSource(),
	}

	metadata, err := dd.Client.GetTaskMetadata(ctx, def.GetSlug())
	if err != nil {
		var merr *api.TaskMissingError
		if !errors.As(err, &merr) {
			return nil, errors.Wrap(err, "unable to get task metadata")
		}

		if dd.MissingTaskHandler == nil {
			return nil, nil
		}

		mptr, err := dd.MissingTaskHandler(ctx, def)
		if err != nil {
			return nil, err
		} else if mptr == nil {
			dd.Logger.Warning(`Task with slug %s does not exist, skipping deploy.`, def.GetSlug())
			return nil, nil
		}
		metadata = *mptr
	}
	tc.TaskID = metadata.ID

	entrypoint, err := def.GetAbsoluteEntrypoint()
	if err == definitions.ErrNoEntrypoint {
		return &tc, nil
	} else if err != nil {
		return nil, err
	} else {
		tc.TaskEntrypoint = entrypoint

		kind, _, err := def.GetKindAndOptions()
		if err != nil {
			return nil, err
		}

		r, err := runtime.Lookup(entrypoint, kind)
		if err != nil {
			return nil, err
		}

		taskroot, err := r.Root(entrypoint)
		if err != nil {
			return nil, err
		}
		tc.TaskRoot = taskroot

		wd, err := r.Workdir(entrypoint)
		if err != nil {
			return nil, err
		}
		if err := def.SetWorkdir(taskroot, wd); err != nil {
			return nil, err
		}

		// Entrypoint for builder needs to be relative to taskroot, not definition directory.
		defnDir := filepath.Dir(dir.DefinitionPath())
		if defnDir != taskroot {
			ep, err := filepath.Rel(taskroot, entrypoint)
			if err != nil {
				return nil, err
			}
			def.SetBuildConfig("entrypoint", ep)
		}
	}

	return &tc, nil
}

func (dd *DefnDiscoverer) TaskConfigSource() TaskConfigSource {
	return TaskConfigSourceDefn
}
