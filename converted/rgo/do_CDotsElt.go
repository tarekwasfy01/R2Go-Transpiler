package rgo

// do_CDotsElt translates batch_002.c (src/main/envir.c); R primitive(s): dotsElt.
func do_CDotsElt(call, op, args, env Value) Value {
	i, e := asInt(arg(args, 0))
	if e != nil || i < 1 || int(i) > nargs(args) {
		return ErrValue("dotsElt: index out of bounds")
	}
	return arg(args, int(i)-1)
}
