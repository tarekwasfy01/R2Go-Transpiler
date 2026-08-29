package cir

// LowerBody retains the universal C control surface. Operand-level lowering is
// deliberately a later pass: this pass must never discard a call, branch or
// assignment from GNU R source.
func LowerBody(name, source string, body []Token) Program {
	p := Program{Name: name, SourceFile: source}
	for i, token := range body {
		if token.Kind == Identifier && (token.Text == "if" || token.Text == "while" || token.Text == "for" || token.Text == "switch") {
			p.Instructions = append(p.Instructions, Instruction{Op: Branch, A: token.Text})
			continue
		}
		if token.Kind == Identifier && token.Text == "return" {
			p.Instructions = append(p.Instructions, Instruction{Op: Return})
			continue
		}
		if token.Kind == Identifier && i+1 < len(body) && body[i+1].Text == "(" {
			p.Instructions = append(p.Instructions, Instruction{Op: Call, A: token.Text})
			continue
		}
		if token.Text == "=" || token.Text == "+=" || token.Text == "-=" {
			p.Instructions = append(p.Instructions, Instruction{Op: Store, A: token.Text})
		}
	}
	if len(p.Instructions) == 0 {
		p.Instructions = []Instruction{{Op: Return}}
	}
	return p
}
