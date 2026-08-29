package rgo

// do_delayedBindingExpr translates batch_057.c (src/main/envir.c); R primitive(s): delayedBindingExpression.
func do_delayedBindingExpr(call, op, args, env Value) Value {
	return envPrimitive("do_delayedBindingExpr", args, env)
}
