package rgo

// do_rawconvalue translates batch_203.c (src/main/connections.c); R primitive(s): rawConnectionValue.
func do_rawconvalue(call, op, args, env Value) Value {
	return connectionPrimitive("do_rawconvalue", args)
}
