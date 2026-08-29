package runtime

import (
	"fmt"
	"r2go/cir"
	"sort"
)

// CIRHost is the bridge between generated C control flow and the existing
// Pure-Go R value kernels. It deliberately rejects unknown C helpers instead
// of silently emulating GNU R or reporting a false capability.
type CIRHost struct {
	Calls    []string
	Stores   []string
	Branches []string
}

var CIRPureGoHelpers = map[string]bool{
	"checkArity": true, "PROTECT": true, "UNPROTECT": true, "PROTECT_WITH_INDEX": true, "REPROTECT": true,
	"length": true, "xlength": true, "asInteger": true, "asReal": true, "asLogical": true, "asChar": true,
	"duplicate": true, "shallow_duplicate": true, "allocVector": true, "allocList": true, "allocMatrix": true,
	"ScalarInteger": true, "ScalarReal": true, "ScalarLogical": true, "mkString": true, "install": true,
	"error": true, "errorcall": true, "warning": true, "warningcall": true, "GetOption1": true,
	"isString": true, "isInteger": true, "isReal": true, "isLogical": true, "isNewList": true,
	"getAttrib": true, "setAttrib": true, "namesgets": true, "classgets": true,
	"STRING_ELT": true, "SET_STRING_ELT": true, "INTEGER_ELT": true, "SET_INTEGER_ELT": true,
	"REAL_ELT": true, "SET_REAL_ELT": true, "VECTOR_ELT": true, "SET_VECTOR_ELT": true,
	"GetRNGstate": true, "PutRNGstate": true, "R_FINITE": true, "ISNA": true, "ISNAN": true,
}

func (h *CIRHost) Call(name string) error {
	h.Calls = append(h.Calls, name)
	if CIRPureGoHelpers[name] {
		return nil
	}
	return fmt.Errorf("unlowered C helper %s", name)
}
func (h *CIRHost) Store(name string) error  { h.Stores = append(h.Stores, name); return nil }
func (h *CIRHost) Branch(kind string) error { h.Branches = append(h.Branches, kind); return nil }

// ExecuteCIRProgram exposes the same strict execution boundary used by the
// primitive dispatcher. A program is executable only once every helper routed
// by its generated CIR can be served by this Pure-Go host.
func ExecuteCIRProgram(program cir.Program) error { return cir.Execute(program, &CIRHost{}) }

func CIRProgramUnsupportedCalls(program cir.Program) []string {
	missing := map[string]bool{}
	for _, instruction := range program.Instructions {
		if instruction.Op == cir.Call && !CIRPureGoHelpers[instruction.A] { missing[instruction.A] = true }
	}
	result := make([]string, 0, len(missing))
	for name := range missing { result = append(result, name) }
	sort.Strings(result)
	return result
}

func CIRProgramExecutable(name string) bool {
	program, ok := GeneratedCIRPrograms[name]
	return ok && len(CIRProgramUnsupportedCalls(program)) == 0
}
