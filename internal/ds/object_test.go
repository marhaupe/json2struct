package ds

import (
	"reflect"
	"testing"
)

func TestAddChild(t *testing.T) {
	obj := &JSONObject{
		Key: "Testobj",
	}
	objToAdd := &JSONObject{
		Key: "ToAdd",
	}

	// Looks like this:
	// {
	// 		"Testobj" : {
	// 			"ToAdd": {}
	// 		}
	// }
	obj.AddChild(objToAdd)
	if !reflect.DeepEqual(obj.Children[0], objToAdd) {
		t.Errorf("Added object to object the wrong way")
	}

	// Looks like this:
	// {
	// 		"Testobj" : {
	// 			"ToAdd": {
	// 					"Testint": 0,
	//      }
	// 		}
	// }
	objToAdd.AddChild(&JSONPrimitive{
		Key:   "Testint",
		Ptype: Int,
	})
	if !reflect.DeepEqual(obj.Children[0], objToAdd) {
		t.Errorf("Added primitive to nested object the wrong way")
	}

	// Looks like this:
	// {
	// 		"Testobj" : {
	// 			"ToAdd": {
	// 					"Testint": 0,
	//      },
	//			"ToAdd": {
	//			    "Testbool": false,
	//		  }
	// 		}
	// }
	secondObjToAdd := &JSONObject{
		Key: "ToAdd",
	}
	secondObjToAdd.AddChild(&JSONPrimitive{
		Key:   "Testbool",
		Ptype: Bool,
	})

	obj.AddChild(secondObjToAdd)

	mergedObj := &JSONObject{
		Key: "ToAdd",
	}
	mergedObj.AddChild(&JSONPrimitive{
		Key:   "Testint",
		Ptype: Int,
	})
	mergedObj.AddChild(&JSONPrimitive{
		Key:   "Testbool",
		Ptype: Bool,
	})

	if !reflect.DeepEqual(obj.Children[0], mergedObj) {
		t.Errorf("Failed to correctly merge objects\ngot:\n%v\nwanted:\n%v\n", mergedObj, obj.Children[0])
	}
}
