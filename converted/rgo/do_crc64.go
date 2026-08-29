package rgo

// do_crc64 translates batch_049.c (src/main/util.c); R primitive(s): crc64.
func do_crc64(call, op, args, env Value) Value {
	return crc64Impl(arg(args, 0))
}
