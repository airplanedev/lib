package mock

import (
	"context"

	"github.com/airplanedev/lib/pkg/api"
)

type MockClient struct {
	Tasks     map[string]api.Task
	Resources []api.Resource
}

var _ api.IAPIClient = &MockClient{}

func (mc *MockClient) GetTask(ctx context.Context, req api.GetTaskRequest) (res api.Task, err error) {
	task, ok := mc.Tasks[req.Slug]
	if !ok {
		return api.Task{}, &api.TaskMissingError{AppURL: "api/", Slug: req.Slug}
	}
	return task, nil
}

func (mc *MockClient) ListResources(ctx context.Context) (res api.ListResourcesResponse, err error) {
	return api.ListResourcesResponse{
		Resources: mc.Resources,
	}, nil
}

func (mc *MockClient) CreateBuildUpload(ctx context.Context, req api.CreateBuildUploadRequest) (res api.CreateBuildUploadResponse, err error) {
	return api.CreateBuildUploadResponse{
		WriteOnlyURL: "writeOnlyURL",
	}, nil
}
