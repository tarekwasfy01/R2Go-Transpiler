package rgo

// do_address translates batch_016.c (src/main/inspect.c); R primitive(s): address.
func do_address(call, op, args, env Value) Value {
	return addressImpl(arg(args, 0))
}
