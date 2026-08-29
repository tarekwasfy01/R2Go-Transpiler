package rgo

// do_pos2env translates batch_188.c (src/main/envir.c); R primitive(s): pos.to.env.
func do_pos2env(call, op, args, env Value) Value {
	return envPrimitive("do_pos2env", args, env)
}
