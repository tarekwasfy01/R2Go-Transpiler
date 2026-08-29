package rgo

// do_dump translates batch_075.c (src/main/deparse.c); R primitive(s): dump.
func do_dump(call, op, args, env Value) Value {
	return unsupported("do_dump", "algorithm depends on GNU R internals not included in this batch corpus")
}
