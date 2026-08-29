package rgo

// do_pretty translates batch_189.c (src/main/util.c); R primitive(s): pretty.
func do_pretty(call, op, args, env Value) Value {
	return prettyImpl(arg(args, 0), arg(args, 1), arg(args, 2))
}
