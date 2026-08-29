package rgo

// do_gcinfo translates batch_105.c (src/main/memory.c); R primitive(s): gcinfo.
func do_gcinfo(call, op, args, env Value) Value {
	return statePrimitive("do_gcinfo", args)
}
