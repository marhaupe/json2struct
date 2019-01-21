package tree

import (
	"reflect"
	"testing"
)

func TestJSONPrimitive_Panic(t *testing.T) {

	type fields struct {
		Key  string
		Type Datatype
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Basic Int",
			fields: fields{
				Key:  "key",
				Type: -1,
			},
			want: "Key int `json:\"key\"`\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jp := &JSONPrimitive{
				Datatype: tt.fields.Type,
				Key:      tt.fields.Key,
			}
			defer func() {
				r := recover()
				if r == nil {
					t.Error("Program should have panicked because of invalid Ptype")
				}
			}()
			jp.String()
		})
	}
}
func TestJSONPrimitive(t *testing.T) {
	type fields struct {
		Key  string
		Type Datatype
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Basic Int",
			fields: fields{
				Key:  "key",
				Type: Int,
			},
			want: "Key int `json:\"key\"`\n",
		},
		{
			name: "Basic Bool",
			fields: fields{
				Key:  "key",
				Type: Bool,
			},
			want: "Key bool `json:\"key\"`\n",
		},
		{
			name: "Basic String",
			fields: fields{
				Key:  "key",
				Type: String,
			},
			want: "Key string `json:\"key\"`\n",
		},
		{
			name: "Basic Null",
			fields: fields{
				Key:  "key",
				Type: Null,
			},
			want: "Key interface{} `json:\"key\"`\n",
		},
		{
			name: "Basic Float",
			fields: fields{
				Key:  "key",
				Type: Float,
			},
			want: "Key float64 `json:\"key\"`\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jp := &JSONPrimitive{
				Datatype: tt.fields.Type,
				Key:      tt.fields.Key,
			}
			if got := jp.String(); got != tt.want {
				t.Errorf("JSONPrimitive.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetParentPrimitives(t *testing.T) {
	prim := &JSONPrimitive{
		Key:      "Test",
		Datatype: String,
	}
	parent := &JSONObject{
		Key: "Testobject",
	}
	parent.AddChild(prim)
	if !reflect.DeepEqual(prim.GetParent(), parent) {
		t.Error("Parent has not been set or retrieved properly")
	}
}
