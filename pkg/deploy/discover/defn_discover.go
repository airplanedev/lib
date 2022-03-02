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

	// MissingTaskHandler is fired if `GetTaskConfig` is called when a task ID cannot be found
	// for a definition file. The handler should either create the task and return the created
	// task's TaskMetadata, or it should return `nil` to signal that the definition should be
	// ignored.
	MissingTaskHandler func(context.Context, definitions.DefinitionInterface) (*api.TaskMetadata, error)
}

var _ TaskDiscoverer = &DefnDiscoverer{}

func (dd *DefnDiscoverer) GetTaskConfig(ctx context.Context, file string) (*TaskConfig, error) {
	if !definitions.IsTaskDef(file) {
		return nil, nil
	}

	dir, err := taskdir.Open(file, true)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	def, err := dir.ReadDefinition_0_3()
	if err != nil {
		return nil, err
	}

	tc := TaskConfig{
		Def:    &def,
		Source: dd.TaskConfigSource(),
		// By definition, all task definitions use JSTs:
		TaskInterpolationMode: "jst",
	}

	metadata, err := dd.Client.GetTaskMetadata(ctx, def.GetSlug())
	if err != nil {
		var merr *api.TaskMissingError
		if !errors.As(err, &merr) {
			return nil, err
		}

		mptr, err := dd.MissingTaskHandler(ctx, &def)
		if err != nil {
			return nil, err
		} else if mptr == nil {
			return nil, nil
		}
		metadata = *mptr
	}
	tc.TaskID = metadata.ID

	entrypoint, err := def.Entrypoint()
	if err == definitions.ErrNoEntrypoint {
		return &tc, nil
	} else if err != nil {
		return nil, err
	}

	defnDir := filepath.Dir(dir.DefinitionPath())
	absEntrypoint, err := filepath.Abs(filepath.Join(defnDir, entrypoint))
	if err != nil {
		return nil, err
	}
	tc.TaskEntrypoint = absEntrypoint

	kind, err := def.Kind()
	if err != nil {
		return nil, err
	}

	r, err := runtime.Lookup(entrypoint, kind)
	if err != nil {
		return nil, err
	}

	taskroot, err := r.Root(absEntrypoint)
	if err != nil {
		return nil, err
	}
	tc.TaskRoot = taskroot

	wd, err := r.Workdir(absEntrypoint)
	if err != nil {
		return nil, err
	}
	if err := def.SetWorkdir(taskroot, wd); err != nil {
		return nil, err
	}

	// Entrypoint for builder needs to be relative to taskroot, not definition directory.
	if defnDir != taskroot {
		ep, err := filepath.Rel(taskroot, absEntrypoint)
		if err != nil {
			return nil, err
		}
		def.SetBuildConfig("entrypoint", ep)
	}

	return &tc, nil
}

func (dd *DefnDiscoverer) TaskConfigSource() TaskConfigSource {
	return TaskConfigSourceDefn
}
