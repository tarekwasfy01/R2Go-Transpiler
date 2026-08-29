package rgo

// do_tempfile translates batch_265.c (src/main/sysutils.c); R primitive(s): tempfile.
func do_tempfile(call, op, args, env Value) Value {
	return systemPrimitive("do_tempfile", args)
}
