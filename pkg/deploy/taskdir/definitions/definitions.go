package definitions

import (
	"context"
	"regexp"
	"strings"

	"github.com/airplanedev/lib/pkg/api"
	"gopkg.in/yaml.v3"
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

func tryOlderDefinitions(buf []byte) (DefinitionInterface, error) {
	var err error
	if err = validateYAML(buf, Definition_0_1{}); err == nil {
		var def Definition_0_1
		if e := yaml.Unmarshal(buf, &def); e != nil {
			return nil, err
		}
		return def.upgrade()
	}
	return nil, err
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
