package rgo

// do_unregNS translates batch_277.c (src/main/envir.c); R primitive(s): unregisterNamespace.
func do_unregNS(call, op, args, env Value) Value {
	return envPrimitive("do_unregNS", args, env)
}
