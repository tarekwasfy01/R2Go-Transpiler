package rgo

// do_delayedBindingEnv translates batch_056.c (src/main/envir.c); R primitive(s): delayedBindingEnvironment.
func do_delayedBindingEnv(call, op, args, env Value) Value {
	return envPrimitive("do_delayedBindingEnv", args, env)
}
