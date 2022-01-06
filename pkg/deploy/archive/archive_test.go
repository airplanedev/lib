package archive

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/airplanedev/lib/pkg/deploy/api/mock"
	"github.com/airplanedev/lib/pkg/utils/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArchive(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	fixturesPath, _ := filepath.Abs("./fixtures")

	testCases := []struct {
		desc        string
		roots       []string
		numUploaded int
	}{
		{
			desc:        "Archives one task",
			roots:       []string{fixturesPath},
			numUploaded: 1,
		},
		{
			desc:        "Archives two tasks, same root",
			roots:       []string{fixturesPath, fixturesPath},
			numUploaded: 1,
		},
		{
			desc:        "Archives two tasks, diff roots",
			roots:       []string{fixturesPath, fixturesPath + "/nested"},
			numUploaded: 2,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			l := &logger.MockLogger{}
			client := &mock.MockClient{}
			uploader := &MockUploader{}
			archiver := NewAPIArchiver(l, client, uploader)

			var numUploaded int
			for _, root := range tC.roots {
				_, sizeBytes, err := archiver.Archive(context.Background(), root)
				require.NoError(err)
				numUploaded++
				if numUploaded > tC.numUploaded {
					assert.Empty(sizeBytes)
				}
			}
			assert.Equal(tC.numUploaded, uploader.UploadCount)
		})
	}
}
