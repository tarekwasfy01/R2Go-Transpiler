package rgo

// do_as_environment translates batch_023.c (src/main/envir.c); R primitive(s): as.environment.
func do_as_environment(call, op, args, env Value) Value {
	return envPrimitive("do_as_environment", args, env)
}
