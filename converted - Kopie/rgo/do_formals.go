package rgo

// do_formals translates batch_103.c (src/main/builtin.c); R primitive(s): formals.
func do_formals(call, op, args, env Value) Value {
	f := arg(args, 0)
	if f.Kind != Function {
		return Nil
	}
	return attr(f, "formals")
}
