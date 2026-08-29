package rgo

// do_dotDelayedEnv translates batch_064.c (src/main/envir.c); R primitive(s): dotDelayedEnvironment.
func do_dotDelayedEnv(call, op, args, env Value) Value {
	return envPrimitive("do_dotDelayedEnv", args, env)
}
