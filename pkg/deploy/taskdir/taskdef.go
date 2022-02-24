package taskdir

import (
	"fmt"
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
		switch err := errors.Cause(err).(type) {
		case definitions.ErrSchemaValidation:
			errorMsgs := []string{}
			for _, verr := range err.Errors {
				errorMsgs = append(errorMsgs, fmt.Sprintf("%s: %s", verr.Field(), verr.Description()))
			}
			return nil, definitions.NewErrReadDefinition(fmt.Sprintf("Error reading %s", defPath), errorMsgs...)
		default:
			return nil, errors.Wrap(err, "unmarshalling task definition")
		}
	}
	def.SetDefinitionPath(defPath)
	return &def, nil
}
