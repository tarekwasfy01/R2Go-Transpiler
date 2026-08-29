package rgo

// do_onexit translates batch_180.c (src/main/builtin.c); R primitive(s): on.exit.
func do_onexit(call, op, args, env Value) Value {
	return unsupported("do_onexit", "algorithm depends on GNU R internals not included in this batch corpus")
}
