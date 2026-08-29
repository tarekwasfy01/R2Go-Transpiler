package rgo

// do_numToBits translates batch_177.c (src/main/raw.c); R primitive(s): numToBits.
func do_numToBits(call, op, args, env Value) Value {
	return numToBitsImpl(arg(args, 0))
}
