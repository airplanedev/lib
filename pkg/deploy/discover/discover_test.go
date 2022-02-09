package discover

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/airplanedev/lib/pkg/api"
	"github.com/airplanedev/lib/pkg/api/mock"
	"github.com/airplanedev/lib/pkg/build"
	"github.com/airplanedev/lib/pkg/deploy/taskdir/definitions"
	"github.com/airplanedev/lib/pkg/utils/logger"
	"github.com/stretchr/testify/require"
)

func TestDiscoverTasks(t *testing.T) {
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
				"my_task": {Slug: "my_task", Kind: build.TaskKindNode},
			},
			want: []TaskConfig{
				{
					TaskRoot:         fixturesPath,
					WorkingDirectory: fixturesPath,
					TaskEntrypoint:   fixturesPath + "/single_task.js",
					Def: &definitions.Definition{
						Slug: "my_task",
						Node: &definitions.NodeDefinition{Entrypoint: "single_task.js"},
					},
					Task: api.Task{Slug: "my_task", Kind: build.TaskKindNode},
					From: TaskConfigSourceScript,
				},
			},
		},
		{
			name:  "multiple scripts",
			paths: []string{"./fixtures/single_task.js", "./fixtures/single_task2.js"},
			existingTasks: map[string]api.Task{
				"my_task":  {Slug: "my_task", Kind: build.TaskKindNode},
				"my_task2": {Slug: "my_task2", Kind: build.TaskKindNode},
			},
			want: []TaskConfig{
				{
					TaskRoot:         fixturesPath,
					WorkingDirectory: fixturesPath,
					TaskEntrypoint:   fixturesPath + "/single_task.js",
					Def: &definitions.Definition{
						Slug: "my_task",
						Node: &definitions.NodeDefinition{Entrypoint: "single_task.js"},
					},
					Task: api.Task{Slug: "my_task", Kind: build.TaskKindNode},
					From: TaskConfigSourceScript,
				},
				{
					TaskRoot:         fixturesPath,
					WorkingDirectory: fixturesPath,
					TaskEntrypoint:   fixturesPath + "/single_task2.js",
					Def: &definitions.Definition{
						Slug: "my_task2",
						Node: &definitions.NodeDefinition{Entrypoint: "single_task2.js"},
					},
					Task: api.Task{Slug: "my_task2", Kind: build.TaskKindNode},
					From: TaskConfigSourceScript,
				},
			},
		},
		{
			name:  "nested scripts",
			paths: []string{"./fixtures/nestedScripts"},
			existingTasks: map[string]api.Task{
				"my_task":  {Slug: "my_task", Kind: build.TaskKindNode},
				"my_task2": {Slug: "my_task2", Kind: build.TaskKindNode},
			},
			want: []TaskConfig{
				{
					TaskRoot:         fixturesPath + "/nestedScripts",
					WorkingDirectory: fixturesPath + "/nestedScripts",
					TaskEntrypoint:   fixturesPath + "/nestedScripts/single_task.js",
					Def: &definitions.Definition{
						Slug: "my_task",
						Node: &definitions.NodeDefinition{Entrypoint: "single_task.js"},
					},
					Task: api.Task{Slug: "my_task", Kind: build.TaskKindNode},
					From: TaskConfigSourceScript,
				},
				{
					TaskRoot:         fixturesPath + "/nestedScripts",
					WorkingDirectory: fixturesPath + "/nestedScripts",
					TaskEntrypoint:   fixturesPath + "/nestedScripts/single_task2.js",
					Def: &definitions.Definition{
						Slug: "my_task2",
						Node: &definitions.NodeDefinition{Entrypoint: "single_task2.js"},
					},
					Task: api.Task{Slug: "my_task2", Kind: build.TaskKindNode},
					From: TaskConfigSourceScript,
				},
			},
		},
		{
			name:  "single defn",
			paths: []string{"./fixtures/defn.task.yaml"},
			existingTasks: map[string]api.Task{
				"my_task": {Slug: "my_task", Kind: build.TaskKindNode},
			},
			want: []TaskConfig{
				{
					TaskRoot: fixturesPath,
					// TODO adjust to be fixturesPath when entrypoint is relative to task defn
					TaskEntrypoint: discoverPath + "/single_task.js",
					Def: &definitions.Definition_0_3{
						Name:        "sunt in tempor eu",
						Slug:        "my_task",
						Description: "ut dolor sit officia ea",
						// TODO adjust to be absolute path
						Node: &definitions.NodeDefinition_0_3{Entrypoint: "./single_task.js", NodeVersion: "14"},
					},
					Task: api.Task{Slug: "my_task", Kind: build.TaskKindNode},
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
		{
			name:  "same task, multiple discoverers",
			paths: []string{"./fixtures/defn.task.yaml", "./fixtures/single_task.js"},
			existingTasks: map[string]api.Task{
				"my_task": {Slug: "my_task", Kind: build.TaskKindNode},
			},
			want: []TaskConfig{
				{
					TaskRoot: fixturesPath,
					// TODO adjust to be fixturesPath when entrypoint is relative to task defn
					TaskEntrypoint: discoverPath + "/single_task.js",
					Def: &definitions.Definition_0_3{
						Name:        "sunt in tempor eu",
						Slug:        "my_task",
						Description: "ut dolor sit officia ea",
						// TODO adjust to be absolute path
						Node: &definitions.NodeDefinition_0_3{Entrypoint: "./single_task.js", NodeVersion: "14"},
					},
					Task: api.Task{Slug: "my_task", Kind: build.TaskKindNode},
					From: TaskConfigSourceDefn,
				},
			},
		},
	}
	for _, tC := range tests {
		t.Run(tC.name, func(t *testing.T) {
			require := require.New(t)
			apiClient := &mock.MockClient{
				Tasks: tC.existingTasks,
			}
			scriptDiscoverer := &ScriptDiscoverer{}
			defnDiscoverer := &DefnDiscoverer{
				Client: apiClient,
			}
			d := &Discoverer{
				TaskDiscoverers: []TaskDiscoverer{defnDiscoverer, scriptDiscoverer},
				Client: &mock.MockClient{
					Tasks: tC.existingTasks,
				},
				Logger: &logger.MockLogger{},
			}
			got, err := d.DiscoverTasks(context.Background(), tC.paths...)
			if tC.expectedErr {
				require.NotNil(err)
				return
			}
			require.NoError(err)

			require.Equal(tC.want, got)
		})
	}
}
