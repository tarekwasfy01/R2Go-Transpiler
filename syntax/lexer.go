package syntax

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Lexer struct {
	src          string
	offset       int
	line, column int
}

func Lex(src string) ([]Token, error) {
	// UTF-8 BOMs are common in R scripts written by Windows editors and are not
	// part of the R program. Keep the original source buffer intact so byte
	// offsets in token/AST spans continue to address the caller's exact source.
	// Starting after the BOM also keeps the first real token at line 1, column 1.
	startOffset := 0
	if strings.HasPrefix(src, "\uFEFF") {
		startOffset = len("\uFEFF")
	}
	l := &Lexer{src: src, offset: startOffset, line: 1, column: 1}
	var out []Token
	for {
		tok, err := l.next()
		if err != nil {
			return nil, err
		}
		out = append(out, tok)
		if tok.Kind == EOF {
			return out, nil
		}
	}
}

func (l *Lexer) pos() Position { return Position{Offset: l.offset, Line: l.line, Column: l.column} }

func (l *Lexer) peek() rune {
	if l.offset >= len(l.src) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.src[l.offset:])
	return r
}

func (l *Lexer) take() rune {
	r, n := utf8.DecodeRuneInString(l.src[l.offset:])
	l.offset += n
	if r == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
	return r
}

func (l *Lexer) token(kind TokenKind, start Position) Token {
	return Token{Kind: kind, Text: l.src[start.Offset:l.offset], At: Span{Start: start, End: l.pos()}}
}

func (l *Lexer) next() (Token, error) {
	for {
		r := l.peek()
		if r == ' ' || r == '\t' || r == '\r' || r == '\f' {
			l.take()
			continue
		}
		if r == '#' {
			for l.peek() != 0 && l.peek() != '\n' {
				l.take()
			}
			continue
		}
		break
	}
	start := l.pos()
	r := l.peek()
	if r == 0 {
		return Token{Kind: EOF, At: Span{Start: start, End: start}}, nil
	}
	if r == '\n' {
		l.take()
		return l.token(Newline, start), nil
	}
	if unicode.IsLetter(r) || r == '.' || r == '_' {
		l.take()
		for {
			r = l.peek()
			if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '_') {
				break
			}
			l.take()
		}
		return l.token(Identifier, start), nil
	}
	if unicode.IsDigit(r) || (r == '.' && l.offset+1 < len(l.src) && l.src[l.offset+1] >= '0' && l.src[l.offset+1] <= '9') {
		l.take()
		seenExp := false
		for {
			r = l.peek()
			if unicode.IsDigit(r) || r == '.' || r == 'x' || r == 'X' || r == 'p' || r == 'P' {
				l.take()
				continue
			}
			if r == 'e' || r == 'E' {
				seenExp = true
				l.take()
				continue
			}
			if (r == '+' || r == '-') && seenExp {
				seenExp = false
				l.take()
				continue
			}
			break
		}
		if l.peek() == 'L' || l.peek() == 'i' {
			l.take()
		}
		return l.token(Number, start), nil
	}
	if r == '\'' || r == '"' || r == '`' {
		quote := l.take()
		for {
			r = l.peek()
			if r == 0 {
				return Token{}, &ParseError{Message: "unterminated string", At: Span{Start: start, End: l.pos()}}
			}
			l.take()
			if r == '\\' && quote != '`' {
				if l.peek() != 0 {
					l.take()
				}
				continue
			}
			if r == quote {
				break
			}
		}
		kind := String
		if quote == '`' {
			kind = Identifier
		}
		return l.token(kind, start), nil
	}
	single := map[rune]TokenKind{'(': LParen, ')': RParen, '{': LBrace, '}': RBrace, '[': LBracket, ']': RBracket, ',': Comma, ';': Semicolon}
	if kind, ok := single[r]; ok {
		l.take()
		return l.token(kind, start), nil
	}
	if r == '%' {
		l.take()
		for l.peek() != 0 && l.peek() != '%' && l.peek() != '\n' {
			l.take()
		}
		if l.peek() != '%' {
			return Token{}, &ParseError{Message: "unterminated special operator", At: Span{Start: start, End: l.pos()}}
		}
		l.take()
		return l.token(Operator, start), nil
	}
	for _, op := range []string{"<<-", "->>", ":::", "::", "<=", ">=", "==", "!=", "&&", "||", "<-", "->", "|>"} {
		if strings.HasPrefix(l.src[l.offset:], op) {
			for range op {
				l.take()
			}
			kind := Operator
			return l.token(kind, start), nil
		}
	}
	if strings.ContainsRune("+-*/^:~!?<>=&|$@\\", r) {
		l.take()
		return l.token(Operator, start), nil
	}
	return Token{}, &ParseError{Message: fmt.Sprintf("unexpected character %q", r), At: Span{Start: start, End: l.pos()}}
}
