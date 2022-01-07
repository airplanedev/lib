package definitions

import (
	"context"

	"github.com/airplanedev/lib/pkg/api"
	"github.com/airplanedev/lib/pkg/build"
)

type DefinitionInterface interface {
	GetKindAndOptions() (build.TaskKind, build.KindOptions, error)
	GetEnv() (api.TaskEnv, error)
	GetSlug() string
	UpgradeJST() error
	GetUpdateTaskRequest(ctx context.Context, client api.IAPIClient, currentTask *api.Task) (api.UpdateTaskRequest, error)
}
