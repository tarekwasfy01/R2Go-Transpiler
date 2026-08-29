package rgo

// do_regNS translates batch_214.c (src/main/envir.c); R primitive(s): registerNamespace.
func do_regNS(call, op, args, env Value) Value {
	return envPrimitive("do_regNS", args, env)
}
