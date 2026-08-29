package rgo

// do_refcnt translates batch_212.c (src/main/inspect.c); R primitive(s): refcnt.
func do_refcnt(call, op, args, env Value) Value {
	return Ints(boolInt(arg(args, 0).shared))
}
