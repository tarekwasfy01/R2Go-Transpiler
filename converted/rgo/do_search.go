package rgo

// do_search translates batch_225.c (src/main/envir.c); R primitive(s): search.
func do_search(call, op, args, env Value) Value {
	return envPrimitive("do_search", args, env)
}
