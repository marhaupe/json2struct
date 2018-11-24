package cmd

import (
	"fmt"

	"github.com/marhaupe/json-to-struct/internal/json"
)

func Start() {

	e := json.JSONObject{
		Root: true,
		Children: []json.JSONElement{
			&json.JSONObject{
				Key: "Intstringobj",
				Children: []json.JSONElement{
					&json.JSONPrimitive{
						Ptype: json.Int,
						Key:   "testint",
					},
					&json.JSONPrimitive{
						Ptype: json.String,
						Key:   "teststring",
					},
				},
			},
			&json.JSONObject{
				Key: "Boolobj",
				Children: []json.JSONElement{
					&json.JSONPrimitive{
						Ptype: json.Bool,
						Key:   "testbool",
					},
				},
			},
		},
	}
	fmt.Print(e.String())
}
