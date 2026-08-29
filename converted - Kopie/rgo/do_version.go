package rgo

// do_version translates batch_285.c (src/main/version.c); R primitive(s): Version.
func do_version(call, op, args, env Value) Value {
	return systemPrimitive("do_version", args)
}
