package rgo

// do_returnValue translates batch_220.c (src/main/eval.c); R primitive(s): returnValue.
func do_returnValue(call, op, args, env Value) Value {
	return arg(args, 0)
}
