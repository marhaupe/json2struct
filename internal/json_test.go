package internal

import (
	"testing"
)

func TestJSONPrimitiveInt_String(t *testing.T) {
	type fields struct {
		JSONElement JSONElement
		Key         string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Basic Int One",
			fields: fields{
				Key: "key",
			},
			want: "Key int `json:\"key\"`\n",
		},
		{
			name: "Basic Int Two",
			fields: fields{
				Key: "anotherKey",
			},
			want: "AnotherKey int `json:\"anotherKey\"`\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jp := &JSONPrimitive{
				Ptype:       Int,
				JSONElement: tt.fields.JSONElement,
				Key:         tt.fields.Key,
			}
			if got := jp.String(); got != tt.want {
				t.Errorf("JSONPrimitive.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONPrimitiveString_String(t *testing.T) {
	type fields struct {
		JSONElement JSONElement
		Key         string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Basic String One",
			fields: fields{
				Key: "key",
			},
			want: "Key string `json:\"key\"`\n",
		},
		{
			name: "Basic String Two",
			fields: fields{
				Key: "anotherKey",
			},
			want: "AnotherKey string `json:\"anotherKey\"`\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jp := &JSONPrimitive{
				Ptype:       String,
				JSONElement: tt.fields.JSONElement,
				Key:         tt.fields.Key,
			}
			if got := jp.String(); got != tt.want {
				t.Errorf("JSONPrimitive.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONObject_String(t *testing.T) {
	type fields struct {
		Key      string
		Children []JSONElement
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Basic Object Test One",
			fields: fields{
				Key: "testobject",
				Children: []JSONElement{
					&JSONPrimitive{
						Ptype: Int,
						Key:   "testint",
					},
					&JSONPrimitive{
						Ptype: String,
						Key:   "teststring",
					},
				},
			},
			want: "Testobject struct{\nTestint int `json:\"testint\"`\nTeststring string `json:\"teststring\"`\n} `json:\"testobject\"`\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jp := &JSONObject{
				Key:      tt.fields.Key,
				Children: tt.fields.Children,
			}
			if got := jp.String(); got != tt.want {
				t.Errorf("JSONObject.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONArray_String(t *testing.T) {
	type fields struct {
		Key      string
		Children []JSONElement
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Array Test Only Bools",
			fields: fields{
				Key: "testarray",
				Children: []JSONElement{
					&JSONPrimitive{
						Ptype: Bool,
						Key:   "testbool",
					},
					&JSONPrimitive{
						Ptype: Bool,
						Key:   "testbool2",
					},
				},
			},
			want: "Testarray []bool `json:\"testarray\"`\n",
		},
		{
			name: "Array Test Only Strings",
			fields: fields{
				Key: "testarray",
				Children: []JSONElement{
					&JSONPrimitive{
						Ptype: String,
						Key:   "teststring1",
					},
					&JSONPrimitive{
						Ptype: String,
						Key:   "teststring2",
					},
				},
			},
			want: "Testarray []string `json:\"testarray\"`\n",
		},
		{
			name: "Array Test Only Ints",
			fields: fields{
				Key: "testarray",
				Children: []JSONElement{
					&JSONPrimitive{
						Ptype: Int,
						Key:   "testint1",
					},
					&JSONPrimitive{
						Ptype: Int,
						Key:   "testint2",
					},
				},
			},
			want: "Testarray []int `json:\"testarray\"`\n",
		},
		{
			name: "Array Test Arrays",
			fields: fields{
				Key: "testarray",
				Children: []JSONElement{
					&JSONArray{
						Key: "Doesn't matter",
						Children: []JSONElement{
							&JSONPrimitive{
								Ptype: Bool,
								Key:   "testbool",
							},
							&JSONPrimitive{
								Ptype: Int,
								Key:   "testint",
							},
						}},
					&JSONArray{
						Key: "Doesn't matter aswell",
						Children: []JSONElement{
							&JSONPrimitive{
								Ptype: Bool,
								Key:   "testbool",
							},
							&JSONPrimitive{
								Ptype: Int,
								Key:   "testint",
							},
						},
					},
				},
			},
			want: "Testarray [][]interface{} `json:\"testarray\"`\n",
		},
		{
			name: "Array Test Different Objects",
			fields: fields{
				Key: "testarray",
				Children: []JSONElement{
					&JSONObject{
						Key: "Intstringobj",
						Children: []JSONElement{
							&JSONPrimitive{
								Ptype: Int,
								Key:   "testint",
							},
							&JSONPrimitive{
								Ptype: String,
								Key:   "teststring",
							},
						},
					},
					&JSONObject{
						Key: "Boolobj",
						Children: []JSONElement{
							&JSONPrimitive{
								Ptype: Bool,
								Key:   "testbool",
							},
						},
					},
				},
			},
			want: "Testarray []struct{\nIntstringobj struct{\nTestint int `json:\"testint\"`\nTeststring string `json:\"teststring\"`\n} `json:\"Intstringobj,omitempty\"`\nBoolobj struct{\nTestbool bool `json:\"testbool\"`\n} `json:\"Boolobj,omitempty\"`\n} `json:\"testarray\"`\n",
		},
		{
			name: "Array Test Different Types Primitive",
			fields: fields{
				Key: "testarray",
				Children: []JSONElement{
					&JSONPrimitive{
						Ptype: Bool,
						Key:   "testbool",
					},
					&JSONPrimitive{
						Ptype: Int,
						Key:   "testint",
					},
				},
			},
			want: "Testarray []interface{} `json:\"testarray\"`\n",
		},
		{
			name: "Array Test Different Types Mixed",
			fields: fields{
				Key: "testarray",
				Children: []JSONElement{
					&JSONPrimitive{
						Ptype: Bool,
						Key:   "testbool",
					},
					&JSONObject{
						Key: "Doesn't matter",
						Children: []JSONElement{
							&JSONPrimitive{
								Ptype: Int,
								Key:   "testint",
							},
							&JSONPrimitive{
								Ptype: String,
								Key:   "teststring",
							},
						},
					},
				},
			},
			want: "Testarray []interface{} `json:\"testarray\"`\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jp := &JSONArray{
				Key:      tt.fields.Key,
				Children: tt.fields.Children,
			}
			if got := jp.String(); got != tt.want {
				t.Errorf("JSONArray.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
func Test_appendOmitEmptyToRootElement(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Simple",
			args: args{
				"Test []struct{} `json:\"test\"`\n",
			},
			want: "Test []struct{} `json:\"test,omitempty\"`\n",
		},
		{
			name: "Nested",
			args: args{
				"Intstringobj struct{\nTestint int `json:\"testint\"`\nTeststring string `json:\"teststring\"`\n} `json:\"Intstringobj\"`\n",
			},
			want: "Intstringobj struct{\nTestint int `json:\"testint\"`\nTeststring string `json:\"teststring\"`\n} `json:\"Intstringobj,omitempty\"`\n",
		},
		{
			name: "Nested 2",
			args: args{
				"Test []struct{\nTestint int `json:\"testint\"`\n} `json:\"test\"`\n",
			},
			want: "Test []struct{\nTestint int `json:\"testint\"`\n} `json:\"test,omitempty\"`\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := appendOmitEmptyToRootElement(tt.args.s); got != tt.want {
				t.Errorf("appendOmitEmptyToRootElement() = %v, want %v", got, tt.want)
			}
		})
	}
}
