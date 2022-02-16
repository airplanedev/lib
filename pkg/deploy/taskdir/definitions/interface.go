package definitions

import (
	"context"

	"github.com/airplanedev/lib/pkg/api"
	"github.com/airplanedev/lib/pkg/build"
)

type DefinitionInterface interface {
	// GetBuildConfig gets the full build config, synthesized from KindOptions and explicitly set
	// BuildConfig. KindOptions are unioned with BuildConfig; non-nil values in BuildConfig take
	// precedence, and a nil BuildConfig value removes the key from the final build config.
	GetBuildConfig() (build.BuildConfig, error)

	// SetBuildConfig sets a build config option. A value of nil means that the key will be
	// excluded from GetBuildConfig; used to mask values that exist in KindOptions.
	SetBuildConfig(key string, value interface{})

	GetKindAndOptions() (build.TaskKind, build.KindOptions, error)
	GetEnv() (api.TaskEnv, error)
	GetSlug() string
	UpgradeJST() error
	GetUpdateTaskRequest(ctx context.Context, client api.IAPIClient, currentTask *api.Task) (api.UpdateTaskRequest, error)
	SetWorkdir(taskroot, workdir string) error
	Entrypoint() (string, error)
}
