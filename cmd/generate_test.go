package cmd

import (
	"bytes"
	"testing"
)

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
	var expect1 bytes.Buffer
	expect1.WriteString("type JSONToStruct struct{\n")
	expect1.WriteString("Src string `json:\"src\"`\n")
	expect1.WriteString("Name string `json:\"name\"`\n")
	expect1.WriteString("HOffset int `json:\"hOffset\"`\n")
	expect1.WriteString("VOffset int `json:\"vOffset\"`\n")
	expect1.WriteString("Alignment string `json:\"alignment\"`\n")
	expect1.WriteString("Visible bool `json:\"visible\"`\n")
	expect1.WriteString("}\n")
	cases["expect1"] = param1
	cases["param1"] = expect1.String()

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
	var expect2 bytes.Buffer
	expect2.WriteString("type JSONToStruct struct{\n")
	expect2.WriteString("Glossary struct{\n")
	expect2.WriteString("Title string `json:\"title\"`\n")
	expect2.WriteString("GlossDiv struct{\n")
	expect2.WriteString("Title string `json:\"title\"`\n")
	expect2.WriteString("GlossList struct{\n")
	expect2.WriteString("GlossEntry struct{\n")
	expect2.WriteString("ID string `json:\"ID\"`\n")
	expect2.WriteString("SortAs string `json:\"SortAs\"`\n")
	expect2.WriteString("GlossTerm string `json:\"GlossTerm\"`\n")
	expect2.WriteString("Acronym string `json:\"Acronym\"`\n")
	expect2.WriteString("Abbrev string `json:\"Abbrev\"`\n")
	expect2.WriteString("GlossDef struct{\n")
	expect2.WriteString("Para string `json:\"para\"`\n")
	expect2.WriteString("GlossSeeAlso []string `json:\"GlossSeeAlso\"`\n")
	expect2.WriteString("} `json:\"GlossDef\"`\n")
	expect2.WriteString("GlossSee string `json:\"GlossSee\"`\n")
	expect2.WriteString("} `json:\"GlossEntry\"`\n")
	expect2.WriteString("} `json:\"GlossList\"`\n")
	expect2.WriteString("} `json:\"GlossDiv\"`\n")
	expect2.WriteString("} `json:\"glossary\"`\n")
	expect2.WriteString("}\n")
	cases["expect2"] = param2
	cases["param2"] = expect2.String()

	param3 := `{
		"Testobj": {
			"Teststring": "Hey"
		},
		"Testobj": {
			"Testbool": true
		}
	}`

	var expect3 bytes.Buffer
	expect3.WriteString("type JSONToStruct struct{\n")
	expect3.WriteString("Testobj struct{\n")
	expect3.WriteString("Teststring string `json:\"Teststring\"`\n")
	expect3.WriteString("Testbool bool `json:\"Testbool\"`\n")
	expect3.WriteString("} `json:\"Testobj\"`\n")
	expect3.WriteString("}\n")

	cases["expect3"] = param3
	cases["param3"] = expect3.String()

	param4 := `[
		{
			"teststring": "Hey"
		},
		{
			"testbool": true
		}
	]`
	var expect4 bytes.Buffer
	expect4.WriteString("type JSONToStruct []struct{\n")
	expect4.WriteString("Teststring string `json:\"teststring,omitempty\"`\n")
	expect4.WriteString("Testbool bool `json:\"testbool,omitempty\"`\n")
	expect4.WriteString("}\n")

	cases["expect4"] = param4
	cases["param4"] = expect4.String()

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

	var expect5 bytes.Buffer
	expect5.WriteString("type JsonToStruct []struct{\n")
	expect5.WriteString("Thissucks bool `json:\"thissucks,omitempty\"`\n")
	expect5.WriteString("Thisdoesntsuck struct{\n")
	expect5.WriteString("Value bool `json:\"value\"`\n")
	expect5.WriteString("} `json:\"thisdoesntsuck,omitempty\"`\n")
	expect5.WriteString("}\n")

	cases["expect5"] = param5
	cases["param5"] = expect5.String()

	return cases
}
func TestGenerate(t *testing.T) {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Generate(tt.args.s); got != tt.want {
				t.Errorf("Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}
