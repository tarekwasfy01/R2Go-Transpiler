package rgo

// do_copyDFattr translates batch_048.c (src/main/attrib.c); R primitive(s): copyDFattr.
func do_copyDFattr(call, op, args, env Value) Value {
	return copyAttrs(arg(args, 0), arg(args, 1))
}
