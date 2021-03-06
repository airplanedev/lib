package build

// KindOptions are part of the task definition, supplied by the user.
type KindOptions map[string]interface{}

// BuildConfig is a collection of build-specific configuration options based on a task's
// KindOptions.
type BuildConfig map[string]interface{}

type TaskKind string

const (
	TaskKindImage  TaskKind = "image"
	TaskKindNode   TaskKind = "node"
	TaskKindPython TaskKind = "python"
	TaskKindShell  TaskKind = "shell"
	TaskKindApp    TaskKind = "app"

	TaskKindSQL  TaskKind = "sql"
	TaskKindREST TaskKind = "rest"
)

type TaskRuntime string

const (
	TaskRuntimeStandard TaskRuntime = ""
	TaskRuntimeWorkflow TaskRuntime = "workflow"
)

// Value represents a value.
type Value interface{}

// Values represent parameters values.
//
// An alias is used because we want the type
// to be `map[string]interface{}` and not a custom one.
//
// They're keyed by the parameter "slug".
type Values = map[string]interface{}
