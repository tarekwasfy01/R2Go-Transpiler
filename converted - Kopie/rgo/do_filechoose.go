package rgo

// do_filechoose translates batch_090.c (src/gnuwin32/extra.c); R primitive(s): file.choose.
func do_filechoose(call, op, args, env Value) Value {
	return filePrimitive("do_filechoose", args)
}
