package rgo

// do_complex translates batch_047.c (src/main/complex.c); R primitive(s): complex.
func do_complex(call, op, args, env Value) Value {
	n, e := asInt(arg(args, 0))
	if e != nil || n < 0 {
		return ErrValue("invalid length")
	}
	re, im := arg(args, 1), arg(args, 2)
	z := complexImpl(re, im)
	if int(n) > length(z) {
		z = lengthGetsImpl(z, Ints(n))
		for i := range z.NA {
			z.NA[i] = false
		}
	}
	return z
}
