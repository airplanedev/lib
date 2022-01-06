package discover

import (
	"context"
	"path/filepath"

	"github.com/airplanedev/lib/pkg/deploy/api"
	"github.com/airplanedev/lib/pkg/deploy/taskdir"
	"github.com/airplanedev/lib/pkg/deploy/definitions"
	"github.com/airplanedev/lib/pkg/utils/logger"
)

type DefnDiscoverer struct {
	Client             api.APIClient
	AssumeYes          bool
	AssumeNo           bool
	Logger             logger.Logger
	MissingTaskHandler func(context.Context, definitions.DefinitionInterface) (*api.Task, error)
}

var _ TaskDiscoverer = &DefnDiscoverer{}

func (dd *DefnDiscoverer) IsAirplaneTask(ctx context.Context, file string) (slug string, err error) {
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

	taskFilePath := ""
	entrypoint, err := def.Entrypoint()
	if err == definitions.ErrNoEntrypoint {
		// nothing
	} else if err != nil {
		return TaskConfig{}, err
	} else {
		taskFilePath, err = filepath.Abs(entrypoint)
		if err != nil {
			return TaskConfig{}, err
		}
	}

	return TaskConfig{
		TaskRoot:       dir.DefinitionRootPath(),
		TaskEntrypoint: taskFilePath,
		Task:           task,
		Def:            &def,
	}, nil
}

func (dd *DefnDiscoverer) TaskConfigSource() TaskConfigSource {
	return TaskConfigSourceDefn
}

func (dd *DefnDiscoverer) HandleMissingTask(ctx context.Context, file string) (*api.Task, error) {
	def, err := getDef(file)
	if err != nil {
		return nil, err
	}
	if dd.MissingTaskHandler != nil {
		return dd.MissingTaskHandler(ctx, def)
	}
	return nil, nil
}

func getDef(file string) (definitions.DefinitionInterface, error) {
	dir, err := taskdir.Open(file, true)
	if err != nil {
		return &definitions.Definition_0_3{}, nil
	}
	defer dir.Close()

	def, err := dir.ReadDefinition_0_3()
	if err != nil {
		return &definitions.Definition_0_3{}, err
	}

	return &def, nil
}
