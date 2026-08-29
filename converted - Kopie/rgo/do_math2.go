package rgo

// do_math2 translates batch_160.c (src/main/arithmetic.c); R primitive(s): besselJ, besselY, psigamma.
func do_math2(call, op, args, env Value) Value {
	off, _ := asInt(op)
	return math2Impl(arg(args, 0), arg(args, 1), int(off))
}
