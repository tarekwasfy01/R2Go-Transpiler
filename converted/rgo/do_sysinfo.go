package rgo

// do_sysinfo translates batch_257.c (src/gnuwin32/extra.c); R primitive(s): Sys.info.
func do_sysinfo(call, op, args, env Value) Value {
	return systemPrimitive("do_sysinfo", args)
}
