package discover

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/airplanedev/lib/pkg/api"
	"github.com/airplanedev/lib/pkg/deploy/taskdir/definitions"
	"github.com/airplanedev/lib/pkg/runtime/typescript"
	"github.com/airplanedev/lib/pkg/utils/logger"
	"github.com/pkg/errors"
)

type ViewDefnDiscoverer struct {
	Client api.IAPIClient
	Logger logger.Logger
}

var _ ViewDiscoverer = &ViewDefnDiscoverer{}

func (dd *ViewDefnDiscoverer) GetViewConfig(ctx context.Context, file string) (*ViewConfig, error) {
	if !definitions.IsViewDef(file) {
		return nil, nil
	}

	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrap(err, "reading view definition")
	}

	format := definitions.GetViewDefFormat(file)
	d := definitions.ViewDefinition{}

	if err = d.Unmarshal(format, buf); err != nil {
		switch err := errors.Cause(err).(type) {
		case definitions.ErrSchemaValidation:
			errorMsgs := []string{}
			for _, verr := range err.Errors {
				errorMsgs = append(errorMsgs, fmt.Sprintf("%s: %s", verr.Field(), verr.Description()))
			}
			return nil, definitions.NewErrReadDefinition(fmt.Sprintf("Error reading %s", file), errorMsgs...)
		default:
			return nil, errors.Wrap(err, "unmarshalling view definition")
		}
	}

	typeScriptRuntime := typescript.Runtime{}
	absFile, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}
	root, err := typeScriptRuntime.Root(absFile)
	if err != nil {
		return nil, errors.Wrap(err, "getting absolute view definition root")
	}

	view, err := dd.Client.GetApp(ctx, api.GetAppRequest{Slug: d.Slug})
	if err != nil {
		var merr *api.AppMissingError
		if !errors.As(err, &merr) {
			return nil, errors.Wrap(err, "unable to get view")
		}
		// TODO offer to create the view.
		if dd.Logger != nil {
			dd.Logger.Warning(`View with slug %s does not exist, skipping deploy.`, d.Slug)
		}
		return nil, nil
	}
	if view.ArchivedAt != nil {
		dd.Logger.Warning(`View with slug %s is archived, skipping deployment.`, view.Slug)
		return nil, nil
	}

	if !filepath.IsAbs(d.Entrypoint) {
		defnDir := filepath.Dir(file)
		d.Entrypoint, err = filepath.Abs(filepath.Join(defnDir, d.Entrypoint))
		if err != nil {
			return nil, err
		}
	}

	return &ViewConfig{
		ID:     view.ID,
		Def:    d,
		Source: dd.ConfigSource(),
		Root:   root,
	}, nil
}

func (dd *ViewDefnDiscoverer) ConfigSource() ConfigSource {
	return ConfigSourceDefn
}
