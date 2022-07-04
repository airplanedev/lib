package typescript

import (
	"bytes"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/airplanedev/lib/pkg/api"
	"github.com/airplanedev/lib/pkg/runtime"
	"github.com/airplanedev/lib/pkg/runtime/javascript"
)

// Init register the runtime.
func init() {
	runtime.Register(".ts", Runtime{})
}

// Data represents the data template.
type data struct {
	Comment string
	Params  string
}

// Param represents the parameter.
type param struct {
	Name string
	Type string
}

// Runtime implementaton.
type Runtime struct {
	javascript.Runtime
}

// Generate implementation.
func (r Runtime) Generate(t *runtime.Task, opts runtime.GenerateOpts) ([]byte, fs.FileMode, error) {
	d := data{}
	if t != nil {
		if t.URL != "" && opts.GenerateComment {
			d.Comment = runtime.Comment(r, t.URL)
		}
		var params map[string]api.Type
		for _, p := range t.Parameters {
			params[p.Slug] = p.Type

		}
		typescriptType, err := CreateParamsType(params, "")
		if err != nil {
			return nil, 0, err
		}
		d.Params = typescriptType
	}

	var buf bytes.Buffer
	if err := code.Execute(&buf, d); err != nil {
		return nil, 0, fmt.Errorf("typescript: template execute - %w", err)
	}

	return buf.Bytes(), 0644, nil
}

// CreateParamsType returns a string representation of a TypeScript type.
func CreateParamsType(parameters map[string]api.Type, typePrefix string) (string, error) {
	var params []string
	for slug, paramType := range parameters {
		slug = slugToTypeKey(slug)
		paramType := typeof(paramType)
		params = append(params, fmt.Sprintf("%s: %s;", slug, paramType))
	}

	sort.Strings(params)
	tc := paramsTemplateConfig{
		TaskName:   strings.Title(strings.ReplaceAll(typePrefix, " ", "")),
		TaskParams: strings.Join(params, "\n  "),
	}

	t := paramsTemplate
	if len(params) == 0 {
		t = paramTypesTemplateNoParams
	}

	var buff bytes.Buffer
	err := t.Execute(&buff, tc)
	if err != nil {
		return "", err
	}
	return buff.String(), nil
}

// typeof translates the given type to TypeScript.
func typeof(t api.Type) string {
	switch t {
	case api.TypeBoolean:
		return "boolean"
	case api.TypeString:
		return "string"
	case api.TypeInteger, api.TypeFloat:
		return "number"
	case api.TypeDate, api.TypeDatetime:
		return "string"
	case api.TypeUpload:
		return fileType
	case api.TypeConfigVar:
		return configType
	}
	return "unknown"
}

// slugToTypeKey converts a slug to a key that can be used in a TypeScript type.
func slugToTypeKey(slug string) string {
	if strings.Contains(slug, "-") {
		return fmt.Sprintf("\"%s\"", slug)
	}
	return slug
}
