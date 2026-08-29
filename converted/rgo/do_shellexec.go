package rgo

// do_shellexec translates batch_236.c (src/gnuwin32/extra.c); R primitive(s): shell.exec.
func do_shellexec(call, op, args, env Value) Value {
	return unsupported("do_shellexec", "algorithm depends on GNU R internals not included in this batch corpus")
}
