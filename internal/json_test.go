package internal

import "testing"

func TestJSONPrimitiveInteger_String(t *testing.T) {
	type fields struct {
		JSONData JSONData
		Key      string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Basic Integer One",
			fields: fields{
				Key: "key",
			},
			want: "Key int `json:\"key\"`",
		},
		{
			name: "Basic Integer Two",
			fields: fields{
				Key: "anotherKey",
			},
			want: "AnotherKey int `json:\"anotherKey\"`",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jp := &JSONPrimitive{
				Ptype:    Integer,
				JSONData: tt.fields.JSONData,
				Key:      tt.fields.Key,
			}
			if got := jp.String(); got != tt.want {
				t.Errorf("JSONPrimitive.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONPrimitiveString_String(t *testing.T) {
	type fields struct {
		JSONData JSONData
		Key      string
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
			want: "Key string `json:\"key\"`",
		},
		{
			name: "Basic String Two",
			fields: fields{
				Key: "anotherKey",
			},
			want: "AnotherKey string `json:\"anotherKey\"`",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jp := &JSONPrimitive{
				Ptype:    String,
				JSONData: tt.fields.JSONData,
				Key:      tt.fields.Key,
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
		Children []JSONData
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
				Children: []JSONData{
					&JSONPrimitive{
						Ptype: Integer,
						Key:   "testint",
					},
					&JSONPrimitive{
						Ptype: String,
						Key:   "teststring",
					},
				},
			},
			want: "Testobject struct { Testint int `json:\"testint\"`Teststring string `json:\"teststring\"` } `json:\"testobject\"`",
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
		Children []JSONData
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Basic Array Test One",
			fields: fields{
				Key: "testarray",
				Children: []JSONData{
					&JSONPrimitive{
						Ptype: Integer,
						Key:   "testint1",
					},
					&JSONPrimitive{
						Ptype: Integer,
						Key:   "testint2",
					},
				},
			},
			want: ""
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

func TestJSONRoot_String(t *testing.T) {
	type fields struct {
		JSONData JSONData
		Children []JSONData
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jp := &JSONRoot{
				JSONData: tt.fields.JSONData,
				Children: tt.fields.Children,
			}
			if got := jp.String(); got != tt.want {
				t.Errorf("JSONRoot.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
