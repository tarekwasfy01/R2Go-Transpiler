package rgo

// do_args translates batch_022.c (src/main/builtin.c); R primitive(s): args.
func do_args(call, op, args, env Value) Value {
	f := arg(args, 0)
	if f.Kind != Function {
		return Nil
	}
	return withAttr(f, "body", Nil)
}
