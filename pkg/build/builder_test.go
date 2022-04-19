package build

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"testing"

	"github.com/airplanedev/dlog"
	"github.com/airplanedev/lib/pkg/examples"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/require"
)

type Test struct {
	// Root is the task root to perform a build inside of.
	Root        string
	Kind        TaskKind
	Options     KindOptions
	ParamValues Values
	// SearchString is a string to look for in the example's output
	// to validate that the task completed successfully. If not set,
	// defaults to a random value which is passed into the example
	// via the `id` parameter.
	SearchString string
}

// RunTests performs a series of builder tests and looks for a given SearchString
// in the task's output to validate that the task built + ran correctly.
func RunTests(tt *testing.T, ctx context.Context, tests []Test) {
	for _, test := range tests {
		test := test // loop local reference
		tt.Run(filepath.Base(test.Root), func(t *testing.T) {
			// These tests can run in parallel, but it may exhaust all memory
			// allocated to the Docker daemon on your computer. For that reason,
			// we don't currently run them in parallel. We could gate parallel
			// execution to CI via `os.Getenv("CI") == "true"`, but that may
			// lead to scenarios where tests break in CI but not locally. If
			// test performance in CI becomes an issue, we should look into caching
			// Docker builds in CI since (locally) that appears to have a significant
			// impact on e2e times for this test suite.
			//
			// t.Parallel()

			require := require.New(t)

			b, err := New(LocalConfig{
				Root:    examples.Path(t, test.Root),
				Builder: string(test.Kind),
				Options: test.Options,
			})
			require.NoError(err)
			t.Cleanup(func() {
				require.NoError(b.Close())
			})

			// Perform the docker build:
			resp, err := b.Build(ctx, "builder-tests", ksuid.New().String())
			require.NoError(err)
			defer func() {
				_, err := b.client.ImageRemove(ctx, resp.ImageURL, types.ImageRemoveOptions{})
				require.NoError(err)
			}()

			if test.ParamValues == nil {
				test.ParamValues = Values{}
			}
			if test.SearchString == "" {
				test.SearchString = ksuid.New().String()
				test.ParamValues["id"] = test.SearchString
			}

			// Run the produced docker image:
			runTask(t, ctx, b.client, resp.ImageURL, test.ParamValues)
			// require.True(strings.Contains(string(out), test.SearchString), "unable to find %q in output:\n%s", test.SearchString, string(out))
		})
	}
}

func runTask(t *testing.T, ctx context.Context, dclient *client.Client, image string, paramValues Values) []byte {
	require := require.New(t)

	pv, err := json.Marshal(paramValues)
	require.NoError(err)

	resp, err := dclient.ContainerCreate(ctx, &container.Config{
		Image: image,
		Tty:   false,
		Cmd:   strslice.StrSlice{string(pv)},
	}, nil, nil, nil, "")
	require.NoError(err)
	containerID := resp.ID
	defer func() {
		// Cleanup this container when we complete these tests:
		require.NoError(dclient.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
			Force: true,
		}))
	}()

	require.NoError(dclient.ContainerStart(ctx, containerID, types.ContainerStartOptions{}))

	resultC, errC := dclient.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)

	logr, err := dclient.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Follow:     true,
		Tail:       "all",
	})
	require.NoError(err)
	defer logr.Close()

	logs, err := io.ReadAll(dlog.NewReader(logr, dlog.Options{
		AppendNewline: true,
	}))
	require.NoError(err)
	fmt.Println("ef", string(logs))

	select {
	case result := <-resultC:
		require.Nil(result.Error)
		require.Equal(int64(0), result.StatusCode, "container exited with non-zero status code: %v", result.StatusCode)
	case err := <-errC:
		require.NoError(err)
	}

	return logs
}

func TestInlineString(t *testing.T) {
	require := require.New(t)

	require.Equal(
		`printf 'The sheep couldn'"'"'t sleep, no matter how many humans he counted.'`,
		inlineString(`The sheep couldn't sleep, no matter how many humans he counted.`),
	)
	require.Equal(
		`printf ''"'"''"'"''"'"''`,
		inlineString(`'''`),
	)
	require.Equal(
		`printf 'hi\nline'`,
		inlineString(`hi\nline`),
	)
}
