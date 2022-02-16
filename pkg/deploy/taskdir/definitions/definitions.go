package definitions

import (
	"context"
	"regexp"
	"strings"

	"github.com/airplanedev/lib/pkg/api"
)

var jsonRegex = regexp.MustCompile(`{{ *JSON *}}`)

func upgradeArguments(args []string) []string {
	upgraded := make([]string, len(args))
	for i, arg := range args {
		jstArg := jsonRegex.ReplaceAllString(arg, "{{JSON.stringify(params)}}")
		upgraded[i] = jstArg
	}
	return upgraded
}

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
	if strings.HasSuffix(fn, ".task.yaml") || strings.HasSuffix(fn, ".task.yml") {
		return TaskDefFormatYAML
	}
	if strings.HasSuffix(fn, ".task.json") {
		return TaskDefFormatJSON
	}
	return TaskDefFormatUnknown
}
