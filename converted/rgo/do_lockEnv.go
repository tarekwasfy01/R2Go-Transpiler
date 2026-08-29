package rgo

// do_lockEnv translates batch_153.c (src/main/envir.c); R primitive(s): lockEnvironment.
func do_lockEnv(call, op, args, env Value) Value {
	return envPrimitive("do_lockEnv", args, env)
}
