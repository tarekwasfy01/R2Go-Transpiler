package rgo

// do_tempdir translates batch_264.c (src/main/sysutils.c); R primitive(s): tempdir.
func do_tempdir(call, op, args, env Value) Value {
	return systemPrimitive("do_tempdir", args)
}
