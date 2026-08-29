package rgo

// do_AT translates batch_001.c (src/main/attrib.c); R primitive(s): @.
func do_AT(call, op, args, env Value) Value {
	return s4Impl("do_AT", args)
}
