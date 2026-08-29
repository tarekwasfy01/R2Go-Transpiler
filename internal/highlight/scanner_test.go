package highlight

import (
	"context"
	"testing"
)

func TestRHighlightingProducesSemanticScopes(t *testing.T) {
	tokens, err := Tokens(context.Background(), R, "# note\nf <- function(x) if (x > 1) print('yes')\ncustom(2)")
	if err != nil {
		t.Fatal(err)
	}
	wanted := map[string]bool{"comment": false, "keyword": false, "name.function": false, "name.builtin": false, "literal.string": false, "literal.number": false, "operator": false}
	for _, token := range tokens {
		if _, ok := wanted[string(token.Scope)]; ok {
			wanted[string(token.Scope)] = true
		}
	}
	for scope, found := range wanted {
		if !found {
			t.Errorf("missing scope %s in %#v", scope, tokens)
		}
	}
}

func TestGoHighlightingWorksForGeneratedFallback(t *testing.T) {
	tokens, err := Tokens(context.Background(), Go, "package main\n// R| x <- 1\nfunc main(){fmt.Println(2)}")
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) < 8 {
		t.Fatalf("too few tokens: %#v", tokens)
	}
}
