package runtime

// Task represents a task.
type Task struct {
	URL        string
	Parameters Parameters
}

// Type enumerates parameter types.
type Type string

// All Parameter types.
const (
	TypeString    Type = "string"
	TypeBoolean   Type = "boolean"
	TypeUpload    Type = "upload"
	TypeInteger   Type = "integer"
	TypeFloat     Type = "float"
	TypeDate      Type = "date"
	TypeDatetime  Type = "datetime"
	TypeConfigVar Type = "configvar"
)

type Parameters []Parameter

// Parameter represents a task parameter.
type Parameter struct {
	Name string
	Slug string
	Type Type
}

// Values represent parameters values.
//
// An alias is used because we want the type
// to be `map[string]interface{}` and not a custom one.
//
// They're keyed by the parameter "slug".
type Values = map[string]interface{}
