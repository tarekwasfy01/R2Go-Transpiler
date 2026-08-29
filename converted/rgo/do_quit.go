package rgo

// do_quit translates batch_196.c (src/main/main.c); R primitive(s): quit.
func do_quit(call, op, args, env Value) Value {
	return ErrValue("quit requested")
}
