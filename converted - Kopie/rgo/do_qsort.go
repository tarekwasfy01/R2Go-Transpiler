package rgo

// do_qsort translates batch_195.c (src/main/qsort.c); R primitive(s): qsort.
func do_qsort(call, op, args, env Value) Value {
	return qsortImpl(arg(args, 0), arg(args, 1), arg(args, 2))
}
