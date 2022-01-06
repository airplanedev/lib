package archive

import (
	"context"
	"io/ioutil"
	"os"
	"path"

	"github.com/airplanedev/archiver"
	"github.com/airplanedev/lib/pkg/build/ignore"
	"github.com/airplanedev/lib/pkg/deploy/api"
	"github.com/airplanedev/lib/pkg/utils/logger"
	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"
)

type Archiver interface {
	Archive(ctx context.Context, root string) (uploadID string, size int, err error)
}

type apiArchiver struct {
	logger   logger.Logger
	client   api.IAPIClient
	uploader Uploader

	uploadArchiveSingleFlightGroup singleflight.Group
	uploadedArchives               map[string]string
}

var _ Archiver = &apiArchiver{}

func NewAPIArchiver(logger logger.Logger, client api.IAPIClient, uploader Uploader) Archiver {
	return &apiArchiver{
		uploadedArchives: make(map[string]string),
		logger:           logger,
		client:           client,
		uploader:         uploader,
	}
}

func (d *apiArchiver) Archive(ctx context.Context, root string) (string, int, error) {
	tmpdir, err := ioutil.TempDir("", "airplane-builds-")
	if err != nil {
		return "", 0, errors.Wrap(err, "creating temporary directory for remote build")
	}
	defer os.RemoveAll(tmpdir)

	archivePath := path.Join(tmpdir, "archive.tar.gz")
	if err := archiveTaskDir(root, archivePath); err != nil {
		return "", 0, err
	}

	uploadIDRes, err, _ := d.uploadArchiveSingleFlightGroup.Do(root, func() (interface{}, error) {
		return d.uploadArchive(ctx, archivePath, root)
	})
	if err != nil {
		return "", 0, err
	}
	upload := uploadIDRes.(uploadRes)
	return upload.uploadID, upload.sizeBytes, nil
}

type uploadRes struct {
	uploadID  string
	sizeBytes int
}

func (d *apiArchiver) uploadArchive(ctx context.Context, archivePath, rootPath string) (uploadRes, error) {
	// Check if anyone has uploaded an archive for this path.
	uid, ok := d.uploadedArchives[rootPath]
	if ok {
		// Somebody has already uploaded the path. Re-use the upload ID.
		return uploadRes{uploadID: uid}, nil
	}

	archive, err := os.OpenFile(archivePath, os.O_RDONLY, 0)
	if err != nil {
		return uploadRes{}, errors.Wrap(err, "opening archive file")
	}
	defer archive.Close()

	info, err := archive.Stat()
	if err != nil {
		return uploadRes{}, errors.Wrap(err, "stat on archive file")
	}
	sizeBytes := int(info.Size())

	upload, err := d.client.CreateBuildUpload(ctx, api.CreateBuildUploadRequest{
		SizeBytes: sizeBytes,
	})
	if err != nil {
		return uploadRes{}, errors.Wrap(err, "creating upload")
	}

	if err := d.uploader.Upload(ctx, upload.WriteOnlyURL, archive); err != nil {
		return uploadRes{}, err
	}

	uploadID := upload.Upload.ID

	// Populate the cache so that we can reuse the upload.
	d.uploadedArchives[rootPath] = uploadID

	return uploadRes{uploadID: uploadID, sizeBytes: sizeBytes}, nil
}

func archiveTaskDir(root string, archivePath string) error {
	// mholt/archiver takes a list of "sources" (files/directories) that will
	// be included in the root of the archive. In our case, we want the root of
	// the archive to be the contents of the task directory, rather than the
	// task directory itself.
	var sources []string
	if files, err := ioutil.ReadDir(root); err != nil {
		return errors.Wrap(err, "inspecting files in task root")
	} else {
		for _, f := range files {
			sources = append(sources, path.Join(root, f.Name()))
		}
	}

	var err error
	arch := archiver.NewTarGz()
	arch.Tar.IncludeFunc, err = ignore.Func(root)
	if err != nil {
		return err
	}

	if err := arch.Archive(sources, archivePath); err != nil {
		return errors.Wrap(err, "building archive")
	}

	return nil
}
