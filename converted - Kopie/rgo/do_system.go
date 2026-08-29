package rgo

// do_system translates batch_258.c (src/gnuwin32/sys-win32.c); R primitive(s): system.
func do_system(call, op, args, env Value) Value {
	return systemImpl(arg(args, 0), arg(args, 1))
}
