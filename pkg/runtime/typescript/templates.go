package typescript

import "text/template"

var code = template.Must(template.New("ts").Parse(`{{with .Comment -}}
{{.}}

{{end -}}
{{.Params}}

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

type paramsTemplateConfig struct {
	TaskName   string
	TaskParams string
}

var paramsTemplate = template.Must(template.New("tsParams").Parse(`export type {{.TaskName}}Params = {
  {{.TaskParams}}
};
`))

var paramTypesTemplateNoParams = template.Must(template.New("tsParamsEmpty").Parse(`export type {{.TaskName}}Params = {};
`))

var fileType = `{ __airplaneType: "upload"; id: string; url: string }`

var configType = `{ name: string; value: string }`
