package rgo

// do_systime translates batch_259.c (src/main/times.c); R primitive(s): Sys.time.
func do_systime(call, op, args, env Value) Value {
	return systemPrimitive("do_systime", args)
}
