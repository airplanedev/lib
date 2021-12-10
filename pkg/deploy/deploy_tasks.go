package deploy

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	libBuild "github.com/airplanedev/lib/pkg/build"
	"github.com/airplanedev/lib/pkg/common/api"
	"github.com/airplanedev/lib/pkg/common/logger"
	"github.com/airplanedev/lib/pkg/deploy/build"
	"github.com/airplanedev/lib/pkg/deploy/discover"
	"github.com/airplanedev/lib/pkg/runtime"
	"github.com/airplanedev/lib/pkg/utils/pointers"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type deployer struct {
	buildCreator build.BuildCreator
	logger       logger.Logger
	changedFiles []string
	client       api.APIClient

	erroredTaskSlugs  map[string]error
	deployedTaskSlugs []string
	mu                sync.Mutex
}

type DeployerOpts struct {
	ChangedFiles []string
}

// Set of properties to track when deploying
type taskDeployedProps struct {
	from       string
	kind       libBuild.TaskKind
	taskID     string
	taskSlug   string
	taskName   string
	buildLocal bool
	buildID    string
}

func NewDeployer(logger logger.Logger, client api.APIClient, bc build.BuildCreator, opts DeployerOpts) *deployer {
	return &deployer{
		buildCreator:     bc,
		erroredTaskSlugs: make(map[string]error),
		logger:           logger,
		changedFiles:     opts.ChangedFiles,
		client:           client,
	}
}

// DeployTasks deploys all taskConfigs as Airplane tasks.
// It concurrently builds (if needed) and updates tasks.
func (d *deployer) DeployTasks(ctx context.Context, taskConfigs []discover.TaskConfig) error {
	if len(d.changedFiles) > 0 {
		// Filter out any tasks that don't have changed files.
		var filteredTaskConfigs []discover.TaskConfig
		for _, tc := range taskConfigs {
			contains, err := containsFile(tc.TaskRoot, d.changedFiles)
			if err != nil {
				return err
			}
			if contains {
				filteredTaskConfigs = append(filteredTaskConfigs, tc)
			}
		}
		if len(taskConfigs) != len(filteredTaskConfigs) {
			d.logger.Log("Changed files specified. Filtered %d task(s) to %d affected task(s)", len(taskConfigs), len(filteredTaskConfigs))
		}
		taskConfigs = filteredTaskConfigs
	}

	if len(taskConfigs) == 0 {
		d.logger.Log("No tasks to deploy")
		return nil
	}

	// Print out a summary before deploying.
	noun := "task"
	if len(taskConfigs) > 1 {
		noun = fmt.Sprintf("%ss", noun)
	}
	d.logger.Log("Deploying %v %v:\n", len(taskConfigs), noun)
	for _, tc := range taskConfigs {
		d.logger.Log(logger.Bold(tc.Task.Slug))
		d.logger.Log("Type: %s", tc.Task.Kind)
		d.logger.Log("Root directory: %s", relpath(tc.TaskRoot))
		if tc.WorkingDirectory != tc.TaskRoot {
			d.logger.Log("Working directory: %s", relpath(tc.WorkingDirectory))
		}
		d.logger.Log("URL: %s", d.client.TaskURL(tc.Task.Slug))
		d.logger.Log("")
	}

	g := new(errgroup.Group)
	// Concurrently deploy the tasks.
	for _, tc := range taskConfigs {
		tc := tc
		g.Go(func() error {
			err := d.deployTask(ctx, tc)
			d.mu.Lock()
			defer d.mu.Unlock()
			if err != nil {
				if !errors.As(err, &runtime.ErrNotLinked{}) {
					d.erroredTaskSlugs[tc.Task.Slug] = err
					return err
				}
			} else {
				d.deployedTaskSlugs = append(d.deployedTaskSlugs, tc.Task.Slug)
			}
			return nil
		})
	}

	groupErr := g.Wait()

	// All of the deploys have finished.
	for taskSlug, err := range d.erroredTaskSlugs {
		d.logger.Log("\n" + logger.Bold(taskSlug))
		d.logger.Log("Status: " + logger.Bold(logger.Red("failed")))
		d.logger.Error(err.Error())
	}
	for _, slug := range d.deployedTaskSlugs {
		d.logger.Log("\n" + logger.Bold(slug))
		d.logger.Log("Status: %s", logger.Bold(logger.Green("succeeded")))
		d.logger.Log("Execute the task: %s", d.client.TaskURL(slug))
	}

	return groupErr
}

func (d *deployer) deployTask(ctx context.Context, tc discover.TaskConfig) (rErr error) {
	tp := taskDeployedProps{
		from:     string(tc.From),
		kind:     tc.Kind,
		taskID:   tc.Task.ID,
		taskSlug: tc.Task.Slug,
		taskName: tc.Task.Name,
	}

	kind, _, err := tc.Def.GetKindAndOptions()
	if err != nil {
		return err
	}
	var image *string
	if ok, err := libBuild.NeedsBuilding(kind); err != nil {
		return err
	} else if ok {
		env, err := tc.Def.GetEnv()
		if err != nil {
			return err
		}
		resp, err := d.buildCreator.CreateBuild(ctx, build.Request{
			TaskID:  tc.Task.ID,
			Root:    tc.TaskRoot,
			Def:     tc.Def,
			TaskEnv: env,
			Shim:    true,
		})
		if err != nil {
			return err
		}
		tp.buildID = resp.BuildID
		image = &resp.ImageURL
	}

	utr, err := tc.Def.GetUpdateTaskRequest(ctx, d.client, image)
	if err != nil {
		return err
	}

	utr.BuildID = pointers.String(tp.buildID)
	utr.InterpolationMode = tc.Task.InterpolationMode
	utr.RequireExplicitPermissions = tc.Task.RequireExplicitPermissions
	utr.Permissions = tc.Task.Permissions

	_, err = d.client.UpdateTask(ctx, utr)
	return err
}

// containsFile returns true if the directory contains at least one of the files.
func containsFile(dir string, filePaths []string) (bool, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return false, errors.Wrapf(err, "calculating absolute path of directory %s", dir)
	}
	for _, cf := range filePaths {
		absCF, err := filepath.Abs(cf)
		if err != nil {
			return false, errors.Wrapf(err, "calculating absolute path of file %s", cf)
		}
		changedFileDir := filepath.Dir(absCF)
		if strings.HasPrefix(changedFileDir, absDir) {
			return true, nil
		}
	}
	return false, nil
}

// Relpath returns the relative using root and the cwd.
func relpath(root string) string {
	if path, err := os.Getwd(); err == nil {
		if rp, err := filepath.Rel(path, root); err == nil {
			if len(rp) == 0 || rp == "." {
				// "." can be missed easily, change it to ./
				return "./"
			}
			return "./" + rp
		}
	}
	return root
}
