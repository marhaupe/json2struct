package ds

import "testing"

func TestJSONPrimitiveInt_String(t *testing.T) {
	type fields struct {
		Key string
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
				Ptype: Int,
				Key:   tt.fields.Key,
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
				Ptype: String,
				Key:   tt.fields.Key,
			}
			if got := jp.String(); got != tt.want {
				t.Errorf("JSONPrimitive.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
