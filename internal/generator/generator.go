package generator

import (
	"strings"

	"github.com/marhaupe/json2struct/internal/parse"

	"github.com/dave/jennifer/jen"
)

func GenerateFile(tree parse.Node) *jen.File {
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

func (g Generator) start() *jen.File {
	rootStmt := g.file.Type().Id("JSONToStruct")

	switch g.Tree.Type() {
	case parse.NodeTypeArray:
		casted := g.Tree.(*parse.ArrayNode)
		rootStmt.Add(makeArray(casted))
	case parse.NodeTypeObject:
		casted := g.Tree.(*parse.ObjectNode)
		rootStmt.Add(makeStruct(casted))
	default:
		panic("temp 1")
	}

	return g.file

}

func countChildrenTypes(node *parse.ArrayNode) int {
	// If there are only zero or one children, then there are zero or one
	// different types of children aswell.
	childrenCount := len(node.Children)
	if childrenCount == 0 || childrenCount == 1 {
		return childrenCount
	}

	foundTypes := make(map[parse.NodeType]bool, 0)
	for _, child := range node.Children {
		foundTypes[child.Type()] = true
	}
	return len(foundTypes)
}

// addJSONTag adds the json-tag, e.g. `json:"title"`. This has to match the original varname from the json file
func makeJSONTag(varname string) *jen.Statement {
	return jen.Tag(map[string]string{"json": varname})
}

func makeArray(arr *parse.ArrayNode) *jen.Statement {

	childrenTypeCount := countChildrenTypes(arr)
	if childrenTypeCount != 1 || arr.Children[0].Type() == parse.NodeTypeArray {
		return jen.Index().Interface()
	}

	if arr.Children[0].Type() == parse.NodeTypeObject {
		mergedChildren := make(map[string][]parse.Node)
		for _, child := range arr.Children {
			childObj := child.(*parse.ObjectNode)

			for key, val := range childObj.Children {
				mergedChildren[key] = val
			}
		}
		compositeObj := &parse.ObjectNode{
			NodeType: parse.NodeTypeObject,
			Children: mergedChildren,
		}
		return jen.Index().Add(makeStruct(compositeObj))
	}

	return jen.Index().Add(makePrimTypedef(arr.Children[0].Type()))
}

func makeStruct(obj *parse.ObjectNode) *jen.Statement {

	var children []jen.Code
	for varname, valueArray := range obj.Children {
		if len(valueArray) != 1 {
			children = append(children, makeVarname(varname).Interface().Add(makeJSONTag(varname)))
		} else {
			switch valueArray[0].Type() {
			case parse.NodeTypeArray:
				childArr := valueArray[0].(*parse.ArrayNode)
				children = append(children, makeVarname(varname).Add(makeArray(childArr)).Add(makeJSONTag(varname)))
			case parse.NodeTypeObject:
				childObj := valueArray[0].(*parse.ObjectNode)
				children = append(children, makeVarname(varname).Add(makeStruct(childObj)).Add(makeJSONTag(varname)))
			default:
				children = append(children, makeVarname(varname).Add(makePrimTypedef(valueArray[0].Type())).Add(makeJSONTag(varname)))
			}
		}
	}
	return jen.Struct(children...)
}

// We have to return e.g. `Title string `json:"title"``.
// Start with the identifier, e.g. `Title`. This has to be uppercase.
func makeVarname(varname string) *jen.Statement {
	upperCaseVarname := strings.Title(strings.ToLower(varname))
	return jen.Id(upperCaseVarname)
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
	case parse.NodeTypeNumber:
		return jen.Float64()
	default:
		panic("temp 2")
	}
}
