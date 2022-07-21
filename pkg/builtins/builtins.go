package builtins

import (
	"fmt"
	"strings"
)

const builtinsSlugPrefix = "airplane"

// Specification of a function via namespace and name. Each function should have a unique namespace + name combination.
type FunctionSpecification struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

func IsBuiltinTaskSlug(slug string) bool {
	_, err := GetBuiltinFunctionSpecification(slug)
	return err == nil
}

// Returns an unvalidated function specification if it matches the format of
// airplane:<namespace>_<name>.
func GetBuiltinFunctionSpecification(slug string) (FunctionSpecification, error) {
	// Validate function specification and schema
	slugParts := strings.Split(slug, ":")
	if len(slugParts) != 2 || slugParts[0] != builtinsSlugPrefix {
		return FunctionSpecification{}, fmt.Errorf("unknown builtin task slug: %s", slug)
	}
	namespaceParts := strings.Split(slugParts[1], "_")
	if len(namespaceParts) != 2 { // This may not be true in the future
		return FunctionSpecification{}, fmt.Errorf("unknown builtin task slug: %s", slug)
	}
	return FunctionSpecification{
		Namespace: namespaceParts[0],
		Name:      namespaceParts[1],
	}, nil
}
