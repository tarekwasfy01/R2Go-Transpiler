package rgo

// do_rawToBits translates batch_201.c (src/main/raw.c); R primitive(s): rawToBits.
func do_rawToBits(call, op, args, env Value) Value {
	return rawToBitsImpl(arg(args, 0))
}
