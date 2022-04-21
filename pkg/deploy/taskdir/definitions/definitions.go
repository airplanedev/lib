package definitions

import (
	"context"
	"strings"

	"github.com/airplanedev/lib/pkg/api"
)

var (
	YamlDefExtensions = []string{".task.yaml", ".task.yml"}
	JSONDefExtensions = []string{".task.json"}
	TaskDefExtensions = append(YamlDefExtensions, JSONDefExtensions...)
)

func NewDefinitionFromTask(ctx context.Context, client api.IAPIClient, t api.Task) (DefinitionInterface, error) {
	def, err := NewDefinitionFromTask_0_3(ctx, client, t)
	if err != nil {
		return nil, err
	}
	return &def, nil
}

type TaskDefFormat string

const (
	TaskDefFormatUnknown TaskDefFormat = ""
	TaskDefFormatYAML    TaskDefFormat = "yaml"
	TaskDefFormatJSON    TaskDefFormat = "json"
)

func IsTaskDef(fn string) bool {
	return GetTaskDefFormat(fn) != TaskDefFormatUnknown
}

func GetTaskDefFormat(fn string) TaskDefFormat {
	for _, de := range YamlDefExtensions {
		if strings.HasSuffix(fn, de) {
			return TaskDefFormatYAML
		}
	}
	for _, de := range JSONDefExtensions {
		if strings.HasSuffix(fn, de) {
			return TaskDefFormatJSON
		}
	}
	return TaskDefFormatUnknown
}
