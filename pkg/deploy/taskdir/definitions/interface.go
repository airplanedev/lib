package definitions

import (
	"context"

	"github.com/airplanedev/lib/pkg/api"
	"github.com/airplanedev/lib/pkg/build"
	"github.com/pkg/errors"
)

type DefinitionInterface interface {
	// GetBuildConfig gets the full build config, synthesized from KindOptions and explicitly set
	// BuildConfig. KindOptions are unioned with BuildConfig; non-nil values in BuildConfig take
	// precedence, and a nil BuildConfig value removes the key from the final build config.
	GetBuildConfig() (build.BuildConfig, error)

	// SetBuildConfig sets a build config option. A value of nil means that the key will be
	// excluded from GetBuildConfig; used to mask values that exist in KindOptions.
	SetBuildConfig(key string, value interface{})

	// SetAbsoluteEntrypoint sets the absolute entrypoint for this definition. Does not change the
	// result of calling Entrypoint(). Returns ErrNoEntrypoint if the task kind definition requires
	// no entrypoint.
	SetAbsoluteEntrypoint(string) error

	// GetAbsoluteEntrypoint gets the absolute entrypoint for this definition. Returns
	// ErrNoEntrypoint if the task kind definition requires no entrypoint. If SetAbsoluteEntrypoint
	// has not been set, returns ErrNoAbsoluteEntrypoint.
	GetAbsoluteEntrypoint() (string, error)

	GetKindAndOptions() (build.TaskKind, build.KindOptions, error)
	GetEnv() (api.TaskEnv, error)
	GetConfigAttachments() ([]api.ConfigAttachment, error)
	GetSlug() string
	GetName() string
	UpgradeJST() error
	GetUpdateTaskRequest(ctx context.Context, client api.IAPIClient) (api.UpdateTaskRequest, error)
	SetWorkdir(taskroot, workdir string) error

	// Entrypoint returns ErrNoEntrypoint if the task kind definition requires no entrypoint. May be
	// empty. May be absolute or relative; if relative, it is relative to the defn file.
	Entrypoint() (string, error)
}

var ErrNoEntrypoint = errors.New("No entrypoint")
var ErrNoAbsoluteEntrypoint = errors.New("No absolute entrypoint")
