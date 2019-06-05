// Package generate contains fields to generate a Go struct from a given JSON string
package generate

import (
	"go/format"
)

// Generate takes a JSON string s and parses it accordingly
// The resulting string is generated code for a Go struct that e.g. can be used to
// marshal and unmarshal JSON
func Generate(s string) (string, error) {
	return "", nil
}

func GenerateWithFormatting(s string) (string, error) {
	res, err := Generate(s)
	if err != nil {
		return "", err
	}
	formattedBytes, err := format.Source([]byte(res))
	if err != nil {
		return "", err
	}
	return string(formattedBytes), nil
}
