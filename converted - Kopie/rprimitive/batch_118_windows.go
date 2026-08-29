//go:build windows

package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_url(call, op, args, env Value) Value {
	var (
		scmd         Value
		sopen        Value
		ans          Value
		class        Value
		enc          Value
		headers      Value
		headers_flat Value
		class2       Value
		url          Value
		open         Value
		ncon         Value
		block        Value
		defmeth      Value
		meth         Value
		winmeth      Value
		ienc         Value
		con          Value
		raw          Value
		type_v       Value
		inet         Value
		cmeth        Value
		lheaders     Value
		nh           Value
		efn          Value
		subtype      Value
		compress     Value
		ct           Value
	)
	Assign(LocalRef(&headers), RT.Symbol("R_NilValue"))
	Assign(LocalRef(&headers_flat), RT.Symbol("R_NilValue"))
	Assign(LocalRef(&class2), RT.Const("string", "\"url\""))
	Assign(LocalRef(&meth), RT.Const("int", "0"))
	Assign(LocalRef(&winmeth), RT.Const("int", "0"))
	Assign(LocalRef(&ienc), RT.Symbol("CE_NATIVE"))
	Assign(LocalRef(&con), RT.Symbol("NULL"))
	Assign(LocalRef(&raw), RT.Symbol("FALSE"))
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&scmd), RT.Call("CAR", args))
	if RT.Truth(func() Value {
		if RT.Truth(func() Value {
			if RT.Truth(RT.Unary("!", RT.Call("isString", scmd))) {
				return true
			}
			return RT.Truth(RT.Binary("!=", RT.Call("LENGTH", scmd), RT.Const("int", "1")))
		}()) {
			return true
		}
		return RT.Truth(RT.Binary("==", RT.Call("STRING_ELT", scmd, RT.Const("int", "0")), RT.Symbol("NA_STRING")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"description\""))
	}
	if RT.Truth(RT.Binary(">", RT.Call("LENGTH", scmd), RT.Const("int", "1"))) {
		RT.Call("warning", RT.Call("_", RT.Const("string", "\"only first element of 'description' argument used\"")))
	}
	Assign(LocalRef(&winmeth), RT.Const("int", "1"))
	if RT.Truth(func() Value {
		if !RT.Truth(RT.Binary("==", RT.Call("PRIMVAL", op), RT.Const("int", "1"))) {
			return false
		}
		return RT.Truth(RT.Unary("!", RT.Call("IS_ASCII", RT.Call("STRING_ELT", scmd, RT.Const("int", "0")))))
	}()) {
		Assign(LocalRef(&ienc), RT.Symbol("CE_UTF8"))
		Assign(LocalRef(&url), RT.Call("trCharUTF8", RT.Call("STRING_ELT", scmd, RT.Const("int", "0"))))
	} else {
		Assign(LocalRef(&ienc), RT.Call("getCharCE", RT.Call("STRING_ELT", scmd, RT.Const("int", "0"))))
		if RT.Truth(RT.Binary("==", ienc, RT.Symbol("CE_UTF8"))) {
			Assign(LocalRef(&url), RT.Call("CHAR", RT.Call("STRING_ELT", scmd, RT.Const("int", "0"))))
		} else {
			Assign(LocalRef(&url), RT.Call("translateCharFP", RT.Call("STRING_ELT", scmd, RT.Const("int", "0"))))
		}
	}
	Assign(LocalRef(&type_v), RT.Symbol("HTTPsh"))
	Assign(LocalRef(&inet), RT.Symbol("TRUE"))
	if RT.Truth(RT.Binary("==", RT.Call("strncmp", url, RT.Const("string", "\"http://\""), RT.Const("int", "7")), RT.Const("int", "0"))) {
		Assign(LocalRef(&type_v), RT.Symbol("HTTPsh"))
	} else {
		if RT.Truth(RT.Binary("==", RT.Call("strncmp", url, RT.Const("string", "\"ftp://\""), RT.Const("int", "6")), RT.Const("int", "0"))) {
			Assign(LocalRef(&type_v), RT.Symbol("FTPsh"))
		} else {
			if RT.Truth(RT.Binary("==", RT.Call("strncmp", url, RT.Const("string", "\"https://\""), RT.Const("int", "8")), RT.Const("int", "0"))) {
				Assign(LocalRef(&type_v), RT.Symbol("HTTPSsh"))
			} else {
				if RT.Truth(RT.Binary("==", RT.Call("strncmp", url, RT.Const("string", "\"ftps://\""), RT.Const("int", "7")), RT.Const("int", "0"))) {
					Assign(LocalRef(&type_v), RT.Symbol("FTPSsh"))
				} else {
					Assign(LocalRef(&inet), RT.Symbol("FALSE"))
				}
			}
		}
	}
	Assign(LocalRef(&sopen), RT.Call("CADR", args))
	if RT.Truth(func() Value {
		if RT.Truth(RT.Unary("!", RT.Call("isString", sopen))) {
			return true
		}
		return RT.Truth(RT.Binary("!=", RT.Call("LENGTH", sopen), RT.Const("int", "1")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"open\""))
	}
	Assign(LocalRef(&open), RT.Call("CHAR", RT.Call("STRING_ELT", sopen, RT.Const("int", "0"))))
	Assign(LocalRef(&block), RT.Call("asLogical", RT.Call("CADDR", args)))
	if RT.Truth(RT.Binary("==", block, RT.Symbol("NA_LOGICAL"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"blocking\""))
	}
	Assign(LocalRef(&enc), RT.Call("CADDDR", args))
	if RT.Truth(func() Value {
		if RT.Truth(func() Value {
			if RT.Truth(RT.Unary("!", RT.Call("isString", enc))) {
				return true
			}
			return RT.Truth(RT.Binary("!=", RT.Call("LENGTH", enc), RT.Const("int", "1")))
		}()) {
			return true
		}
		return RT.Truth(RT.Binary(">", RT.Call("strlen", RT.Call("CHAR", RT.Call("STRING_ELT", enc, RT.Const("int", "0")))), RT.Const("int", "100")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"encoding\""))
	}
	Assign(LocalRef(&cmeth), RT.Call("CHAR", RT.Call("asChar", RT.Call("CAD4R", args))))
	Assign(LocalRef(&meth), RT.Call("streql", cmeth, RT.Const("string", "\"libcurl\"")))
	Assign(LocalRef(&defmeth), RT.Call("streql", cmeth, RT.Const("string", "\"default\"")))
	if RT.Truth(defmeth) {
		Assign(LocalRef(&meth), RT.Const("int", "1"))
	}
	if RT.Truth(RT.Call("streql", cmeth, RT.Const("string", "\"wininet\""))) {
		Assign(LocalRef(&winmeth), RT.Const("int", "1"))
	} else {
		if RT.Truth(RT.Call("streql", cmeth, RT.Const("string", "\"internal\""))) {
			Assign(LocalRef(&winmeth), RT.Const("int", "0"))
		}
	}
	if RT.Truth(RT.Binary("==", RT.Call("PRIMVAL", op), RT.Const("int", "1"))) {
		Assign(LocalRef(&raw), RT.Call("asRbool", RT.Call("CAD5R", args), call))
		if RT.Truth(RT.Binary("==", raw, RT.Symbol("NA_LOGICAL"))) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"raw\""))
		}
	}
	if RT.Truth(RT.Binary("==", RT.Call("PRIMVAL", op), RT.Const("int", "0"))) {
		Assign(LocalRef(&lheaders), RT.Call("CAD5R", args))
		if RT.Truth(RT.Unary("!", RT.Call("isNull", lheaders))) {
			Assign(LocalRef(&headers), RT.Call("VECTOR_ELT", lheaders, RT.Const("int", "0")))
			Assign(LocalRef(&headers_flat), RT.Call("VECTOR_ELT", lheaders, RT.Const("int", "1")))
		}
	}
	if RT.Truth(RT.Unary("!", meth)) {
		if RT.Truth(RT.Binary("==", RT.Call("strncmp", url, RT.Const("string", "\"ftps://\""), RT.Const("int", "7")), RT.Const("int", "0"))) {
			if RT.Truth(defmeth) {
				Assign(LocalRef(&meth), RT.Const("int", "1"))
			} else {
				RT.Call("error", RT.Const("string", "\"ftps:// URLs are not supported by this method\""))
			}
		}
		if RT.Truth(func() Value {
			if !RT.Truth(RT.Unary("!", winmeth)) {
				return false
			}
			return RT.Truth(RT.Binary("==", RT.Call("strncmp", url, RT.Const("string", "\"https://\""), RT.Const("int", "8")), RT.Const("int", "0")))
		}()) {
			if RT.Truth(defmeth) {
				Assign(LocalRef(&meth), RT.Const("int", "1"))
			} else {
				RT.Call("error", RT.Const("string", "\"https:// URLs are not supported by this method\""))
			}
		}
	}
	Assign(LocalRef(&ncon), RT.Call("NextConnection"))
	if RT.Truth(RT.Binary("==", RT.Call("strncmp", url, RT.Const("string", "\"file://\""), RT.Const("int", "7")), RT.Const("int", "0"))) {
		Assign(LocalRef(&nh), RT.Const("int", "7"))
		if RT.Truth(func() Value {
			if !RT.Truth(func() Value {
				if !RT.Truth(RT.Binary(">", RT.Call("strlen", url), RT.Const("int", "9"))) {
					return false
				}
				return RT.Truth(RT.Binary("==", RT.Index(url, RT.Const("int", "7")), RT.Const("char", "'/'")))
			}()) {
				return false
			}
			return RT.Truth(RT.Binary("==", RT.Index(url, RT.Const("int", "9")), RT.Const("char", "':'")))
		}()) {
			Assign(LocalRef(&nh), RT.Const("int", "8"))
		}
		Assign(LocalRef(&con), RT.Call("newfile", RT.Binary("+", url, nh), ienc, func() Value {
			if RT.Truth(RT.Call("strlen", open)) {
				return open
			}
			return RT.Const("string", "\"r\"")
		}(), raw))
		Assign(LocalRef(&class2), RT.Const("string", "\"file\""))
	} else {
		if RT.Truth(inet) {
			if RT.Truth(meth) {
				Assign(LocalRef(&con), RT.Call("R_newCurlUrl", url, func() Value {
					if RT.Truth(RT.Call("strlen", open)) {
						return open
					}
					return RT.Const("string", "\"r\"")
				}(), headers, RT.Const("int", "0")))
			} else {
				if RT.Truth(RT.Unary("!", winmeth)) {
					RT.Call("error", RT.Call("_", RT.Const("string", "\"the 'internal' method of url() is defunct for http:// and ftp:// URLs\"")))
				}
				Assign(LocalRef(&con), RT.Call("R_newurl", url, func() Value {
					if RT.Truth(RT.Call("strlen", open)) {
						return open
					}
					return RT.Const("string", "\"r\"")
				}(), headers_flat, winmeth))
				RT.AssignField(RT.Cast("Rurlconn", RT.Field(con, "private")), "type", type_v)
			}
		} else {
			if RT.Truth(RT.Binary("==", RT.Call("PRIMVAL", op), RT.Const("int", "1"))) {
				if RT.Truth(RT.Binary("==", RT.Call("strlen", url), RT.Const("int", "0"))) {
					if RT.Truth(RT.Unary("!", RT.Call("strlen", open))) {
						Assign(LocalRef(&open), RT.Const("string", "\"w+\""))
					}
					if RT.Truth(func() Value {
						if !RT.Truth(RT.Binary("!=", RT.Call("strcmp", open, RT.Const("string", "\"w+\"")), RT.Const("int", "0"))) {
							return false
						}
						return RT.Truth(RT.Binary("!=", RT.Call("strcmp", open, RT.Const("string", "\"w+b\"")), RT.Const("int", "0")))
					}()) {
						Assign(LocalRef(&open), RT.Const("string", "\"w+\""))
						RT.Call("warning", RT.Call("_", RT.Const("string", "\"file(\\\"\\\") only supports open = \\\"w+\\\" and open = \\\"w+b\\\": using the former\"")))
					}
				}
				if RT.Truth(func() Value {
					if RT.Truth(RT.Binary("==", RT.Call("strcmp", url, RT.Const("string", "\"clipboard\"")), RT.Const("int", "0"))) {
						return true
					}
					return RT.Truth(RT.Binary("==", RT.Call("strncmp", url, RT.Const("string", "\"clipboard-\""), RT.Const("int", "10")), RT.Const("int", "0")))
				}()) {
					Assign(LocalRef(&con), RT.Call("newclp", url, func() Value {
						if RT.Truth(RT.Call("strlen", open)) {
							return open
						}
						return RT.Const("string", "\"r\"")
					}()))
				} else {
					Assign(LocalRef(&efn), RT.Call("R_ExpandFileName", url))
					if RT.Truth(func() Value {
						if !RT.Truth(RT.Unary("!", raw)) {
							return false
						}
						return RT.Truth(func() Value {
							if RT.Truth(func() Value {
								if RT.Truth(RT.Unary("!", RT.Call("strlen", open))) {
									return true
								}
								return RT.Truth(RT.Call("streql", open, RT.Const("string", "\"r\"")))
							}()) {
								return true
							}
							return RT.Truth(RT.Call("streql", open, RT.Const("string", "\"rt\"")))
						}())
					}()) {
						Assign(LocalRef(&subtype), RT.Const("int", "0"))
						Assign(LocalRef(&compress), RT.Const("int", "0"))
						Assign(LocalRef(&ct), RT.Call("comp_type_from_file", efn, RT.Symbol("FALSE"), LocalRef(&subtype)))
						switch RT.Key(ct) {
						case RT.Key(RT.Symbol("COMP_UNKNOWN")):
							Assign(LocalRef(&con), RT.Call("newfile", url, ienc, func() Value {
								if RT.Truth(RT.Call("strlen", open)) {
									return open
								}
								return RT.Const("string", "\"r\"")
							}(), raw))
							break
						case RT.Key(RT.Symbol("COMP_GZ")):
							Assign(LocalRef(&con), RT.Call("newgzfile", url, func() Value {
								if RT.Truth(RT.Call("strlen", open)) {
									return open
								}
								return RT.Const("string", "\"rt\"")
							}(), compress))
							break
						case RT.Key(RT.Symbol("COMP_BZ")):
							Assign(LocalRef(&con), RT.Call("newbzfile", url, func() Value {
								if RT.Truth(RT.Call("strlen", open)) {
									return open
								}
								return RT.Const("string", "\"rt\"")
							}(), compress))
							break
						case RT.Key(RT.Symbol("COMP_XZ")):
							Assign(LocalRef(&con), RT.Call("newxzfile", url, func() Value {
								if RT.Truth(RT.Call("strlen", open)) {
									return open
								}
								return RT.Const("string", "\"rt\"")
							}(), subtype, compress))
							break
						case RT.Key(RT.Symbol("COMP_ZSTD")):
							Assign(LocalRef(&con), RT.Call("newzstdfile", url, func() Value {
								if RT.Truth(RT.Call("strlen", open)) {
									return open
								}
								return RT.Const("string", "\"rt\"")
							}(), compress))
							break
						}
					} else {
						Assign(LocalRef(&con), RT.Call("newfile", url, ienc, func() Value {
							if RT.Truth(RT.Call("strlen", open)) {
								return open
							}
							return RT.Const("string", "\"r\"")
						}(), raw))
					}
				}
				Assign(LocalRef(&class2), RT.Const("string", "\"file\""))
			} else {
				RT.Call("error", RT.Call("_", RT.Const("string", "\"URL scheme unsupported by this method\"")))
			}
		}
	}
	RT.AssignIndex(RT.Symbol("Connections"), ncon, con)
	RT.AssignField(con, "blocking", RT.Cast("Rboolean", block))
	RT.Call("strncpy", RT.Field(con, "encname"), RT.Call("CHAR", RT.Call("STRING_ELT", enc, RT.Const("int", "0"))), RT.Const("int", "100"))
	RT.AssignIndex(RT.Field(con, "encname"), RT.Binary("-", RT.Const("int", "100"), RT.Const("int", "1")), RT.Const("char", "'\\0'"))
	if RT.Truth(func() Value {
		if !RT.Truth(RT.Index(RT.Field(con, "encname"), RT.Const("int", "0"))) {
			return false
		}
		return RT.Truth(RT.Unary("!", RT.Call("streql", RT.Field(con, "encname"), RT.Const("string", "\"native.enc\""))))
	}()) {
		RT.AssignField(con, "canseek", RT.Const("int", "0"))
	}
	RT.AssignField(con, "ex_ptr", RT.Call("PROTECT", RT.Call("R_MakeExternalPtr", RT.Field(con, "id"), RT.Call("install", RT.Const("string", "\"connection\"")), RT.Symbol("R_NilValue"))))
	if RT.Truth(RT.Call("strlen", open)) {
		RT.Call("checked_open", ncon)
	}
	RT.Call("PROTECT", Assign(LocalRef(&ans), RT.Call("ScalarInteger", ncon)))
	RT.Call("PROTECT", Assign(LocalRef(&class), RT.Call("allocVector", RT.Symbol("STRSXP"), RT.Const("int", "2"))))
	RT.Call("SET_STRING_ELT", class, RT.Const("int", "0"), RT.Call("mkChar", class2))
	RT.Call("SET_STRING_ELT", class, RT.Const("int", "1"), RT.Call("mkChar", RT.Const("string", "\"connection\"")))
	RT.Call("classgets", ans, class)
	RT.Call("setAttrib", ans, RT.Symbol("R_ConnIdSymbol"), RT.Field(con, "ex_ptr"))
	RT.Call("R_RegisterCFinalizerEx", RT.Field(con, "ex_ptr"), RT.Symbol("conFinalizer"), RT.Symbol("FALSE"))
	RT.Call("UNPROTECT", RT.Const("int", "3"))
	return ans
}
