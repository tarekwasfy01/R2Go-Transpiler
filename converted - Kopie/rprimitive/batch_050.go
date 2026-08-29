package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_getNSRegistry(call, op, args, rho Value) Value {
	RT.Call("checkArity", op, args)
	return RT.Symbol("R_NamespaceRegistry")
}
