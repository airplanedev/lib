package taskdir

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/airplanedev/lib/pkg/deploy/taskdir/definitions"
	"github.com/pkg/errors"
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
	def.SetDefinitionPath(defPath)
	return &def, nil
}
