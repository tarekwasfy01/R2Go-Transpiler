package rgo

// do_tzone_name translates batch_276.c (src/gnuwin32/extra.c); R primitive(s): tzone_name.
func do_tzone_name(call, op, args, env Value) Value {
	return systemPrimitive("do_tzone_name", args)
}
