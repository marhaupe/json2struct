package lex

import (
	"testing"
)

type lexTest struct {
	name  string
	input string
	items []Item
}

func mkItem(typ ItemType, text string) Item {
	return Item{
		Typ:   typ,
		Value: text,
	}
}

var (
	iLeftBrace     = mkItem(ItemLeftBrace, "{")
	iRightBrace    = mkItem(ItemRightBrace, "}")
	iLeftSqrBrace  = mkItem(ItemLeftSqrBrace, "[")
	iRightSqrBrace = mkItem(ItemRightSqrBrace, "]")
	iEOF           = mkItem(ItemEOF, "")
	iColon         = mkItem(ItemColon, ":")
	iComma         = mkItem(ItemComma, ",")
)

// collect gathers the emitted items into a slice.
func collect(t *lexTest) (items []Item) {
	l := Lex(t.name, t.input)
	for {
		item := l.NextItem()
		items = append(items, item)
		if item.Typ == ItemEOF {
			break
		}
	}
	return items
}

func equal(i1, i2 []Item, checkPos bool) bool {
	if len(i1) != len(i2) {
		return false
	}
	for k := range i1 {
		if i1[k].Typ != i2[k].Typ {
			return false
		}
		if i1[k].Value != i2[k].Value {
			return false
		}
		if checkPos && i1[k].Pos != i2[k].Pos {
			return false
		}
	}
	return true
}

var lexTests = []lexTest{
	{"empty file", "", []Item{iEOF}},
	{"empty object root", "{}", []Item{
		iLeftBrace,
		iRightBrace,
		iEOF,
	}},
	{"bools in object root", `{
		"bool1": true, 
		"bool2": false
		}`, []Item{
		iLeftBrace,
		mkItem(ItemString, "bool1"),
		iColon,
		mkItem(ItemBool, "true"),
		iComma,
		mkItem(ItemString, "bool2"),
		iColon,
		mkItem(ItemBool, "false"),
		iRightBrace,
		iEOF,
	}},
	{"valid numbers in object root", `{
		"number1": 1,
		"number2": 1.0,
		"number3": -1,
		"number4": -1.0
		}`, []Item{
		iLeftBrace,
		mkItem(ItemString, "number1"),
		iColon,
		mkItem(ItemInteger, "1"),
		iComma,
		mkItem(ItemString, "number2"),
		iColon,
		mkItem(ItemFloat, "1.0"),
		iComma,
		mkItem(ItemString, "number3"),
		iColon,
		mkItem(ItemInteger, "-1"),
		iComma,
		mkItem(ItemString, "number4"),
		iColon,
		mkItem(ItemFloat, "-1.0"),
		iRightBrace,
		iEOF,
	}},
	{"valid strings in object root", `{
		"string1": "value1",
		"string2": "true",
		"string3": "false",
		"string4": "null",
		"string5": "1"
		}`, []Item{
		iLeftBrace,
		mkItem(ItemString, "string1"),
		iColon,
		mkItem(ItemString, "value1"),
		iComma,
		mkItem(ItemString, "string2"),
		iColon,
		mkItem(ItemString, "true"),
		iComma,
		mkItem(ItemString, "string3"),
		iColon,
		mkItem(ItemString, "false"),
		iComma,
		mkItem(ItemString, "string4"),
		iColon,
		mkItem(ItemString, "null"),
		iComma,
		mkItem(ItemString, "string5"),
		iColon,
		mkItem(ItemString, "1"),
		iRightBrace,
		iEOF,
	}},
	{"null in object root", `{ "null1": null }`, []Item{
		iLeftBrace,
		mkItem(ItemString, "null1"),
		iColon,
		mkItem(ItemNil, "null"),
		iRightBrace,
		iEOF,
	}},
	{"nested objects with multiple keys and values", `
	{
		"obj1": {
		"null1": null,
		"bool1": true
		}
	}`, []Item{
		iLeftBrace,
		mkItem(ItemString, "obj1"),
		iColon,
		iLeftBrace,
		mkItem(ItemString, "null1"),
		iColon,
		mkItem(ItemNil, "null"),
		iComma,
		mkItem(ItemString, "bool1"),
		iColon,
		mkItem(ItemBool, "true"),
		iRightBrace,
		iRightBrace,
		iEOF,
	}},
	{"empty array root", "[]", []Item{
		iLeftSqrBrace,
		iRightSqrBrace,
		iEOF,
	}},
	{"bools in array root", `[ true, false, false, true ]`, []Item{
		iLeftSqrBrace,
		mkItem(ItemBool, "true"),
		iComma,
		mkItem(ItemBool, "false"),
		iComma,
		mkItem(ItemBool, "false"),
		iComma,
		mkItem(ItemBool, "true"),
		iRightSqrBrace,
		iEOF,
	}},
}

func TestLex(t *testing.T) {
	for _, test := range lexTests {
		items := collect(&test)
		if !equal(items, test.items, false) {
			t.Errorf("%s: got\n\t%+v\nexpected\n\t%v", test.name, items, test.items)
		}
	}
}
