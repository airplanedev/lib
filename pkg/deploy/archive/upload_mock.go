package archive

import (
	"context"
	"os"
)

type MockUploader struct {
	UploadCount int
}

var _ Uploader = &MockUploader{}

func (d *MockUploader) Upload(ctx context.Context, url string, archive *os.File) error {
	d.UploadCount++
	return nil
}
