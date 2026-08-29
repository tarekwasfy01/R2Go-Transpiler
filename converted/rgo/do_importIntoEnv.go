package rgo

// do_importIntoEnv translates batch_126.c (src/main/envir.c); R primitive(s): importIntoEnv.
func do_importIntoEnv(call, op, args, env Value) Value {
	return envPrimitive("do_importIntoEnv", args, env)
}
