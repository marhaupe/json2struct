package ds

import (
	"strings"
)

// func appendOmitEmptyToRootElement(s string) string {
// 	re := regexp.MustCompile("(?s)`json:\"(.*)\"`\n$")
// 	return re.ReplaceAllString(s, "`json:\"$1,omitempty\"`\n")
// }

func capitalizeKey(k string) string {
	return strings.Title(k)
}
