package rgo

// do_inspect translates batch_127.c (src/main/inspect.c); R primitive(s): inspect.
func do_inspect(call, op, args, env Value) Value {
	return unsupported("do_inspect", "algorithm depends on GNU R internals not included in this batch corpus")
}
