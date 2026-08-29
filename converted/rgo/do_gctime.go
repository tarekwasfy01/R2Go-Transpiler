package rgo

// do_gctime translates batch_106.c (src/main/memory.c); R primitive(s): gc.time.
func do_gctime(call, op, args, env Value) Value {
	return gcPrimitive("do_gctime", args)
}
