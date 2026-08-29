package rgo

// do_dotType translates batch_067.c (src/main/envir.c); R primitive(s): getDotType.
func do_dotType(call, op, args, env Value) Value {
	return envPrimitive("do_dotType", args, env)
}
