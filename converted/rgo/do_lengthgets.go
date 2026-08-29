package rgo

// do_lengthgets translates batch_146.c (src/main/builtin.c); R primitive(s): length<-.
func do_lengthgets(call, op, args, env Value) Value {
	return lengthGetsImpl(arg(args, 0), arg(args, 1))
}
