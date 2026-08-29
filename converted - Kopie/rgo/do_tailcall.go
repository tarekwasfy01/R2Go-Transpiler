package rgo

// do_tailcall translates batch_263.c (src/main/eval.c); R primitive(s): Exec, Tailcall.
func do_tailcall(call, op, args, env Value) Value {
	return unsupported("do_tailcall", "requires the evaluator/context stack; isolated batch corpus does not define enough evaluator semantics")
}
