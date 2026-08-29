package rgo

// do_intToUtf8 translates batch_129.c (src/main/raw.c); R primitive(s): intToUtf8.
func do_intToUtf8(call, op, args, env Value) Value {
	return intToUtf8Impl(arg(args, 0), arg(args, 1), arg(args, 2))
}
