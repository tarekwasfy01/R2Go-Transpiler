package rgo

// do_sockselect translates batch_242.c (src/main/connections.c); R primitive(s): sockSelect.
func do_sockselect(call, op, args, env Value) Value {
	return connectionPrimitive("do_sockselect", args)
}
