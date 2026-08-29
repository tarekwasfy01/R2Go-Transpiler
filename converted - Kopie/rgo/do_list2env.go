package rgo

// do_list2env translates batch_147.c (src/main/envir.c); R primitive(s): list2env.
func do_list2env(call, op, args, env Value) Value {
	return envPrimitive("do_list2env", args, env)
}
