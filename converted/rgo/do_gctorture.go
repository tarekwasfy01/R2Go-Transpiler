package rgo

// do_gctorture translates batch_107.c (src/main/memory.c); R primitive(s): gctorture.
func do_gctorture(call, op, args, env Value) Value {
	return gcPrimitive("do_gctorture", args)
}
