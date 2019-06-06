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
	iLeftBrace  = mkItem(ItemLeftBrace, "{")
	iRightBrace = mkItem(ItemRightBrace, "}")
	iEOF        = mkItem(ItemEOF, "")
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
	{"empty file", "", []Item{mkItem(ItemEOF, "")}},
	{"empty json", "{}", []Item{
		iLeftBrace,
		iRightBrace,
		iEOF,
	}},
	{"only bool", `{"bool1": true}`, []Item{
		iLeftBrace,
		mkItem(ItemString, "bool1"),
		mkItem(ItemColon, ":"),
		mkItem(ItemBool, "true"),
		iRightBrace,
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
