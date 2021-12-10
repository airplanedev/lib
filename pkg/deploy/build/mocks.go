package build

import (
	"context"

	"github.com/airplanedev/lib/pkg/build"
)

type MockBuildCreator struct{}

var _ BuildCreator = &MockBuildCreator{}

func (mbc *MockBuildCreator) CreateBuild(ctx context.Context, req Request) (*build.Response, error) {
	return &build.Response{
		BuildID:  "buildID",
		ImageURL: "imageURL",
	}, nil
}
