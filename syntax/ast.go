package syntax

// Position is a byte-oriented source location. Line and Column are one-based.
type Position struct {
	Offset int `json:"offset"`
	Line   int `json:"line"`
	Column int `json:"column"`
}

type Span struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type Expr interface {
	exprNode()
	SourceSpan() Span
}

type Program struct {
	Expressions []Expr `json:"expressions"`
}

type LiteralKind string

const (
	NumberLiteral  LiteralKind = "number"
	StringLiteral  LiteralKind = "string"
	LogicalLiteral LiteralKind = "logical"
	NullLiteral    LiteralKind = "null"
	NALiteral      LiteralKind = "na"
)

type Literal struct {
	Kind LiteralKind `json:"kind"`
	Text string      `json:"text"`
	At   Span        `json:"span"`
}

func (*Literal) exprNode()          {}
func (e *Literal) SourceSpan() Span { return e.At }

type Symbol struct {
	Name string `json:"name"`
	At   Span   `json:"span"`
}

func (*Symbol) exprNode()          {}
func (e *Symbol) SourceSpan() Span { return e.At }

// Call is deliberately close to R's language-object model: operators and
// control forms are calls too, so quote/substitute/eval can preserve them.
type Call struct {
	Function  Expr       `json:"function"`
	Arguments []Argument `json:"arguments"`
	At        Span       `json:"span"`
}

type Argument struct {
	Name  string `json:"name,omitempty"`
	Value Expr   `json:"value,omitempty"`
	At    Span   `json:"span"`
}

func (*Call) exprNode()          {}
func (e *Call) SourceSpan() Span { return e.At }

type Block struct {
	Expressions []Expr `json:"expressions"`
	At          Span   `json:"span"`
}

func (*Block) exprNode()          {}
func (e *Block) SourceSpan() Span { return e.At }

type Function struct {
	Parameters []Parameter `json:"parameters"`
	Body       Expr        `json:"body"`
	At         Span        `json:"span"`
}

type Parameter struct {
	Name    string `json:"name"`
	Default Expr   `json:"default,omitempty"`
	At      Span   `json:"span"`
}

func (*Function) exprNode()          {}
func (e *Function) SourceSpan() Span { return e.At }

type If struct {
	Condition Expr `json:"condition"`
	Then      Expr `json:"then"`
	Else      Expr `json:"else,omitempty"`
	At        Span `json:"span"`
}

func (*If) exprNode()          {}
func (e *If) SourceSpan() Span { return e.At }

type While struct {
	Condition Expr `json:"condition"`
	Body      Expr `json:"body"`
	At        Span `json:"span"`
}

func (*While) exprNode()          {}
func (e *While) SourceSpan() Span { return e.At }

type For struct {
	Variable string `json:"variable"`
	Sequence Expr   `json:"sequence"`
	Body     Expr   `json:"body"`
	At       Span   `json:"span"`
}

func (*For) exprNode()          {}
func (e *For) SourceSpan() Span { return e.At }

type Repeat struct {
	Body Expr `json:"body"`
	At   Span `json:"span"`
}

func (*Repeat) exprNode()          {}
func (e *Repeat) SourceSpan() Span { return e.At }
