package discover

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/airplanedev/lib/pkg/api"
	"github.com/airplanedev/lib/pkg/deploy/taskdir/definitions"
	"github.com/airplanedev/lib/pkg/utils/logger"
	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

type AppDefinition struct {
	Slug       string `json:"slug"`
	Entrypoint string `json:"entrypoint"`
}

type AppDefnDiscoverer struct {
	Client api.IAPIClient
	Logger logger.Logger
}

var _ AppDiscoverer = &AppDefnDiscoverer{}

func (dd *AppDefnDiscoverer) IsAirplaneApp(ctx context.Context, file string) (*AppConfig, error) {
	fmt.Println("hi")

	if !definitions.IsAppDef(file) {
		return nil, nil
	}
	fmt.Println("hi")

	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrap(err, "reading app definition")
	}

	format := definitions.GetAppDefFormat(file)
	switch format {
	case definitions.DefFormatYAML:
		buf, err = yaml.YAMLToJSON(buf)
		if err != nil {
			return nil, err
		}
	case definitions.DefFormatJSON:
		// nothing
	default:
		return nil, errors.Errorf("unknown format: %s", format)
	}

	var d AppDefinition
	if err = json.Unmarshal(buf, &d); err != nil {
		return nil, err
	}

	return &AppConfig{
		Slug:       d.Slug,
		Entrypoint: d.Entrypoint,
		Source:     dd.AppConfigSource(),
	}, nil
}

func (dd *AppDefnDiscoverer) AppConfigSource() AppConfigSource {
	return AppConfigSourceDefn
}
