package archive

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

type Uploader interface {
	Upload(ctx context.Context, url string, archive *os.File) error
}

type HttpUploader struct {
}

var _ Uploader = &HttpUploader{}

func (d *HttpUploader) Upload(ctx context.Context, url string, archive *os.File) error {
	info, err := archive.Stat()
	if err != nil {
		return errors.Wrap(err, "stat on archive file")
	}
	sizeBytes := int(info.Size())

	req, err := http.NewRequestWithContext(ctx, "PUT", url, archive)
	if err != nil {
		return errors.Wrap(err, "creating GCS upload request")
	}
	req.Header.Add("X-Goog-Content-Length-Range", fmt.Sprintf("0,%d", sizeBytes))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "uploading to GCS")
	}
	defer resp.Body.Close()

	return nil
}
