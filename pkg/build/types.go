package build

type KindOptions map[string]interface{}

type TaskKind string

const (
	TaskKindDeno       TaskKind = "deno"
	TaskKindDockerfile TaskKind = "dockerfile"
	TaskKindGo         TaskKind = "go"
	TaskKindImage      TaskKind = "image"
	TaskKindNode       TaskKind = "node"
	TaskKindPython     TaskKind = "python"
	TaskKindShell      TaskKind = "shell"

	TaskKindSQL  TaskKind = "sql"
	TaskKindREST TaskKind = "rest"
)

// Type enumerates parameter types.
type Type string

// Value represents a value.
type Value interface{}

// Values represent parameters values.
//
// An alias is used because we want the type
// to be `map[string]interface{}` and not a custom one.
//
// They're keyed by the parameter "slug".
type Values = map[string]interface{}
