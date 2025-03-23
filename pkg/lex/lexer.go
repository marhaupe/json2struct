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
	items chan *Item
}

func (l *Lexer) NextItem() *Item {
	return <-l.items
}

func Lex(json string) *Lexer {
	l := &Lexer{
		input: json,
		pos:   0,
		start: 0,
		width: 0,
		items: make(chan *Item, 32),
	}
	go l.run()
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

func (l *Lexer) emit(t ItemType) {
	l.items <- &Item{
		Typ:   t,
		Pos:   l.pos,
		Value: l.input[l.start:l.pos],
	}
	l.start = l.pos
}

func (l *Lexer) emitError(msg string) {
	l.items <- &Item{
		Typ:   ItemError,
		Pos:   l.pos,
		Value: msg,
	}
	l.start = l.pos
}

func (l *Lexer) ignore() {
	l.start = l.pos
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
	defer close(l.items)
}

type stateFn func(l *Lexer) stateFn

func lexWhitespace(l *Lexer) stateFn {
	for r := l.next(); isSpace(r) || r == newline; r = l.next() {
	}
	l.backup()
	l.ignore()

	switch r := l.next(); {
	case r == EOF:
		l.emit(ItemEOF)
		return nil
	case r == leftBrace:
		return lexLeftBrace
	case r == leftSqrBrace:
		return lexLeftSqrBrace
	case r == rightBrace:
		return lexRightBrace
	case r == rightSqrBrace:
		return lexRightSqrBrace
	case r == colon:
		return lexColon
	case r == comma:
		return lexComma
	case r == quote:
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
	l.pos += wordNullLength
	l.emit(ItemNil)
	return lexWhitespace
}

func lexBool(l *Lexer) stateFn {
	for {
		r := l.next()
		if r == EOF {
			panic("unexpected eof while lexing bool")
		}
		word := l.input[l.start:l.pos]
		if word == wordTrue || word == wordFalse {
			l.emit(ItemBool)
			break
		}
	}
	return lexWhitespace
}

func lexNumber(l *Lexer) stateFn {
	l.acceptRun(validNumberPrefixes)
	l.acceptRun(validNumberDigits)

	count := l.acceptRun(validNumberSuffixes)
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
	l.emit(ItemString)
	l.pos++
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
