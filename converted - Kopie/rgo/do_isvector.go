package rgo

// do_isvector translates batch_141.c (src/main/coerce.c); R primitive(s): is.vector.
func do_isvector(call, op, args, env Value) Value {
	return isVectorImpl(arg(args, 0), arg(args, 1))
}
