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

	// TODO: set tc.TaskID from the slug and handle missing tasks

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
