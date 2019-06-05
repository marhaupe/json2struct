package lex

import (
	"unicode/utf8"
)

type Item struct {
	Typ   ItemType
	Pos   int
	Value string
}

type ItemType int

const (
	ItemString ItemType = iota
	ItemBool
	ItemFloat
	ItemNil
	ItemKey
	ItemLeftBrace
	ItemRightBrace
	ItemLeftSqrBrace
	ItemRightSqrBrace
	ItemComma
	ItemColon
)

const EOF = -1

type Lexer struct {
	Name  string
	Input string
	Pos   int
	Start int
	Width int
	Items chan Item
}

func (l *Lexer) NextItem() Item {
	return <-l.Items
}

func Lex(name, json string) *Lexer {
	l := &Lexer{
		Name:  name,
		Input: json,
		Pos:   0,
		Start: 0,
		Width: 0,
		Items: make(chan Item),
	}
	go l.run()
	return l
}

func (l *Lexer) next() rune {
	if l.Pos > len(l.Input) {
		l.Width = 0
		return -1
	}
	rn, width := utf8.DecodeLastRuneInString(l.Input[l.Pos:])
	l.Width = width
	l.Pos += l.Width
	return rn
}

func (l *Lexer) peek() rune {
	rn := l.next()
	l.backup()
	return rn
}

func (l *Lexer) backup() {
	l.Pos -= l.Width
}

func (l *Lexer) emit(t ItemType) {
	l.Items <- Item{
		Typ:   t,
		Pos:   l.Pos,
		Value: l.Input[l.Start:l.Pos],
	}
	l.Start = l.Pos
}

func (l *Lexer) ignore() {
	l.Start = l.Pos
}

func (l *Lexer) run() {
	for state := lexWhitespace; state != nil; {
		state = state(l)
	}
	close(l.Items)
}

const (
	errorInvalidStart = "json started with neither [ nor {"
)

type stateFn func(l *Lexer) stateFn

func lexWhitespace(l *Lexer) stateFn {
	for r := l.next(); isSpace(r) || r == '\n'; l.next() {
		l.peek()
	}
	l.backup()
	l.ignore()

	switch r := l.next(); {
	case r == '{':
		return lexLeftBrace
	case r == '[':
		return lexLeftSqrBrace
	case r == '}':
		return lexRightBrace
	case r == ']':
		return lexRightSqrBrace
	case r == ':':
		return lexColon
	case r == ',':
		return lexComma
	case r == '"':
		return lexString
	case isNumber(r):
		l.backup()
		return lexNumber
	case isNull(r):
		l.backup()
		return lexNull
	case isBool(r):
		l.backup()
		return lexBool
	default:
		panic("unexpected character " + string(r))
	}
}

func lexNull(l *Lexer) stateFn {
	l.Pos += len("null")
	l.emit(ItemNil)
	return nil
}

func lexBool(l *Lexer) stateFn {
	for {
		r := l.next()
		if r == EOF {
			panic("unexpected eof while lexing bool")
		}
		word := l.Input[l.Start:l.Pos]
		if word == "true" || word == "false" {
			l.emit(ItemBool)
			break
		}
	}
	return lexWhitespace
}

func lexNumber(l *Lexer) stateFn {
	return nil
}

func lexString(l *Lexer) stateFn {
	for r := l.next(); r != '"'; r = l.next() {
		if r == EOF {
			panic("unexpected eof while lexing string")
		}
	}
	l.emit(ItemString)
	return lexWhitespace
}

func lexComma(l *Lexer) stateFn {
	return nil
}

func lexColon(l *Lexer) stateFn {
	return nil
}

func lexLeftBrace(l *Lexer) stateFn {
	return nil
}

func lexRightBrace(l *Lexer) stateFn {
	return nil
}

func lexLeftSqrBrace(l *Lexer) stateFn {
	return nil
}

func lexRightSqrBrace(l *Lexer) stateFn {
	return nil
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

func isNumber(r rune) bool {
	return r == '+' || r == '-' || ('0' <= r && r <= '9')
}

func isNull(r rune) bool {
	return r == 'n'
}

func isBool(r rune) bool {
	return r == 't' || r == 'f'
}
