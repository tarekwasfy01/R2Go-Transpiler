package rgo

// do_packBits translates batch_182.c (src/main/raw.c); R primitive(s): packBits.
func do_packBits(call, op, args, env Value) Value {
	return packBitsImpl(arg(args, 0), arg(args, 1))
}
