package cir

// Op is a portable, side-effect-explicit C subset instruction.  The complete
// 355-entry lowering matrix will reference programs made from these ops.
type Op uint8

const (
	Load Op = iota + 1
	Store
	Constant
	Call
	Binary
	Unary
	Branch
	Jump
	Return
	Allocate
	Index
	SetIndex
)

type Instruction struct {
	Op      Op
	A, B, C string
	Target  int
}
type Program struct {
	Name         string
	SourceFile   string
	Instructions []Instruction
}

// FunctionBody returns the balanced token range following the chosen C
// function declaration.  It deliberately accepts attributes/prototypes before
// the opening brace and is reusable by the later statement parser.
func FunctionBody(tokens []Token, name string) ([]Token, bool) {
	for i, t := range tokens {
		if t.Kind != Identifier || t.Text != name {
			continue
		}
		depth := 0
		open := -1
		for j := i + 1; j < len(tokens); j++ {
			if tokens[j].Text == "(" {
				depth++
			}
			if tokens[j].Text == ")" {
				depth--
			}
			if depth == 0 && tokens[j].Text == "{" {
				open = j
				break
			}
			if depth == 0 && tokens[j].Text == ";" {
				break
			}
		}
		if open < 0 {
			continue
		}
		depth = 1
		for j := open + 1; j < len(tokens); j++ {
			if tokens[j].Text == "{" {
				depth++
			}
			if tokens[j].Text == "}" {
				depth--
				if depth == 0 {
					return tokens[open+1 : j], true
				}
			}
		}
	}
	return nil, false
}
