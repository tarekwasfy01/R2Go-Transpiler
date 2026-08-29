package rgo

// do_commandArgs translates batch_044.c (src/main/CommandLineArgs.c); R primitive(s): commandArgs.
func do_commandArgs(call, op, args, env Value) Value {
	return systemPrimitive("do_commandArgs", args)
}
