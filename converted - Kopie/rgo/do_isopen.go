package rgo

// do_isopen translates batch_139.c (src/main/connections.c); R primitive(s): isOpen.
func do_isopen(call, op, args, env Value) Value {
	return connectionPrimitive("do_isopen", args)
}
