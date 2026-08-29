package rgo

// do_mget translates batch_168.c (src/main/envir.c); R primitive(s): mget.
func do_mget(call, op, args, env Value) Value {
	return envPrimitive("do_mget", args, env)
}
