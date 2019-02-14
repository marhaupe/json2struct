// json2struct is a simple cli utility that can be used to parse JSON into Go structs.
// Behind the scenes, the tool builds an ast of the JSON and generates the
// resulting Go code accordingly.
package main

import (
	"github.com/marhaupe/json2struct/cmd"
)

func main() {
	cmd.Execute()
}
