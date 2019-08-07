package generator

import "testing"

func Test_isNumber(t *testing.T) {
	type args struct {
		varname string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "floating",
			args: args{"1.1"},
			want: true,
		},
		{
			name: "negative floating",
			args: args{"-1.1"},
			want: true,
		},
		{
			name: "int",
			args: args{"1"},
			want: true,
		},
		{
			name: "negative int",
			args: args{"-1"},
			want: true,
		},
		{
			name: "not number",
			args: args{"xyz"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNumber(tt.args.varname); got != tt.want {
				t.Errorf("isNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
