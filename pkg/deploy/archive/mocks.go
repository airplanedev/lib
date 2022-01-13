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

type MockArchiver struct {
	UploadCount int
}

var _ Archiver = &MockArchiver{}

func (d *MockArchiver) Archive(ctx context.Context, root string) (uploadID string, size int, err error) {
	return "uploadID", 65, nil
}
