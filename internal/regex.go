package internal

import (
	"regexp"
)

var ObjectRoot *regexp.Regexp
var ArrayRoot *regexp.Regexp
var ObjectValue *regexp.Regexp
var ArrayValue *regexp.Regexp
var IntValue *regexp.Regexp
var StringValue *regexp.Regexp
var BoolValue *regexp.Regexp
var NullValue *regexp.Regexp

func init() {
	ObjectRoot := regexp.MustCompile("(?s)^{(.*?)}$")
	ArrayRoot := regexp.MustCompile("(?s)^\\[(.*?)\\]$")
	ObjectValue := regexp.MustCompile("(?s)\"(.*?)\":\\s*({.*})")
	ArrayValue := regexp.MustCompile("(?s)\"(.*?)\":\\s*([.*])")
	IntValue := regexp.MustCompile("(?s)\"(.*?)\":\\s*(\\d*)")
	StringValue := regexp.MustCompile("(?s)\"(.*?)\":\\s*\"(.*)\"")
	BoolValue := regexp.MustCompile("(?s)\"(.*?)\":\\s*(true|false)")
	NullValue := regexp.MustCompile("(?s)\"(.*?)\":\\s*(null)")
}
