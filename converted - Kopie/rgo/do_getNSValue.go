package rgo

// do_getNSValue translates batch_111.c (src/main/envir.c); R primitive(s): getNamespaceValue.
func do_getNSValue(call, op, args, env Value) Value {
	return envPrimitive("do_getNSValue", args, env)
}
