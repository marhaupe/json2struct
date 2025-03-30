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
	ItemLeftBrace
	ItemRightBrace
	ItemLeftSqrBrace
	ItemRightSqrBrace
	ItemComma
	ItemColon
	ItemEOF
	ItemError
)

const (
	whitespace          = ' '
	tab                 = '\t'
	newline             = '\n'
	comma               = ','
	colon               = ':'
	leftBrace           = '{'
	rightBrace          = '}'
	leftSqrBrace        = '['
	rightSqrBrace       = ']'
	quote               = '"'
	plus                = '+'
	minus               = '-'
	zero                = '0'
	nine                = '9'
	n                   = 'n'
	t                   = 't'
	f                   = 'f'
	backslash           = '\\'
	validNumberPrefixes = "+-"
	validNumberDigits   = "0123456789"
	validNumberSuffixes = ".eE+-0123456789"
	wordTrue            = "true"
	wordFalse           = "false"
	wordNullLength      = len("null")
)

const EOF = -1

type Lexer struct {
	input string
	pos   int
	start int
	width int
	state stateFn
}

func (l *Lexer) NextItem() *Item {
	for {
		newState, item := l.state(l)
		l.state = newState
		if item != nil {
			return item
		}
	}
}

func Lex(json string) *Lexer {
	l := &Lexer{
		input: json,
		pos:   0,
		start: 0,
		width: 0,
		state: lexWhitespace,
	}
	return l
}

func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return EOF
	}
	rn, width := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = width
	l.pos += l.width
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
	l.pos -= l.width
}

func (l *Lexer) emit(t ItemType) *Item {
	item := &Item{
		Typ:   t,
		Pos:   l.pos,
		Value: l.input[l.start:l.pos],
	}
	l.start = l.pos
	return item
}

func (l *Lexer) emitError(msg string) *Item {
	item := &Item{
		Typ:   ItemError,
		Pos:   l.pos,
		Value: msg,
	}
	l.start = l.pos
	return item
}

func (l *Lexer) ignore() {
	l.start = l.pos
}

// func (l *Lexer) run() {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			l.emitError(fmt.Sprint(r))
// 		}
// 	}()

// 	for state := lexWhitespace; state != nil; {
// 		state = state(l)
// 	}
// 	defer close(l.items)
// }

type stateFn func(l *Lexer) (stateFn, *Item)

func lexWhitespace(l *Lexer) (stateFn, *Item) {
	for r := l.next(); isSpace(r) || r == newline; r = l.next() {
	}
	l.backup()
	l.ignore()

	switch r := l.next(); {
	case r == EOF:
		item := l.emit(ItemEOF)
		return nil, item
	case r == leftBrace:
		return lexLeftBrace, nil
	case r == leftSqrBrace:
		return lexLeftSqrBrace, nil
	case r == rightBrace:
		return lexRightBrace, nil
	case r == rightSqrBrace:
		return lexRightSqrBrace, nil
	case r == colon:
		return lexColon, nil
	case r == comma:
		return lexComma, nil
	case r == quote:
		return lexString, nil
	case isNumber(r):
		l.backup()
		return lexNumber, nil
	case isNull(r):
		l.backup()
		return lexNull, nil
	case isBool(r):
		l.backup()
		return lexBool, nil
	default:
		panic(fmt.Sprintf("unexpected character %v", string(r)))
	}
}

func lexNull(l *Lexer) (stateFn, *Item) {
	l.pos += wordNullLength
	item := l.emit(ItemNil)
	return lexWhitespace, item
}

func lexBool(l *Lexer) (stateFn, *Item) {
	var item *Item
	for {
		r := l.next()
		if r == EOF {
			panic("unexpected eof while lexing bool")
		}
		word := l.input[l.start:l.pos]
		if word == wordTrue || word == wordFalse {
			item = l.emit(ItemBool)
			break
		}
	}
	return lexWhitespace, item
}

func lexNumber(l *Lexer) (stateFn, *Item) {
	l.acceptRun(validNumberPrefixes)
	l.acceptRun(validNumberDigits)

	count := l.acceptRun(validNumberSuffixes)
	var item *Item
	if count > 0 {
		item = l.emit(ItemFloat)
	} else {
		item = l.emit(ItemInteger)
	}

	return lexWhitespace, item
}

func lexString(l *Lexer) (stateFn, *Item) {
	// the current rune is known to be `"`. Throw it away in order to emit the raw string value
	// e.g. example instead of "example".
	l.ignore()
	for r := l.next(); r != quote; r = l.next() {
		if r == backslash {
			r = l.next()
		}
		if r == EOF {
			panic("unterminated string")
		}
	}

	// the current rune is the closing `"`. Throw this one away as well by decrementing the position. Reset the correct state after emitting the lexem.
	l.pos--
	item := l.emit(ItemString)
	l.pos++
	return lexWhitespace, item
}

func lexComma(l *Lexer) (stateFn, *Item) {
	item := l.emit(ItemComma)
	return lexWhitespace, item
}

func lexColon(l *Lexer) (stateFn, *Item) {
	item := l.emit(ItemColon)
	return lexWhitespace, item
}

func lexLeftBrace(l *Lexer) (stateFn, *Item) {
	item := l.emit(ItemLeftBrace)
	return lexWhitespace, item
}

func lexRightBrace(l *Lexer) (stateFn, *Item) {
	item := l.emit(ItemRightBrace)
	return lexWhitespace, item
}

func lexLeftSqrBrace(l *Lexer) (stateFn, *Item) {
	item := l.emit(ItemLeftSqrBrace)
	return lexWhitespace, item
}

func lexRightSqrBrace(l *Lexer) (stateFn, *Item) {
	item := l.emit(ItemRightSqrBrace)
	return lexWhitespace, item
}

const ()

func isSpace(r rune) bool {
	return r == whitespace || r == tab
}

func isNumber(r rune) bool {
	return r == plus || r == minus || (zero <= r && r <= nine)
}

func isNull(r rune) bool {
	return r == n
}

func isBool(r rune) bool {
	return r == t || r == f
}
