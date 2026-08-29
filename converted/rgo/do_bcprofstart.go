package rgo

// do_bcprofstart translates batch_029.c (src/main/eval.c); R primitive(s): bcprofstart.
func do_bcprofstart(call, op, args, env Value) Value {
	return unsupported("do_bcprofstart", "requires the evaluator/context stack; isolated batch corpus does not define enough evaluator semantics")
}
