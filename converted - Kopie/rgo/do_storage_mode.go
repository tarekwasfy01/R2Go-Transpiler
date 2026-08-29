package rgo

// do_storage_mode translates batch_249.c (src/main/coerce.c); R primitive(s): storage.mode<-.
func do_storage_mode(call, op, args, env Value) Value {
	x := clone(arg(args, 0))
	mode, _ := asString(arg(args, 1))
	return coerceMode(x, mode)
}
