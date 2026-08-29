package syntax

import "fmt"

type TokenKind int

const (
	EOF TokenKind = iota
	Newline
	Identifier
	Number
	String
	Operator
	LParen
	RParen
	LBrace
	RBrace
	LBracket
	RBracket
	Comma
	Semicolon
)

type Token struct {
	Kind TokenKind
	Text string
	At   Span
}

func (t Token) String() string {
	return fmt.Sprintf("%q at %d:%d", t.Text, t.At.Start.Line, t.At.Start.Column)
}

type ParseError struct {
	Message string
	At      Span
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("%d:%d: %s", e.At.Start.Line, e.At.Start.Column, e.Message)
}
