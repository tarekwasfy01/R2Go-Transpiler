package rgo

// do_setS4Object translates batch_232.c (src/main/objects.c); R primitive(s): setS4Object.
func do_setS4Object(call, op, args, env Value) Value {
	return s4Impl("do_setS4Object", args)
}
