package rgo

// do_math3 translates batch_161.c (src/main/arithmetic.c); R primitive(s): besselI, besselK.
func do_math3(call, op, args, env Value) Value {
	off, _ := asInt(op)
	return math3Impl(arg(args, 0), arg(args, 1), arg(args, 2), int(off))
}
