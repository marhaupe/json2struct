package lex

import (
	"fmt"
	"strings"
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
	ItemInteger
	ItemFloat
	ItemNil
	ItemKey
	ItemLeftBrace
	ItemRightBrace
	ItemLeftSqrBrace
	ItemRightSqrBrace
	ItemComma
	ItemColon
	ItemEOF
	ItemError
)

const EOF = -1

type Lexer struct {
	Name  string
	Input string
	Pos   int
	Start int
	Width int
	Items chan *Item
}

func (l *Lexer) NextItem() *Item {
	return <-l.Items
}

func Lex(name, json string) *Lexer {
	l := &Lexer{
		Name:  name,
		Input: json,
		Pos:   0,
		Start: 0,
		Width: 0,
		Items: make(chan *Item, 32),
	}
	go l.run()
	return l
}

func (l *Lexer) next() rune {
	if l.Pos >= len(l.Input) {
		l.Width = 0
		return EOF
	}
	rn, width := utf8.DecodeRuneInString(l.Input[l.Pos:])
	l.Width = width
	l.Pos += l.Width
	return rn
}

func (l *Lexer) acceptRun(valid string) int {
	validCount := 0
	for strings.ContainsRune(valid, l.next()) {
		validCount++
	}
	l.backup()
	return validCount
}

func (l *Lexer) backup() {
	l.Pos -= l.Width
}

func (l *Lexer) emit(t ItemType) {
	l.Items <- &Item{
		Typ:   t,
		Pos:   l.Pos,
		Value: l.Input[l.Start:l.Pos],
	}
	l.Start = l.Pos
}

func (l *Lexer) emitError(msg string) {
	l.Items <- &Item{
		Typ:   ItemError,
		Pos:   l.Pos,
		Value: msg,
	}
	l.Start = l.Pos
}

func (l *Lexer) ignore() {
	l.Start = l.Pos
}

func (l *Lexer) run() {
	defer func() {
		if r := recover(); r != nil {
			l.emitError(fmt.Sprint(r))
		}
	}()

	for state := lexWhitespace; state != nil; {
		state = state(l)
	}
	defer close(l.Items)
}

type stateFn func(l *Lexer) stateFn

func lexWhitespace(l *Lexer) stateFn {
	for r := l.next(); isSpace(r) || r == '\n'; r = l.next() {
	}
	l.backup()
	l.ignore()

	switch r := l.next(); {
	case r == EOF:
		l.emit(ItemEOF)
		return nil
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
		panic(fmt.Sprintf("unexpected character %v", string(r)))
	}
}

func lexNull(l *Lexer) stateFn {
	l.Pos += len("null")
	l.emit(ItemNil)
	return lexWhitespace
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
	l.acceptRun("+-")
	l.acceptRun("0123456789")

	count := l.acceptRun(".eE+-0123456789")
	if count > 0 {
		l.emit(ItemFloat)
	} else {
		l.emit(ItemInteger)
	}

	return lexWhitespace
}

func lexString(l *Lexer) stateFn {
	// the current rune is known to be `"`. Throw it away in order to emit the raw string value
	// e.g. example instead of "example".
	l.ignore()
	for r := l.next(); r != '"'; r = l.next() {
		if r == '\\' {
			r = l.next()
		}
		if r == EOF {
			panic("unterminated string")
		}
	}

	// the current rune is the closing `"`. Throw this one away as well by decrementing the position. Reset the correct state after emitting the lexem.
	l.Pos--
	l.emit(ItemString)
	l.Pos++
	return lexWhitespace
}

func lexComma(l *Lexer) stateFn {
	l.emit(ItemComma)
	return lexWhitespace
}

func lexColon(l *Lexer) stateFn {
	l.emit(ItemColon)
	return lexWhitespace
}

func lexLeftBrace(l *Lexer) stateFn {
	l.emit(ItemLeftBrace)
	return lexWhitespace
}

func lexRightBrace(l *Lexer) stateFn {
	l.emit(ItemRightBrace)
	return lexWhitespace
}

func lexLeftSqrBrace(l *Lexer) stateFn {
	l.emit(ItemLeftSqrBrace)
	return lexWhitespace
}

func lexRightSqrBrace(l *Lexer) stateFn {
	l.emit(ItemRightSqrBrace)
	return lexWhitespace
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
