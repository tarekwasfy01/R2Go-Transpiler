package rgo

// do_maxNSize translates batch_162.c (src/main/memory.c); R primitive(s): mem.maxNSize.
func do_maxNSize(call, op, args, env Value) Value {
	return systemPrimitive("do_maxNSize", args)
}
