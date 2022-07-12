package outputs

import (
	"fmt"

	"github.com/airplanedev/ojson"
	"github.com/airplanedev/path"
	"github.com/pkg/errors"
)

func ApplyOutputCommand(cmd *ParsedLine, o *ojson.Value) error {
	fmt.Printf(">>>>> %+v\n", cmd)

	switch cmd.Command {
	case "":
		if err := applyLegacy(cmd.Name, cmd.Value.V, o); err != nil {
			return err
		}

	case "set":
		if err := applySet(cmd.JsonPath, cmd.Value.V, o); err != nil {
			return err
		}

	case "append":
		if err := applyAppend(cmd.JsonPath, cmd.Value.V, o); err != nil {
			return err
		}

	default:
		return errors.New("unknown command")
	}

	return nil
}

func applyLegacy(name string, v interface{}, o *ojson.Value) error {
	fmt.Printf(">>>> Applying legacy with with name %s and value %+v", name, v)

	if o.V == nil {
		o.V = ojson.NewObject()
	}

	obj, ok := o.V.(*ojson.Object)
	if !ok {
		return errors.New("expected json object at top level")
	}

	target, ok := obj.Get(name)
	if !ok {
		target = []interface{}{}
	}

	arr, ok := target.([]interface{})
	if !ok {
		return errors.New("expected array")
	}

	obj.Set(name, append(arr, v))
	return nil
}

func applySet(jsPath string, v interface{}, o *ojson.Value) error {
	p, err := path.FromJS(jsPath)
	if err != nil {
		return err
	}

	loc, err := getLocation(p, o)
	if err != nil {
		return err
	}

	updateLocation(loc, v)
	return nil
}

func applyAppend(jsPath string, v interface{}, o *ojson.Value) error {
	p, err := path.FromJS(jsPath)
	if err != nil {
		return err
	}

	loc, err := getLocation(p, o)
	if err != nil {
		return err
	}

	var locArr []interface{}
	locVal := getAtLocation(loc)
	// If the append point is a null, we effectively "insert" an empty array at
	// that point first before appending to it.
	if locVal == nil {
		locVal = []interface{}{}
	}
	locArr, ok := locVal.([]interface{})
	if !ok {
		return errors.New("expected array at append point")
	}
	updateLocation(loc, append(locArr, v))
	return nil
}
