package build

import (
	"context"

	"github.com/airplanedev/lib/pkg/build"
	"github.com/airplanedev/lib/pkg/common/api"
	"github.com/airplanedev/lib/pkg/deploy/taskdir/definitions"
)

type BuildCreator interface {
	CreateBuild(ctx context.Context, req Request) (*build.Response, error)
}

// Request represents a build request.
type Request struct {
	Root    string
	Def     definitions.DefinitionInterface
	TaskID  string
	TaskEnv api.TaskEnv
	Shim    bool
	GitMeta api.BuildGitMeta
}

// Response represents a build response.
type Response struct {
	ImageURL string
	// Optional, only if applicable
	BuildID string
}
