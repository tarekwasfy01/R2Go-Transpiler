package rgo

// do_setmaxnumthreads translates batch_234.c (src/main/eval.c); R primitive(s): setMaxNumMathThreads.
func do_setmaxnumthreads(call, op, args, env Value) Value {
	return systemPrimitive("do_setmaxnumthreads", args)
}
