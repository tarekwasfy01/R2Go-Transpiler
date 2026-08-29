package rgo

// do_builtins translates batch_039.c (src/main/envir.c); R primitive(s): builtins.
func do_builtins(call, op, args, env Value) Value {
	return envPrimitive("do_builtins", args, env)
}
