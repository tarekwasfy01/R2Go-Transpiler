package rgo

// do_bincode translates batch_032.c (src/main/util.c); R primitive(s): bincode.
func do_bincode(call, op, args, env Value) Value {
	return bincodeImpl(arg(args, 0), arg(args, 1), arg(args, 2), arg(args, 3))
}
