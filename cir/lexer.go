// Package cir provides the C intermediate representation used to lower GNU R
// entry points into data-driven Pure-Go programs.
package cir

import (
	"fmt"
	"unicode"
)

type TokenKind uint8

const (
	EOF TokenKind = iota
	Identifier
	Number
	String
	Character
	Operator
	Punctuation
)

type Token struct {
	Kind         TokenKind
	Text         string
	Line, Column int
}

// Lex is intentionally independent of a C preprocessor. It preserves macro
// identifiers as normal tokens and skips directive lines; this is sufficient
// for lowering already-expanded R entry-point bodies in src/main.
func Lex(source string) ([]Token, error) {
	b := []rune(source)
	out := make([]Token, 0, len(b)/3)
	line, col := 1, 1
	advance := func() {
		if b[0] == '\n' {
			line++
			col = 1
		} else {
			col++
		}
		b = b[1:]
	}
	for len(b) > 0 {
		if unicode.IsSpace(b[0]) {
			advance()
			continue
		}
		if b[0] == '#' && col == 1 {
			for len(b) > 0 && b[0] != '\n' {
				advance()
			}
			continue
		}
		if len(b) > 1 && b[0] == '/' && b[1] == '/' {
			for len(b) > 0 && b[0] != '\n' {
				advance()
			}
			continue
		}
		if len(b) > 1 && b[0] == '/' && b[1] == '*' {
			advance()
			advance()
			for len(b) > 1 && !(b[0] == '*' && b[1] == '/') {
				advance()
			}
			if len(b) < 2 {
				return nil, fmt.Errorf("unterminated comment at %d", line)
			}
			advance()
			advance()
			continue
		}
		startLine, startCol := line, col
		start := 0
		if unicode.IsLetter(b[0]) || b[0] == '_' {
			for start < len(b) && (unicode.IsLetter(b[start]) || unicode.IsDigit(b[start]) || b[start] == '_') {
				start++
			}
			out = append(out, Token{Identifier, string(b[:start]), startLine, startCol})
			for i := 0; i < start; i++ {
				advance()
			}
			continue
		}
		if unicode.IsDigit(b[0]) || b[0] == '.' && len(b) > 1 && unicode.IsDigit(b[1]) {
			for start < len(b) && (unicode.IsLetter(b[start]) || unicode.IsDigit(b[start]) || b[start] == '.' || b[start] == '+' || b[start] == '-') {
				start++
			}
			out = append(out, Token{Number, string(b[:start]), startLine, startCol})
			for i := 0; i < start; i++ {
				advance()
			}
			continue
		}
		if b[0] == '\'' || b[0] == '"' {
			quote := b[0]
			start = 1
			for start < len(b) {
				if b[start] == '\\' {
					start += 2
					continue
				}
				if b[start] == quote {
					start++
					break
				}
				start++
			}
			if start > len(b) || b[start-1] != quote {
				return nil, fmt.Errorf("unterminated literal at %d", line)
			}
			kind := String
			if quote == '\'' {
				kind = Character
			}
			out = append(out, Token{kind, string(b[:start]), startLine, startCol})
			for i := 0; i < start; i++ {
				advance()
			}
			continue
		}
		operators := []string{"<<=", ">>=", "...", "->", "++", "--", "==", "!=", "<=", ">=", "&&", "||", "+=", "-=", "*=", "/=", "%=", "&=", "|=", "^=", "<<", ">>"}
		matched := ""
		for _, op := range operators {
			r := []rune(op)
			if len(b) >= len(r) && string(b[:len(r)]) == op {
				matched = op
				break
			}
		}
		if matched != "" {
			out = append(out, Token{Operator, matched, startLine, startCol})
			for range matched {
				advance()
			}
			continue
		}
		kind := Punctuation
		if string(b[0]) != "+-*/%=!<>&|^~?" {
			kind = Punctuation
		} else {
			kind = Operator
		}
		out = append(out, Token{kind, string(b[0]), startLine, startCol})
		advance()
	}
	return append(out, Token{Kind: EOF, Line: line, Column: col}), nil
}
