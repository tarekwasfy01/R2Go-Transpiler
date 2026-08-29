package rgo

// do_dotDelayedExpr translates batch_065.c (src/main/envir.c); R primitive(s): dotDelayedExpression.
func do_dotDelayedExpr(call, op, args, env Value) Value {
	return envPrimitive("do_dotDelayedExpr", args, env)
}
