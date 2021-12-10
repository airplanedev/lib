package definitions

import (
	"context"

	"github.com/airplanedev/lib/pkg/build"
	"github.com/airplanedev/lib/pkg/common/api"
)

type DefinitionInterface interface {
	GetKindAndOptions() (build.TaskKind, build.KindOptions, error)
	GetEnv() (api.TaskEnv, error)
	GetSlug() string
	UpgradeJST() error
	GetUpdateTaskRequest(context.Context, api.APIClient, *string) (api.UpdateTaskRequest, error)
}
