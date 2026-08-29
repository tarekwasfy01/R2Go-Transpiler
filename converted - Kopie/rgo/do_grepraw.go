package rgo

// do_grepraw translates batch_122.c (src/main/grep.c); R primitive(s): grepRaw.
func do_grepraw(call, op, args, env Value) Value {
	return grepRawImpl(arg(args, 0), arg(args, 1), arg(args, 2), arg(args, 3), arg(args, 4), arg(args, 5), arg(args, 6))
}
