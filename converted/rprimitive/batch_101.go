package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_serialize(call, op, args, env Value) Value {
	var (
		object Value
		icon   Value
		type_v Value
		ver    Value
		fun    Value
	)
	RT.Call("checkArity", op, args)
	if RT.Truth(RT.Binary("==", RT.Call("PRIMVAL", op), RT.Const("int", "2"))) {
		return RT.Call("checkNotPromise", RT.Call("R_unserialize", RT.Call("CAR", args), RT.Call("CADR", args)))
	}
	Assign(LocalRef(&object), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&icon), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&type_v), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&ver), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&fun), RT.Call("CAR", args))
	if RT.Truth(RT.Binary("==", RT.Call("PRIMVAL", op), RT.Const("int", "1"))) {
		return RT.Call("R_serializeb", object, icon, type_v, ver, fun)
	} else {
		return RT.Call("R_serialize", object, icon, type_v, ver, fun)
	}
}
