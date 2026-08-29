package rgo

// do_envprofile translates batch_085.c (src/main/envir.c); R primitive(s): env.profile.
func do_envprofile(call, op, args, env Value) Value {
	return envPrimitive("do_envprofile", args, env)
}
