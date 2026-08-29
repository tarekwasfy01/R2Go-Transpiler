package rgo

// do_sequence translates batch_227.c (src/main/seq.c); R primitive(s): sequence.
func do_sequence(call, op, args, env Value) Value {
	return sequenceImpl(arg(args, 0), arg(args, 1), arg(args, 2))
}
