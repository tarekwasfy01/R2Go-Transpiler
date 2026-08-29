package rgo

// do_textconvalue translates batch_267.c (src/main/connections.c); R primitive(s): textConnectionValue.
func do_textconvalue(call, op, args, env Value) Value {
	return connectionPrimitive("do_textconvalue", args)
}
