# Missing / constrained helpers

This file records functionality that cannot be made GNU-R-equivalent from the isolated batch corpus alone, or that conflicts with the no-cgo constraint. The Go package still compiles; affected entry points return explicit R-style error Values rather than silently succeeding.

## Runtime model provided

`runtime.go` provides Pure-Go `Value`, `Env`, `Conn`, attributes, names/dim metadata, NA masks, cloning/copy-on-write helpers, vector recycling helpers, coercions, filesystem/connection adapters, regexp/string helpers, numeric/raw algorithms, and structured error Values. `PROTECT`/`UNPROTECT` are intentionally absent because Go GC owns object lifetime. CAR/CADR/CDR-style access is represented by list indexing helpers instead of C cons-cell macros.

## Native dynamic loading boundary

These C entry points fundamentally call dynamically loaded native routines or R DLL registries. Under `no cgo` they cannot execute arbitrary C/Fortran symbols. Their Pure-Go implementation is the defined error path; a future Pure-Go plugin ABI would be needed for equivalent extensibility.

`do_External`, `do_Externalgr`, `do_dotCode`, `do_dotcall`, `do_dotcallgr`, `do_dynload`, `do_dynunload`, `do_getDllTable`, `do_getRegisteredRoutines`, `do_getSymbolInfo`, `do_isloaded`

## Transport/compression boundary

Sockets, FIFOs, URL and compressed R connection implementations need protocol/compression state that is not present in the extracted batch corpus. Basic file/raw/text/stdin/stdout/stderr connection mechanics are implemented; these specialized constructors return explicit errors.

`do_fifo`, `do_gzcon`, `do_gzfile`, `do_pipe`, `do_serversocket`, `do_sockconn`, `do_unz`, `do_url`

## GNU R serialization/lazy-load wire format

A deterministic internal Pure-Go serializer exists for basic values, but GNU R XDR/native serialization and lazy-load database framing are not reconstructed from these isolated functions. Exact interoperability needs a full R serialization codec.

`do_getVarsFromFrame`, `do_lazyLoadDBfetch`, `do_lazyLoadDBflush`, `do_lazyLoadDBinsertValue`, `do_serialize`, `do_serializeToConn`, `do_unserializeFromConn`

## Evaluator/environment integration

The package includes lexical `Env` lookup/mutation, locks, list conversion, parent links and basic function callbacks. Entries that depend on promises, active bindings, namespace registries, search path semantics, or full R evaluation may return explicit evaluator-state errors for unsupported branches.

`do_activeBndFun`, `do_as_environment`, `do_attach`, `do_bindingType`, `do_bndIsActive`, `do_bndIsLocked`, `do_builtins`, `do_delayedBindingEnv`, `do_delayedBindingExpr`, `do_detach`, `do_dotDelayedEnv`, `do_dotDelayedExpr`, `do_dotForcedExpr`, `do_dotType`, `do_eapply`, `do_env2list`, `do_envIsLocked`, `do_envprofile`, `do_forcedBindingExpr`, `do_getNSRegistry`, `do_getNSValue`, `do_getRegNS`, `do_importIntoEnv`, `do_isNSEnv`, `do_list2env`, `do_lockBnd`, `do_lockEnv`, `do_ls`, `do_mget`, `do_mkActiveBnd`, `do_mkUnbound`, `do_pos2env`, `do_regNS`, `do_remove`, `do_search`, `do_unregNS`

## Corpus-external R internals

These entry points still depend on algorithms or global runtime machinery not contained in the extracted function body (for example bytecode evaluator/context stacks, full LAPACK dispatch, native R save/load formats, or other R-internal helpers). They compile and fail explicitly rather than using pseudo-code or an empty body.

`do_addCondHands`, `do_addGlobHands`, `do_addRestart`, `do_addTryHandlers`, `do_asfunction`, `do_bcclose`, `do_bcprofcounts`, `do_bcprofstart`, `do_bcprofstop`, `do_bodyCode`, `do_curlDownload`, `do_curlGetHeaders`, `do_curlVersion`, `do_debug`, `do_declare`, `do_delayed`, `do_disassemble`, `do_docall`, `do_dump`, `do_envirName`, `do_envirgets`, `do_getconst`, `do_growconst`, `do_inspect`, `do_invokeRestart`, `do_is_builtin_internal`, `do_lapack`, `do_load`, `do_loadFromConn2`, `do_loadfile`, `do_makelazy`, `do_matchcall`, `do_mkcode`, `do_mmap_file`, `do_munmap_file`, `do_onexit`, `do_parentenvgets`, `do_putconst`, `do_readln`, `do_recall`, `do_recordGraphics`, `do_regFinaliz`, `do_save`, `do_saveToConn`, `do_savefile`, `do_scan`, `do_shellexec`, `do_signalCondition`, `do_standardGeneric`, `do_substitute`, `do_sys`, `do_sysbrowser`, `do_tailcall`, `do_traceback`, `do_tryCatchHelper`, `do_tryWrap`, `do_wrap_meta`

## Deliberate C-API mappings

- `SEXP` -> `Value`; vector payloads are typed Go slices and NA is tracked separately from IEEE NaN.
- `allocVector` -> typed constructors/slice allocation.
- `CAR`/`CADR`/... -> `arg(args, i)` / `Value.V`.
- attributes (`getAttrib`/`setAttrib`) -> `Attr` map with shape/name preservation helpers.
- `PROTECT`/`UNPROTECT` -> no direct equivalent; Go GC provides reachability-based lifetime.
- R errors/warnings -> `Value{Kind: Error}` for errors; warning-only behavior is generally folded into deterministic return behavior because the isolated runtime has no warning sink.
