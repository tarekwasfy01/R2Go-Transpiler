package highlight

import (
	"context"
	"unicode"

	"github.com/oligo/gvcode/textstyle/syntax"
)

type Language string

const (
	R  Language = "R"
	Go Language = "Go"
)

var rKeywords = wordSet("if", "else", "repeat", "while", "function", "for", "in", "next", "break", "return", "TRUE", "FALSE", "T", "F", "NULL", "NA", "Inf", "NaN")
var rBuiltins = wordSet("c", "list", "matrix", "array", "data.frame", "print", "length", "names", "dim", "nrow", "ncol", "sum", "mean", "min", "max", "range", "sort", "order", "match", "lapply", "sapply", "stop", "stopifnot", "library", "require", "requireNamespace")
var goKeywords = wordSet("break", "default", "func", "interface", "select", "case", "defer", "go", "map", "struct", "chan", "else", "goto", "package", "switch", "const", "fallthrough", "if", "range", "type", "continue", "for", "import", "return", "var", "true", "false", "nil")

// Tokens is a deterministic Pure-Go lexer for editor highlighting. It accepts
// incomplete code while the user types and cannot panic the GUI event loop.
func Tokens(ctx context.Context, language Language, text string) ([]syntax.Token, error) {
	runes := []rune(text)
	tokens := make([]syntax.Token, 0, max(32, len(runes)/8))
	for i := 0; i < len(runes); {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		start := i
		r := runes[i]
		if language == R && r == '#' || language == Go && r == '/' && i+1 < len(runes) && runes[i+1] == '/' {
			for i < len(runes) && runes[i] != '\n' {
				i++
			}
			tokens = appendToken(tokens, start, i, "comment")
			continue
		}
		if language == Go && r == '/' && i+1 < len(runes) && runes[i+1] == '*' {
			i += 2
			for i+1 < len(runes) && !(runes[i] == '*' && runes[i+1] == '/') {
				i++
			}
			if i+1 < len(runes) {
				i += 2
			}
			tokens = appendToken(tokens, start, i, "comment")
			continue
		}
		if r == '\'' || r == '"' || r == '`' {
			quote := r
			i++
			for i < len(runes) {
				if runes[i] == '\\' && quote != '`' && i+1 < len(runes) {
					i += 2
					continue
				}
				current := runes[i]
				i++
				if current == quote {
					break
				}
			}
			tokens = appendToken(tokens, start, i, "literal.string")
			continue
		}
		if unicode.IsDigit(r) || r == '.' && i+1 < len(runes) && unicode.IsDigit(runes[i+1]) {
			i = scanNumber(runes, i)
			tokens = appendToken(tokens, start, i, "literal.number")
			continue
		}
		if identStart(r, language) {
			i++
			for i < len(runes) && identContinue(runes[i], language) {
				i++
			}
			word := string(runes[start:i])
			scope := syntax.StyleScope("")
			if language == R {
				if rKeywords[word] {
					scope = "keyword"
				} else if rBuiltins[word] {
					scope = "name.builtin"
				} else if nextNonSpace(runes, i) == '(' {
					scope = "name.function"
				}
			} else if goKeywords[word] {
				scope = "keyword"
			} else if nextNonSpace(runes, i) == '(' {
				scope = "name.function"
			}
			tokens = appendToken(tokens, start, i, scope)
			continue
		}
		if containsRune("+-*/^:~!?<>=&|$@%", r) {
			i++
			for i < len(runes) && containsRune("+-*/^:~!?<>=&|$@%", runes[i]) {
				i++
			}
			tokens = appendToken(tokens, start, i, "operator")
			continue
		}
		if containsRune("(){}[],;", r) {
			i++
			tokens = appendToken(tokens, start, i, "punctuation")
			continue
		}
		i++
	}
	return tokens, nil
}

func scanNumber(runes []rune, i int) int {
	if i+1 < len(runes) && runes[i] == '0' && (runes[i+1] == 'x' || runes[i+1] == 'X') {
		i += 2
		for i < len(runes) && (unicode.IsDigit(runes[i]) || runes[i] >= 'a' && runes[i] <= 'f' || runes[i] >= 'A' && runes[i] <= 'F') {
			i++
		}
	} else {
		for i < len(runes) && unicode.IsDigit(runes[i]) {
			i++
		}
		if i < len(runes) && runes[i] == '.' {
			i++
			for i < len(runes) && unicode.IsDigit(runes[i]) {
				i++
			}
		}
		if i < len(runes) && (runes[i] == 'e' || runes[i] == 'E') {
			i++
			if i < len(runes) && (runes[i] == '+' || runes[i] == '-') {
				i++
			}
			for i < len(runes) && unicode.IsDigit(runes[i]) {
				i++
			}
		}
	}
	if i < len(runes) && (runes[i] == 'L' || runes[i] == 'i') {
		i++
	}
	return i
}

func appendToken(tokens []syntax.Token, start, end int, scope syntax.StyleScope) []syntax.Token {
	if scope == "" || end <= start {
		return tokens
	}
	return append(tokens, syntax.Token{Start: start, End: end, Scope: scope})
}

func wordSet(words ...string) map[string]bool {
	set := make(map[string]bool, len(words))
	for _, word := range words {
		set[word] = true
	}
	return set
}

func identStart(r rune, language Language) bool {
	return unicode.IsLetter(r) || r == '_' || language == R && r == '.'
}

func identContinue(r rune, language Language) bool {
	return identStart(r, language) || unicode.IsDigit(r)
}

func nextNonSpace(runes []rune, i int) rune {
	for i < len(runes) && unicode.IsSpace(runes[i]) {
		i++
	}
	if i < len(runes) {
		return runes[i]
	}
	return 0
}

func containsRune(set string, wanted rune) bool {
	for _, r := range set {
		if r == wanted {
			return true
		}
	}
	return false
}
