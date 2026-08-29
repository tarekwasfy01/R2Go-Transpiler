package rgo

// do_savefile translates batch_223.c (src/main/saveload.c); R primitive(s): save.to.file.
func do_savefile(call, op, args, env Value) Value {
	return unsupported("do_savefile", "algorithm depends on GNU R internals not included in this batch corpus")
}
