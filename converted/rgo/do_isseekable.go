package rgo

// do_isseekable translates batch_140.c (src/main/connections.c); R primitive(s): isSeekable.
func do_isseekable(call, op, args, env Value) Value {
	return connectionPrimitive("do_isseekable", args)
}
