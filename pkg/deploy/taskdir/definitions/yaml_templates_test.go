package definitions

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/airplanedev/lib/pkg/build"
	"github.com/stretchr/testify/require"
)

func TestYAMLComments(t *testing.T) {
	fixturesPath, _ := filepath.Abs("./fixtures")
	for _, test := range []struct {
		name       string
		file       string
		kind       build.TaskKind
		entrypoint string
	}{
		{
			name:       "python",
			file:       fixturesPath + "/python.task.yaml",
			kind:       build.TaskKindPython,
			entrypoint: "my_task.py",
		},
		{
			name:       "node",
			file:       fixturesPath + "/node.task.yaml",
			kind:       build.TaskKindNode,
			entrypoint: "my_task.ts",
		},
		{
			name:       "shell",
			file:       fixturesPath + "/shell.task.yaml",
			kind:       build.TaskKindShell,
			entrypoint: "my_task.sh",
		},
		{
			name: "docker",
			file: fixturesPath + "/docker.task.yaml",
			kind: build.TaskKindImage,
		},
		{
			name:       "sql",
			file:       fixturesPath + "/sql.task.yaml",
			kind:       build.TaskKindSQL,
			entrypoint: "my_task.sql",
		},
		{
			name: "rest",
			file: fixturesPath + "/rest.task.yaml",
			kind: build.TaskKindREST,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)
			def, err := NewDefinition_0_3("My Task", "my_task", test.kind, test.entrypoint)
			require.NoError(err)

			got, err := def.GenerateCommentedFile(TaskDefFormatYAML)
			require.NoError(err)

			expected, err := os.ReadFile(test.file)
			require.NoError(err)

			require.Equal(string(expected), string(got))

			unmarshalled := Definition_0_3{}
			err = unmarshalled.Unmarshal(TaskDefFormatYAML, got)
			require.NoError(err)
			require.Equal(def, unmarshalled)
		})
	}
}
