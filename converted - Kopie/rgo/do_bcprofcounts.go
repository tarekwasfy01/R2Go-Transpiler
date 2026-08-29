package rgo

// do_bcprofcounts translates batch_028.c (src/main/eval.c); R primitive(s): bcprofcounts.
func do_bcprofcounts(call, op, args, env Value) Value {
	return unsupported("do_bcprofcounts", "requires the evaluator/context stack; isolated batch corpus does not define enough evaluator semantics")
}
