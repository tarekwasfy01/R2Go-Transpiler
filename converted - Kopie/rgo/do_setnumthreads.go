package rgo

// do_setnumthreads translates batch_235.c (src/main/eval.c); R primitive(s): setNumMathThreads.
func do_setnumthreads(call, op, args, env Value) Value {
	return systemPrimitive("do_setnumthreads", args)
}
