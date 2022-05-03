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

func TestDiscover(t *testing.T) {
	fixturesPath, _ := filepath.Abs("./fixtures")
	tests := []struct {
		name                string
		paths               []string
		existingTasks       map[string]api.Task
		existingApps        map[string]api.App
		expectedErr         bool
		expectedTaskConfigs []TaskConfig
		expectedAppConfigs  []AppConfig
		buildConfigs        []build.BuildConfig
		defnFilePath        string
		absEntrypoints      []string
	}{
		{
			name:  "single script",
			paths: []string{"./fixtures/single_task.js"},
			existingTasks: map[string]api.Task{
				"my_task": {ID: "tsk123", Slug: "my_task", Kind: build.TaskKindNode, InterpolationMode: "handlebars"},
			},
			expectedTaskConfigs: []TaskConfig{
				{
					TaskID:         "tsk123",
					TaskRoot:       fixturesPath,
					TaskEntrypoint: fixturesPath + "/single_task.js",
					Def: &definitions.Definition_0_3{
						Slug: "my_task",
						Node: &definitions.NodeDefinition_0_3{},
					},
					Source: ConfigSourceScript,
				},
			},
			buildConfigs: []build.BuildConfig{
				{
					"workdir":    "",
					"entrypoint": "single_task.js",
				},
			},
		},
		{
			name:  "multiple scripts",
			paths: []string{"./fixtures/single_task.js", "./fixtures/single_task2.js"},
			existingTasks: map[string]api.Task{
				"my_task":  {ID: "tsk123", Slug: "my_task", Kind: build.TaskKindNode, InterpolationMode: "jst"},
				"my_task2": {ID: "tsk456", Slug: "my_task2", Kind: build.TaskKindNode, InterpolationMode: "handlebars"},
			},
			expectedTaskConfigs: []TaskConfig{
				{
					TaskID:         "tsk123",
					TaskRoot:       fixturesPath,
					TaskEntrypoint: fixturesPath + "/single_task.js",
					Def: &definitions.Definition_0_3{
						Slug: "my_task",
						Node: &definitions.NodeDefinition_0_3{},
					},
					Source: ConfigSourceScript,
				},
				{
					TaskID:         "tsk456",
					TaskRoot:       fixturesPath,
					TaskEntrypoint: fixturesPath + "/single_task2.js",
					Def: &definitions.Definition_0_3{
						Slug: "my_task2",
						Node: &definitions.NodeDefinition_0_3{},
					},
					Source: ConfigSourceScript,
				},
			},
			buildConfigs: []build.BuildConfig{
				{
					"workdir":    "",
					"entrypoint": "single_task.js",
				},
				{
					"workdir":    "",
					"entrypoint": "single_task2.js",
				},
			},
		},
		{
			name:  "nested scripts",
			paths: []string{"./fixtures/nestedScripts"},
			existingTasks: map[string]api.Task{
				"my_task":  {ID: "tsk123", Slug: "my_task", Kind: build.TaskKindNode, InterpolationMode: "jst"},
				"my_task2": {ID: "tsk456", Slug: "my_task2", Kind: build.TaskKindNode, InterpolationMode: "jst"},
			},
			expectedTaskConfigs: []TaskConfig{
				{
					TaskID:         "tsk123",
					TaskRoot:       fixturesPath + "/nestedScripts",
					TaskEntrypoint: fixturesPath + "/nestedScripts/single_task.js",
					Def: &definitions.Definition_0_3{
						Slug: "my_task",
						Node: &definitions.NodeDefinition_0_3{},
					},
					Source: ConfigSourceScript,
				},
				{
					TaskID:         "tsk456",
					TaskRoot:       fixturesPath + "/nestedScripts",
					TaskEntrypoint: fixturesPath + "/nestedScripts/single_task2.js",
					Def: &definitions.Definition_0_3{
						Slug: "my_task2",
						Node: &definitions.NodeDefinition_0_3{},
					},
					Source: ConfigSourceScript,
				},
			},
			buildConfigs: []build.BuildConfig{
				{
					"workdir":    "",
					"entrypoint": "single_task.js",
				},
				{
					"workdir":    "",
					"entrypoint": "single_task2.js",
				},
			},
		},
		{
			name:  "single defn",
			paths: []string{"./fixtures/defn.task.yaml"},
			existingTasks: map[string]api.Task{
				"my_task": {ID: "tsk123", Slug: "my_task", Kind: build.TaskKindNode, InterpolationMode: "jst"},
			},
			expectedTaskConfigs: []TaskConfig{
				{
					TaskID:         "tsk123",
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
					Source: ConfigSourceDefn,
				},
			},
			buildConfigs: []build.BuildConfig{
				{
					"workdir": "",
				},
			},
			defnFilePath: fixturesPath + "/defn.task.yaml",
			absEntrypoints: []string{
				fixturesPath + "/single_task.js",
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
				"my_task": {ID: "tsk123", Slug: "my_task", Kind: build.TaskKindNode, InterpolationMode: "jst"},
			},
			expectedTaskConfigs: []TaskConfig{
				{
					TaskID:         "tsk123",
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
					Source: ConfigSourceDefn,
				},
			},
			buildConfigs: []build.BuildConfig{
				{
					"workdir": "",
				},
			},
			defnFilePath: fixturesPath + "/defn.task.yaml",
			absEntrypoints: []string{
				fixturesPath + "/single_task.js",
			},
		},
		{
			name:  "different working directory",
			paths: []string{"./fixtures/subdir/single_task.js"},
			existingTasks: map[string]api.Task{
				"my_task": {ID: "tsk123", Slug: "my_task", Kind: build.TaskKindNode, InterpolationMode: "jst"},
			},
			expectedTaskConfigs: []TaskConfig{
				{
					TaskID:         "tsk123",
					TaskRoot:       fixturesPath,
					TaskEntrypoint: fixturesPath + "/subdir/single_task.js",
					Def: &definitions.Definition_0_3{
						Slug: "my_task",
						Node: &definitions.NodeDefinition_0_3{},
					},
					Source: ConfigSourceScript,
				},
			},
			buildConfigs: []build.BuildConfig{
				{
					"workdir":    "/subdir",
					"entrypoint": "subdir/single_task.js",
				},
			},
		},
		{
			name:  "non linked script with def in same directory",
			paths: []string{"./fixtures/nonlinkedscript/single_task.js"},
			existingTasks: map[string]api.Task{
				"my_task": {ID: "tsk123", Slug: "my_task", Kind: build.TaskKindNode, InterpolationMode: "jst"},
			},
			expectedTaskConfigs: []TaskConfig{
				{
					TaskID:         "tsk123",
					TaskRoot:       fixturesPath,
					TaskEntrypoint: fixturesPath + "/nonlinkedscript/single_task.js",
					Def: &definitions.Definition_0_3{
						Name:        "sunt in tempor eu",
						Slug:        "my_task",
						Description: "ut dolor sit officia ea",
						Node: &definitions.NodeDefinition_0_3{
							Entrypoint:  "./single_task.js",
							NodeVersion: "14",
						},
					},
					Source: ConfigSourceDefn,
				},
			},
			buildConfigs: []build.BuildConfig{
				{
					"workdir":    "/nonlinkedscript",
					"entrypoint": "nonlinkedscript/single_task.js",
				},
			},
			defnFilePath: fixturesPath + "/nonlinkedscript/single_task.task.yaml",
			absEntrypoints: []string{
				fixturesPath + "/nonlinkedscript/single_task.js",
			},
		},
		{
			name:  "non linked script with def in same directory - entire directory deployed",
			paths: []string{"./fixtures/nonlinkedscript"},
			existingTasks: map[string]api.Task{
				"my_task": {ID: "tsk123", Slug: "my_task", Kind: build.TaskKindNode, InterpolationMode: "jst"},
			},
			expectedTaskConfigs: []TaskConfig{
				{
					TaskID:         "tsk123",
					TaskRoot:       fixturesPath,
					TaskEntrypoint: fixturesPath + "/nonlinkedscript/single_task.js",
					Def: &definitions.Definition_0_3{
						Name:        "sunt in tempor eu",
						Slug:        "my_task",
						Description: "ut dolor sit officia ea",
						Node: &definitions.NodeDefinition_0_3{
							Entrypoint:  "./single_task.js",
							NodeVersion: "14",
						},
					},
					Source: ConfigSourceDefn,
				},
			},
			buildConfigs: []build.BuildConfig{
				{
					"workdir":    "/nonlinkedscript",
					"entrypoint": "nonlinkedscript/single_task.js",
				},
			},
			defnFilePath: fixturesPath + "/nonlinkedscript/single_task.task.yaml",
			absEntrypoints: []string{
				fixturesPath + "/nonlinkedscript/single_task.js",
			},
		},
		{
			name:  "discovers definition when script is deployed",
			paths: []string{"./fixtures/subdir/defn.task.yaml"},
			existingTasks: map[string]api.Task{
				"my_task": {ID: "tsk123", Slug: "my_task", Kind: build.TaskKindNode, InterpolationMode: "jst"},
			},
			expectedTaskConfigs: []TaskConfig{
				{
					TaskID:         "tsk123",
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
					Source: ConfigSourceDefn,
				},
			},
			buildConfigs: []build.BuildConfig{
				{
					"workdir":    "/subdir",
					"entrypoint": "subdir/single_task.js",
				},
			},
			defnFilePath: fixturesPath + "/subdir/defn.task.yaml",
			absEntrypoints: []string{
				fixturesPath + "/subdir/single_task.js",
			},
		},
		{
			name:  "defn - entrypoint does not exist",
			paths: []string{"./fixtures/defn_incorrect_entrypoint.task.yaml"},
			existingTasks: map[string]api.Task{
				"incorrect_entrypoint": {ID: "tsk123", Slug: "incorrect_entrypoint", Kind: build.TaskKindNode, InterpolationMode: "jst"},
			},
			expectedErr: true,
		},
		{
			name:  "app defn",
			paths: []string{"./fixtures/app/defn.app.yaml"},
			existingApps: map[string]api.App{
				"my_app": {ID: "app123", Slug: "my_app", Name: "My App"},
			},
			expectedAppConfigs: []AppConfig{
				{
					ID:         "app123",
					Slug:       "my_app",
					Root:       fixturesPath + "/app",
					Entrypoint: fixturesPath + "/app/foo.js",
					Source:     ConfigSourceDefn,
				},
			},
		},
	}
	for _, tC := range tests {
		t.Run(tC.name, func(t *testing.T) {
			require := require.New(t)
			apiClient := &mock.MockClient{
				Tasks: tC.existingTasks,
				Apps:  tC.existingApps,
			}
			scriptDiscoverer := &ScriptDiscoverer{
				Client: apiClient,
				Logger: &logger.MockLogger{},
			}
			defnDiscoverer := &DefnDiscoverer{
				Client: apiClient,
				Logger: &logger.MockLogger{},
			}
			appDefnDiscoverer := &AppDefnDiscoverer{
				Client: apiClient,
				Logger: &logger.MockLogger{},
			}
			d := &Discoverer{
				TaskDiscoverers: []TaskDiscoverer{defnDiscoverer, scriptDiscoverer},
				AppDiscoverers:  []AppDiscoverer{appDefnDiscoverer},
				Client:          apiClient,
				Logger:          &logger.MockLogger{},
			}
			taskConfigs, appConfigs, err := d.Discover(context.Background(), tC.paths...)
			if tC.expectedErr {
				require.NotNil(err)
				return
			}
			require.NoError(err)

			require.Equal(len(tC.expectedTaskConfigs), len(taskConfigs))
			for i := range tC.expectedTaskConfigs {
				for k, v := range tC.buildConfigs[i] {
					tC.expectedTaskConfigs[i].Def.SetBuildConfig(k, v)
				}
				if i < len(tC.absEntrypoints) {
					err := tC.expectedTaskConfigs[i].Def.SetAbsoluteEntrypoint(tC.absEntrypoints[i])
					require.NoError(err)
				}
				tC.expectedTaskConfigs[i].Def.SetDefnFilePath(tC.defnFilePath)
				require.Equal(tC.expectedTaskConfigs[i], taskConfigs[i])
			}

			require.Equal(len(tC.expectedAppConfigs), len(appConfigs))
			for i := range tC.expectedAppConfigs {
				require.Equal(tC.expectedAppConfigs[i], appConfigs[i])
			}
		})
	}
}
