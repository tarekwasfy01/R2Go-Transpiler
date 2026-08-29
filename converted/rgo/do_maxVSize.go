package rgo

// do_maxVSize translates batch_163.c (src/main/memory.c); R primitive(s): mem.maxVSize.
func do_maxVSize(call, op, args, env Value) Value {
	return systemPrimitive("do_maxVSize", args)
}
