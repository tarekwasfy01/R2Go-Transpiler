package rgo

// do_is_builtin_internal translates batch_134.c (src/main/eval.c); R primitive(s): is.builtin.internal.
func do_is_builtin_internal(call, op, args, env Value) Value {
	return unsupported("do_is_builtin_internal", "requires the evaluator/context stack; isolated batch corpus does not define enough evaluator semantics")
}
