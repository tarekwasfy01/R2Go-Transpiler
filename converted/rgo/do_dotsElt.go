package rgo

// do_dotsElt translates batch_070.c (src/main/envir.c); R primitive(s): ...elt.
func do_dotsElt(call, op, args, env Value) Value {
	i, e := asInt(arg(args, 0))
	if e != nil || i < 1 || int(i) > nargs(args) {
		return ErrValue("...elt: index out of bounds")
	}
	return arg(args, int(i)-1)
}
