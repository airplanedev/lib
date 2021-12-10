package api

import "context"

type APIClient interface {
	GetTask(ctx context.Context, slug string) (res Task, err error)
	ListResources(ctx context.Context) (res ListResourcesResponse, err error)
	TaskURL(slug string) string
	UpdateTask(ctx context.Context, req UpdateTaskRequest) (res UpdateTaskResponse, err error)
	CreateTask(ctx context.Context, req CreateTaskRequest) (res CreateTaskResponse, err error)
	CreateBuild(ctx context.Context, req CreateBuildRequest) (res CreateBuildResponse, err error)
	GetRegistryToken(ctx context.Context) (res RegistryTokenResponse, err error)
	CreateBuildUpload(ctx context.Context, req CreateBuildUploadRequest) (res CreateBuildUploadResponse, err error)
}
