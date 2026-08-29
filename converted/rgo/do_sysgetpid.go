package rgo

// do_sysgetpid translates batch_256.c (src/main/platform.c); R primitive(s): Sys.getpid.
func do_sysgetpid(call, op, args, env Value) Value {
	return systemPrimitive("do_sysgetpid", args)
}
