package build

import (
	"io/ioutil"
	"path/filepath"

	"github.com/airplanedev/lib/pkg/build/fsx"
	"github.com/pkg/errors"
)

func dockerfile(root string, options KindOptions) (string, error) {
	dockerfile, _ := options["dockerfile"].(string)
	dockerfilePath := filepath.Join(root, dockerfile)
	if err := fsx.AssertExistsAll(dockerfilePath); err != nil {
		return "", err
	}

	contents, err := ioutil.ReadFile(dockerfilePath)
	if err != nil {
		return "", errors.Wrap(err, "opening dockerfile")
	}

	return string(contents), nil
}
