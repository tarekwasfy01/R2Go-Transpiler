package rgo

// do_printdefault translates batch_191.c (src/main/print.c); R primitive(s): print.default.
func do_printdefault(call, op, args, env Value) Value {
	return printDefaultImpl(arg(args, 0))
}
