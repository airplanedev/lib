package discover

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/airplanedev/lib/pkg/build"
	"github.com/airplanedev/lib/pkg/common/api"
	"github.com/airplanedev/lib/pkg/common/logger"
	"github.com/airplanedev/lib/pkg/deploy/taskdir/definitions"
	_ "github.com/airplanedev/lib/pkg/runtime/javascript"
	_ "github.com/airplanedev/lib/pkg/runtime/python"
	_ "github.com/airplanedev/lib/pkg/runtime/shell"
	_ "github.com/airplanedev/lib/pkg/runtime/typescript"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscoverTasks(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	fixturesPath, _ := filepath.Abs("./fixtures")
	discoverPath, _ := filepath.Abs("./")
	tests := []struct {
		name          string
		paths         []string
		existingTasks map[string]api.Task
		expectedErr   bool
		want          []TaskConfig
	}{
		{
			name:  "single script",
			paths: []string{"./fixtures/single_task.js"},
			existingTasks: map[string]api.Task{
				"my_task": {Kind: build.TaskKindNode},
			},
			want: []TaskConfig{
				{
					TaskRoot:         fixturesPath,
					WorkingDirectory: fixturesPath,
					TaskFilePath:     fixturesPath + "/single_task.js",
					Kind:             build.TaskKindNode,
					Def: &definitions.Definition{
						Node: &definitions.NodeDefinition{Entrypoint: "single_task.js"},
					},
					Task: api.Task{Kind: build.TaskKindNode},
					KindOptions: build.KindOptions{
						"entrypoint":  "single_task.js",
						"language":    "",
						"nodeVersion": "",
						"workdir":     "",
					},
					From: TaskConfigSourceScript,
				},
			},
		},
		{
			name:  "multiple scripts",
			paths: []string{"./fixtures/single_task.js", "./fixtures/single_task2.js"},
			existingTasks: map[string]api.Task{
				"my_task":  {Kind: build.TaskKindNode},
				"my_task2": {Kind: build.TaskKindNode},
			},
			want: []TaskConfig{
				{
					TaskRoot:         fixturesPath,
					WorkingDirectory: fixturesPath,
					TaskFilePath:     fixturesPath + "/single_task.js",
					Kind:             build.TaskKindNode,
					Def: &definitions.Definition{
						Node: &definitions.NodeDefinition{Entrypoint: "single_task.js"},
					},
					Task: api.Task{Kind: build.TaskKindNode},
					KindOptions: build.KindOptions{
						"entrypoint":  "single_task.js",
						"language":    "",
						"nodeVersion": "",
						"workdir":     "",
					},
					From: TaskConfigSourceScript,
				},
				{
					TaskRoot:         fixturesPath,
					WorkingDirectory: fixturesPath,
					TaskFilePath:     fixturesPath + "/single_task2.js",
					Kind:             build.TaskKindNode,
					Def: &definitions.Definition{
						Node: &definitions.NodeDefinition{Entrypoint: "single_task2.js"},
					},
					Task: api.Task{Kind: build.TaskKindNode},
					KindOptions: build.KindOptions{
						"entrypoint":  "single_task2.js",
						"language":    "",
						"nodeVersion": "",
						"workdir":     "",
					},
					From: TaskConfigSourceScript,
				},
			},
		},
		{
			name:  "nested scripts",
			paths: []string{"./fixtures/nestedScripts"},
			existingTasks: map[string]api.Task{
				"my_task":  {Kind: build.TaskKindNode},
				"my_task2": {Kind: build.TaskKindNode},
			},
			want: []TaskConfig{
				{
					TaskRoot:         fixturesPath + "/nestedScripts",
					WorkingDirectory: fixturesPath + "/nestedScripts",
					TaskFilePath:     fixturesPath + "/nestedScripts/single_task.js",
					Kind:             build.TaskKindNode,
					Def: &definitions.Definition{
						Node: &definitions.NodeDefinition{Entrypoint: "single_task.js"},
					},
					Task: api.Task{Kind: build.TaskKindNode},
					KindOptions: build.KindOptions{
						"entrypoint":  "single_task.js",
						"language":    "",
						"nodeVersion": "",
						"workdir":     "",
					},
					From: TaskConfigSourceScript,
				},
				{
					TaskRoot:         fixturesPath + "/nestedScripts",
					WorkingDirectory: fixturesPath + "/nestedScripts",
					TaskFilePath:     fixturesPath + "/nestedScripts/single_task2.js",
					Kind:             build.TaskKindNode,
					Def: &definitions.Definition{
						Node: &definitions.NodeDefinition{Entrypoint: "single_task2.js"},
					},
					Task: api.Task{Kind: build.TaskKindNode},
					KindOptions: build.KindOptions{
						"entrypoint":  "single_task2.js",
						"language":    "",
						"nodeVersion": "",
						"workdir":     "",
					},
					From: TaskConfigSourceScript,
				},
			},
		},
		{
			name:  "single defn",
			paths: []string{"./fixtures/defn.task.yaml"},
			existingTasks: map[string]api.Task{
				"my_task": {Kind: build.TaskKindNode},
			},
			want: []TaskConfig{
				{
					TaskRoot: fixturesPath,
					// TODO adjust to be fixturesPath when entrypoint is relative to task defn
					TaskFilePath: discoverPath + "/single_task.js",
					Kind:         build.TaskKindNode,
					Def: &definitions.Definition_0_3{
						Name:        "sunt in tempor eu",
						Slug:        "my_task",
						Description: "ut dolor sit officia ea",
						Node:        &definitions.NodeDefinition_0_3{Entrypoint: "./single_task.js", NodeVersion: "14"},
					},
					Task: api.Task{Kind: build.TaskKindNode},
					KindOptions: build.KindOptions{
						// TODO adjust to be absolute path
						"entrypoint":  "./single_task.js",
						"nodeVersion": "14",
					},
					From: TaskConfigSourceDefn,
				},
			},
		},
		{
			name:          "task not returned by api - deploy skipped",
			paths:         []string{"./fixtures/single_task.js", "./fixtures/defn.task.yaml"},
			existingTasks: map[string]api.Task{},
			expectedErr:   false,
		},
	}
	for _, tC := range tests {
		t.Run(tC.name, func(t *testing.T) {
			apiClient := &api.MockClient{
				Tasks: tC.existingTasks,
			}
			scriptDiscoverer := &ScriptDiscoverer{}
			defnDiscoverer := &DefnDiscoverer{
				Client: apiClient,
			}
			d := &Discoverer{
				TaskDiscoverers: []TaskDiscoverer{scriptDiscoverer, defnDiscoverer},
				Client: &api.MockClient{
					Tasks: tC.existingTasks,
				},
				Logger: &logger.MockLogger{},
			}
			got, err := d.DiscoverTasks(context.Background(), tC.paths...)
			if tC.expectedErr {
				assert.NotNil(err)
				return
			}
			require.NoError(err)

			assert.Equal(tC.want, got)
		})
	}
}
