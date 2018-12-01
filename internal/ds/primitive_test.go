package ds

import "testing"

func TestJSONPrimitive_Panic(t *testing.T) {

	type fields struct {
		Key   string
		Ptype PrimitiveType
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Basic Int",
			fields: fields{
				Key:   "key",
				Ptype: -1,
			},
			want: "Key int `json:\"key\"`\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jp := &JSONPrimitive{
				Ptype: tt.fields.Ptype,
				Key:   tt.fields.Key,
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
		Key   string
		Ptype PrimitiveType
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Basic Int",
			fields: fields{
				Key:   "key",
				Ptype: Int,
			},
			want: "Key int `json:\"key\"`\n",
		},
		{
			name: "Basic Bool",
			fields: fields{
				Key:   "key",
				Ptype: Bool,
			},
			want: "Key bool `json:\"key\"`\n",
		},
		{
			name: "Basic String",
			fields: fields{
				Key:   "key",
				Ptype: String,
			},
			want: "Key string `json:\"key\"`\n",
		},
		{
			name: "Basic Float",
			fields: fields{
				Key:   "key",
				Ptype: Float,
			},
			want: "Key float64 `json:\"key\"`\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jp := &JSONPrimitive{
				Ptype: tt.fields.Ptype,
				Key:   tt.fields.Key,
			}
			if got := jp.String(); got != tt.want {
				t.Errorf("JSONPrimitive.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
