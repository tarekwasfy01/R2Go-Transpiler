package rgo

// do_url translates batch_281.c (src/main/connections.c); R primitive(s): file, url.
func do_url(call, op, args, env Value) Value {
	return connectionPrimitive("do_url", args)
}
