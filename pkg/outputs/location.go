package outputs

import (
	"github.com/airplanedev/ojson"
	"github.com/airplanedev/path"
	"github.com/pkg/errors"
)

type rootLocation ojson.Value

type objLocation struct {
	Key string
	Obj *ojson.Object
}
type arrLocation struct {
	Key int
	Arr *[]interface{}
}

type location struct {
	Root *rootLocation
	Obj  *objLocation
	Arr  *arrLocation
}

func getAtLocation(loc location) interface{} {
	if loc.Root != nil {
		return loc.Root.V
	} else if loc.Obj != nil {
		v, _ := loc.Obj.Obj.Get(loc.Obj.Key)
		return v
	} else if loc.Arr != nil {
		return (*loc.Arr.Arr)[loc.Arr.Key]
	} else {
		return nil
	}
}

func updateLocation(loc location, v interface{}) interface{} {
	if loc.Root != nil {
		loc.Root.V = v
	} else if loc.Obj != nil {
		loc.Obj.Obj.Set(loc.Obj.Key, v)
	} else if loc.Arr != nil {
		(*loc.Arr.Arr)[loc.Arr.Key] = v
	}
	return v
}

// getLocation returns a location at path p, usually for modifying.
// Note that this will create new `ojson.Object`s along the path if nulls are
// encountered. For example, trying to access "a.b" in an empty object would
// insert another empty object at key "a".
func getLocation(p path.P, root *ojson.Value) (location, error) {
	loc := location{
		Root: (*rootLocation)(root),
	}
	var cur interface{}
	cur = root.V
	for _, component := range p.Components() {
		switch c := component.(type) {
		case string:
			obj, ok := cur.(*ojson.Object)
			if !ok {
				if cur == nil {
					updateLocation(loc, ojson.NewObject())
					obj = getAtLocation(loc).(*ojson.Object)
				} else {
					return location{}, errors.New("expected *ojson.Object")
				}
			}
			cur, _ = obj.Get(c)
			loc = location{
				Obj: &objLocation{
					Key: c,
					Obj: obj,
				},
			}

		case int:
			arr, ok := cur.([]interface{})
			if !ok {
				return location{}, errors.New("expected array")
			}
			if c >= len(arr) {
				return location{}, errors.New("array had too few elements")
			}
			cur = arr[c]
			loc = location{
				Arr: &arrLocation{
					Key: c,
					Arr: &arr,
				},
			}

		default:
			return location{}, errors.New("unexpected component type")
		}
	}
	return loc, nil
}
