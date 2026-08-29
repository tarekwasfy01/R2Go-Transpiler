package rgo

// do_ascall translates batch_024.c (src/main/coerce.c); R primitive(s): as.call.
func do_ascall(call, op, args, env Value) Value {
	x := arg(args, 0)
	if x.Kind != List {
		return ErrValue("invalid argument list")
	}
	return withAttr(x, "class", Strings("call"))
}
