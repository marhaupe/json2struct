package json

import (
	"regexp"
	"strings"
)

type JSONElement interface {
	String() string
	Datatype() string
}

func listChildrenTypes(c []JSONElement) []string {
	foundChildrenTypes := make(map[string]bool)
	var foundChildren []string
	for _, entry := range c {
		foundChildrenTypes[entry.Datatype()] = true
	}
	for k := range foundChildrenTypes {
		foundChildren = append(foundChildren, k)
	}
	return foundChildren
}

func appendOmitEmptyToRootElement(s string) string {
	re := regexp.MustCompile("(?s)`json:\"(.*)\"`\n$")
	return re.ReplaceAllString(s, "`json:\"$1,omitempty\"`\n")
}

func capitalizeKey(k string) string {
	return strings.Title(k)
}
