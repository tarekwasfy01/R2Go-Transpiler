package rgo

// do_syswhich translates batch_261.c (src/gnuwin32/run.c); R primitive(s): Sys.which.
func do_syswhich(call, op, args, env Value) Value {
	return systemPrimitive("do_syswhich", args)
}
