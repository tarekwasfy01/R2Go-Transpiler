package rgo

// do_allnames translates batch_019.c (src/main/list.c); R primitive(s): all.names.
func do_allnames(call, op, args, env Value) Value {
	return allNamesImpl(arg(args, 0), arg(args, 1), arg(args, 2), arg(args, 3))
}
