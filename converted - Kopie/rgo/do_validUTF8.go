package rgo

// do_validUTF8 translates batch_284.c (src/main/util.c); R primitive(s): validUTF8.
func do_validUTF8(call, op, args, env Value) Value {
	return validUTF8Impl(arg(args, 0))
}
