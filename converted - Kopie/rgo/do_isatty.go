package rgo

// do_isatty translates batch_135.c (src/main/connections.c); R primitive(s): isatty.
func do_isatty(call, op, args, env Value) Value {
	return systemPrimitive("do_isatty", args)
}
