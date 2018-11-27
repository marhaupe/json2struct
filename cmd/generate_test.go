package cmd

import (
	"bytes"
	"testing"
)

func setupCases() map[string]string {
	cases := make(map[string]string)

	want1 := `{
					"src": "Images/Sun.png",
					"name": "sun1", 
					"hOffset": 250, 
					"vOffset": 250, 
					"alignment": "center",
					"visible": true
				}`
	var have1 bytes.Buffer
	have1.WriteString("type JSONToStruct struct{\n")
	have1.WriteString("Src string `json:\"src\"`\n")
	have1.WriteString("Name string `json:\"name\"`\n")
	have1.WriteString("HOffset int `json:\"hOffset\"`\n")
	have1.WriteString("VOffset int `json:\"vOffset\"`\n")
	have1.WriteString("Alignment string `json:\"alignment\"`\n")
	have1.WriteString("Visible bool `json:\"visible\"`\n")
	have1.WriteString("}\n")
	cases["want1"] = want1
	cases["have1"] = have1.String()

	want2 := `{
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
	var have2 bytes.Buffer
	have2.WriteString("type JSONToStruct struct{\n")
	have2.WriteString("Glossary struct{\n")
	have2.WriteString("Title string `json:\"title\"`\n")
	have2.WriteString("GlossDiv struct{\n")
	have2.WriteString("Title string `json:\"title\"`\n")
	have2.WriteString("GlossList struct{\n")
	have2.WriteString("GlossEntry struct{\n")
	have2.WriteString("ID string `json:\"ID\"`\n")
	have2.WriteString("SortAs string `json:\"SortAs\"`\n")
	have2.WriteString("GlossTerm string `json:\"GlossTerm\"`\n")
	have2.WriteString("Acronym string `json:\"Acronym\"`\n")
	have2.WriteString("Abbrev string `json:\"Abbrev\"`\n")
	have2.WriteString("GlossDef struct{\n")
	have2.WriteString("Para string `json:\"para\"`\n")
	have2.WriteString("GlossSeeAlso []string `json:\"GlossSeeAlso\"`\n")
	have2.WriteString("} `json:\"GlossDef\"`\n")
	have2.WriteString("GlossSee string `json:\"GlossSee\"`\n")
	have2.WriteString("} `json:\"GlossEntry\"`\n")
	have2.WriteString("} `json:\"GlossList\"`\n")
	have2.WriteString("} `json:\"GlossDiv\"`\n")
	have2.WriteString("} `json:\"glossary\"`\n")
	have2.WriteString("}\n")
	cases["want2"] = want2
	cases["have2"] = have2.String()

	want3 := `{
		"Testobj": {
			"Teststring": "Hey"
		},
		"Testobj": {
			"Testbool": true
		}
	}`
	var have3 bytes.Buffer
	have3.WriteString("type JSONToStruct struct{\n")
	have3.WriteString("Testobj struct{\n")
	have3.WriteString("Teststring string `json:\"Teststring\"`\n")
	have3.WriteString("} `json:\"Testobj\"`\n")
	have3.WriteString("}\n")

	cases["want3"] = want3
	cases["have3"] = have3.String()

	want4 := `[
		{
			"teststring": "Hey"
		},
		{
			"testbool": true
		}
	]`
	var have4 bytes.Buffer
	have4.WriteString("type JSONToStruct []struct{\n")
	have4.WriteString("Teststring string `json:\"teststring,omitempty\"`\n")
	have4.WriteString("Testbool bool `json:\"testbool,omitempty\"`\n")
	have4.WriteString("}\n")

	cases["want4"] = want4
	cases["have4"] = have4.String()

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
			want: cases["have1"],
			args: args{
				s: cases["want1"],
			},
		},
		{
			name: "Case Two",
			want: cases["have2"],
			args: args{
				s: cases["want2"],
			},
		},
		{
			name: "Case Three",
			want: cases["have3"],
			args: args{
				s: cases["want3"],
			},
		},
		{
			name: "Case Four",
			want: cases["have4"],
			args: args{
				s: cases["want4"],
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
