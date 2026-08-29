package rgo

// do_nargs translates batch_176.c (src/main/util.c); R primitive(s): nargs.
func do_nargs(call, op, args, env Value) Value {
	return Ints(int64(nargs(args)))
}
