package rgo

// do_memoryprofile translates batch_166.c (src/main/memory.c); R primitive(s): memory.profile.
func do_memoryprofile(call, op, args, env Value) Value {
	return gcPrimitive("do_memoryprofile", args)
}
