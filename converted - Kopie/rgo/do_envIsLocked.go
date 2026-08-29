package rgo

// do_envIsLocked translates batch_082.c (src/main/envir.c); R primitive(s): environmentIsLocked.
func do_envIsLocked(call, op, args, env Value) Value {
	return envPrimitive("do_envIsLocked", args, env)
}
