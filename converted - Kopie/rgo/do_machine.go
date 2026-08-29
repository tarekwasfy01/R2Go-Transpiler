package rgo

// do_machine translates batch_155.c (src/gnuwin32/sys-win32.c); R primitive(s): machine.
func do_machine(call, op, args, env Value) Value {
	return systemPrimitive("do_machine", args)
}
