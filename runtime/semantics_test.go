package runtime

import (
	"math"
	"testing"
)

func TestCommonKindMatrix(t *testing.T) {
	cases := []struct{ a, b, want Kind }{
		{LogicalKind, IntegerKind, IntegerKind},
		{IntegerKind, DoubleKind, DoubleKind},
		{DoubleKind, ComplexKind, ComplexKind},
		{ComplexKind, CharacterKind, CharacterKind},
		{RawKind, LogicalKind, LogicalKind},
		{CharacterKind, ListKind, ListKind},
	}
	for _, tc := range cases {
		got, err := CommonKind(tc.a, tc.b)
		if err != nil || got != tc.want {
			t.Fatalf("CommonKind(%s,%s)=%s,%v want %s", tc.a, tc.b, got, err, tc.want)
		}
	}
}

func TestReplacePromotesTarget(t *testing.T) {
	src := &IntegerVector{Data: []int64{1, 2, 3}, Missing: []bool{false, false, false}}
	out, err := replacePositions(src, []int{1}, &DoubleVector{Data: []float64{2.5}})
	if err != nil {
		t.Fatal(err)
	}
	d, ok := out.(*DoubleVector)
	if !ok {
		t.Fatalf("expected promoted double vector, got %T", out)
	}
	if d.Data[0] != 1 || d.Data[1] != 2.5 || d.Data[2] != 3 {
		t.Fatalf("unexpected promoted data: %#v", d.Data)
	}
}

func TestMatchUsesCommonNumericKind(t *testing.T) {
	got := matchValues(&IntegerVector{Data: []int64{1}}, &DoubleVector{Data: []float64{1}}).(*IntegerVector)
	if len(got.Data) != 1 || got.Data[0] != 1 || got.Missing[0] {
		t.Fatalf("integer 1 should match double 1: %#v", got)
	}
}

func TestMatchKeyDistinguishesNAAndNaN(t *testing.T) {
	na := &DoubleVector{Data: []float64{NAReal()}, Missing: []bool{true}}
	nan := &DoubleVector{Data: []float64{math.NaN()}, Missing: []bool{false}}
	if MatchKey(na) == MatchKey(nan) {
		t.Fatalf("NA and NaN must use different match keys")
	}
}
