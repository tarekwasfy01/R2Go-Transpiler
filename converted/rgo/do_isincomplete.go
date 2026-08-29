package rgo

// do_isincomplete translates batch_136.c (src/main/connections.c); R primitive(s): isIncomplete.
func do_isincomplete(call, op, args, env Value) Value {
	return connectionPrimitive("do_isincomplete", args)
}
