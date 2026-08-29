package rgo

// do_standardGeneric translates batch_245.c (src/main/objects.c); R primitive(s): standardGeneric.
func do_standardGeneric(call, op, args, env Value) Value {
	return unsupported("do_standardGeneric", "requires the evaluator/context stack; isolated batch corpus does not define enough evaluator semantics")
}
