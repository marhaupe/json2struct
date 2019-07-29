package generator

import (
	"strings"

	"github.com/marhaupe/json2struct/internal/parse"

	"github.com/dave/jennifer/jen"
)

func GenerateFile(tree parse.Node) *jen.File {
	g := Generator{
		Tree:        tree,
		CurrentNode: tree,
	}
	return g.run()
}

type Generator struct {
	Tree        parse.Node
	CurrentNode parse.Node

	file        *jen.File
	currentStmt *jen.Statement
}

func (g Generator) run() *jen.File {
	g.file = jen.NewFile("generated")
	g.currentStmt = g.file.Type().Id("JSONToStruct")

	switch g.Tree.Type() {
	case parse.NodeTypeArray:
		g.currentStmt.Index().Struct()
	case parse.NodeTypeObject:
		g.currentStmt.Struct()
	default:
		panic("temp 1")
	}

	return g.file

}

func makePrimitiveDecl(typ parse.NodeType, varname string) *jen.Statement {

	// We have to return e.g. `Title string `json:"title"``.
	// Start with the identifier, e.g. `Title`. This has to be uppercase.
	upperCaseVarname := strings.Title(strings.ToLower(varname))
	stmt := jen.Id(upperCaseVarname)

	// Depending on `typ`, add the type of the identifier, e.g. `Title string`.
	switch typ {
	case parse.NodeTypeBool:
		stmt = stmt.Bool()
	case parse.NodeTypeNil:
		stmt = stmt.Interface()
	case parse.NodeTypeString:
		stmt = stmt.String()
	case parse.NodeTypeNumber:
		stmt = stmt.Float64()
	default:
		panic("temp 2")
	}

	// Add the json-tag, e.g. `json:"title"`. This has to match the original varname from the json file
	stmt = stmt.Tag(map[string]string{"json": varname})

	return stmt
}
