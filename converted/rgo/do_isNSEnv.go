package rgo

// do_isNSEnv translates batch_133.c (src/main/envir.c); R primitive(s): isNamespaceEnv.
func do_isNSEnv(call, op, args, env Value) Value {
	return envPrimitive("do_isNSEnv", args, env)
}
