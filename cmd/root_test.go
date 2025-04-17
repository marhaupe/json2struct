package cmd

import "testing"

func TestRun(t *testing.T) {
	type args struct {
		inputString        string
		inputFile          string
		shouldBenchmark    bool
		shouldUseClipboard bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "should parse input string",
			args: args{
				inputString:        `{"name": "John", "age": 30}`,
				inputFile:          "",
				shouldBenchmark:    false,
				shouldUseClipboard: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Run(tt.args.inputString, tt.args.inputFile, tt.args.shouldBenchmark, tt.args.shouldUseClipboard)
		})
	}
}
