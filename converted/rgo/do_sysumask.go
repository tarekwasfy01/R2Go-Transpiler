package rgo

// do_sysumask translates batch_260.c (src/main/platform.c); R primitive(s): Sys.umask.
func do_sysumask(call, op, args, env Value) Value {
	return systemPrimitive("do_sysumask", args)
}
