package rgo

// do_split translates batch_244.c (src/main/split.c); R primitive(s): split.
func do_split(call, op, args, env Value) Value {
	return splitImpl(arg(args, 0), arg(args, 1))
}
