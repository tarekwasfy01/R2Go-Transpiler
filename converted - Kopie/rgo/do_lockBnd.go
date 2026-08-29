package rgo

// do_lockBnd translates batch_152.c (src/main/envir.c); R primitive(s): lockBinding, unlockBinding.
func do_lockBnd(call, op, args, env Value) Value {
	return envPrimitive("do_lockBnd", args, env)
}
