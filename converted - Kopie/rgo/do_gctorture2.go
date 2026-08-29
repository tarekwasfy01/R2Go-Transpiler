package rgo

// do_gctorture2 translates batch_108.c (src/main/memory.c); R primitive(s): gctorture2.
func do_gctorture2(call, op, args, env Value) Value {
	return gcPrimitive("do_gctorture2", args)
}
