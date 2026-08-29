package rgo

// do_encodeString translates batch_080.c (src/main/util.c); R primitive(s): encodeString.
func do_encodeString(call, op, args, env Value) Value {
	return encodeStringImpl(arg(args, 0), arg(args, 2), arg(args, 4), Nil, arg(args, 3), arg(args, 1))
}
