package rgo

// do_env2list translates batch_081.c (src/main/envir.c); R primitive(s): env2list.
func do_env2list(call, op, args, env Value) Value {
	return envPrimitive("do_env2list", args, env)
}
