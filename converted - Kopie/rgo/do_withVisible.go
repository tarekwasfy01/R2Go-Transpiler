package rgo

// do_withVisible translates batch_286.c (src/main/eval.c); R primitive(s): withVisible.
func do_withVisible(call, op, args, env Value) Value {
	return Lists(arg(args, 0), Bool(true))
}
