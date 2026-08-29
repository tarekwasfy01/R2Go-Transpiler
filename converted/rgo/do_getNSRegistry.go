package rgo

// do_getNSRegistry translates batch_110.c (src/main/envir.c); R primitive(s): getNamespaceRegistry.
func do_getNSRegistry(call, op, args, env Value) Value {
	return envPrimitive("do_getNSRegistry", args, env)
}
