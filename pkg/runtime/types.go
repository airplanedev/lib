package runtime

import "github.com/airplanedev/lib/pkg/api"

// Task represents a task.
type Task struct {
	URL        string
	Parameters Parameters
}

type Parameters []Parameter

// Parameter represents a task parameter.
type Parameter struct {
	Name string
	Slug string
	Type api.Type
}

// Values represent parameters values.
//
// An alias is used because we want the type
// to be `map[string]interface{}` and not a custom one.
//
// They're keyed by the parameter "slug".
type Values = map[string]interface{}
