package api

import (
	"context"
)

type MockClient struct {
	Tasks map[string]Task
}

var _ IAPIClient = &MockClient{}

func (mc *MockClient) GetTask(ctx context.Context, slug string) (res Task, err error) {
	task, ok := mc.Tasks[slug]
	if !ok {
		return Task{}, &TaskMissingError{AppURL: "api/", Slug: slug}
	}
	return task, nil
}

func (mc *MockClient) ListResources(ctx context.Context) (res ListResourcesResponse, err error) {
	return ListResourcesResponse{}, nil
}

func (mc *MockClient) CreateBuildUpload(ctx context.Context, req CreateBuildUploadRequest) (res CreateBuildUploadResponse, err error) {
	return CreateBuildUploadResponse{
		WriteOnlyURL: "writeOnlyURL",
	}, nil
}
