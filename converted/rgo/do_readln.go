package rgo

// do_readln translates batch_209.c (src/main/scan.c); R primitive(s): readline.
func do_readln(call, op, args, env Value) Value {
	return unsupported("do_readln", "algorithm depends on GNU R internals not included in this batch corpus")
}
