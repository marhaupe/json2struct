package json

import (
	"testing"
)

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
