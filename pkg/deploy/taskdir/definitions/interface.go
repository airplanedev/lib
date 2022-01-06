package definitions

import (
	"context"

	"github.com/airplanedev/lib/pkg/build"
	"github.com/airplanedev/lib/pkg/deploy/api"
)

type DefinitionInterface interface {
	GetKindAndOptions() (build.TaskKind, build.KindOptions, error)
	GetEnv() (api.TaskEnv, error)
	GetSlug() string
	UpgradeJST() error
	GetUpdateTaskRequest(ctx context.Context, client api.APIClient, currentTask *api.Task) (api.UpdateTaskRequest, error)
}
