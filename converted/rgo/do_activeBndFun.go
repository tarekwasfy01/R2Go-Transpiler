package rgo

// do_activeBndFun translates batch_011.c (src/main/envir.c); R primitive(s): activeBindingFunction.
func do_activeBndFun(call, op, args, env Value) Value {
	return envPrimitive("do_activeBndFun", args, env)
}
