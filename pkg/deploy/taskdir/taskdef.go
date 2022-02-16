package taskdir

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/airplanedev/lib/pkg/deploy/taskdir/definitions"
	yamlutils "github.com/airplanedev/lib/pkg/utils/yaml"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func (td TaskDirectory) ReadDefinition() (definitions.DefinitionInterface, error) {
	buf, err := ioutil.ReadFile(td.defPath)
	if err != nil {
		return nil, errors.Wrap(err, "reading task definition")
	}

	defPath := td.defPath
	// Attempt to set a prettier defPath, best effort
	if wd, err := os.Getwd(); err != nil {
	} else if path, err := filepath.Rel(wd, defPath); err != nil {
	} else {
		defPath = path
	}

	def := definitions.Definition_0_3{}
	if err := def.Unmarshal(definitions.GetTaskDefFormat(defPath), buf); err != nil {
		return nil, errors.Wrap(err, "unmarshalling task definition")
	}
	return &def, nil
}

// WriteSlug updates the slug of a task definition and persists td to disk.
//
// It attempts to retain the existing file's formatting (comments, etc.) where possible.
func (td TaskDirectory) WriteSlug(slug string) error {
	if err := yamlutils.SetYAMLField(td.defPath, "slug", slug); err != nil {
		return errors.Wrap(err, "setting slug")
	}

	return nil
}

func (td TaskDirectory) WriteDefinition(def definitions.Definition) error {
	data, err := yaml.Marshal(def)
	if err != nil {
		return errors.Wrap(err, "marshalling definition")
	}

	if err := ioutil.WriteFile(td.defPath, data, 0664); err != nil {
		return errors.Wrap(err, "writing file")
	}

	return nil
}
