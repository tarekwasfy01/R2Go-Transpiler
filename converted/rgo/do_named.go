package rgo

// do_named translates batch_175.c (src/main/inspect.c); R primitive(s): named.
func do_named(call, op, args, env Value) Value {
	return Ints(boolInt(arg(args, 0).shared))
}
