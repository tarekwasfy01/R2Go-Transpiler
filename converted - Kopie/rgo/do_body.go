package rgo

// do_body translates batch_037.c (src/main/builtin.c); R primitive(s): body.
func do_body(call, op, args, env Value) Value {
	f := arg(args, 0)
	if f.Kind != Function {
		return Nil
	}
	return attr(f, "body")
}
