package definitions

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

type ErrSchemaValidation struct {
	Errors []gojsonschema.ResultError
}

func (err ErrSchemaValidation) Error() string {
	return fmt.Sprintf("invalid format: %v", err.Errors)
}
