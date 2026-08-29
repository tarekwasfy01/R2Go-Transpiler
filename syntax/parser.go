package syntax

import (
	"fmt"
	"strconv"
	"strings"
)

type Parser struct {
	tokens []Token
	pos    int
}

func Parse(src string) (*Program, error) {
	tokens, err := Lex(src)
	if err != nil {
		return nil, err
	}
	p := &Parser{tokens: tokens}
	program := &Program{}
	p.separators()
	for p.current().Kind != EOF {
		e, err := p.expression(0)
		if err != nil {
			return nil, err
		}
		program.Expressions = append(program.Expressions, e)
		if p.current().Kind != EOF && p.current().Kind != RBrace && !p.isSeparator() {
			return nil, p.errorf(p.current(), "expected newline or ';', got %q", p.current().Text)
		}
		p.separators()
	}
	return program, nil
}

func (p *Parser) current() Token { return p.tokens[p.pos] }
func (p *Parser) advance() Token {
	t := p.current()
	if p.pos < len(p.tokens)-1 {
		p.pos++
	}
	return t
}
func (p *Parser) match(kind TokenKind, texts ...string) (Token, bool) {
	t := p.current()
	if t.Kind != kind {
		return Token{}, false
	}
	if len(texts) > 0 && texts[0] != "" {
		ok := false
		for _, s := range texts {
			if t.Text == s {
				ok = true
			}
		}
		if !ok {
			return Token{}, false
		}
	}
	p.advance()
	return t, true
}
func (p *Parser) expect(kind TokenKind, text string) (Token, error) {
	t := p.current()
	if t.Kind != kind || (text != "" && t.Text != text) {
		return Token{}, p.errorf(t, "expected %q, got %q", text, t.Text)
	}
	p.advance()
	return t, nil
}
func (p *Parser) errorf(t Token, f string, a ...any) error {
	return &ParseError{Message: fmt.Sprintf(f, a...), At: t.At}
}
func (p *Parser) isSeparator() bool { k := p.current().Kind; return k == Newline || k == Semicolon }
func (p *Parser) separators() {
	for p.isSeparator() {
		p.advance()
	}
}
func span(a, b Span) Span { return Span{Start: a.Start, End: b.End} }

var precedence = map[string]int{
	"<-": 1, "<<-": 1, "=": 1, "->": 1, "->>": 1,
	"~": 2, "|>": 3, "||": 4, "|": 5, "&&": 6, "&": 7,
	"==": 8, "!=": 8, "<": 8, "<=": 8, ">": 8, ">=": 8,
	"+": 9, "-": 9, "*": 10, "/": 10, "%%": 10, "%/%": 10,
	":": 11, "^": 12, "$": 14, "@": 14, "::": 15, ":::": 15,
}

func (p *Parser) expression(min int) (Expr, error) {
	left, err := p.prefix()
	if err != nil {
		return nil, err
	}
	for {
		// A newline is whitespace when the following token continues the
		// current expression.  This is grammar-driven and works for every
		// operator instead of patching individual Base-R source files.
		if p.current().Kind == Newline {
			n := p.pos
			for n < len(p.tokens) && p.tokens[n].Kind == Newline {
				n++
			}
			infix := false
			if n < len(p.tokens) {
				_, infix = precedence[p.tokens[n].Text]
			}
			if n < len(p.tokens) && p.tokens[n].Kind == Operator && (infix || strings.HasPrefix(p.tokens[n].Text, "%")) {
				p.pos = n
			} else {
				return left, nil
			}
		}
		if p.current().Kind == LParen {
			if 16 < min {
				return left, nil
			}
			left, err = p.call(left)
			if err != nil {
				return nil, err
			}
			continue
		}
		if p.current().Kind == LBracket {
			if 16 < min {
				return left, nil
			}
			left, err = p.index(left)
			if err != nil {
				return nil, err
			}
			continue
		}
		t := p.current()
		if t.Kind != Operator {
			return left, nil
		}
		prec, ok := precedence[t.Text]
		if !ok {
			if strings.HasPrefix(t.Text, "%") {
				prec = 10
				ok = true
			}
		}
		if !ok || prec < min {
			return left, nil
		}
		p.advance()
		p.skipSoftNewlines()
		next := prec + 1
		if t.Text == "^" || t.Text == "<-" || t.Text == "<<-" || t.Text == "=" {
			next = prec
		}
		var right Expr
		var e error
		if t.Text == "$" || t.Text == "@" || t.Text == "::" || t.Text == ":::" {
			r := p.current()
			if r.Kind != Identifier && r.Kind != String {
				return nil, p.errorf(r, "expected name after %s", t.Text)
			}
			p.advance()
			name := strings.Trim(r.Text, "`")
			if r.Kind == String {
				name, e = unquoteRString(r.Text)
				if e != nil {
					return nil, e
				}
			}
			right = &Symbol{Name: name, At: r.At}
		} else {
			right, e = p.expression(next)
		}
		if e != nil {
			return nil, e
		}
		if t.Text == "->" || t.Text == "->>" {
			left, right = right, left
			if t.Text == "->" {
				t.Text = "<-"
			} else {
				t.Text = "<<-"
			}
		}
		if t.Text == "|>" {
			call, ok := right.(*Call)
			if !ok {
				return nil, p.errorf(t, "pipe target must be a function call")
			}
			call.Arguments = append([]Argument{{Value: left, At: left.SourceSpan()}}, call.Arguments...)
			call.At = span(left.SourceSpan(), call.SourceSpan())
			left = call
			continue
		}
		left = &Call{Function: &Symbol{Name: t.Text, At: t.At}, Arguments: []Argument{{Value: left, At: left.SourceSpan()}, {Value: right, At: right.SourceSpan()}}, At: span(left.SourceSpan(), right.SourceSpan())}
	}
}

func (p *Parser) prefix() (Expr, error) {
	t := p.advance()
	switch t.Kind {
	case Number:
		return &Literal{Kind: NumberLiteral, Text: t.Text, At: t.At}, nil
	case String:
		text, err := unquoteRString(t.Text)
		if err != nil {
			return nil, p.errorf(t, "invalid string: %v", err)
		}
		return &Literal{Kind: StringLiteral, Text: text, At: t.At}, nil
	case Identifier:
		name := strings.Trim(t.Text, "`")
		switch name {
		case "TRUE", "T", "FALSE", "F":
			return &Literal{Kind: LogicalLiteral, Text: name, At: t.At}, nil
		case "NULL":
			return &Literal{Kind: NullLiteral, Text: name, At: t.At}, nil
		case "NA", "NA_integer_", "NA_real_", "NA_character_", "NA_complex_":
			return &Literal{Kind: NALiteral, Text: name, At: t.At}, nil
		case "Inf", "NaN":
			return &Literal{Kind: NumberLiteral, Text: name, At: t.At}, nil
		case "function":
			return p.function(t)
		case "if":
			return p.ifExpr(t)
		case "while":
			return p.whileExpr(t)
		case "for":
			return p.forExpr(t)
		case "repeat":
			body, e := p.expression(0)
			if e != nil {
				return nil, e
			}
			return &Repeat{Body: body, At: span(t.At, body.SourceSpan())}, nil
		case "break", "next":
			return &Call{Function: &Symbol{Name: name, At: t.At}, At: t.At}, nil
		case "return":
			args := []Argument{}
			end := t.At
			if p.current().Kind == LParen {
				c, e := p.call(&Symbol{Name: "return", At: t.At})
				return c, e
			}
			if !p.isSeparator() && p.current().Kind != RBrace && p.current().Kind != EOF {
				v, e := p.expression(0)
				if e != nil {
					return nil, e
				}
				args = append(args, Argument{Value: v, At: v.SourceSpan()})
				end = v.SourceSpan()
			}
			return &Call{Function: &Symbol{Name: "return", At: t.At}, Arguments: args, At: span(t.At, end)}, nil
		default:
			return &Symbol{Name: name, At: t.At}, nil
		}
	case Operator:
		if t.Text == "\\" {
			return p.function(t)
		}
		if t.Text == "+" || t.Text == "-" || t.Text == "!" || t.Text == "~" {
			v, e := p.expression(13)
			if e != nil {
				return nil, e
			}
			return &Call{Function: &Symbol{Name: t.Text, At: t.At}, Arguments: []Argument{{Value: v, At: v.SourceSpan()}}, At: span(t.At, v.SourceSpan())}, nil
		}
	case LParen:
		p.skipSoftNewlines()
		e, err := p.expression(0)
		if err != nil {
			return nil, err
		}
		p.skipSoftNewlines()
		end, err := p.expect(RParen, ")")
		if err != nil {
			return nil, err
		}
		return &Call{Function: &Symbol{Name: "(", At: t.At}, Arguments: []Argument{{Value: e, At: e.SourceSpan()}}, At: span(t.At, end.At)}, nil
	case LBrace:
		return p.block(t)
	}
	return nil, p.errorf(t, "expected expression, got %q", t.Text)
}

func unquoteRString(text string) (string, error) {
	if len(text) < 2 {
		return "", fmt.Errorf("invalid quoted string")
	}
	quote := text[0]
	rest := text[1 : len(text)-1]
	var out strings.Builder
	for rest != "" {
		if rest[0] == '\n' {
			out.WriteByte('\n')
			rest = rest[1:]
			continue
		}
		// R accepts an escaped single quote in a double-quoted string and
		// vice versa; strconv.UnquoteChar intentionally rejects that form.
		if len(rest) >= 2 && rest[0] == '\\' && ((quote == '"' && rest[1] == '\'') || (quote == '\'' && rest[1] == '"')) {
			out.WriteByte(rest[1])
			rest = rest[2:]
			continue
		}
		r, multibyte, tail, err := strconv.UnquoteChar(rest, quote)
		if err != nil {
			return "", err
		}
		if multibyte {
			out.WriteRune(r)
		} else {
			out.WriteByte(byte(r))
		}
		rest = tail
	}
	return out.String(), nil
}

func (p *Parser) block(start Token) (Expr, error) {
	b := &Block{At: start.At}
	p.separators()
	for p.current().Kind != RBrace {
		if p.current().Kind == EOF {
			return nil, p.errorf(p.current(), "unterminated block")
		}
		e, err := p.expression(0)
		if err != nil {
			return nil, err
		}
		b.Expressions = append(b.Expressions, e)
		if p.current().Kind != RBrace && !p.isSeparator() {
			return nil, p.errorf(p.current(), "expected separator in block")
		}
		p.separators()
	}
	end := p.advance()
	b.At = span(start.At, end.At)
	return b, nil
}

func (p *Parser) call(fn Expr) (Expr, error) {
	start := fn.SourceSpan()
	p.advance()
	args, err := p.arguments(RParen)
	if err != nil {
		return nil, err
	}
	p.skipSoftNewlines()
	end, err := p.expect(RParen, ")")
	if err != nil {
		return nil, err
	}
	return &Call{Function: fn, Arguments: args, At: span(start, end.At)}, nil
}

func (p *Parser) index(object Expr) (Expr, error) {
	open := p.advance()
	name := "["
	double := false
	if p.current().Kind == LBracket {
		p.advance()
		name = "[["
		double = true
	}
	args, err := p.arguments(RBracket)
	if err != nil {
		return nil, err
	}
	end, err := p.expect(RBracket, "]")
	if err != nil {
		return nil, err
	}
	if double {
		end, err = p.expect(RBracket, "]")
		if err != nil {
			return nil, err
		}
	}
	args = append([]Argument{{Value: object, At: object.SourceSpan()}}, args...)
	return &Call{Function: &Symbol{Name: name, At: open.At}, Arguments: args, At: span(object.SourceSpan(), end.At)}, nil
}

func (p *Parser) arguments(end TokenKind) ([]Argument, error) {
	var args []Argument
	p.skipSoftNewlines()
	if p.current().Kind == end {
		return args, nil
	}
	for {
		if p.current().Kind == Comma {
			t := p.advance()
			args = append(args, Argument{At: t.At})
			p.skipSoftNewlines()
			if p.current().Kind == end {
				break
			}
			continue
		} else {
			var name string
			named := false
			start := p.current().At
			if (p.current().Kind == Identifier || p.current().Kind == String) && p.pos+1 < len(p.tokens) && p.tokens[p.pos+1].Kind == Operator && p.tokens[p.pos+1].Text == "=" {
				nameToken := p.advance()
				name = strings.Trim(nameToken.Text, "`")
				if nameToken.Kind == String {
					decoded, e := unquoteRString(nameToken.Text)
					if e != nil {
						return nil, e
					}
					name = decoded
				}
				p.advance()
				named = true
			}
			if named && (p.current().Kind == Comma || p.current().Kind == end) {
				args = append(args, Argument{Name: name, At: start})
			} else {
				p.skipSoftNewlines()
				v, err := p.expression(0)
				if err != nil {
					return nil, err
				}
				args = append(args, Argument{Name: name, Value: v, At: span(start, v.SourceSpan())})
			}
		}
		p.skipSoftNewlines()
		if _, ok := p.match(Comma, ""); !ok {
			break
		}
		p.skipSoftNewlines()
		if p.current().Kind == end {
			args = append(args, Argument{At: p.current().At})
			break
		}
	}
	return args, nil
}
func (p *Parser) skipSoftNewlines() {
	for p.current().Kind == Newline {
		p.advance()
	}
}

func (p *Parser) function(start Token) (Expr, error) {
	if _, e := p.expect(LParen, "("); e != nil {
		return nil, e
	}
	var params []Parameter
	p.skipSoftNewlines()
	for p.current().Kind != RParen {
		t, e := p.expect(Identifier, "")
		if e != nil {
			return nil, e
		}
		param := Parameter{Name: strings.Trim(t.Text, "`"), At: t.At}
		if _, ok := p.match(Operator, "="); ok {
			v, e := p.expression(2)
			if e != nil {
				return nil, e
			}
			param.Default = v
			param.At = span(t.At, v.SourceSpan())
		}
		params = append(params, param)
		p.skipSoftNewlines()
		if _, ok := p.match(Comma, ""); !ok {
			break
		}
		p.skipSoftNewlines()
	}
	if _, e := p.expect(RParen, ")"); e != nil {
		return nil, e
	}
	p.skipSoftNewlines()
	body, e := p.expression(0)
	if e != nil {
		return nil, e
	}
	return &Function{Parameters: params, Body: body, At: span(start.At, body.SourceSpan())}, nil
}
func (p *Parser) ifExpr(start Token) (Expr, error) {
	if _, e := p.expect(LParen, "("); e != nil {
		return nil, e
	}
	p.skipSoftNewlines()
	c, e := p.expression(0)
	if e != nil {
		return nil, e
	}
	p.skipSoftNewlines()
	if _, e = p.expect(RParen, ")"); e != nil {
		return nil, e
	}
	p.skipSoftNewlines()
	then, e := p.expression(0)
	if e != nil {
		return nil, e
	}
	n := &If{Condition: c, Then: then, At: span(start.At, then.SourceSpan())}
	save := p.pos
	p.skipSoftNewlines()
	if p.current().Kind == Identifier && p.current().Text == "else" {
		p.advance()
		p.skipSoftNewlines()
		other, e := p.expression(0)
		if e != nil {
			return nil, e
		}
		n.Else = other
		n.At = span(start.At, other.SourceSpan())
	} else {
		p.pos = save
	}
	return n, nil
}
func (p *Parser) whileExpr(start Token) (Expr, error) {
	if _, e := p.expect(LParen, "("); e != nil {
		return nil, e
	}
	c, e := p.expression(0)
	if e != nil {
		return nil, e
	}
	if _, e = p.expect(RParen, ")"); e != nil {
		return nil, e
	}
	p.skipSoftNewlines()
	body, e := p.expression(0)
	if e != nil {
		return nil, e
	}
	return &While{Condition: c, Body: body, At: span(start.At, body.SourceSpan())}, nil
}
func (p *Parser) forExpr(start Token) (Expr, error) {
	if _, e := p.expect(LParen, "("); e != nil {
		return nil, e
	}
	v, e := p.expect(Identifier, "")
	if e != nil {
		return nil, e
	}
	in, e := p.expect(Identifier, "in")
	_ = in
	if e != nil {
		return nil, e
	}
	seq, e := p.expression(0)
	if e != nil {
		return nil, e
	}
	if _, e = p.expect(RParen, ")"); e != nil {
		return nil, e
	}
	p.skipSoftNewlines()
	body, e := p.expression(0)
	if e != nil {
		return nil, e
	}
	return &For{Variable: v.Text, Sequence: seq, Body: body, At: span(start.At, body.SourceSpan())}, nil
}
