package rgo

// do_getlocale translates batch_121.c (src/main/platform.c); R primitive(s): Sys.getlocale.
func do_getlocale(call, op, args, env Value) Value {
	return systemPrimitive("do_getlocale", args)
}
