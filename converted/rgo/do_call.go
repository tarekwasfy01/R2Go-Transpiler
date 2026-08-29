package rgo

// do_call translates batch_040.c (src/main/coerce.c); R primitive(s): call.
func do_call(call, op, args, env Value) Value {
	name, e := asString(arg(args, 0))
	if e != nil {
		return ErrValue("first argument must be a character string")
	}
	vals := append([]Value{Sym(name)}, arg(args, 1).V...)
	return withAttr(Lists(vals...), "class", Strings("call"))
}
