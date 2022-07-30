package generator

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/marhaupe/json2struct/pkg/parse"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/dave/jennifer/jen"
)

func GenerateOutputFromString(s string) (string, error) {
	generatedFile, err := GenerateFileFromString(s)
	if err != nil {
		return "", err
	}

	buf := &bytes.Buffer{}
	err = generatedFile.Render(buf)
	if err != nil {
		return "", fmt.Errorf("error rendering file: %v", err)
	}
	return buf.String(), nil
}

func GenerateFileFromString(s string) (*jen.File, error) {
	node, err := parse.ParseFromString("json2struct", s)
	if err != nil {
		return nil, err
	}
	return GenerateFileFromAST(node)
}

func GenerateFileFromAST(tree parse.Node) (*jen.File, error) {
	g := Generator{
		Tree:        tree,
		currentNode: tree,
		file:        jen.NewFile("generated"),
	}
	return g.start()
}

type Generator struct {
	Tree        parse.Node
	currentNode parse.Node

	file        *jen.File
	currentStmt *jen.Statement
}

func (g Generator) start() (file *jen.File, err error) {
	defer func() {
		if r := recover(); r != nil {
			file = nil
			err = errors.New(fmt.Sprint(r))
			return
		}
	}()

	rootStmt := g.file.Type().Id("JSONToStruct")

	switch g.Tree.Type() {
	case parse.NodeTypeArray:
		casted := g.Tree.(*parse.ArrayNode)
		rootStmt.Add(makeArray(casted))
	case parse.NodeTypeObject:
		casted := g.Tree.(*parse.ObjectNode)
		rootStmt.Add(makeStruct(casted))
	default:
		panic("invalid json. expected { or [ as initial node but received something else")
	}

	return g.file, nil
}

func makeArray(arr *parse.ArrayNode) *jen.Statement {
	// 	Many different datatypes e.g. strings and objects,
	// 	or no datatypes at all (empty array)
	//	-> The generated code is []interface{}
	childrenTypeCount := countNodeTypes(arr.Children)
	if childrenTypeCount != 1 {
		return jen.Index().Interface()
	}

	// 	Only arrays as children
	//	-> The generated code is []interface{}
	if arr.Children[0].Type() == parse.NodeTypeArray {
		return jen.Index().Interface()
	}

	// 	Only structs as children
	//	-> We have to merge the structs since we don't
	// 			want to "lose" data when ultimately parsing with
	//			the generated code.
	if arr.Children[0].Type() == parse.NodeTypeObject {

		// Only merge the children objects if the size is greater than one.
		if len(arr.Children) > 1 {
			mergedObject := mergeObjects(castToObjectArr(arr.Children))
			return jen.Index().Add(makeStruct(mergedObject))
		}

		// At this point, we are sure that there is only one object as a child
		return jen.Index().Add(makeStruct(arr.Children[0].(*parse.ObjectNode)))
	}

	// 	Only one primitive datatype e.g. only strings
	//	-> The generated code is []string
	return jen.Index().Add(makePrimTypedef(arr.Children[0].Type()))
}

func makeStruct(obj *parse.ObjectNode) *jen.Statement {
	var children []jen.Code

	var sortedVarnames []string
	for varname := range obj.Children {
		sortedVarnames = append(sortedVarnames, varname)
	}
	sort.Strings(sortedVarnames)

	for _, varname := range sortedVarnames {

		valueArray := obj.Children[varname]
		childrenWithSharedKey := len(valueArray)
		typeCount := countNodeTypes(valueArray)

		if childrenWithSharedKey == 0 {
			continue
		}

		if childrenWithSharedKey == 1 {
			switch valueArray[0].Type() {
			case parse.NodeTypeArray:
				childArr := valueArray[0].(*parse.ArrayNode)
				children = append(
					children,
					makeVarname(varname).
						Add(makeArray(childArr)).
						Add(makeJSONTag(varname)),
				)
			case parse.NodeTypeObject:
				childObj := valueArray[0].(*parse.ObjectNode)
				children = append(
					children,
					makeVarname(varname).
						Add(makeStruct(childObj)).
						Add(makeJSONTag(varname)),
				)
			default:
				children = append(
					children,
					makeVarname(varname).
						Add(makePrimTypedef(valueArray[0].Type())).
						Add(makeJSONTag(varname)),
				)
			}

		} else if childrenWithSharedKey > 1 {

			if typeCount > 1 {
				children = append(
					children,
					makeVarname(varname).
						Add(jen.Interface()).
						Add(makeJSONTag(varname)),
				)

				// If there are different objects for the same key, we should merge them together.
			} else if typeCount == 1 {

				switch typ := valueArray[0].Type(); typ {
				case parse.NodeTypeObject:
					compositeObj := mergeObjects(castToObjectArr(valueArray))
					children = append(children,
						makeVarname(varname).
							Add(makeStruct(compositeObj)).
							Add(makeJSONTag(varname)),
					)
				case parse.NodeTypeBool:
					fallthrough
				case parse.NodeTypeFloat:
					fallthrough
				case parse.NodeTypeString:
					fallthrough
				case parse.NodeTypeInteger:
					children = append(children,
						makeVarname(varname).
							Add(makePrimTypedef(typ)).
							Add(makeJSONTag(varname)),
					)
				default:
					children = append(
						children,
						makeVarname(varname).
							Add(jen.Interface()).
							Add(makeJSONTag(varname)),
					)

				}
			}
		}
	}

	return jen.Struct(children...)
}

func mergeObjects(children []*parse.ObjectNode) *parse.ObjectNode {
	mergedChildren := make(map[string][]parse.Node)

	for _, object := range children {
		for varname, valueArray := range object.Children {

			if mergedChildren[varname] == nil {
				mergedChildren[varname] = valueArray

			} else {

				typeCount := countNodeTypes(mergedChildren[varname])
				// We want to merge nested objects aswell.
				// For that, we need to check if
				// 1) the type for the values in mergedChildren is object
				// 2) the type for the values in valueArray is object
				if typeCount == 1 &&
					mergedChildren[varname][0].Type() == parse.NodeTypeObject &&
					valueArray[0].Type() == parse.NodeTypeObject {

					var objectsToBeMerged []*parse.ObjectNode
					objectsToBeMerged = append(objectsToBeMerged, castToObjectArr(mergedChildren[varname])...)
					objectsToBeMerged = append(objectsToBeMerged, castToObjectArr(valueArray)...)

					mergedObj := mergeObjects(objectsToBeMerged)

					mergedChildren[varname] = []parse.Node{mergedObj}

				}
			}
		}
	}

	return &parse.ObjectNode{
		NodeType: parse.NodeTypeObject,
		Children: mergedChildren,
	}
}

var caser = cases.Title(language.English)

func makeVarname(varname string) *jen.Statement {
	upperCaseVarname := caser.String(strings.ToLower(varname))

	// Maybe we can find a shortcut? This is getting a bit computation heavy
	if isValid, validIdentifier := identifierIsValid(varname); !isValid {
		upperCaseVarname = validIdentifier
		fmt.Println("Invalid identifier found. We cleaned that up for you, but you might want to double check if you're happy with our naming: ", validIdentifier)
	}

	if identifierIsPredeclared(varname) {
		upperCaseVarname = "_" + upperCaseVarname
		fmt.Println("Predeclared identifier found. We cleaned that up for you, but you might want to double check if you're happy with our naming: ", upperCaseVarname)
	}

	return jen.Id(upperCaseVarname)
}

// Valid identifiers are specified here: https://go.dev/ref/spec#Identifiers.
// Summarized we can say:
// identifier 	  = letter { letter | unicode_digit } .
// letter         = unicode_letter | "_" .
// unicode_letter = /* a Unicode code point classified as "Letter" */ .
// unicode_digit  = /* a Unicode code point classified as "Number, decimal digit" */ .
func identifierIsValid(identifier string) (bool, string) {
	characters := []rune(identifier)
	if !isLetter(characters[0]) {
		characters[0] = '_'
	}
	for i, char := range characters[1:] {
		if !(isLetter(char) || unicode.IsDigit(char)) {
			characters[i] = '_'
		}
	}
	validIdentifier := string(characters)
	return validIdentifier == identifier, validIdentifier
}

func isLetter(letter rune) bool {
	return unicode.IsLetter(letter) || letter == '_'
}

var predeclaredIdentifiers = []string{"any", "bool", "byte", "comparable", "complex64", "complex128", "error", "float32", "float64", "int", "int8", "int16", "int32", "int64", "rune", "string", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "true", "false", "iota", "nil", "append", "cap", "close", "complex", "copy", "delete", "imag", "len", "make", "new", "panic", "print", "println", "real", "recover"}

func identifierIsPredeclared(identifier string) bool {
	for _, predeclared := range predeclaredIdentifiers {
		if identifier == predeclared {
			return true
		}
	}
	return false
}

// addJSONTag adds the json-tag, e.g. `json:"title"`. This has to match the original varname from the json file
func makeJSONTag(varname string) *jen.Statement {
	return jen.Tag(map[string]string{"json": varname})
}

// Depending on `typ`, add the type of the identifier, e.g. `Title string`.
func makePrimTypedef(typ parse.NodeType) *jen.Statement {
	switch typ {
	case parse.NodeTypeBool:
		return jen.Bool()
	case parse.NodeTypeNil:
		return jen.Interface()
	case parse.NodeTypeString:
		return jen.String()
	case parse.NodeTypeInteger:
		return jen.Int()
	case parse.NodeTypeFloat:
		return jen.Float64()
	default:
		panic("received unexpected primitive nodetype")
	}
}

func castToObjectArr(arr []parse.Node) []*parse.ObjectNode {
	var objectArr []*parse.ObjectNode
	for _, child := range arr {
		obj, ok := child.(*parse.ObjectNode)
		if !ok {
			panic("casting a node to objectnode failed")
		}
		objectArr = append(objectArr, obj)
	}
	return objectArr
}

func countNodeTypes(children []parse.Node) int {
	// If there are only zero or one children, then there are zero or one
	// different types of children aswell.
	childrenCount := len(children)
	if childrenCount == 0 || childrenCount == 1 {
		return childrenCount
	}

	foundTypes := make(map[parse.NodeType]bool, 0)
	for _, child := range children {
		foundTypes[child.Type()] = true
	}
	return len(foundTypes)
}
