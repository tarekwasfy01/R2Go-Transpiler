package rgo

// do_forcedBindingExpr translates batch_102.c (src/main/envir.c); R primitive(s): forcedBindingExpression.
func do_forcedBindingExpr(call, op, args, env Value) Value {
	return envPrimitive("do_forcedBindingExpr", args, env)
}
