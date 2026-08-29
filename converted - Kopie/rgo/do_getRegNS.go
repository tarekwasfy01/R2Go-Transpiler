package rgo

// do_getRegNS translates batch_112.c (src/main/envir.c); R primitive(s): getRegisteredNamespace, isRegisteredNamespace.
func do_getRegNS(call, op, args, env Value) Value {
	return envPrimitive("do_getRegNS", args, env)
}
