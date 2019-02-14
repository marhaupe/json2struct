package generate

import (
	"testing"
)

func TestDetectInvalidJson(t *testing.T) {
	json := `{ "Key": value }`
	_, err := Generate(json)
	if err == nil {
		t.Error("Supplied an invalid JSON that has not been detected as such")
	}
}

func setupCases() map[string]string {
	cases := make(map[string]string)

	param1 := `{
					"src": "Images/Sun.png",
					"name": "sun1", 
					"hOffset": 250, 
					"vOffset": 250, 
					"alignment": "center",
					"visible": true
				}`
	expect1 := "type JSONToStruct struct{\n" +
		"Src string `json:\"src\"`\n" +
		"Name string `json:\"name\"`\n" +
		"HOffset int `json:\"hOffset\"`\n" +
		"VOffset int `json:\"vOffset\"`\n" +
		"Alignment string `json:\"alignment\"`\n" +
		"Visible bool `json:\"visible\"`\n" +
		"}\n"
	cases["expect1"] = param1
	cases["param1"] = expect1

	param2 := `{
    "glossary": {
        "title": "example glossary",
				"GlossDiv": {
            "title": "S",
						"GlossList": {
                "GlossEntry": {
                    "ID": "SGML",
										"SortAs": "SGML",
										"GlossTerm": "Standard Generalized Markup Language",
										"Acronym": "SGML",
										"Abbrev": "ISO 8879:1986",
										"GlossDef": {
                        "para": "A meta-markup language, used to create markup languages such as DocBook.",
												"GlossSeeAlso": ["GML", "XML"]
										},
										"GlossSee": "markup"
								}
            }
        }
    }
	}`
	expect2 := "type JSONToStruct struct{\n" +
		"Glossary struct{\n" +
		"Title string `json:\"title\"`\n" +
		"GlossDiv struct{\n" +
		"Title string `json:\"title\"`\n" +
		"GlossList struct{\n" +
		"GlossEntry struct{\n" +
		"ID string `json:\"ID\"`\n" +
		"SortAs string `json:\"SortAs\"`\n" +
		"GlossTerm string `json:\"GlossTerm\"`\n" +
		"Acronym string `json:\"Acronym\"`\n" +
		"Abbrev string `json:\"Abbrev\"`\n" +
		"GlossDef struct{\n" +
		"Para string `json:\"para\"`\n" +
		"GlossSeeAlso []string `json:\"GlossSeeAlso\"`\n" +
		"} `json:\"GlossDef\"`\n" +
		"GlossSee string `json:\"GlossSee\"`\n" +
		"} `json:\"GlossEntry\"`\n" +
		"} `json:\"GlossList\"`\n" +
		"} `json:\"GlossDiv\"`\n" +
		"} `json:\"glossary\"`\n" +
		"}\n"
	cases["expect2"] = param2
	cases["param2"] = expect2

	param3 := `{
		"Testobj": {
			"Teststring": "Hey"
		},
		"Testobj": {
			"Testbool": true
		}
	}`
	expect3 := "type JSONToStruct struct{\n" +
		"Testobj struct{\n" +
		"Teststring string `json:\"Teststring\"`\n" +
		"Testbool bool `json:\"Testbool\"`\n" +
		"} `json:\"Testobj\"`\n" +
		"}\n"
	cases["expect3"] = param3
	cases["param3"] = expect3

	param4 := `[
		{
			"teststring": "Hey"
		},
		{
			"testbool": true
		}
	]`
	expect4 := "type JSONToStruct []struct{\n" +
		"Teststring string `json:\"teststring\"`\n" +
		"Testbool bool `json:\"testbool\"`\n" +
		"}\n"
	cases["expect4"] = param4
	cases["param4"] = expect4

	param5 := `[
  {
    "thissucks": true
  },
  {
    "thisdoesntsuck": {
          "value": false
    }
  }
	]`
	expect5 := "type JSONToStruct []struct{\n" +
		"Thissucks bool `json:\"thissucks\"`\n" +
		"Thisdoesntsuck struct{\n" +
		"Value bool `json:\"value\"`\n" +
		"} `json:\"thisdoesntsuck\"`\n" +
		"}\n"
	cases["expect5"] = param5
	cases["param5"] = expect5

	param6 := `{
		"NestedObj": {
			"Testbool": true
		}
	}`
	expect6 := "type JSONToStruct struct{\n" +
		"NestedObj struct{\n" +
		"Testbool bool `json:\"Testbool\"`\n" +
		"} `json:\"NestedObj\"`\n" +
		"}\n"
	cases["expect6"] = param6
	cases["param6"] = expect6

	param7 := `{
		"Testfloat": 56.7,
		"Testint": 56
	}`
	expect7 := "type JSONToStruct struct{\n" +
		"Testfloat float64 `json:\"Testfloat\"`\n" +
		"Testint int `json:\"Testint\"`\n" +
		"}\n"
	cases["expect7"] = param7
	cases["param7"] = expect7
	return cases
}
func TestVariousCombinations(t *testing.T) {
	type args struct {
		s string
	}
	cases := setupCases()
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Case One",
			want: cases["param1"],
			args: args{
				s: cases["expect1"],
			},
		},
		{
			name: "Case Two",
			want: cases["param2"],
			args: args{
				s: cases["expect2"],
			},
		},
		{
			name: "Case Three",
			want: cases["param3"],
			args: args{
				s: cases["expect3"],
			},
		},
		{
			name: "Case Four",
			want: cases["param4"],
			args: args{
				s: cases["expect4"],
			},
		},
		{
			name: "Case Five",
			want: cases["param5"],
			args: args{
				s: cases["expect5"],
			},
		},
		{
			name: "Case Six",
			want: cases["param6"],
			args: args{
				s: cases["expect6"],
			},
		},
		{
			name: "Case Seven",
			want: cases["param7"],
			args: args{
				s: cases["expect7"],
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := Generate(tt.args.s)
			if got != tt.want {
				t.Errorf("Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestObjectRootWithPrimitives(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "With strings",
			args: args{
				`{
				"Test": "Here"
			 }`,
			},
			want: "type JSONToStruct struct{\n" +
				"Test string `json:\"Test\"`\n" +
				"}\n",
		},
		{
			name: "With ints",
			args: args{
				`{
				"Test": 500
			 }`,
			},
			want: "type JSONToStruct struct{\n" +
				"Test int `json:\"Test\"`\n" +
				"}\n",
		},
		{
			name: "With floats",
			args: args{
				`{
				"Test": 500.1
			 }`,
			},
			want: "type JSONToStruct struct{\n" +
				"Test float64 `json:\"Test\"`\n" +
				"}\n",
		},
		{
			name: "With Nulls",
			args: args{
				`{
				"Test": null
			 }`,
			},
			want: "type JSONToStruct struct{\n" +
				"Test interface{} `json:\"Test\"`\n" +
				"}\n",
		},
		{
			name: "With bools",
			args: args{
				`{
				"Test": true
			 }`,
			},
			want: "type JSONToStruct struct{\n" +
				"Test bool `json:\"Test\"`\n" +
				"}\n",
		},
		{
			name: "With every primitive",
			args: args{
				`{
				"Testbool": true,
				"Teststring": "here",
				"Testfloat": 500.1,
				"Testint": 500
			 }`,
			},
			want: "type JSONToStruct struct{\n" +
				"Testbool bool `json:\"Testbool\"`\n" +
				"Teststring string `json:\"Teststring\"`\n" +
				"Testfloat float64 `json:\"Testfloat\"`\n" +
				"Testint int `json:\"Testint\"`\n" +
				"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := Generate(tt.args.s)
			if got != tt.want {
				t.Errorf("Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestObjectRootWithObjects(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "With one object",
			args: args{
				`{
					"Testobj": {
						"Teststring": "Here",
						"Testbool": true,
						"Testint": 500,
						"Testfloat": 500.1
					}
				}`,
			},
			want: "type JSONToStruct struct{\n" +
				"Testobj struct{\n" +
				"Teststring string `json:\"Teststring\"`\n" +
				"Testbool bool `json:\"Testbool\"`\n" +
				"Testint int `json:\"Testint\"`\n" +
				"Testfloat float64 `json:\"Testfloat\"`\n" +
				"} `json:\"Testobj\"`\n" +
				"}\n",
		},
		{
			name: "With two objects",
			args: args{
				`{
					"Testobj": {
						"Teststring": "Here",
						"Testbool": true
					},
					"Testobj": {
						"Testint": 500,
						"Testfloat": 500.1
					}
				}`,
			},
			want: "type JSONToStruct struct{\n" +
				"Testobj struct{\n" +
				"Teststring string `json:\"Teststring\"`\n" +
				"Testbool bool `json:\"Testbool\"`\n" +
				"Testint int `json:\"Testint\"`\n" +
				"Testfloat float64 `json:\"Testfloat\"`\n" +
				"} `json:\"Testobj\"`\n" +
				"}\n",
		},
		{
			name: "With two objects with the same key",
			args: args{
				`{
					"Testobj": {
						"Teststring": "Here",
						"Testbool": true,
						"Testint": 500,
						"Testfloat": 500.1
					}
				}`,
			},
			want: "type JSONToStruct struct{\n" +
				"Testobj struct{\n" +
				"Teststring string `json:\"Teststring\"`\n" +
				"Testbool bool `json:\"Testbool\"`\n" +
				"Testint int `json:\"Testint\"`\n" +
				"Testfloat float64 `json:\"Testfloat\"`\n" +
				"} `json:\"Testobj\"`\n" +
				"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := Generate(tt.args.s)
			if got != tt.want {
				t.Errorf("Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayRootWithObjects(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "With Similar Objects",
			args: args{
				`[
					{
						"Testbool": true,
						"Testint": 500
					},
					{
						"Testbool": true,
						"Testint": 500
					}
				 ]`,
			},
			want: "type JSONToStruct []struct{\n" +
				"Testbool bool `json:\"Testbool\"`\n" +
				"Testint int `json:\"Testint\"`\n" +
				"}\n",
		},
		{
			name: "With Floats",
			args: args{
				`[
					500.1,
					600.1,
					300.1
				 ]`,
			},
			want: "type JSONToStruct []float64",
		},
		{
			name: "With Bools",
			args: args{
				`[
					true,
					false,
					true
				 ]`,
			},
			want: "type JSONToStruct []bool",
		},
		{
			name: "With Strings",
			args: args{
				`[
					"500",
					"600",
					"300"
				 ]`,
			},
			want: "type JSONToStruct []string",
		},
		{
			name: "With Nulls",
			args: args{
				`[
					null,
					null
				 ]`,
			},
			want: "type JSONToStruct []interface{}",
		},
		{
			name: "With Ints",
			args: args{
				`[
					500,
					600,
					300
				 ]`,
			},
			want: "type JSONToStruct []int",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := Generate(tt.args.s)
			if got != tt.want {
				t.Errorf("Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestArrayRootWithMultipleTypes(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "With Multiples",
			args: args{
				`[
					"500",
					600
				 ]`,
			},
			want: "type JSONToStruct []interface{}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := Generate(tt.args.s)
			if got != tt.want {
				t.Errorf("Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestArrayRootWithPrimitives(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "With Strings",
			args: args{
				`[
					"500",
					"600",
					"300"
				 ]`,
			},
			want: "type JSONToStruct []string",
		},
		{
			name: "With Floats",
			args: args{
				`[
					500.1,
					600.1,
					300.1
				 ]`,
			},
			want: "type JSONToStruct []float64",
		},
		{
			name: "With Bools",
			args: args{
				`[
					true,
					false,
					true
				 ]`,
			},
			want: "type JSONToStruct []bool",
		},
		{
			name: "With Ints",
			args: args{
				`[
					500,
					600,
					300
				 ]`,
			},
			want: "type JSONToStruct []int",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := Generate(tt.args.s)
			if got != tt.want {
				t.Errorf("Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}
