package rgo

// do_intToBits translates batch_128.c (src/main/raw.c); R primitive(s): intToBits.
func do_intToBits(call, op, args, env Value) Value {
	return intToBitsImpl(arg(args, 0))
}
