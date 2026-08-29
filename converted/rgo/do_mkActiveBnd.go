package rgo

// do_mkActiveBnd translates batch_169.c (src/main/envir.c); R primitive(s): makeActiveBinding.
func do_mkActiveBnd(call, op, args, env Value) Value {
	return envPrimitive("do_mkActiveBnd", args, env)
}
