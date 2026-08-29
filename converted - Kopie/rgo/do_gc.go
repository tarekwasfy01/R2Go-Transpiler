package rgo

// do_gc translates batch_104.c (src/main/memory.c); R primitive(s): gc.
func do_gc(call, op, args, env Value) Value {
	return gcPrimitive("do_gc", args)
}
