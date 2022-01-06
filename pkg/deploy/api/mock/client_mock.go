package mock

import (
	"context"

	"github.com/airplanedev/lib/pkg/deploy/api"
)

type MockClient struct {
	Tasks map[string]api.Task
}

var _ api.IAPIClient = &MockClient{}

func (mc *MockClient) GetTask(ctx context.Context, slug string) (res api.Task, err error) {
	task, ok := mc.Tasks[slug]
	if !ok {
		return api.Task{}, &api.TaskMissingError{AppURL: "api/", Slug: slug}
	}
	return task, nil
}

func (mc *MockClient) ListResources(ctx context.Context) (res api.ListResourcesResponse, err error) {
	return api.ListResourcesResponse{}, nil
}

func (mc *MockClient) CreateBuildUpload(ctx context.Context, req api.CreateBuildUploadRequest) (res api.CreateBuildUploadResponse, err error) {
	return api.CreateBuildUploadResponse{
		WriteOnlyURL: "writeOnlyURL",
	}, nil
}
