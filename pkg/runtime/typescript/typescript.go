package typescript

import (
	"bytes"
	"fmt"
	"io/fs"
	"text/template"

	"github.com/airplanedev/lib/pkg/build"
	"github.com/airplanedev/lib/pkg/runtime"
	"github.com/airplanedev/lib/pkg/runtime/javascript"
)

// Init register the runtime.
func init() {
	runtime.Register(".ts", Runtime{})
}

// Code template.
var code = template.Must(template.New("ts").Parse(`{{.Comment}}

type Params = {
  {{- range .Params }}
  {{ .Name }}: {{ .Type }}
  {{- end }}
}

export default async function(params: Params) {
  console.log('parameters:', params);
}
`))

// Data represents the data template.
type data struct {
	Comment string
	Params  []param
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
func (r Runtime) Generate(t build.Task) ([]byte, fs.FileMode, error) {
	var args = data{Comment: runtime.Comment(r, t)}
	var params = t.Parameters
	var buf bytes.Buffer

	for _, p := range params {
		args.Params = append(args.Params, param{
			Name: p.Slug,
			Type: typeof(p.Type),
		})
	}

	if err := code.Execute(&buf, args); err != nil {
		return nil, 0, fmt.Errorf("typescript: template execute - %w", err)
	}

	return buf.Bytes(), 0644, nil
}

// Typeof translates the given type to typescript.
func typeof(t build.Type) string {
	switch t {
	case build.TypeInteger, TypeFloat:
		return "number"
	case build.TypeDate, TypeDatetime:
		return "string"
	case build.TypeBoolean:
		return "boolean"
	case build.TypeString:
		return "string"
	case build.TypeUpload:
		return "string"
	default:
		return "unknown"
	}
}
