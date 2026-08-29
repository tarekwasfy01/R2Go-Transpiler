package ui

import (
	"testing"

	"github.com/oligo/gvcode/textstyle/syntax"
)

func TestColorSchemeIsOpaqueBlackAndColorful(t *testing.T) {
	scheme := codeColorScheme(true, true)
	if got := scheme.Foreground.NRGBA(); got.R != 0 || got.G != 0 || got.B != 0 || got.A != 255 {
		t.Fatalf("foreground=%v, want opaque black", got)
	}
	seen := map[[3]uint8]bool{}
	for _, scope := range []syntax.StyleScope{"keyword", "name.function", "name.builtin", "literal.string", "literal.number", "comment", "operator", "punctuation"} {
		style := scheme.GetTokenStyle(scope)
		got := scheme.GetColor(style.Foreground()).NRGBA()
		if got.A != 255 {
			t.Errorf("scope %s has alpha %d, want 255", scope, got.A)
		}
		seen[[3]uint8{got.R, got.G, got.B}] = true
	}
	if len(seen) < 6 {
		t.Fatalf("syntax palette has only %d distinct colors", len(seen))
	}
}
