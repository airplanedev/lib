package viewdir

import (
	"context"
	"path/filepath"

	"github.com/airplanedev/lib/pkg/api"
	"github.com/airplanedev/lib/pkg/deploy/discover"
	"github.com/airplanedev/lib/pkg/deploy/taskdir/definitions"
	"github.com/airplanedev/lib/pkg/utils/logger"
	"github.com/pkg/errors"
)

type ViewDirectoryInterface interface {
	Root() string
	EntrypointPath() string
}

type ViewDirectory struct {
	root           string
	entrypointPath string
	logger         logger.Logger
}

func (this *ViewDirectory) Root() string {
	return this.root
}

func (this *ViewDirectory) EntrypointPath() string {
	return this.entrypointPath
}

func missingViewHandler(ctx context.Context, defn definitions.ViewDefinition) (*api.View, error) {
	// TODO(zhan): generate view?
	return &api.View{
		ID:   "temp",
		Slug: defn.Slug,
	}, nil
}

func NewViewDirectory(
	ctx context.Context,
	api *api.Client,
	logger logger.Logger,
	root string,
	envSlug string,
) (ViewDirectory, error) {
	// Discover local views in the directory of the file.
	d := &discover.Discoverer{
		ViewDiscoverers: []discover.ViewDiscoverer{
			&discover.ViewDefnDiscoverer{
				Client:             api,
				MissingViewHandler: missingViewHandler,
			},
		},
		EnvSlug: envSlug,
		Client:  api,
	}
	_, viewConfigs, err := d.Discover(ctx, root)
	if err != nil {
		return ViewDirectory{}, errors.Wrap(err, "discovering view configs")
	}
	if len(viewConfigs) != 1 {
		return ViewDirectory{}, errors.New("currently can only have one view!")
	}
	vc := viewConfigs[0]
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return ViewDirectory{}, errors.Wrap(err, "getting absolute root filepath")
	}

	return ViewDirectory{
		root:           absRoot,
		entrypointPath: vc.Def.Entrypoint,
		logger:         logger,
	}, nil
}
