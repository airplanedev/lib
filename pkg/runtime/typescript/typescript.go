package typescript

import (
	"bytes"
	"fmt"
	"io/fs"
	"text/template"

	"github.com/airplanedev/lib/pkg/runtime"
	"github.com/airplanedev/lib/pkg/runtime/javascript"
)

// Init register the runtime.
func init() {
	runtime.Register(".ts", Runtime{})
}

// Code template.
var code = template.Must(template.New("ts").Parse(`{{with .Comment -}}
{{.}}

{{end -}}
type Params = {
  {{- range .Params }}
  {{ .NameTag }}: {{ .Type }}
  {{- end }}
}

// Put the main logic of the task in this function.
export default async function(params: Params) {
  console.log('parameters:', params);

  // You can return data to show outputs to users.
  // Outputs documentation: https://docs.airplane.dev/tasks/outputs
  return [
    {element: 'hydrogen', weight: 1.008},
    {element: 'helium', weight: 4.0026},
  ];
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
func (r Runtime) Generate(t *runtime.Task) ([]byte, fs.FileMode, error) {
	d := data{}
	if t != nil {
		d.Comment = runtime.Comment(r, t.URL)
		for _, p := range t.Parameters {
			d.Params = append(d.Params, param{
				Name: p.Slug,
				Type: typeof(p.Type),
			})
		}
	}

	var buf bytes.Buffer
	if err := code.Execute(&buf, d); err != nil {
		return nil, 0, fmt.Errorf("typescript: template execute - %w", err)
	}

	return buf.Bytes(), 0644, nil
}

// Typeof translates the given type to typescript.
func typeof(t runtime.Type) string {
	switch t {
	case runtime.TypeInteger, runtime.TypeFloat:
		return "number"
	case runtime.TypeDate, runtime.TypeDatetime:
		return "string"
	case runtime.TypeBoolean:
		return "boolean"
	case runtime.TypeString:
		return "string"
	case runtime.TypeUpload:
		return "string"
	default:
		return "unknown"
	}
}
