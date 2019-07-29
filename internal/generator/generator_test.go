package generator

import (
	"reflect"
	"testing"

	"github.com/dave/jennifer/jen"
	"github.com/marhaupe/json2struct/internal/parse"
)

func Test_makePrimitiveDecl(t *testing.T) {
	type args struct {
		typ     parse.NodeType
		varname string
	}
	tests := []struct {
		name string
		args args
		want *jen.Statement
	}{
		{
			name: "string",
			args: args{
				typ:     parse.NodeTypeString,
				varname: "teststring",
			},
			want: jen.Id("Teststring").String().Tag(map[string]string{"json": "teststring"}),
		},
		{
			name: "number",
			args: args{
				typ:     parse.NodeTypeNumber,
				varname: "testnumber",
			},
			want: jen.Id("Testnumber").Float64().Tag(map[string]string{"json": "testnumber"}),
		},
		{
			name: "nil",
			args: args{
				typ:     parse.NodeTypeNil,
				varname: "testnil",
			},
			want: jen.Id("Testnil").Interface().Tag(map[string]string{"json": "testnil"}),
		},
		{
			name: "bool",
			args: args{
				typ:     parse.NodeTypeBool,
				varname: "testbool",
			},
			want: jen.Id("Testbool").Bool().Tag(map[string]string{"json": "testbool"}),
		},
		{
			name: "weird identifier",
			args: args{
				typ:     parse.NodeTypeString,
				varname: "tEsTStRiNG",
			},
			want: jen.Id("Teststring").String().Tag(map[string]string{"json": "tEsTStRiNG"}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makePrimitiveDecl(tt.args.typ, tt.args.varname); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makePrimitiveDecl() = %v, want %v", got, tt.want)
			}
		})
	}
}
