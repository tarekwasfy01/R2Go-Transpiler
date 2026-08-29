package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_dotCode(call, op, args, env Value) Value {
	var (
		cargs      Value
		cargs0     Value
		naok       Value
		na         Value
		nargs      Value
		Fort       Value
		copy       Value
		fun        Value
		ans        Value
		pa         Value
		s          Value
		symbol     Value
		checkTypes Value
		vmax       Value
		symName    Value
		havenames  Value
		names      Value
		nprotect   Value
		targetType Value
		n          Value
		t          Value
		ptr        Value
		ss         Value
		iptr       Value
		i          Value
		rptr       Value
		sptr       Value
		zptr       Value
		fptr       Value
		cptr       Value
		cptr0      Value
		nn         Value
		lptr       Value
		p          Value
		arg        Value
		type_v     Value
		tmp        Value
		buf        Value
		z          Value
		j          Value
		k          Value
	)
	Assign(LocalRef(&cargs0), RT.Symbol("NULL"))
	Assign(LocalRef(&copy), RT.Symbol("R_CBoundsCheck"))
	Assign(LocalRef(&fun), RT.Symbol("NULL"))
	Assign(LocalRef(&symbol), RT.List(RT.Symbol("R_C_SYM"), RT.List(RT.Symbol("NULL")), RT.Symbol("NULL")))
	Assign(LocalRef(&checkTypes), RT.Symbol("NULL"))
	Assign(LocalRef(&symName), RT.NewArray(RT.Symbol("MaxSymbolBytes")))
	if RT.Truth(RT.Binary("<", RT.Call("length", args), RT.Const("int", "1"))) {
		RT.Call("errorcall", call, RT.Call("_", RT.Const("string", "\"'.NAME' is missing\"")))
	}
	RT.Call("check1arg2", args, call, RT.Const("string", "\".NAME\""))
	if RT.Truth(func() Value {
		if RT.Truth(func() Value {
			if RT.Truth(RT.Binary("==", RT.Symbol("NaokSymbol"), RT.Symbol("NULL"))) {
				return true
			}
			return RT.Truth(RT.Binary("==", RT.Symbol("DupSymbol"), RT.Symbol("NULL")))
		}()) {
			return true
		}
		return RT.Truth(RT.Binary("==", RT.Symbol("PkgSymbol"), RT.Symbol("NULL")))
	}()) {
		RT.AssignSymbol("NaokSymbol", RT.Call("install", RT.Const("string", "\"NAOK\"")))
		RT.AssignSymbol("DupSymbol", RT.Call("install", RT.Const("string", "\"DUP\"")))
		RT.AssignSymbol("PkgSymbol", RT.Call("install", RT.Const("string", "\"PACKAGE\"")))
	}
	if RT.Truth(RT.Binary("==", RT.Symbol("EncSymbol"), RT.Symbol("NULL"))) {
		RT.AssignSymbol("EncSymbol", RT.Call("install", RT.Const("string", "\"ENCODING\"")))
	}
	if RT.Truth(RT.Binary("==", RT.Symbol("CSingSymbol"), RT.Symbol("NULL"))) {
		RT.AssignSymbol("CSingSymbol", RT.Call("install", RT.Const("string", "\"Csingle\"")))
	}
	Assign(LocalRef(&vmax), RT.Call("vmaxget"))
	Assign(LocalRef(&Fort), RT.Call("PRIMVAL", op))
	if RT.Truth(Fort) {
		RT.AssignField(symbol, "type", RT.Symbol("R_FORTRAN_SYM"))
	}
	Assign(LocalRef(&args), RT.Call("enctrim", args))
	Assign(LocalRef(&args), RT.Call("resolveNativeRoutine", args, LocalRef(&fun), LocalRef(&symbol), symName, LocalRef(&nargs), LocalRef(&naok), call, env))
	if RT.Truth(func() Value {
		if !RT.Truth(RT.Field(RT.Field(symbol, "symbol"), "c")) {
			return false
		}
		return RT.Truth(RT.Binary(">", RT.Field(RT.Field(RT.Field(symbol, "symbol"), "c"), "numArgs"), RT.Unary("-", RT.Const("int", "1"))))
	}()) {
		if RT.Truth(RT.Binary("!=", RT.Field(RT.Field(RT.Field(symbol, "symbol"), "c"), "numArgs"), nargs)) {
			RT.Call("errorcall", call, RT.Call("_", RT.Const("string", "\"Incorrect number of arguments (%d), expecting %d for '%s'\"")), nargs, RT.Field(RT.Field(RT.Field(symbol, "symbol"), "c"), "numArgs"), symName)
		}
		Assign(LocalRef(&checkTypes), RT.Field(RT.Field(RT.Field(symbol, "symbol"), "c"), "types"))
	}
	Assign(LocalRef(&nargs), RT.Const("int", "0"))
	Assign(LocalRef(&havenames), RT.Symbol("false"))
	for Assign(LocalRef(&pa), args); RT.Truth(RT.Binary("!=", pa, RT.Symbol("R_NilValue"))); Assign(LocalRef(&pa), RT.Call("CDR", pa)) {
		if RT.Truth(RT.Binary("!=", RT.Call("TAG", pa), RT.Symbol("R_NilValue"))) {
			Assign(LocalRef(&havenames), RT.Symbol("true"))
		}
		RT.Inc(LocalRef(&nargs), 1, true)
	}
	RT.Call("PROTECT", Assign(LocalRef(&ans), RT.Call("allocVector", RT.Symbol("VECSXP"), nargs)))
	if RT.Truth(havenames) {
		RT.Call("PROTECT", Assign(LocalRef(&names), RT.Call("allocVector", RT.Symbol("STRSXP"), nargs)))
		for RT.Sequence(Assign(LocalRef(&na), RT.Const("int", "0")), Assign(LocalRef(&pa), args)); RT.Truth(RT.Binary("!=", pa, RT.Symbol("R_NilValue"))); RT.Sequence(Assign(LocalRef(&pa), RT.Call("CDR", pa)), RT.Inc(LocalRef(&na), 1, true)) {
			if RT.Truth(RT.Binary("==", RT.Call("TAG", pa), RT.Symbol("R_NilValue"))) {
				RT.Call("SET_STRING_ELT", names, na, RT.Symbol("R_BlankString"))
			} else {
				RT.Call("SET_STRING_ELT", names, na, RT.Call("PRINTNAME", RT.Call("TAG", pa)))
			}
		}
		RT.Call("setAttrib", ans, RT.Symbol("R_NamesSymbol"), names)
		RT.Call("UNPROTECT", RT.Const("int", "1"))
	}
	Assign(LocalRef(&cargs), RT.Cast("void **", RT.Call("R_alloc", nargs, RT.SizeOfType("void *"))))
	if RT.Truth(copy) {
		Assign(LocalRef(&cargs0), RT.Cast("void **", RT.Call("R_alloc", nargs, RT.SizeOfType("void *"))))
	}
	for RT.Sequence(Assign(LocalRef(&na), RT.Const("int", "0")), Assign(LocalRef(&pa), args)); RT.Truth(RT.Binary("!=", pa, RT.Symbol("R_NilValue"))); RT.Sequence(Assign(LocalRef(&pa), RT.Call("CDR", pa)), RT.Inc(LocalRef(&na), 1, true)) {
		if RT.Truth(func() Value {
			if !RT.Truth(checkTypes) {
				return false
			}
			return RT.Truth(RT.Unary("!", RT.Call("comparePrimitiveTypes", RT.Index(checkTypes, na), RT.Call("CAR", pa))))
		}()) {
			RT.Call("errorcall", call, RT.Call("_", RT.Const("string", "\"wrong type for argument %d in call to %s\"")), RT.Binary("+", na, RT.Const("int", "1")), symName)
		}
		Assign(LocalRef(&nprotect), RT.Const("int", "0"))
		Assign(LocalRef(&targetType), func() Value {
			if RT.Truth(checkTypes) {
				return RT.Index(checkTypes, na)
			}
			return RT.Const("int", "0")
		}())
		Assign(LocalRef(&s), RT.Call("CAR", pa))
		RT.Call("SET_VECTOR_ELT", ans, na, s)
		if RT.Truth(func() Value {
			if !RT.Truth(RT.Binary("==", RT.Call("checkNativeType", targetType, RT.Call("TYPEOF", s)), RT.Symbol("false"))) {
				return false
			}
			return RT.Truth(RT.Binary("!=", targetType, RT.Symbol("SINGLESXP")))
		}()) {
			RT.Call("PROTECT", Assign(LocalRef(&s), RT.Call("coerceVector", s, targetType)))
			RT.Inc(LocalRef(&nprotect), 1, true)
		}
		if RT.Truth(func() Value {
			if !RT.Truth(RT.Call("isVector", s)) {
				return false
			}
			return RT.Truth(RT.Call("IS_LONG_VEC", s))
		}()) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"long vectors (argument %d) are not supported in %s\"")), RT.Binary("+", na, RT.Const("int", "1")), func() Value {
				if RT.Truth(Fort) {
					return RT.Const("string", "\".Fortran\"")
				}
				return RT.Const("string", "\".C\"")
			}())
		}
		Assign(LocalRef(&t), RT.Call("TYPEOF", s))
		switch RT.Key(t) {
		case RT.Key(RT.Symbol("RAWSXP")):
			if RT.Truth(copy) {
				Assign(LocalRef(&n), RT.Call("XLENGTH", s))
				Assign(LocalRef(&ptr), RT.Call("R_alloc", RT.Binary("+", RT.Binary("*", n, RT.SizeOfType("Rbyte")), RT.Binary("*", RT.Const("int", "2"), RT.Symbol("NG"))), RT.Const("int", "1")))
				RT.Call("memset", ptr, RT.Symbol("FILL"), RT.Binary("+", RT.Binary("*", n, RT.SizeOfType("Rbyte")), RT.Binary("*", RT.Const("int", "2"), RT.Symbol("NG"))))
				Assign(LocalRef(&ptr), RT.Binary("+", ptr, RT.Symbol("NG")))
				if RT.Truth(n) {
					RT.Call("memcpy", ptr, RT.Call("RAW", s), n)
				}
				RT.AssignIndex(cargs, na, RT.Cast("void *", ptr))
			} else {
				if RT.Truth(RT.Call("MAYBE_REFERENCED", s)) {
					Assign(LocalRef(&n), RT.Call("XLENGTH", s))
					Assign(LocalRef(&ss), RT.Call("allocVector", t, n))
					if RT.Truth(n) {
						RT.Call("memcpy", RT.Call("RAW", ss), RT.Call("RAW", s), RT.Binary("*", n, RT.SizeOfType("Rbyte")))
					}
					RT.Call("SET_VECTOR_ELT", ans, na, ss)
					RT.AssignIndex(cargs, na, RT.Cast("void *", RT.Call("RAW", ss)))
				} else {
					RT.AssignIndex(cargs, na, RT.Cast("void *", RT.Call("RAW", s)))
				}
			}
			break
		case RT.Key(RT.Symbol("LGLSXP")), RT.Key(RT.Symbol("INTSXP")):
			Assign(LocalRef(&n), RT.Call("XLENGTH", s))
			Assign(LocalRef(&iptr), RT.Call("INTEGER", s))
			if RT.Truth(RT.Unary("!", naok)) {
				for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, n)); RT.Inc(LocalRef(&i), 1, true) {
					if RT.Truth(RT.Binary("==", RT.Index(iptr, i), RT.Symbol("NA_INTEGER"))) {
						RT.Call("error", RT.Call("_", RT.Const("string", "\"NAs in foreign function call (arg %d)\"")), RT.Binary("+", na, RT.Const("int", "1")))
					}
				}
			}
			if RT.Truth(copy) {
				Assign(LocalRef(&ptr), RT.Call("R_alloc", RT.Binary("+", RT.Binary("*", n, RT.SizeOfType("int")), RT.Binary("*", RT.Const("int", "2"), RT.Symbol("NG"))), RT.Const("int", "1")))
				RT.Call("memset", ptr, RT.Symbol("FILL"), RT.Binary("+", RT.Binary("*", n, RT.SizeOfType("int")), RT.Binary("*", RT.Const("int", "2"), RT.Symbol("NG"))))
				Assign(LocalRef(&ptr), RT.Binary("+", ptr, RT.Symbol("NG")))
				if RT.Truth(n) {
					RT.Call("memcpy", ptr, RT.Call("INTEGER", s), RT.Binary("*", n, RT.SizeOfType("int")))
				}
				RT.AssignIndex(cargs, na, RT.Cast("void *", ptr))
			} else {
				if RT.Truth(RT.Call("MAYBE_REFERENCED", s)) {
					Assign(LocalRef(&ss), RT.Call("allocVector", t, n))
					if RT.Truth(n) {
						RT.Call("memcpy", RT.Call("INTEGER", ss), RT.Call("INTEGER", s), RT.Binary("*", n, RT.SizeOfType("int")))
					}
					RT.Call("SET_VECTOR_ELT", ans, na, ss)
					RT.AssignIndex(cargs, na, RT.Cast("void *", RT.Call("INTEGER", ss)))
				} else {
					RT.AssignIndex(cargs, na, RT.Cast("void *", iptr))
				}
			}
			break
		case RT.Key(RT.Symbol("REALSXP")):
			Assign(LocalRef(&n), RT.Call("XLENGTH", s))
			Assign(LocalRef(&rptr), RT.Call("REAL", s))
			if RT.Truth(RT.Unary("!", naok)) {
				for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, n)); RT.Inc(LocalRef(&i), 1, true) {
					if RT.Truth(RT.Unary("!", RT.Call("R_FINITE", RT.Index(rptr, i)))) {
						RT.Call("error", RT.Call("_", RT.Const("string", "\"NA/NaN/Inf in foreign function call (arg %d)\"")), RT.Binary("+", na, RT.Const("int", "1")))
					}
				}
			}
			if RT.Truth(RT.Binary("==", RT.Call("asLogical", RT.Call("getAttrib", s, RT.Symbol("CSingSymbol"))), RT.Const("int", "1"))) {
				Assign(LocalRef(&sptr), RT.Cast("float *", RT.Call("R_alloc", n, RT.SizeOfType("float"))))
				for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, n)); RT.Inc(LocalRef(&i), 1, true) {
					RT.AssignIndex(sptr, i, RT.Cast("float", RT.Index(RT.Call("REAL", s), i)))
				}
				RT.AssignIndex(cargs, na, RT.Cast("void *", sptr))
			} else {
				if RT.Truth(copy) {
					Assign(LocalRef(&ptr), RT.Call("R_alloc", RT.Binary("+", RT.Binary("*", n, RT.SizeOfType("double")), RT.Binary("*", RT.Const("int", "2"), RT.Symbol("NG"))), RT.Const("int", "1")))
					RT.Call("memset", ptr, RT.Symbol("FILL"), RT.Binary("+", RT.Binary("*", n, RT.SizeOfType("double")), RT.Binary("*", RT.Const("int", "2"), RT.Symbol("NG"))))
					Assign(LocalRef(&ptr), RT.Binary("+", ptr, RT.Symbol("NG")))
					if RT.Truth(n) {
						RT.Call("memcpy", ptr, RT.Call("REAL", s), RT.Binary("*", n, RT.SizeOfType("double")))
					}
					RT.AssignIndex(cargs, na, RT.Cast("void *", ptr))
				} else {
					if RT.Truth(RT.Call("MAYBE_REFERENCED", s)) {
						Assign(LocalRef(&ss), RT.Call("allocVector", t, n))
						if RT.Truth(n) {
							RT.Call("memcpy", RT.Call("REAL", ss), RT.Call("REAL", s), RT.Binary("*", n, RT.SizeOfType("double")))
						}
						RT.Call("SET_VECTOR_ELT", ans, na, ss)
						RT.AssignIndex(cargs, na, RT.Cast("void *", RT.Call("REAL", ss)))
					} else {
						RT.AssignIndex(cargs, na, RT.Cast("void *", rptr))
					}
				}
			}
			break
		case RT.Key(RT.Symbol("CPLXSXP")):
			Assign(LocalRef(&n), RT.Call("XLENGTH", s))
			Assign(LocalRef(&zptr), RT.Call("COMPLEX", s))
			if RT.Truth(RT.Unary("!", naok)) {
				for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, n)); RT.Inc(LocalRef(&i), 1, true) {
					if RT.Truth(func() Value {
						if RT.Truth(RT.Unary("!", RT.Call("R_FINITE", RT.Field(RT.Index(zptr, i), "r")))) {
							return true
						}
						return RT.Truth(RT.Unary("!", RT.Call("R_FINITE", RT.Field(RT.Index(zptr, i), "i"))))
					}()) {
						RT.Call("error", RT.Call("_", RT.Const("string", "\"complex NA/NaN/Inf in foreign function call (arg %d)\"")), RT.Binary("+", na, RT.Const("int", "1")))
					}
				}
			}
			if RT.Truth(copy) {
				Assign(LocalRef(&ptr), RT.Call("R_alloc", RT.Binary("+", RT.Binary("*", n, RT.SizeOfType("Rcomplex")), RT.Binary("*", RT.Const("int", "2"), RT.Symbol("NG"))), RT.Const("int", "1")))
				RT.Call("memset", ptr, RT.Symbol("FILL"), RT.Binary("+", RT.Binary("*", n, RT.SizeOfType("Rcomplex")), RT.Binary("*", RT.Const("int", "2"), RT.Symbol("NG"))))
				Assign(LocalRef(&ptr), RT.Binary("+", ptr, RT.Symbol("NG")))
				if RT.Truth(n) {
					RT.Call("memcpy", ptr, RT.Call("COMPLEX", s), RT.Binary("*", n, RT.SizeOfType("Rcomplex")))
				}
				RT.AssignIndex(cargs, na, RT.Cast("void *", ptr))
			} else {
				if RT.Truth(RT.Call("MAYBE_REFERENCED", s)) {
					Assign(LocalRef(&ss), RT.Call("allocVector", t, n))
					if RT.Truth(n) {
						RT.Call("memcpy", RT.Call("COMPLEX", ss), RT.Call("COMPLEX", s), RT.Binary("*", n, RT.SizeOfType("Rcomplex")))
					}
					RT.Call("SET_VECTOR_ELT", ans, na, ss)
					RT.AssignIndex(cargs, na, RT.Cast("void *", RT.Call("COMPLEX", ss)))
				} else {
					RT.AssignIndex(cargs, na, RT.Cast("void *", zptr))
				}
			}
			break
		case RT.Key(RT.Symbol("STRSXP")):
			Assign(LocalRef(&n), RT.Call("XLENGTH", s))
			if RT.Truth(Fort) {
				Assign(LocalRef(&ss), RT.Call("translateChar", RT.Call("STRING_ELT", s, RT.Const("int", "0"))))
				if RT.Truth(RT.Binary(">", n, RT.Const("int", "1"))) {
					RT.Call("warning", RT.Const("string", "\"only the first string in a char vector used in .Fortran\""))
				} else {
					RT.Call("warning", RT.Const("string", "\"passing a char vector to .Fortran is not portable\""))
				}
				Assign(LocalRef(&fptr), RT.Cast("char *", RT.Call("R_alloc", RT.Binary("+", RT.Call("max", RT.Const("int", "255"), RT.Call("strlen", ss)), RT.Const("int", "1")), RT.SizeOfType("char"))))
				RT.Call("strcpy", fptr, ss)
				RT.AssignIndex(cargs, na, RT.Cast("void *", fptr))
			} else {
				if RT.Truth(copy) {
					Assign(LocalRef(&cptr), RT.Cast("char **", RT.Call("R_alloc", n, RT.SizeOfType("char *"))))
					Assign(LocalRef(&cptr0), RT.Cast("char **", RT.Call("R_alloc", n, RT.SizeOfType("char *"))))
					for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, n)); RT.Inc(LocalRef(&i), 1, true) {
						Assign(LocalRef(&ss), RT.Call("translateChar", RT.Call("STRING_ELT", s, i)))
						Assign(LocalRef(&nn), RT.Binary("+", RT.Binary("+", RT.Call("strlen", ss), RT.Const("int", "1")), RT.Binary("*", RT.Const("int", "2"), RT.Symbol("NG"))))
						Assign(LocalRef(&ptr), RT.Cast("char *", RT.Call("R_alloc", nn, RT.SizeOfType("char"))))
						RT.Call("memset", ptr, RT.Symbol("FILL"), nn)
						RT.AssignIndex(cptr, i, RT.AssignIndex(cptr0, i, RT.Binary("+", ptr, RT.Symbol("NG"))))
						RT.Call("strcpy", RT.Index(cptr, i), ss)
					}
					RT.AssignIndex(cargs, na, RT.Cast("void *", cptr))
					RT.AssignIndex(cargs0, na, RT.Cast("void *", cptr0))
				} else {
					Assign(LocalRef(&cptr), RT.Cast("char **", RT.Call("R_alloc", n, RT.SizeOfType("char *"))))
					for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, n)); RT.Inc(LocalRef(&i), 1, true) {
						Assign(LocalRef(&ss), RT.Call("translateChar", RT.Call("STRING_ELT", s, i)))
						Assign(LocalRef(&nn), RT.Binary("+", RT.Call("strlen", ss), RT.Const("int", "1")))
						if RT.Truth(RT.Binary(">", nn, RT.Const("int", "1"))) {
							RT.AssignIndex(cptr, i, RT.Cast("char *", RT.Call("R_alloc", nn, RT.SizeOfType("char"))))
							RT.Call("strcpy", RT.Index(cptr, i), ss)
						} else {
							Assign(LocalRef(&nn), RT.Const("int", "128"))
							RT.AssignIndex(cptr, i, RT.Cast("char *", RT.Call("R_alloc", nn, RT.SizeOfType("char"))))
							RT.Call("memset", RT.Index(cptr, i), RT.Const("int", "0"), nn)
						}
					}
					RT.AssignIndex(cargs, na, RT.Cast("void *", cptr))
				}
			}
			break
		case RT.Key(RT.Symbol("VECSXP")):
			if RT.Truth(Fort) {
				RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid mode (%s) to pass to Fortran (arg %d)\"")), RT.Call("R_typeToChar", s), RT.Binary("+", na, RT.Const("int", "1")))
			}
			Assign(LocalRef(&n), RT.Call("XLENGTH", s))
			Assign(LocalRef(&lptr), RT.Cast("SEXP *", RT.Call("R_alloc", n, RT.SizeOfType("SEXP"))))
			for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, n)); RT.Inc(LocalRef(&i), 1, true) {
				RT.AssignIndex(lptr, i, RT.Call("VECTOR_ELT", s, i))
			}
			RT.AssignIndex(cargs, na, RT.Cast("void *", lptr))
			break
		case RT.Key(RT.Symbol("CLOSXP")), RT.Key(RT.Symbol("BUILTINSXP")), RT.Key(RT.Symbol("SPECIALSXP")), RT.Key(RT.Symbol("ENVSXP")):
			if RT.Truth(Fort) {
				RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid mode (%s) to pass to Fortran (arg %d)\"")), RT.Call("R_typeToChar", s), RT.Binary("+", na, RT.Const("int", "1")))
			}
			RT.AssignIndex(cargs, na, RT.Cast("void *", s))
			break
		case RT.Key(RT.Symbol("NILSXP")):
			RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid mode (%s) to pass to C or Fortran (arg %d)\"")), RT.Call("R_typeToChar", s), RT.Binary("+", na, RT.Const("int", "1")))
			RT.AssignIndex(cargs, na, RT.Cast("void *", s))
			break
		default:
			if RT.Truth(Fort) {
				RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid mode (%s) to pass to Fortran (arg %d)\"")), RT.Call("R_typeToChar", s), RT.Binary("+", na, RT.Const("int", "1")))
			}
			RT.Call("warning", RT.Const("string", "\"passing an object of type '%s' to .C (arg %d) is deprecated\""), RT.Call("R_typeToChar", s), RT.Binary("+", na, RT.Const("int", "1")))
			if RT.Truth(RT.Binary("==", t, RT.Symbol("LISTSXP"))) {
				RT.Call("warning", RT.Call("_", RT.Const("string", "\"pairlists are passed as SEXP as from R 2.15.0\"")))
			}
			RT.AssignIndex(cargs, na, RT.Cast("void *", s))
			break
		}
		if RT.Truth(nprotect) {
			RT.Call("UNPROTECT", nprotect)
		}
	}
	switch RT.Key(nargs) {
	case RT.Key(RT.Const("int", "0")):
		RT.CallIndirect(RT.Cast("FUNV0", fun))
		break
	case RT.Key(RT.Const("int", "1")):
		RT.CallIndirect(RT.Cast("FUNV1", fun), RT.Index(cargs, RT.Const("int", "0")))
		break
	case RT.Key(RT.Const("int", "2")):
		RT.CallIndirect(RT.Cast("FUNV2", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")))
		break
	case RT.Key(RT.Const("int", "3")):
		RT.CallIndirect(RT.Cast("FUNV3", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")))
		break
	case RT.Key(RT.Const("int", "4")):
		RT.CallIndirect(RT.Cast("FUNV4", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")))
		break
	case RT.Key(RT.Const("int", "5")):
		RT.CallIndirect(RT.Cast("FUNV5", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")))
		break
	case RT.Key(RT.Const("int", "6")):
		RT.CallIndirect(RT.Cast("FUNV6", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")))
		break
	case RT.Key(RT.Const("int", "7")):
		RT.CallIndirect(RT.Cast("FUNV7", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")))
		break
	case RT.Key(RT.Const("int", "8")):
		RT.CallIndirect(RT.Cast("FUNV8", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")))
		break
	case RT.Key(RT.Const("int", "9")):
		RT.CallIndirect(RT.Cast("FUNV9", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")))
		break
	case RT.Key(RT.Const("int", "10")):
		RT.CallIndirect(RT.Cast("FUNV10", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")))
		break
	case RT.Key(RT.Const("int", "11")):
		RT.CallIndirect(RT.Cast("FUNV11", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")))
		break
	case RT.Key(RT.Const("int", "12")):
		RT.CallIndirect(RT.Cast("FUNV12", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")))
		break
	case RT.Key(RT.Const("int", "13")):
		RT.CallIndirect(RT.Cast("FUNV13", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")))
		break
	case RT.Key(RT.Const("int", "14")):
		RT.CallIndirect(RT.Cast("FUNV14", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")))
		break
	case RT.Key(RT.Const("int", "15")):
		RT.CallIndirect(RT.Cast("FUNV15", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")))
		break
	case RT.Key(RT.Const("int", "16")):
		RT.CallIndirect(RT.Cast("FUNV16", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")))
		break
	case RT.Key(RT.Const("int", "17")):
		RT.CallIndirect(RT.Cast("FUNV17", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")))
		break
	case RT.Key(RT.Const("int", "18")):
		RT.CallIndirect(RT.Cast("FUNV18", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")))
		break
	case RT.Key(RT.Const("int", "19")):
		RT.CallIndirect(RT.Cast("FUNV19", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")))
		break
	case RT.Key(RT.Const("int", "20")):
		RT.CallIndirect(RT.Cast("FUNV20", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")))
		break
	case RT.Key(RT.Const("int", "21")):
		RT.CallIndirect(RT.Cast("FUNV21", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")))
		break
	case RT.Key(RT.Const("int", "22")):
		RT.CallIndirect(RT.Cast("FUNV22", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")))
		break
	case RT.Key(RT.Const("int", "23")):
		RT.CallIndirect(RT.Cast("FUNV23", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")))
		break
	case RT.Key(RT.Const("int", "24")):
		RT.CallIndirect(RT.Cast("FUNV24", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")))
		break
	case RT.Key(RT.Const("int", "25")):
		RT.CallIndirect(RT.Cast("FUNV25", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")))
		break
	case RT.Key(RT.Const("int", "26")):
		RT.CallIndirect(RT.Cast("FUNV26", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")))
		break
	case RT.Key(RT.Const("int", "27")):
		RT.CallIndirect(RT.Cast("FUNV27", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")))
		break
	case RT.Key(RT.Const("int", "28")):
		RT.CallIndirect(RT.Cast("FUNV28", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")))
		break
	case RT.Key(RT.Const("int", "29")):
		RT.CallIndirect(RT.Cast("FUNV29", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")))
		break
	case RT.Key(RT.Const("int", "30")):
		RT.CallIndirect(RT.Cast("FUNV30", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")))
		break
	case RT.Key(RT.Const("int", "31")):
		RT.CallIndirect(RT.Cast("FUNV31", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")))
		break
	case RT.Key(RT.Const("int", "32")):
		RT.CallIndirect(RT.Cast("FUNV32", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")))
		break
	case RT.Key(RT.Const("int", "33")):
		RT.CallIndirect(RT.Cast("FUNV33", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")))
		break
	case RT.Key(RT.Const("int", "34")):
		RT.CallIndirect(RT.Cast("FUNV34", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")))
		break
	case RT.Key(RT.Const("int", "35")):
		RT.CallIndirect(RT.Cast("FUNV35", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")), RT.Index(cargs, RT.Const("int", "34")))
		break
	case RT.Key(RT.Const("int", "36")):
		RT.CallIndirect(RT.Cast("FUNV36", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")), RT.Index(cargs, RT.Const("int", "34")), RT.Index(cargs, RT.Const("int", "35")))
		break
	case RT.Key(RT.Const("int", "37")):
		RT.CallIndirect(RT.Cast("FUNV37", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")), RT.Index(cargs, RT.Const("int", "34")), RT.Index(cargs, RT.Const("int", "35")), RT.Index(cargs, RT.Const("int", "36")))
		break
	case RT.Key(RT.Const("int", "38")):
		RT.CallIndirect(RT.Cast("FUNV38", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")), RT.Index(cargs, RT.Const("int", "34")), RT.Index(cargs, RT.Const("int", "35")), RT.Index(cargs, RT.Const("int", "36")), RT.Index(cargs, RT.Const("int", "37")))
		break
	case RT.Key(RT.Const("int", "39")):
		RT.CallIndirect(RT.Cast("FUNV39", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")), RT.Index(cargs, RT.Const("int", "34")), RT.Index(cargs, RT.Const("int", "35")), RT.Index(cargs, RT.Const("int", "36")), RT.Index(cargs, RT.Const("int", "37")), RT.Index(cargs, RT.Const("int", "38")))
		break
	case RT.Key(RT.Const("int", "40")):
		RT.CallIndirect(RT.Cast("FUNV40", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")), RT.Index(cargs, RT.Const("int", "34")), RT.Index(cargs, RT.Const("int", "35")), RT.Index(cargs, RT.Const("int", "36")), RT.Index(cargs, RT.Const("int", "37")), RT.Index(cargs, RT.Const("int", "38")), RT.Index(cargs, RT.Const("int", "39")))
		break
	case RT.Key(RT.Const("int", "41")):
		RT.CallIndirect(RT.Cast("FUNV41", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")), RT.Index(cargs, RT.Const("int", "34")), RT.Index(cargs, RT.Const("int", "35")), RT.Index(cargs, RT.Const("int", "36")), RT.Index(cargs, RT.Const("int", "37")), RT.Index(cargs, RT.Const("int", "38")), RT.Index(cargs, RT.Const("int", "39")), RT.Index(cargs, RT.Const("int", "40")))
		break
	case RT.Key(RT.Const("int", "42")):
		RT.CallIndirect(RT.Cast("FUNV42", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")), RT.Index(cargs, RT.Const("int", "34")), RT.Index(cargs, RT.Const("int", "35")), RT.Index(cargs, RT.Const("int", "36")), RT.Index(cargs, RT.Const("int", "37")), RT.Index(cargs, RT.Const("int", "38")), RT.Index(cargs, RT.Const("int", "39")), RT.Index(cargs, RT.Const("int", "40")), RT.Index(cargs, RT.Const("int", "41")))
		break
	case RT.Key(RT.Const("int", "43")):
		RT.CallIndirect(RT.Cast("FUNV43", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")), RT.Index(cargs, RT.Const("int", "34")), RT.Index(cargs, RT.Const("int", "35")), RT.Index(cargs, RT.Const("int", "36")), RT.Index(cargs, RT.Const("int", "37")), RT.Index(cargs, RT.Const("int", "38")), RT.Index(cargs, RT.Const("int", "39")), RT.Index(cargs, RT.Const("int", "40")), RT.Index(cargs, RT.Const("int", "41")), RT.Index(cargs, RT.Const("int", "42")))
		break
	case RT.Key(RT.Const("int", "44")):
		RT.CallIndirect(RT.Cast("FUNV44", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")), RT.Index(cargs, RT.Const("int", "34")), RT.Index(cargs, RT.Const("int", "35")), RT.Index(cargs, RT.Const("int", "36")), RT.Index(cargs, RT.Const("int", "37")), RT.Index(cargs, RT.Const("int", "38")), RT.Index(cargs, RT.Const("int", "39")), RT.Index(cargs, RT.Const("int", "40")), RT.Index(cargs, RT.Const("int", "41")), RT.Index(cargs, RT.Const("int", "42")), RT.Index(cargs, RT.Const("int", "43")))
		break
	case RT.Key(RT.Const("int", "45")):
		RT.CallIndirect(RT.Cast("FUNV45", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")), RT.Index(cargs, RT.Const("int", "34")), RT.Index(cargs, RT.Const("int", "35")), RT.Index(cargs, RT.Const("int", "36")), RT.Index(cargs, RT.Const("int", "37")), RT.Index(cargs, RT.Const("int", "38")), RT.Index(cargs, RT.Const("int", "39")), RT.Index(cargs, RT.Const("int", "40")), RT.Index(cargs, RT.Const("int", "41")), RT.Index(cargs, RT.Const("int", "42")), RT.Index(cargs, RT.Const("int", "43")), RT.Index(cargs, RT.Const("int", "44")))
		break
	case RT.Key(RT.Const("int", "46")):
		RT.CallIndirect(RT.Cast("FUNV46", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")), RT.Index(cargs, RT.Const("int", "34")), RT.Index(cargs, RT.Const("int", "35")), RT.Index(cargs, RT.Const("int", "36")), RT.Index(cargs, RT.Const("int", "37")), RT.Index(cargs, RT.Const("int", "38")), RT.Index(cargs, RT.Const("int", "39")), RT.Index(cargs, RT.Const("int", "40")), RT.Index(cargs, RT.Const("int", "41")), RT.Index(cargs, RT.Const("int", "42")), RT.Index(cargs, RT.Const("int", "43")), RT.Index(cargs, RT.Const("int", "44")), RT.Index(cargs, RT.Const("int", "45")))
		break
	case RT.Key(RT.Const("int", "47")):
		RT.CallIndirect(RT.Cast("FUNV47", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")), RT.Index(cargs, RT.Const("int", "34")), RT.Index(cargs, RT.Const("int", "35")), RT.Index(cargs, RT.Const("int", "36")), RT.Index(cargs, RT.Const("int", "37")), RT.Index(cargs, RT.Const("int", "38")), RT.Index(cargs, RT.Const("int", "39")), RT.Index(cargs, RT.Const("int", "40")), RT.Index(cargs, RT.Const("int", "41")), RT.Index(cargs, RT.Const("int", "42")), RT.Index(cargs, RT.Const("int", "43")), RT.Index(cargs, RT.Const("int", "44")), RT.Index(cargs, RT.Const("int", "45")), RT.Index(cargs, RT.Const("int", "46")))
		break
	case RT.Key(RT.Const("int", "48")):
		RT.CallIndirect(RT.Cast("FUNV48", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")), RT.Index(cargs, RT.Const("int", "34")), RT.Index(cargs, RT.Const("int", "35")), RT.Index(cargs, RT.Const("int", "36")), RT.Index(cargs, RT.Const("int", "37")), RT.Index(cargs, RT.Const("int", "38")), RT.Index(cargs, RT.Const("int", "39")), RT.Index(cargs, RT.Const("int", "40")), RT.Index(cargs, RT.Const("int", "41")), RT.Index(cargs, RT.Const("int", "42")), RT.Index(cargs, RT.Const("int", "43")), RT.Index(cargs, RT.Const("int", "44")), RT.Index(cargs, RT.Const("int", "45")), RT.Index(cargs, RT.Const("int", "46")), RT.Index(cargs, RT.Const("int", "47")))
		break
	case RT.Key(RT.Const("int", "49")):
		RT.CallIndirect(RT.Cast("FUNV49", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")), RT.Index(cargs, RT.Const("int", "34")), RT.Index(cargs, RT.Const("int", "35")), RT.Index(cargs, RT.Const("int", "36")), RT.Index(cargs, RT.Const("int", "37")), RT.Index(cargs, RT.Const("int", "38")), RT.Index(cargs, RT.Const("int", "39")), RT.Index(cargs, RT.Const("int", "40")), RT.Index(cargs, RT.Const("int", "41")), RT.Index(cargs, RT.Const("int", "42")), RT.Index(cargs, RT.Const("int", "43")), RT.Index(cargs, RT.Const("int", "44")), RT.Index(cargs, RT.Const("int", "45")), RT.Index(cargs, RT.Const("int", "46")), RT.Index(cargs, RT.Const("int", "47")), RT.Index(cargs, RT.Const("int", "48")))
		break
	case RT.Key(RT.Const("int", "50")):
		RT.CallIndirect(RT.Cast("FUNV50", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")), RT.Index(cargs, RT.Const("int", "34")), RT.Index(cargs, RT.Const("int", "35")), RT.Index(cargs, RT.Const("int", "36")), RT.Index(cargs, RT.Const("int", "37")), RT.Index(cargs, RT.Const("int", "38")), RT.Index(cargs, RT.Const("int", "39")), RT.Index(cargs, RT.Const("int", "40")), RT.Index(cargs, RT.Const("int", "41")), RT.Index(cargs, RT.Const("int", "42")), RT.Index(cargs, RT.Const("int", "43")), RT.Index(cargs, RT.Const("int", "44")), RT.Index(cargs, RT.Const("int", "45")), RT.Index(cargs, RT.Const("int", "46")), RT.Index(cargs, RT.Const("int", "47")), RT.Index(cargs, RT.Const("int", "48")), RT.Index(cargs, RT.Const("int", "49")))
		break
	case RT.Key(RT.Const("int", "51")):
		RT.CallIndirect(RT.Cast("FUNV51", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")), RT.Index(cargs, RT.Const("int", "34")), RT.Index(cargs, RT.Const("int", "35")), RT.Index(cargs, RT.Const("int", "36")), RT.Index(cargs, RT.Const("int", "37")), RT.Index(cargs, RT.Const("int", "38")), RT.Index(cargs, RT.Const("int", "39")), RT.Index(cargs, RT.Const("int", "40")), RT.Index(cargs, RT.Const("int", "41")), RT.Index(cargs, RT.Const("int", "42")), RT.Index(cargs, RT.Const("int", "43")), RT.Index(cargs, RT.Const("int", "44")), RT.Index(cargs, RT.Const("int", "45")), RT.Index(cargs, RT.Const("int", "46")), RT.Index(cargs, RT.Const("int", "47")), RT.Index(cargs, RT.Const("int", "48")), RT.Index(cargs, RT.Const("int", "49")), RT.Index(cargs, RT.Const("int", "50")))
		break
	case RT.Key(RT.Const("int", "52")):
		RT.CallIndirect(RT.Cast("FUNV52", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")), RT.Index(cargs, RT.Const("int", "34")), RT.Index(cargs, RT.Const("int", "35")), RT.Index(cargs, RT.Const("int", "36")), RT.Index(cargs, RT.Const("int", "37")), RT.Index(cargs, RT.Const("int", "38")), RT.Index(cargs, RT.Const("int", "39")), RT.Index(cargs, RT.Const("int", "40")), RT.Index(cargs, RT.Const("int", "41")), RT.Index(cargs, RT.Const("int", "42")), RT.Index(cargs, RT.Const("int", "43")), RT.Index(cargs, RT.Const("int", "44")), RT.Index(cargs, RT.Const("int", "45")), RT.Index(cargs, RT.Const("int", "46")), RT.Index(cargs, RT.Const("int", "47")), RT.Index(cargs, RT.Const("int", "48")), RT.Index(cargs, RT.Const("int", "49")), RT.Index(cargs, RT.Const("int", "50")), RT.Index(cargs, RT.Const("int", "51")))
		break
	case RT.Key(RT.Const("int", "53")):
		RT.CallIndirect(RT.Cast("FUNV53", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")), RT.Index(cargs, RT.Const("int", "34")), RT.Index(cargs, RT.Const("int", "35")), RT.Index(cargs, RT.Const("int", "36")), RT.Index(cargs, RT.Const("int", "37")), RT.Index(cargs, RT.Const("int", "38")), RT.Index(cargs, RT.Const("int", "39")), RT.Index(cargs, RT.Const("int", "40")), RT.Index(cargs, RT.Const("int", "41")), RT.Index(cargs, RT.Const("int", "42")), RT.Index(cargs, RT.Const("int", "43")), RT.Index(cargs, RT.Const("int", "44")), RT.Index(cargs, RT.Const("int", "45")), RT.Index(cargs, RT.Const("int", "46")), RT.Index(cargs, RT.Const("int", "47")), RT.Index(cargs, RT.Const("int", "48")), RT.Index(cargs, RT.Const("int", "49")), RT.Index(cargs, RT.Const("int", "50")), RT.Index(cargs, RT.Const("int", "51")), RT.Index(cargs, RT.Const("int", "52")))
		break
	case RT.Key(RT.Const("int", "54")):
		RT.CallIndirect(RT.Cast("FUNV54", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")), RT.Index(cargs, RT.Const("int", "34")), RT.Index(cargs, RT.Const("int", "35")), RT.Index(cargs, RT.Const("int", "36")), RT.Index(cargs, RT.Const("int", "37")), RT.Index(cargs, RT.Const("int", "38")), RT.Index(cargs, RT.Const("int", "39")), RT.Index(cargs, RT.Const("int", "40")), RT.Index(cargs, RT.Const("int", "41")), RT.Index(cargs, RT.Const("int", "42")), RT.Index(cargs, RT.Const("int", "43")), RT.Index(cargs, RT.Const("int", "44")), RT.Index(cargs, RT.Const("int", "45")), RT.Index(cargs, RT.Const("int", "46")), RT.Index(cargs, RT.Const("int", "47")), RT.Index(cargs, RT.Const("int", "48")), RT.Index(cargs, RT.Const("int", "49")), RT.Index(cargs, RT.Const("int", "50")), RT.Index(cargs, RT.Const("int", "51")), RT.Index(cargs, RT.Const("int", "52")), RT.Index(cargs, RT.Const("int", "53")))
		break
	case RT.Key(RT.Const("int", "55")):
		RT.CallIndirect(RT.Cast("FUNV55", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")), RT.Index(cargs, RT.Const("int", "34")), RT.Index(cargs, RT.Const("int", "35")), RT.Index(cargs, RT.Const("int", "36")), RT.Index(cargs, RT.Const("int", "37")), RT.Index(cargs, RT.Const("int", "38")), RT.Index(cargs, RT.Const("int", "39")), RT.Index(cargs, RT.Const("int", "40")), RT.Index(cargs, RT.Const("int", "41")), RT.Index(cargs, RT.Const("int", "42")), RT.Index(cargs, RT.Const("int", "43")), RT.Index(cargs, RT.Const("int", "44")), RT.Index(cargs, RT.Const("int", "45")), RT.Index(cargs, RT.Const("int", "46")), RT.Index(cargs, RT.Const("int", "47")), RT.Index(cargs, RT.Const("int", "48")), RT.Index(cargs, RT.Const("int", "49")), RT.Index(cargs, RT.Const("int", "50")), RT.Index(cargs, RT.Const("int", "51")), RT.Index(cargs, RT.Const("int", "52")), RT.Index(cargs, RT.Const("int", "53")), RT.Index(cargs, RT.Const("int", "54")))
		break
	case RT.Key(RT.Const("int", "56")):
		RT.CallIndirect(RT.Cast("FUNV56", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")), RT.Index(cargs, RT.Const("int", "34")), RT.Index(cargs, RT.Const("int", "35")), RT.Index(cargs, RT.Const("int", "36")), RT.Index(cargs, RT.Const("int", "37")), RT.Index(cargs, RT.Const("int", "38")), RT.Index(cargs, RT.Const("int", "39")), RT.Index(cargs, RT.Const("int", "40")), RT.Index(cargs, RT.Const("int", "41")), RT.Index(cargs, RT.Const("int", "42")), RT.Index(cargs, RT.Const("int", "43")), RT.Index(cargs, RT.Const("int", "44")), RT.Index(cargs, RT.Const("int", "45")), RT.Index(cargs, RT.Const("int", "46")), RT.Index(cargs, RT.Const("int", "47")), RT.Index(cargs, RT.Const("int", "48")), RT.Index(cargs, RT.Const("int", "49")), RT.Index(cargs, RT.Const("int", "50")), RT.Index(cargs, RT.Const("int", "51")), RT.Index(cargs, RT.Const("int", "52")), RT.Index(cargs, RT.Const("int", "53")), RT.Index(cargs, RT.Const("int", "54")), RT.Index(cargs, RT.Const("int", "55")))
		break
	case RT.Key(RT.Const("int", "57")):
		RT.CallIndirect(RT.Cast("FUNV57", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")), RT.Index(cargs, RT.Const("int", "34")), RT.Index(cargs, RT.Const("int", "35")), RT.Index(cargs, RT.Const("int", "36")), RT.Index(cargs, RT.Const("int", "37")), RT.Index(cargs, RT.Const("int", "38")), RT.Index(cargs, RT.Const("int", "39")), RT.Index(cargs, RT.Const("int", "40")), RT.Index(cargs, RT.Const("int", "41")), RT.Index(cargs, RT.Const("int", "42")), RT.Index(cargs, RT.Const("int", "43")), RT.Index(cargs, RT.Const("int", "44")), RT.Index(cargs, RT.Const("int", "45")), RT.Index(cargs, RT.Const("int", "46")), RT.Index(cargs, RT.Const("int", "47")), RT.Index(cargs, RT.Const("int", "48")), RT.Index(cargs, RT.Const("int", "49")), RT.Index(cargs, RT.Const("int", "50")), RT.Index(cargs, RT.Const("int", "51")), RT.Index(cargs, RT.Const("int", "52")), RT.Index(cargs, RT.Const("int", "53")), RT.Index(cargs, RT.Const("int", "54")), RT.Index(cargs, RT.Const("int", "55")), RT.Index(cargs, RT.Const("int", "56")))
		break
	case RT.Key(RT.Const("int", "58")):
		RT.CallIndirect(RT.Cast("FUNV58", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")), RT.Index(cargs, RT.Const("int", "34")), RT.Index(cargs, RT.Const("int", "35")), RT.Index(cargs, RT.Const("int", "36")), RT.Index(cargs, RT.Const("int", "37")), RT.Index(cargs, RT.Const("int", "38")), RT.Index(cargs, RT.Const("int", "39")), RT.Index(cargs, RT.Const("int", "40")), RT.Index(cargs, RT.Const("int", "41")), RT.Index(cargs, RT.Const("int", "42")), RT.Index(cargs, RT.Const("int", "43")), RT.Index(cargs, RT.Const("int", "44")), RT.Index(cargs, RT.Const("int", "45")), RT.Index(cargs, RT.Const("int", "46")), RT.Index(cargs, RT.Const("int", "47")), RT.Index(cargs, RT.Const("int", "48")), RT.Index(cargs, RT.Const("int", "49")), RT.Index(cargs, RT.Const("int", "50")), RT.Index(cargs, RT.Const("int", "51")), RT.Index(cargs, RT.Const("int", "52")), RT.Index(cargs, RT.Const("int", "53")), RT.Index(cargs, RT.Const("int", "54")), RT.Index(cargs, RT.Const("int", "55")), RT.Index(cargs, RT.Const("int", "56")), RT.Index(cargs, RT.Const("int", "57")))
		break
	case RT.Key(RT.Const("int", "59")):
		RT.CallIndirect(RT.Cast("FUNV59", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")), RT.Index(cargs, RT.Const("int", "34")), RT.Index(cargs, RT.Const("int", "35")), RT.Index(cargs, RT.Const("int", "36")), RT.Index(cargs, RT.Const("int", "37")), RT.Index(cargs, RT.Const("int", "38")), RT.Index(cargs, RT.Const("int", "39")), RT.Index(cargs, RT.Const("int", "40")), RT.Index(cargs, RT.Const("int", "41")), RT.Index(cargs, RT.Const("int", "42")), RT.Index(cargs, RT.Const("int", "43")), RT.Index(cargs, RT.Const("int", "44")), RT.Index(cargs, RT.Const("int", "45")), RT.Index(cargs, RT.Const("int", "46")), RT.Index(cargs, RT.Const("int", "47")), RT.Index(cargs, RT.Const("int", "48")), RT.Index(cargs, RT.Const("int", "49")), RT.Index(cargs, RT.Const("int", "50")), RT.Index(cargs, RT.Const("int", "51")), RT.Index(cargs, RT.Const("int", "52")), RT.Index(cargs, RT.Const("int", "53")), RT.Index(cargs, RT.Const("int", "54")), RT.Index(cargs, RT.Const("int", "55")), RT.Index(cargs, RT.Const("int", "56")), RT.Index(cargs, RT.Const("int", "57")), RT.Index(cargs, RT.Const("int", "58")))
		break
	case RT.Key(RT.Const("int", "60")):
		RT.CallIndirect(RT.Cast("FUNV60", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")), RT.Index(cargs, RT.Const("int", "34")), RT.Index(cargs, RT.Const("int", "35")), RT.Index(cargs, RT.Const("int", "36")), RT.Index(cargs, RT.Const("int", "37")), RT.Index(cargs, RT.Const("int", "38")), RT.Index(cargs, RT.Const("int", "39")), RT.Index(cargs, RT.Const("int", "40")), RT.Index(cargs, RT.Const("int", "41")), RT.Index(cargs, RT.Const("int", "42")), RT.Index(cargs, RT.Const("int", "43")), RT.Index(cargs, RT.Const("int", "44")), RT.Index(cargs, RT.Const("int", "45")), RT.Index(cargs, RT.Const("int", "46")), RT.Index(cargs, RT.Const("int", "47")), RT.Index(cargs, RT.Const("int", "48")), RT.Index(cargs, RT.Const("int", "49")), RT.Index(cargs, RT.Const("int", "50")), RT.Index(cargs, RT.Const("int", "51")), RT.Index(cargs, RT.Const("int", "52")), RT.Index(cargs, RT.Const("int", "53")), RT.Index(cargs, RT.Const("int", "54")), RT.Index(cargs, RT.Const("int", "55")), RT.Index(cargs, RT.Const("int", "56")), RT.Index(cargs, RT.Const("int", "57")), RT.Index(cargs, RT.Const("int", "58")), RT.Index(cargs, RT.Const("int", "59")))
		break
	case RT.Key(RT.Const("int", "61")):
		RT.CallIndirect(RT.Cast("FUNV61", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")), RT.Index(cargs, RT.Const("int", "34")), RT.Index(cargs, RT.Const("int", "35")), RT.Index(cargs, RT.Const("int", "36")), RT.Index(cargs, RT.Const("int", "37")), RT.Index(cargs, RT.Const("int", "38")), RT.Index(cargs, RT.Const("int", "39")), RT.Index(cargs, RT.Const("int", "40")), RT.Index(cargs, RT.Const("int", "41")), RT.Index(cargs, RT.Const("int", "42")), RT.Index(cargs, RT.Const("int", "43")), RT.Index(cargs, RT.Const("int", "44")), RT.Index(cargs, RT.Const("int", "45")), RT.Index(cargs, RT.Const("int", "46")), RT.Index(cargs, RT.Const("int", "47")), RT.Index(cargs, RT.Const("int", "48")), RT.Index(cargs, RT.Const("int", "49")), RT.Index(cargs, RT.Const("int", "50")), RT.Index(cargs, RT.Const("int", "51")), RT.Index(cargs, RT.Const("int", "52")), RT.Index(cargs, RT.Const("int", "53")), RT.Index(cargs, RT.Const("int", "54")), RT.Index(cargs, RT.Const("int", "55")), RT.Index(cargs, RT.Const("int", "56")), RT.Index(cargs, RT.Const("int", "57")), RT.Index(cargs, RT.Const("int", "58")), RT.Index(cargs, RT.Const("int", "59")), RT.Index(cargs, RT.Const("int", "60")))
		break
	case RT.Key(RT.Const("int", "62")):
		RT.CallIndirect(RT.Cast("FUNV62", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")), RT.Index(cargs, RT.Const("int", "34")), RT.Index(cargs, RT.Const("int", "35")), RT.Index(cargs, RT.Const("int", "36")), RT.Index(cargs, RT.Const("int", "37")), RT.Index(cargs, RT.Const("int", "38")), RT.Index(cargs, RT.Const("int", "39")), RT.Index(cargs, RT.Const("int", "40")), RT.Index(cargs, RT.Const("int", "41")), RT.Index(cargs, RT.Const("int", "42")), RT.Index(cargs, RT.Const("int", "43")), RT.Index(cargs, RT.Const("int", "44")), RT.Index(cargs, RT.Const("int", "45")), RT.Index(cargs, RT.Const("int", "46")), RT.Index(cargs, RT.Const("int", "47")), RT.Index(cargs, RT.Const("int", "48")), RT.Index(cargs, RT.Const("int", "49")), RT.Index(cargs, RT.Const("int", "50")), RT.Index(cargs, RT.Const("int", "51")), RT.Index(cargs, RT.Const("int", "52")), RT.Index(cargs, RT.Const("int", "53")), RT.Index(cargs, RT.Const("int", "54")), RT.Index(cargs, RT.Const("int", "55")), RT.Index(cargs, RT.Const("int", "56")), RT.Index(cargs, RT.Const("int", "57")), RT.Index(cargs, RT.Const("int", "58")), RT.Index(cargs, RT.Const("int", "59")), RT.Index(cargs, RT.Const("int", "60")), RT.Index(cargs, RT.Const("int", "61")))
		break
	case RT.Key(RT.Const("int", "63")):
		RT.CallIndirect(RT.Cast("FUNV63", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")), RT.Index(cargs, RT.Const("int", "34")), RT.Index(cargs, RT.Const("int", "35")), RT.Index(cargs, RT.Const("int", "36")), RT.Index(cargs, RT.Const("int", "37")), RT.Index(cargs, RT.Const("int", "38")), RT.Index(cargs, RT.Const("int", "39")), RT.Index(cargs, RT.Const("int", "40")), RT.Index(cargs, RT.Const("int", "41")), RT.Index(cargs, RT.Const("int", "42")), RT.Index(cargs, RT.Const("int", "43")), RT.Index(cargs, RT.Const("int", "44")), RT.Index(cargs, RT.Const("int", "45")), RT.Index(cargs, RT.Const("int", "46")), RT.Index(cargs, RT.Const("int", "47")), RT.Index(cargs, RT.Const("int", "48")), RT.Index(cargs, RT.Const("int", "49")), RT.Index(cargs, RT.Const("int", "50")), RT.Index(cargs, RT.Const("int", "51")), RT.Index(cargs, RT.Const("int", "52")), RT.Index(cargs, RT.Const("int", "53")), RT.Index(cargs, RT.Const("int", "54")), RT.Index(cargs, RT.Const("int", "55")), RT.Index(cargs, RT.Const("int", "56")), RT.Index(cargs, RT.Const("int", "57")), RT.Index(cargs, RT.Const("int", "58")), RT.Index(cargs, RT.Const("int", "59")), RT.Index(cargs, RT.Const("int", "60")), RT.Index(cargs, RT.Const("int", "61")), RT.Index(cargs, RT.Const("int", "62")))
		break
	case RT.Key(RT.Const("int", "64")):
		RT.CallIndirect(RT.Cast("FUNV64", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")), RT.Index(cargs, RT.Const("int", "34")), RT.Index(cargs, RT.Const("int", "35")), RT.Index(cargs, RT.Const("int", "36")), RT.Index(cargs, RT.Const("int", "37")), RT.Index(cargs, RT.Const("int", "38")), RT.Index(cargs, RT.Const("int", "39")), RT.Index(cargs, RT.Const("int", "40")), RT.Index(cargs, RT.Const("int", "41")), RT.Index(cargs, RT.Const("int", "42")), RT.Index(cargs, RT.Const("int", "43")), RT.Index(cargs, RT.Const("int", "44")), RT.Index(cargs, RT.Const("int", "45")), RT.Index(cargs, RT.Const("int", "46")), RT.Index(cargs, RT.Const("int", "47")), RT.Index(cargs, RT.Const("int", "48")), RT.Index(cargs, RT.Const("int", "49")), RT.Index(cargs, RT.Const("int", "50")), RT.Index(cargs, RT.Const("int", "51")), RT.Index(cargs, RT.Const("int", "52")), RT.Index(cargs, RT.Const("int", "53")), RT.Index(cargs, RT.Const("int", "54")), RT.Index(cargs, RT.Const("int", "55")), RT.Index(cargs, RT.Const("int", "56")), RT.Index(cargs, RT.Const("int", "57")), RT.Index(cargs, RT.Const("int", "58")), RT.Index(cargs, RT.Const("int", "59")), RT.Index(cargs, RT.Const("int", "60")), RT.Index(cargs, RT.Const("int", "61")), RT.Index(cargs, RT.Const("int", "62")), RT.Index(cargs, RT.Const("int", "63")))
		break
	case RT.Key(RT.Const("int", "65")):
		RT.CallIndirect(RT.Cast("FUNV65", fun), RT.Index(cargs, RT.Const("int", "0")), RT.Index(cargs, RT.Const("int", "1")), RT.Index(cargs, RT.Const("int", "2")), RT.Index(cargs, RT.Const("int", "3")), RT.Index(cargs, RT.Const("int", "4")), RT.Index(cargs, RT.Const("int", "5")), RT.Index(cargs, RT.Const("int", "6")), RT.Index(cargs, RT.Const("int", "7")), RT.Index(cargs, RT.Const("int", "8")), RT.Index(cargs, RT.Const("int", "9")), RT.Index(cargs, RT.Const("int", "10")), RT.Index(cargs, RT.Const("int", "11")), RT.Index(cargs, RT.Const("int", "12")), RT.Index(cargs, RT.Const("int", "13")), RT.Index(cargs, RT.Const("int", "14")), RT.Index(cargs, RT.Const("int", "15")), RT.Index(cargs, RT.Const("int", "16")), RT.Index(cargs, RT.Const("int", "17")), RT.Index(cargs, RT.Const("int", "18")), RT.Index(cargs, RT.Const("int", "19")), RT.Index(cargs, RT.Const("int", "20")), RT.Index(cargs, RT.Const("int", "21")), RT.Index(cargs, RT.Const("int", "22")), RT.Index(cargs, RT.Const("int", "23")), RT.Index(cargs, RT.Const("int", "24")), RT.Index(cargs, RT.Const("int", "25")), RT.Index(cargs, RT.Const("int", "26")), RT.Index(cargs, RT.Const("int", "27")), RT.Index(cargs, RT.Const("int", "28")), RT.Index(cargs, RT.Const("int", "29")), RT.Index(cargs, RT.Const("int", "30")), RT.Index(cargs, RT.Const("int", "31")), RT.Index(cargs, RT.Const("int", "32")), RT.Index(cargs, RT.Const("int", "33")), RT.Index(cargs, RT.Const("int", "34")), RT.Index(cargs, RT.Const("int", "35")), RT.Index(cargs, RT.Const("int", "36")), RT.Index(cargs, RT.Const("int", "37")), RT.Index(cargs, RT.Const("int", "38")), RT.Index(cargs, RT.Const("int", "39")), RT.Index(cargs, RT.Const("int", "40")), RT.Index(cargs, RT.Const("int", "41")), RT.Index(cargs, RT.Const("int", "42")), RT.Index(cargs, RT.Const("int", "43")), RT.Index(cargs, RT.Const("int", "44")), RT.Index(cargs, RT.Const("int", "45")), RT.Index(cargs, RT.Const("int", "46")), RT.Index(cargs, RT.Const("int", "47")), RT.Index(cargs, RT.Const("int", "48")), RT.Index(cargs, RT.Const("int", "49")), RT.Index(cargs, RT.Const("int", "50")), RT.Index(cargs, RT.Const("int", "51")), RT.Index(cargs, RT.Const("int", "52")), RT.Index(cargs, RT.Const("int", "53")), RT.Index(cargs, RT.Const("int", "54")), RT.Index(cargs, RT.Const("int", "55")), RT.Index(cargs, RT.Const("int", "56")), RT.Index(cargs, RT.Const("int", "57")), RT.Index(cargs, RT.Const("int", "58")), RT.Index(cargs, RT.Const("int", "59")), RT.Index(cargs, RT.Const("int", "60")), RT.Index(cargs, RT.Const("int", "61")), RT.Index(cargs, RT.Const("int", "62")), RT.Index(cargs, RT.Const("int", "63")), RT.Index(cargs, RT.Const("int", "64")))
		break
	default:
		RT.Call("errorcall", call, RT.Call("_", RT.Const("string", "\"too many arguments, sorry\"")))
	}
	for RT.Sequence(Assign(LocalRef(&na), RT.Const("int", "0")), Assign(LocalRef(&pa), args)); RT.Truth(RT.Binary("!=", pa, RT.Symbol("R_NilValue"))); RT.Sequence(Assign(LocalRef(&pa), RT.Call("CDR", pa)), RT.Inc(LocalRef(&na), 1, true)) {
		Assign(LocalRef(&p), RT.Index(cargs, na))
		Assign(LocalRef(&arg), RT.Call("CAR", pa))
		Assign(LocalRef(&s), RT.Call("VECTOR_ELT", ans, na))
		Assign(LocalRef(&type_v), func() Value {
			if RT.Truth(checkTypes) {
				return RT.Index(checkTypes, na)
			}
			return RT.Call("TYPEOF", arg)
		}())
		Assign(LocalRef(&n), RT.Call("xlength", arg))
		switch RT.Key(type_v) {
		case RT.Key(RT.Symbol("RAWSXP")):
			if RT.Truth(copy) {
				Assign(LocalRef(&s), RT.Call("allocVector", type_v, n))
				Assign(LocalRef(&ptr), RT.Cast("unsigned char *", p))
				if RT.Truth(n) {
					RT.Call("memcpy", RT.Call("RAW", s), ptr, RT.Binary("*", n, RT.SizeOfType("Rbyte")))
				}
				Assign(LocalRef(&ptr), RT.Binary("+", ptr, RT.Binary("*", n, RT.SizeOfType("Rbyte"))))
				for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, RT.Symbol("NG"))); RT.Inc(LocalRef(&i), 1, true) {
					if RT.Truth(RT.Binary("!=", RT.Deref(RT.Inc(LocalRef(&ptr), 1, true)), RT.Symbol("FILL"))) {
						RT.Call("error", RT.Const("string", "\"array over-run in %s(\\\"%s\\\") in %s argument %d\\n\""), func() Value {
							if RT.Truth(Fort) {
								return RT.Const("string", "\".Fortran\"")
							}
							return RT.Const("string", "\".C\"")
						}(), symName, RT.Call("type2char", type_v), RT.Binary("+", na, RT.Const("int", "1")))
					}
				}
				Assign(LocalRef(&ptr), RT.Cast("unsigned char *", p))
				for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, RT.Symbol("NG"))); RT.Inc(LocalRef(&i), 1, true) {
					if RT.Truth(RT.Binary("!=", RT.Deref(RT.Inc(LocalRef(&ptr), -1, false)), RT.Symbol("FILL"))) {
						RT.Call("error", RT.Const("string", "\"array under-run in %s(\\\"%s\\\") in %s argument %d\\n\""), func() Value {
							if RT.Truth(Fort) {
								return RT.Const("string", "\".Fortran\"")
							}
							return RT.Const("string", "\".C\"")
						}(), symName, RT.Call("type2char", type_v), RT.Binary("+", na, RT.Const("int", "1")))
					}
				}
			}
			break
		case RT.Key(RT.Symbol("INTSXP")):
			if RT.Truth(copy) {
				Assign(LocalRef(&s), RT.Call("allocVector", type_v, n))
				Assign(LocalRef(&ptr), RT.Cast("unsigned char *", p))
				if RT.Truth(n) {
					RT.Call("memcpy", RT.Call("INTEGER", s), ptr, RT.Binary("*", n, RT.SizeOfType("int")))
				}
				Assign(LocalRef(&ptr), RT.Binary("+", ptr, RT.Binary("*", n, RT.SizeOfType("int"))))
				for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, RT.Symbol("NG"))); RT.Inc(LocalRef(&i), 1, true) {
					if RT.Truth(RT.Binary("!=", RT.Deref(RT.Inc(LocalRef(&ptr), 1, true)), RT.Symbol("FILL"))) {
						RT.Call("error", RT.Const("string", "\"array over-run in %s(\\\"%s\\\") in %s argument %d\\n\""), func() Value {
							if RT.Truth(Fort) {
								return RT.Const("string", "\".Fortran\"")
							}
							return RT.Const("string", "\".C\"")
						}(), symName, RT.Call("type2char", type_v), RT.Binary("+", na, RT.Const("int", "1")))
					}
				}
				Assign(LocalRef(&ptr), RT.Cast("unsigned char *", p))
				for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, RT.Symbol("NG"))); RT.Inc(LocalRef(&i), 1, true) {
					if RT.Truth(RT.Binary("!=", RT.Deref(RT.Inc(LocalRef(&ptr), -1, false)), RT.Symbol("FILL"))) {
						RT.Call("error", RT.Const("string", "\"array under-run in %s(\\\"%s\\\") in %s argument %d\\n\""), func() Value {
							if RT.Truth(Fort) {
								return RT.Const("string", "\".Fortran\"")
							}
							return RT.Const("string", "\".C\"")
						}(), symName, RT.Call("type2char", type_v), RT.Binary("+", na, RT.Const("int", "1")))
					}
				}
			}
			break
		case RT.Key(RT.Symbol("LGLSXP")):
			if RT.Truth(copy) {
				Assign(LocalRef(&s), RT.Call("allocVector", type_v, n))
				Assign(LocalRef(&ptr), RT.Cast("unsigned char *", p))
				Assign(LocalRef(&iptr), RT.Cast("int *", ptr))
				for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, n)); RT.Inc(LocalRef(&i), 1, true) {
					Assign(LocalRef(&tmp), RT.Index(iptr, i))
					RT.AssignIndex(RT.Call("LOGICAL", s), i, func() Value {
						if RT.Truth(func() Value {
							if RT.Truth(RT.Binary("==", tmp, RT.Symbol("NA_INTEGER"))) {
								return true
							}
							return RT.Truth(RT.Binary("==", tmp, RT.Const("int", "0")))
						}()) {
							return tmp
						}
						return RT.Const("int", "1")
					}())
				}
				Assign(LocalRef(&ptr), RT.Binary("+", ptr, RT.Binary("*", n, RT.SizeOfType("int"))))
				for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, RT.Symbol("NG"))); RT.Inc(LocalRef(&i), 1, true) {
					if RT.Truth(RT.Binary("!=", RT.Deref(RT.Inc(LocalRef(&ptr), 1, true)), RT.Symbol("FILL"))) {
						RT.Call("error", RT.Const("string", "\"array over-run in %s(\\\"%s\\\") in %s argument %d\\n\""), func() Value {
							if RT.Truth(Fort) {
								return RT.Const("string", "\".Fortran\"")
							}
							return RT.Const("string", "\".C\"")
						}(), symName, RT.Call("type2char", type_v), RT.Binary("+", na, RT.Const("int", "1")))
					}
				}
				Assign(LocalRef(&ptr), RT.Cast("unsigned char *", p))
				for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, RT.Symbol("NG"))); RT.Inc(LocalRef(&i), 1, true) {
					if RT.Truth(RT.Binary("!=", RT.Deref(RT.Inc(LocalRef(&ptr), -1, false)), RT.Symbol("FILL"))) {
						RT.Call("error", RT.Const("string", "\"array under-run in %s(\\\"%s\\\") in %s argument %d\\n\""), func() Value {
							if RT.Truth(Fort) {
								return RT.Const("string", "\".Fortran\"")
							}
							return RT.Const("string", "\".C\"")
						}(), symName, RT.Call("type2char", type_v), RT.Binary("+", na, RT.Const("int", "1")))
					}
				}
			} else {
				Assign(LocalRef(&iptr), RT.Cast("int *", p))
				for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, n)); RT.Inc(LocalRef(&i), 1, true) {
					Assign(LocalRef(&tmp), RT.Index(iptr, i))
					RT.AssignIndex(iptr, i, func() Value {
						if RT.Truth(func() Value {
							if RT.Truth(RT.Binary("==", tmp, RT.Symbol("NA_INTEGER"))) {
								return true
							}
							return RT.Truth(RT.Binary("==", tmp, RT.Const("int", "0")))
						}()) {
							return tmp
						}
						return RT.Const("int", "1")
					}())
				}
			}
			break
		case RT.Key(RT.Symbol("REALSXP")), RT.Key(RT.Symbol("SINGLESXP")):
			if RT.Truth(copy) {
				RT.Call("PROTECT", Assign(LocalRef(&s), RT.Call("allocVector", RT.Symbol("REALSXP"), n)))
				if RT.Truth(func() Value {
					if RT.Truth(RT.Binary("==", type_v, RT.Symbol("SINGLESXP"))) {
						return true
					}
					return RT.Truth(RT.Binary("==", RT.Call("asLogical", RT.Call("getAttrib", arg, RT.Symbol("CSingSymbol"))), RT.Const("int", "1")))
				}()) {
					Assign(LocalRef(&sptr), RT.Cast("float *", p))
					for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, n)); RT.Inc(LocalRef(&i), 1, true) {
						RT.AssignIndex(RT.Call("REAL", s), i, RT.Cast("double", RT.Index(sptr, i)))
					}
				} else {
					Assign(LocalRef(&ptr), RT.Cast("unsigned char *", p))
					if RT.Truth(n) {
						RT.Call("memcpy", RT.Call("REAL", s), ptr, RT.Binary("*", n, RT.SizeOfType("double")))
					}
					Assign(LocalRef(&ptr), RT.Binary("+", ptr, RT.Binary("*", n, RT.SizeOfType("double"))))
					for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, RT.Symbol("NG"))); RT.Inc(LocalRef(&i), 1, true) {
						if RT.Truth(RT.Binary("!=", RT.Deref(RT.Inc(LocalRef(&ptr), 1, true)), RT.Symbol("FILL"))) {
							RT.Call("error", RT.Const("string", "\"array over-run in %s(\\\"%s\\\") in %s argument %d\\n\""), func() Value {
								if RT.Truth(Fort) {
									return RT.Const("string", "\".Fortran\"")
								}
								return RT.Const("string", "\".C\"")
							}(), symName, RT.Call("type2char", type_v), RT.Binary("+", na, RT.Const("int", "1")))
						}
					}
					Assign(LocalRef(&ptr), RT.Cast("unsigned char *", p))
					for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, RT.Symbol("NG"))); RT.Inc(LocalRef(&i), 1, true) {
						if RT.Truth(RT.Binary("!=", RT.Deref(RT.Inc(LocalRef(&ptr), -1, false)), RT.Symbol("FILL"))) {
							RT.Call("error", RT.Const("string", "\"array under-run in %s(\\\"%s\\\") in %s argument %d\\n\""), func() Value {
								if RT.Truth(Fort) {
									return RT.Const("string", "\".Fortran\"")
								}
								return RT.Const("string", "\".C\"")
							}(), symName, RT.Call("type2char", type_v), RT.Binary("+", na, RT.Const("int", "1")))
						}
					}
				}
				RT.Call("UNPROTECT", RT.Const("int", "1"))
			} else {
				if RT.Truth(func() Value {
					if RT.Truth(RT.Binary("==", type_v, RT.Symbol("SINGLESXP"))) {
						return true
					}
					return RT.Truth(RT.Binary("==", RT.Call("asLogical", RT.Call("getAttrib", arg, RT.Symbol("CSingSymbol"))), RT.Const("int", "1")))
				}()) {
					Assign(LocalRef(&s), RT.Call("allocVector", RT.Symbol("REALSXP"), n))
					Assign(LocalRef(&sptr), RT.Cast("float *", p))
					for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, n)); RT.Inc(LocalRef(&i), 1, true) {
						RT.AssignIndex(RT.Call("REAL", s), i, RT.Cast("double", RT.Index(sptr, i)))
					}
				}
			}
			break
		case RT.Key(RT.Symbol("CPLXSXP")):
			if RT.Truth(copy) {
				Assign(LocalRef(&s), RT.Call("allocVector", type_v, n))
				Assign(LocalRef(&ptr), RT.Cast("unsigned char *", p))
				if RT.Truth(n) {
					RT.Call("memcpy", RT.Call("COMPLEX", s), p, RT.Binary("*", n, RT.SizeOfType("Rcomplex")))
				}
				Assign(LocalRef(&ptr), RT.Binary("+", ptr, RT.Binary("*", n, RT.SizeOfType("Rcomplex"))))
				for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, RT.Symbol("NG"))); RT.Inc(LocalRef(&i), 1, true) {
					if RT.Truth(RT.Binary("!=", RT.Deref(RT.Inc(LocalRef(&ptr), 1, true)), RT.Symbol("FILL"))) {
						RT.Call("error", RT.Const("string", "\"array over-run in %s(\\\"%s\\\") in %s argument %d\\n\""), func() Value {
							if RT.Truth(Fort) {
								return RT.Const("string", "\".Fortran\"")
							}
							return RT.Const("string", "\".C\"")
						}(), symName, RT.Call("type2char", type_v), RT.Binary("+", na, RT.Const("int", "1")))
					}
				}
				Assign(LocalRef(&ptr), RT.Cast("unsigned char *", p))
				for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, RT.Symbol("NG"))); RT.Inc(LocalRef(&i), 1, true) {
					if RT.Truth(RT.Binary("!=", RT.Deref(RT.Inc(LocalRef(&ptr), -1, false)), RT.Symbol("FILL"))) {
						RT.Call("error", RT.Const("string", "\"array under-run in %s(\\\"%s\\\") in %s argument %d\\n\""), func() Value {
							if RT.Truth(Fort) {
								return RT.Const("string", "\".Fortran\"")
							}
							return RT.Const("string", "\".C\"")
						}(), symName, RT.Call("type2char", type_v), RT.Binary("+", na, RT.Const("int", "1")))
					}
				}
			}
			break
		case RT.Key(RT.Symbol("STRSXP")):
			if RT.Truth(Fort) {
				Assign(LocalRef(&buf), RT.NewArray(RT.Const("int", "256")))
				RT.Call("strncpy", buf, RT.Cast("char *", p), RT.Const("int", "255"))
				RT.AssignIndex(buf, RT.Const("int", "255"), RT.Const("char", "'\\0'"))
				RT.Call("PROTECT", Assign(LocalRef(&s), RT.Call("allocVector", type_v, RT.Const("int", "1"))))
				RT.Call("SET_STRING_ELT", s, RT.Const("int", "0"), RT.Call("mkChar", buf))
				RT.Call("UNPROTECT", RT.Const("int", "1"))
			} else {
				if RT.Truth(copy) {
					Assign(LocalRef(&ss), arg)
					RT.Call("PROTECT", Assign(LocalRef(&s), RT.Call("allocVector", type_v, n)))
					Assign(LocalRef(&cptr), RT.Cast("char **", p))
					Assign(LocalRef(&cptr0), RT.Cast("char **", RT.Index(cargs0, na)))
					for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, n)); RT.Inc(LocalRef(&i), 1, true) {
						Assign(LocalRef(&ptr), RT.Cast("unsigned char *", RT.Index(cptr, i)))
						RT.Call("SET_STRING_ELT", s, i, RT.Call("mkChar", RT.Index(cptr, i)))
						if RT.Truth(RT.Binary("==", RT.Index(cptr, i), RT.Index(cptr0, i))) {
							Assign(LocalRef(&z), RT.Call("translateChar", RT.Call("STRING_ELT", ss, i)))
							for RT.Sequence(Assign(LocalRef(&j), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", j, RT.Symbol("NG"))); RT.Inc(LocalRef(&j), 1, true) {
								if RT.Truth(RT.Binary("!=", RT.Deref(RT.Inc(LocalRef(&ptr), -1, false)), RT.Symbol("FILL"))) {
									RT.Call("error", RT.Const("string", "\"array under-run in .C(\\\"%s\\\") in character argument %d, element %d\""), symName, RT.Binary("+", na, RT.Const("int", "1")), RT.Cast("int", RT.Binary("+", i, RT.Const("int", "1"))))
								}
							}
							Assign(LocalRef(&ptr), RT.Cast("unsigned char *", RT.Index(cptr, i)))
							Assign(LocalRef(&ptr), RT.Binary("+", ptr, RT.Binary("+", RT.Call("strlen", z), RT.Const("int", "1"))))
							for RT.Sequence(Assign(LocalRef(&j), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", j, RT.Symbol("NG"))); RT.Inc(LocalRef(&j), 1, true) {
								if RT.Truth(RT.Binary("!=", RT.Deref(RT.Inc(LocalRef(&ptr), 1, true)), RT.Symbol("FILL"))) {
									Assign(LocalRef(&p), ptr)
									for RT.Sequence(Assign(LocalRef(&k), RT.Const("int", "1"))); RT.Truth(RT.Binary("<", k, RT.Binary("-", RT.Symbol("NG"), j))); RT.Sequence(RT.Inc(LocalRef(&k), 1, true), RT.Inc(LocalRef(&p), 1, true)) {
										if RT.Truth(RT.Binary("==", RT.Deref(p), RT.Symbol("FILL"))) {
											RT.AssignDeref(p, RT.Const("char", "'\\0'"))
										}
									}
									RT.Call("error", RT.Const("string", "\"array over-run in .C(\\\"%s\\\") in character argument %d, element %d\\n'%s'->'%s'\\n\""), symName, RT.Binary("+", na, RT.Const("int", "1")), RT.Cast("int", RT.Binary("+", i, RT.Const("int", "1"))), z, RT.Index(cptr, i))
								}
							}
						}
					}
					RT.Call("UNPROTECT", RT.Const("int", "1"))
				} else {
					RT.Call("PROTECT", Assign(LocalRef(&s), RT.Call("allocVector", type_v, n)))
					Assign(LocalRef(&cptr), RT.Cast("char **", p))
					for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, n)); RT.Inc(LocalRef(&i), 1, true) {
						RT.Call("SET_STRING_ELT", s, i, RT.Call("mkChar", RT.Index(cptr, i)))
					}
					RT.Call("UNPROTECT", RT.Const("int", "1"))
				}
			}
			break
		default:
			break
		}
		if RT.Truth(RT.Binary("!=", s, arg)) {
			RT.Call("PROTECT", s)
			RT.Call("SHALLOW_DUPLICATE_ATTRIB", s, arg)
			RT.Call("SET_VECTOR_ELT", ans, na, s)
			RT.Call("UNPROTECT", RT.Const("int", "1"))
		}
	}
	RT.Call("UNPROTECT", RT.Const("int", "1"))
	RT.Call("vmaxset", vmax)
	return ans
}
