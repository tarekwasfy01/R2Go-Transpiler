package rgo

// do_readchar translates batch_207.c (src/main/connections.c); R primitive(s): readChar.
func do_readchar(call, op, args, env Value) Value {
	return connectionPrimitive("do_readchar", args)
}
