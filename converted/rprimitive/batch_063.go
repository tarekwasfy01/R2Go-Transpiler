package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_isNSEnv(call, op, args, rho Value) Value {
	RT.Call("checkArity", op, args)
	return func() Value {
		if RT.Truth(RT.Call("R_IsNamespaceEnv", RT.Call("CAR", args))) {
			return RT.Call("mkTrue")
		}
		return RT.Call("mkFalse")
	}()
}
