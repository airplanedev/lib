package definitions

import (
	"context"
	"strings"

	"github.com/airplanedev/lib/pkg/api"
)

var (
	YamlTaskDefExtensions = []string{".task.yaml", ".task.yml"}
	JSONTaskDefExtensions = []string{".task.json"}
	TaskDefExtensions     = append(YamlTaskDefExtensions, JSONTaskDefExtensions...)

	YamlAppDefExtensions = []string{".app.yaml", ".app.yml"}
	JSONAppDefExtensions = []string{".app.json"}
	AppDefExtensions     = append(YamlAppDefExtensions, JSONAppDefExtensions...)
)

func NewDefinitionFromTask(ctx context.Context, client api.IAPIClient, t api.Task) (DefinitionInterface, error) {
	def, err := NewDefinitionFromTask_0_3(ctx, client, t)
	if err != nil {
		return nil, err
	}
	return &def, nil
}

type DefFormat string

const (
	DefFormatUnknown DefFormat = ""
	DefFormatYAML    DefFormat = "yaml"
	DefFormatJSON    DefFormat = "json"
)

func IsTaskDef(fn string) bool {
	return GetTaskDefFormat(fn) != DefFormatUnknown
}

func IsAppDef(fn string) bool {
	return GetAppDefFormat(fn) != DefFormatUnknown
}

func GetDefFormat(fn string, yamlExtensions, jsonExtensions []string) DefFormat {
	for _, de := range yamlExtensions {
		if strings.HasSuffix(fn, de) {
			return DefFormatYAML
		}
	}
	for _, de := range jsonExtensions {
		if strings.HasSuffix(fn, de) {
			return DefFormatJSON
		}
	}
	return DefFormatUnknown
}

func GetAppDefFormat(fn string) DefFormat {
	return GetDefFormat(fn, YamlAppDefExtensions, JSONAppDefExtensions)
}

func GetTaskDefFormat(fn string) DefFormat {
	return GetDefFormat(fn, YamlTaskDefExtensions, JSONTaskDefExtensions)
}
