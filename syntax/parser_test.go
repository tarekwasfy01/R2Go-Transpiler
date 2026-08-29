package syntax

import "testing"

func TestUTF8BOMAtStartIsIgnored(t *testing.T) {
	p, err := Parse("\uFEFFx <- 1\nx + 2")
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Expressions) != 2 {
		t.Fatalf("expressions=%d", len(p.Expressions))
	}
}

func TestNamespaceCallBindsBeforeArguments(t *testing.T) {
	program, err := Parse("base::sum(1:4)")
	if err != nil {
		t.Fatal(err)
	}
	outer, ok := program.Expressions[0].(*Call)
	if !ok {
		t.Fatalf("outer expression is %T", program.Expressions[0])
	}
	qualified, ok := outer.Function.(*Call)
	if !ok {
		t.Fatalf("call function is %T", outer.Function)
	}
	op, _ := qualified.Function.(*Symbol)
	name, _ := qualified.Arguments[1].Value.(*Symbol)
	if op == nil || op.Name != "::" || name == nil || name.Name != "sum" {
		t.Fatalf("unexpected namespace AST: %#v", qualified)
	}
}

func TestNativePipeRewritesToCall(t *testing.T) {
	program, err := Parse("1:5 |> sum()")
	if err != nil {
		t.Fatal(err)
	}
	call, ok := program.Expressions[0].(*Call)
	if !ok {
		t.Fatalf("pipe expression is %T", program.Expressions[0])
	}
	fn, _ := call.Function.(*Symbol)
	if fn == nil || fn.Name != "sum" || len(call.Arguments) != 1 {
		t.Fatalf("pipe was not rewritten to sum(lhs): %#v", call)
	}
}

func TestParseUniversalCoreShapes(t *testing.T) {
	src := "f <- function(x, y = 2) {\n  if (x > y) x + y else x * y\n}\nfor (i in 1:3) print(f(i))\nquoted <- quote(a + 1)\n"
	p, err := Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Expressions) != 3 {
		t.Fatalf("expressions=%d", len(p.Expressions))
	}
}

func TestLexDoesNotConsumeAdditionAsNumber(t *testing.T) {
	ts, err := Lex("1+2")
	if err != nil {
		t.Fatal(err)
	}
	if len(ts) < 4 || ts[0].Text != "1" || ts[1].Text != "+" || ts[2].Text != "2" {
		t.Fatalf("tokens=%#v", ts)
	}
}

func TestParseDoubleBracket(t *testing.T) {
	if _, err := Parse("x[[1]]"); err != nil {
		t.Fatal(err)
	}
}

func TestParseMissingArraySubscripts(t *testing.T) {
	p, err := Parse("x[, 2]\nx[1, ]")
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Expressions) != 2 {
		t.Fatalf("expressions=%d", len(p.Expressions))
	}
}

func TestDollarBindsBeforeFollowingSubset(t *testing.T) {
	p, err := Parse("d$id[2]")
	if err != nil {
		t.Fatal(err)
	}
	outer, ok := p.Expressions[0].(*Call)
	if !ok {
		t.Fatal("expected outer call")
	}
	if fn, ok := outer.Function.(*Symbol); !ok || fn.Name != "[" {
		t.Fatalf("outer=%#v", outer.Function)
	}
}
