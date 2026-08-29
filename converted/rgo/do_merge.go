package rgo

// do_merge translates batch_167.c (src/main/util.c); R primitive(s): merge.
func do_merge(call, op, args, env Value) Value {
	return mergeImpl(arg(args, 0), arg(args, 1), arg(args, 2), arg(args, 3))
}
