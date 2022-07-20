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
	"github.com/pkg/errors"
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
		params := make(map[string]api.Type, len(t.Parameters))
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
		paramType, err := typeof(paramType)
		if err != nil {
			return "", err
		}
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
		return "", errors.Wrap(err, "applying parameter template")
	}
	return buff.String(), nil
}

// typeof translates the given type to TypeScript.
func typeof(t api.Type) (string, error) {
	switch t {
	case api.TypeBoolean:
		return "boolean", nil
	case api.TypeString:
		return "string", nil
	case api.TypeInteger, api.TypeFloat:
		return "number", nil
	case api.TypeDate, api.TypeDatetime:
		return "string", nil
	case api.TypeUpload:
		return fileType, nil
	case api.TypeConfigVar:
		return configType, nil
	}
	return "", errors.Errorf("unknown parameter type %s", t)
}
