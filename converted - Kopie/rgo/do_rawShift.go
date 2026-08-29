package rgo

// do_rawShift translates batch_200.c (src/main/raw.c); R primitive(s): rawShift.
func do_rawShift(call, op, args, env Value) Value {
	return rawShiftImpl(arg(args, 0), arg(args, 1))
}
