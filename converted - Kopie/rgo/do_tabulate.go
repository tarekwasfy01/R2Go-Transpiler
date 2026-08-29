package rgo

// do_tabulate translates batch_262.c (src/main/util.c); R primitive(s): tabulate.
func do_tabulate(call, op, args, env Value) Value {
	return tabulateImpl(arg(args, 0), arg(args, 1))
}
