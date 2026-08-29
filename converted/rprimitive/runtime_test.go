package rprimitive

import "testing"

func TestTranslationInventory(t *testing.T) {
	if len(PrimitiveTable) != 158 {
		t.Fatalf("primitive table: got %d want 158", len(PrimitiveTable))
	}
}
func TestRuntimeOperators(t *testing.T) {
	r := NewRuntime()
	r.InstallCoreCompatibility()
	if !r.Truth(r.Binary("==", r.Const("int", "2"), r.Const("int", "2"))) {
		t.Fatal("equality")
	}
	if got := r.Call("PRIMVAL", PrimitiveOp{Offset: 41}); asInt(got) != 41 {
		t.Fatalf("PRIMVAL=%v", got)
	}
}
