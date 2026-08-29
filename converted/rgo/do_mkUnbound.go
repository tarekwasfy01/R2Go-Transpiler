package rgo

// do_mkUnbound translates batch_170.c (src/main/envir.c); R primitive(s): mkUnbound.
func do_mkUnbound(call, op, args, env Value) Value {
	return envPrimitive("do_mkUnbound", args, env)
}
