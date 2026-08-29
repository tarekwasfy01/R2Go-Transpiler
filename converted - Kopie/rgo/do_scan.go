package rgo

// do_scan translates batch_224.c (src/main/scan.c); R primitive(s): scan.
func do_scan(call, op, args, env Value) Value {
	return unsupported("do_scan", "algorithm depends on GNU R internals not included in this batch corpus")
}
