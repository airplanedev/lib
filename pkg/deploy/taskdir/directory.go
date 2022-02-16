package taskdir

import (
	"io"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type TaskDirectory struct {
	// path is the absolute path of the airplane.yml task definition.
	defPath string
	// closer is used to clean up TaskDirectory.
	closer io.Closer
}

// New creates a TaskDirectory struct with the (desired) definition file as input
func New(file string) (TaskDirectory, error) {
	var td TaskDirectory
	var err error
	td.defPath, err = filepath.Abs(file)
	if err != nil {
		return td, errors.Wrap(err, "converting local file path to absolute path")
	}
	return td, nil
}

// Open creates a TaskDirectory struct from a file argument
// Supports file in the form of github.com/path/to/repo/example and will download from GitHub
// Supports file in the form of local_file.yml and will read it to determine the full details
func Open(file string) (TaskDirectory, error) {
	if strings.HasPrefix(file, "http://") {
		return TaskDirectory{}, errors.New("http:// paths are not supported, use https:// instead")
	}

	var td TaskDirectory
	var err error
	if strings.HasPrefix(file, "github.com/") || strings.HasPrefix(file, "https://github.com/") {
		td.defPath, td.closer, err = openGitHubDirectory(file)
		if err != nil {
			return TaskDirectory{}, err
		}
	} else {
		td.defPath, err = filepath.Abs(file)
		if err != nil {
			return TaskDirectory{}, errors.Wrap(err, "converting local file path to absolute path")
		}
	}

	return td, nil
}

func (td TaskDirectory) DefinitionPath() string {
	return td.defPath
}

func (td TaskDirectory) Close() error {
	if td.closer != nil {
		return td.closer.Close()
	}

	return nil
}
