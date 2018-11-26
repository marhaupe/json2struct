package ds

import "testing"

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
			name: "Nested Object Test One",
			fields: fields{
				Key: "testobject",
				Children: []JSONElement{
					&JSONObject{
						Key: "nested_object",
						Children: []JSONElement{
							&JSONPrimitive{
								Ptype: Int,
								Key:   "testint",
							},
						},
					},
				},
			},
			want: "Testobject struct{\nNested_object struct{\nTestint int `json:\"testint\"`\n} `json:\"nested_object\"`\n} `json:\"testobject\"`\n",
		},
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
