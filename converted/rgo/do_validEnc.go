package rgo

// do_validEnc translates batch_283.c (src/main/util.c); R primitive(s): validEnc.
func do_validEnc(call, op, args, env Value) Value {
	return validUTF8Impl(arg(args, 0))
}
