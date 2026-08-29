package rgo

// do_numToInts translates batch_178.c (src/main/raw.c); R primitive(s): numToInts.
func do_numToInts(call, op, args, env Value) Value {
	return numToIntsImpl(arg(args, 0))
}
