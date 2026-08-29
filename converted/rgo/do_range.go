package rgo

// do_range translates batch_198.c (src/main/summary.c); R primitive(s): range.
func do_range(call, op, args, env Value) Value {
	return rangeImpl(arg(args, 0), arg(args, 1))
}
