package rgo

// do_dotForcedExpr translates batch_066.c (src/main/envir.c); R primitive(s): dotForcedExpression.
func do_dotForcedExpr(call, op, args, env Value) Value {
	return envPrimitive("do_dotForcedExpr", args, env)
}
