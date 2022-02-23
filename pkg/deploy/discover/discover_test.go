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
	tests := []struct {
		name          string
		paths         []string
		existingTasks map[string]api.Task
		expectedErr   bool
		want          []TaskConfig
		buildConfigs  []build.BuildConfig
	}{
		{
			name:  "single script",
			paths: []string{"./fixtures/single_task.js"},
			existingTasks: map[string]api.Task{
				"my_task": {Slug: "my_task", Kind: build.TaskKindNode},
			},
			want: []TaskConfig{
				{
					TaskRoot:       fixturesPath,
					TaskEntrypoint: fixturesPath + "/single_task.js",
					Def: &definitions.Definition{
						Slug: "my_task",
						Node: &definitions.NodeDefinition{
							Entrypoint: "single_task.js",
						},
					},
					Task: api.Task{Slug: "my_task", Kind: build.TaskKindNode},
					From: TaskConfigSourceScript,
				},
			},
			buildConfigs: []build.BuildConfig{
				{
					"workdir": "",
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
					TaskRoot:       fixturesPath,
					TaskEntrypoint: fixturesPath + "/single_task.js",
					Def: &definitions.Definition{
						Slug: "my_task",
						Node: &definitions.NodeDefinition{
							Entrypoint: "single_task.js",
						},
					},
					Task: api.Task{Slug: "my_task", Kind: build.TaskKindNode},
					From: TaskConfigSourceScript,
				},
				{
					TaskRoot:       fixturesPath,
					TaskEntrypoint: fixturesPath + "/single_task2.js",
					Def: &definitions.Definition{
						Slug: "my_task2",
						Node: &definitions.NodeDefinition{
							Entrypoint: "single_task2.js",
						},
					},
					Task: api.Task{Slug: "my_task2", Kind: build.TaskKindNode},
					From: TaskConfigSourceScript,
				},
			},
			buildConfigs: []build.BuildConfig{
				{
					"workdir": "",
				},
				{
					"workdir": "",
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
					TaskRoot:       fixturesPath + "/nestedScripts",
					TaskEntrypoint: fixturesPath + "/nestedScripts/single_task.js",
					Def: &definitions.Definition{
						Slug: "my_task",
						Node: &definitions.NodeDefinition{
							Entrypoint: "single_task.js",
						},
					},
					Task: api.Task{Slug: "my_task", Kind: build.TaskKindNode},
					From: TaskConfigSourceScript,
				},
				{
					TaskRoot:       fixturesPath + "/nestedScripts",
					TaskEntrypoint: fixturesPath + "/nestedScripts/single_task2.js",
					Def: &definitions.Definition{
						Slug: "my_task2",
						Node: &definitions.NodeDefinition{
							Entrypoint: "single_task2.js",
						},
					},
					Task: api.Task{Slug: "my_task2", Kind: build.TaskKindNode},
					From: TaskConfigSourceScript,
				},
			},
			buildConfigs: []build.BuildConfig{
				{
					"workdir": "",
				},
				{
					"workdir": "",
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
					TaskRoot:       fixturesPath,
					TaskEntrypoint: fixturesPath + "/single_task.js",
					Def: &definitions.Definition_0_3{
						Name:        "sunt in tempor eu",
						Slug:        "my_task",
						Description: "ut dolor sit officia ea",
						Node: &definitions.NodeDefinition_0_3{
							Entrypoint:  "./single_task.js",
							NodeVersion: "14",
						},
					},
					Task: api.Task{Slug: "my_task", Kind: build.TaskKindNode},
					From: TaskConfigSourceDefn,
				},
			},
			buildConfigs: []build.BuildConfig{
				{
					"workdir": "",
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
					TaskRoot:       fixturesPath,
					TaskEntrypoint: fixturesPath + "/single_task.js",
					Def: &definitions.Definition_0_3{
						Name:        "sunt in tempor eu",
						Slug:        "my_task",
						Description: "ut dolor sit officia ea",
						Node: &definitions.NodeDefinition_0_3{
							Entrypoint:  "./single_task.js",
							NodeVersion: "14",
						},
					},
					Task: api.Task{Slug: "my_task", Kind: build.TaskKindNode},
					From: TaskConfigSourceDefn,
				},
			},
			buildConfigs: []build.BuildConfig{
				{
					"workdir": "",
				},
			},
		},
		{
			name:  "different working directory",
			paths: []string{"./fixtures/subdir/single_task.js"},
			existingTasks: map[string]api.Task{
				"my_task": {Slug: "my_task", Kind: build.TaskKindNode},
			},
			want: []TaskConfig{
				{
					TaskRoot:       fixturesPath,
					TaskEntrypoint: fixturesPath + "/subdir/single_task.js",
					Def: &definitions.Definition{
						Slug: "my_task",
						Node: &definitions.NodeDefinition{
							Entrypoint: "subdir/single_task.js",
						},
					},
					Task: api.Task{Slug: "my_task", Kind: build.TaskKindNode},
					From: TaskConfigSourceScript,
				},
			},
			buildConfigs: []build.BuildConfig{
				{
					"workdir": "/subdir",
				},
			},
		},
		{
			name:  "different working directory, with definition",
			paths: []string{"./fixtures/subdir/defn.task.yaml"},
			existingTasks: map[string]api.Task{
				"my_task": {Slug: "my_task", Kind: build.TaskKindNode},
			},
			want: []TaskConfig{
				{
					TaskRoot:       fixturesPath,
					TaskEntrypoint: fixturesPath + "/subdir/single_task.js",
					Def: &definitions.Definition_0_3{
						Name:        "sunt in tempor eu",
						Slug:        "my_task",
						Description: "ut dolor sit officia ea",
						Node: &definitions.NodeDefinition_0_3{
							Entrypoint:  "./single_task.js",
							NodeVersion: "14",
						},
					},
					Task: api.Task{Slug: "my_task", Kind: build.TaskKindNode},
					From: TaskConfigSourceDefn,
				},
			},
			buildConfigs: []build.BuildConfig{
				{
					"workdir":    "/subdir",
					"entrypoint": "subdir/single_task.js",
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

			require.Equal(len(tC.want), len(got))
			for i := range tC.want {
				for k, v := range tC.buildConfigs[i] {
					tC.want[i].Def.SetBuildConfig(k, v)
				}
				require.Equal(tC.want[i], got[i])
			}
		})
	}
}
