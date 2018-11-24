package cmd

import (
	"fmt"

	"github.com/marhaupe/json-to-struct/internal"
)

func Start() {

	e := internal.JSONObject{
		Root: true,
		Children: []internal.JSONElement{
			&internal.JSONObject{
				Key: "Intstringobj",
				Children: []internal.JSONElement{
					&internal.JSONPrimitive{
						Ptype: internal.Int,
						Key:   "testint",
					},
					&internal.JSONPrimitive{
						Ptype: internal.String,
						Key:   "teststring",
					},
				},
			},
			&internal.JSONObject{
				Key: "Boolobj",
				Children: []internal.JSONElement{
					&internal.JSONPrimitive{
						Ptype: internal.Bool,
						Key:   "testbool",
					},
				},
			},
		},
	}
	fmt.Print(e.String())
}
