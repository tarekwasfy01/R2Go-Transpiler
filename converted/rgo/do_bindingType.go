package rgo

// do_bindingType translates batch_033.c (src/main/envir.c); R primitive(s): getBindingType.
func do_bindingType(call, op, args, env Value) Value {
	return envPrimitive("do_bindingType", args, env)
}
