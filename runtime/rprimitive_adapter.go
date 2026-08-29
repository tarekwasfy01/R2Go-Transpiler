package runtime

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	rprimitive "r2go/converted/rprimitive"
	"r2go/syntax"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// rprimitive uses a package-global DynamicRuntime by design. Serializing the
// temporary installation makes it safe to use from the evaluator while the
// long-term runtime hooks are progressively filled in.
var translatedRuntimeMu sync.Mutex

type translatedRuntimeBridge struct {
	*rprimitive.Runtime
	ctx *Context
	env *Environment

	environmentLocked map[*Environment]bool
	bindingLocked     map[*Environment]map[string]bool
	activeBindings    map[*Environment]map[string]rprimitive.Value
	debugFlags        map[string]bool
	stepFlags         map[string]bool
	connections       *dynamicConnectionTable
}

type dynamicSymbol string
type dynamicNAString struct{}
type dynamicNAReal struct{}
type dynamicUnbound struct{}
type dynamicMissingArg struct{}

type dynamicExternalPtr struct {
	Addr      rprimitive.Value
	Tag       rprimitive.Value
	Protected rprimitive.Value
	Finalizer rprimitive.Value
	OnExit    bool
}

type dynamicConnectionTable struct {
	values map[int]rprimitive.Value
}

type dynamicConnectionRef struct {
	table *dynamicConnectionTable
	index int
}

func (r dynamicConnectionRef) Get() rprimitive.Value {
	if r.table == nil {
		return nil
	}
	return r.table.values[r.index]
}

func (r dynamicConnectionRef) Set(value rprimitive.Value) {
	if r.table == nil {
		return
	}
	if value == nil {
		delete(r.table.values, r.index)
		return
	}
	r.table.values[r.index] = value
}

type dynamicHandlerEntry struct {
	Class   string
	Parent  rprimitive.Value
	Handler rprimitive.Value
	Target  rprimitive.Value
	Result  rprimitive.Value
	Calling bool
}

// The translated C boundary also stores a few sentinels and bookkeeping
// objects inside host Lists.  Lists carry runtime.Value, so give those
// deliberately opaque compatibility values stable host identities instead of
// silently coercing them to NULL.
func (dynamicSymbol) Kind() Kind            { return CharacterKind }
func (s dynamicSymbol) String() string      { return string(s) }
func (*dynamicNAString) Kind() Kind         { return CharacterKind }
func (*dynamicNAString) String() string     { return "NA" }
func (*dynamicNAReal) Kind() Kind           { return DoubleKind }
func (*dynamicNAReal) String() string       { return "NA" }
func (*dynamicUnbound) Kind() Kind          { return NullKind }
func (*dynamicUnbound) String() string      { return "<unbound>" }
func (*dynamicMissingArg) Kind() Kind       { return NullKind }
func (*dynamicMissingArg) String() string   { return "<missing>" }
func (*dynamicExternalPtr) Kind() Kind      { return ListKind }
func (*dynamicExternalPtr) String() string  { return "<externalptr>" }
func (*dynamicHandlerEntry) Kind() Kind     { return ListKind }
func (*dynamicHandlerEntry) String() string { return "<condition-handler>" }

var (
	dynamicNAStringValue = &dynamicNAString{}
	dynamicNARealValue   = &dynamicNAReal{}
	dynamicUnboundValue  = &dynamicUnbound{}
	dynamicMissingValue  = &dynamicMissingArg{}
)

func newTranslatedRuntimeBridge(ctx *Context, env *Environment) *translatedRuntimeBridge {
	base := rprimitive.NewRuntime()
	base.InstallCoreCompatibility()
	base.SetSymbol("R_HandlerStack", NullValue)
	base.SetSymbol("R_RestartStack", NullValue)
	base.SetSymbol("R_HandlerResultToken", NullValue)
	return &translatedRuntimeBridge{
		Runtime:           base,
		ctx:               ctx,
		env:               env,
		environmentLocked: map[*Environment]bool{},
		bindingLocked:     map[*Environment]map[string]bool{},
		activeBindings:    map[*Environment]map[string]rprimitive.Value{},
		debugFlags:        map[string]bool{},
		stepFlags:         map[string]bool{},
		connections:       &dynamicConnectionTable{values: map[int]rprimitive.Value{}},
	}
}

func (b *translatedRuntimeBridge) Symbol(name string) rprimitive.Value {
	switch name {
	case "R_NilValue", "NULL":
		return NullValue
	case "NA_LOGICAL":
		return NA
	case "NA_INTEGER":
		return int64(math.MinInt32)
	case "NA_REAL":
		return dynamicNARealValue
	case "R_NaN":
		return math.NaN()
	case "R_PosInf":
		return math.Inf(1)
	case "R_NegInf":
		return math.Inf(-1)
	case "NA_STRING":
		return dynamicNAStringValue
	case "R_UnboundValue":
		return dynamicUnboundValue
	case "R_MissingArg":
		return dynamicMissingValue
	case "Connections":
		return b.connections
	case "TRUE", "true":
		return True
	case "FALSE", "false":
		return False
	case "R_GlobalEnv", "R_BaseEnv":
		if root := b.rootEnvironment(); root != nil {
			return root
		}
		return NullValue
	case "R_EmptyEnv":
		return NullValue
	case "STRSXP", "VECSXP", "INTSXP", "REALSXP", "LGLSXP", "CPLXSXP", "RAWSXP", "ENVSXP", "CLOSXP", "LISTSXP", "SYMSXP", "LANGSXP", "SPECIALSXP", "BUILTINSXP", "EXTPTRSXP", "EXPRSXP", "BCODESXP", "PROMSXP", "NILSXP":
		return name
	case "R_BindingTypeUnbound", "R_BindingTypeValue", "R_BindingTypeMissing", "R_BindingTypeDelayed", "R_BindingTypeForced", "R_BindingTypeActive", "R_ClassSymbol", "R_NameSymbol", "R_NamesSymbol", "R_DimSymbol", "R_DimNamesSymbol", "R_LevelsSymbol", "R_RowNamesSymbol", "R_ConnIdSymbol":
		return name
	}
	return b.Runtime.Symbol(name)
}

// Call supplies the C-API-independent common layer. Calls not covered here
// retain rprimitive.Runtime's explicit panic, which executeTranslatedEntry
// converts into a normal evaluator error rather than a fake implementation.
func (b *translatedRuntimeBridge) Call(name string, values ...rprimitive.Value) rprimitive.Value {
	// Protection and C-runtime lifetime operations. Go GC owns the storage, but
	// PROTECT/REPROTECT still have an observable return value in GNU R.
	switch name {
	case "PROTECT", "REPROTECT":
		return dynamicArg(values, 0)
	case "PROTECT_WITH_INDEX":
		if ref, ok := dynamicArg(values, 1).(rprimitive.Ref); ok && ref != nil {
			ref.Set(int64(0))
		}
		return dynamicArg(values, 0)
	case "UNPROTECT", "R_PreserveObject", "R_ReleaseObject", "R_args_enable_refcnt", "R_try_clear_args_refcnt", "La_Init":
		return NullValue
	case "vmaxget":
		return NullValue
	case "vmaxset":
		return NullValue
	case "checkArity", "check1arg", "check1arg2", "CheckFormals", "checkNSname", "checkNotPromise", "R_checkConstants", "ENSURE_NAMEDMAX", "INCREMENT_NAMED":
		if name == "checkNotPromise" {
			return dynamicArg(values, 0)
		}
		return NullValue
	case "check_retval":
		if len(values) != 0 {
			return values[len(values)-1]
		}
		return NullValue
	case "MAYBE_REFERENCED", "MAYBE_SHARED":
		return dynamicMutableValue(runtimeValue(dynamicArg(values, 0)))
	case "NAMED":
		if dynamicMutableValue(runtimeValue(dynamicArg(values, 0))) {
			return int64(2)
		}
		return int64(0)
	case "RAISE_NAMED":
		return dynamicArg(values, 0)
	}

	// R errors are non-local exits. Returning an error sentinel lets translated
	// C code continue after error(), which is observably wrong, so use a typed
	// panic that executeTranslatedEntry converts back to an evaluator error.
	switch name {
	case "error", "Rf_error":
		panic(translatedRuntimeError{message: dynamicFormatMessage(values, 0)})
	case "errorcall", "errorcall_dflt":
		panic(translatedRuntimeError{message: dynamicFormatMessage(values, 1)})
	case "R_MissingArgError":
		panic(translatedRuntimeError{message: fmt.Sprintf("argument %q is missing, with no default", dynamicName(dynamicArg(values, 0)))})
	case "R_ObjectNotFoundError":
		panic(translatedRuntimeError{message: fmt.Sprintf("object %q not found", dynamicName(dynamicArg(values, 0)))})
	case "R_Suicide":
		panic(translatedRuntimeError{message: dynamicFormatMessage(values, 0)})
	case "warning":
		return NullValue
	case "warningcall":
		return NullValue
	}

	// Pairlists/language lists are represented by the host List value. The
	// Names slice carries pairlist tags, preserving the information needed by
	// TAG/SET_TAG while keeping CAR/CDR mutation in the same storage model.
	switch name {
	case "CAR", "CADR", "CADDR", "CADDDR", "CAD4R", "CAD5R":
		index := map[string]int{"CAR": 0, "CADR": 1, "CADDR": 2, "CADDDR": 3, "CAD4R": 4, "CAD5R": 5}[name]
		return dynamicListElement(runtimeValue(dynamicArg(values, 0)), index)
	case "CDR":
		return dynamicListTail(runtimeValue(dynamicArg(values, 0)), 1)
	case "CDDR":
		return dynamicListTail(runtimeValue(dynamicArg(values, 0)), 2)
	case "CDDDR":
		return dynamicListTail(runtimeValue(dynamicArg(values, 0)), 3)
	case "nthcdr":
		return dynamicListTail(runtimeValue(dynamicArg(values, 0)), dynamicIndex(dynamicArg(values, 1)))
	case "SETCAR":
		if list, ok := runtimeValue(dynamicArg(values, 0)).(*List); ok && len(list.Data) != 0 {
			list.Data[0] = runtimeValue(dynamicArg(values, 1))
		}
		return dynamicArg(values, 0)
	case "SETCDR":
		return dynamicSetListTail(runtimeValue(dynamicArg(values, 0)), runtimeValue(dynamicArg(values, 1)))
	case "CONS", "LCONS":
		return dynamicCons(runtimeValue(dynamicArg(values, 0)), runtimeValue(dynamicArg(values, 1)))
	case "TAG":
		if list, ok := runtimeValue(dynamicArg(values, 0)).(*List); ok && len(list.Names) != 0 && list.Names[0] != "" {
			return dynamicSymbol(list.Names[0])
		}
		return NullValue
	case "SET_TAG":
		if list, ok := runtimeValue(dynamicArg(values, 0)).(*List); ok {
			dynamicEnsureListNames(list)
			if len(list.Names) != 0 {
				if _, nilTag := dynamicArg(values, 1).(Null); nilTag || dynamicArg(values, 1) == nil {
					list.Names[0] = ""
				} else {
					list.Names[0] = dynamicName(dynamicArg(values, 1))
				}
			}
		}
		return dynamicArg(values, 0)
	case "allocList":
		n := maxDynamic(0, dynamicIndex(dynamicArg(values, 0)))
		return &List{Data: make([]Value, n), Names: make([]string, n)}
	case "allocLang":
		n := maxDynamic(0, dynamicIndex(dynamicArg(values, 0)))
		return &List{Data: make([]Value, n), Names: make([]string, n)}
	case "lang4":
		return &List{Data: []Value{runtimeValue(dynamicArg(values, 0)), runtimeValue(dynamicArg(values, 1)), runtimeValue(dynamicArg(values, 2)), runtimeValue(dynamicArg(values, 3))}}
	case "allocFormalsList2", "allocFormalsList3":
		return dynamicFormals(values)
	case "VectorToPairList":
		return dynamicVectorToList(runtimeValue(dynamicArg(values, 0)))
	case "listAppend":
		return dynamicAppendLists(runtimeValue(dynamicArg(values, 0)), runtimeValue(dynamicArg(values, 1)))
	case "ItemName":
		if list, ok := runtimeValue(dynamicArg(values, 0)).(*List); ok {
			i := dynamicIndex(dynamicArg(values, 1))
			if i >= 0 && i < len(list.Names) {
				return list.Names[i]
			}
		}
		return ""
	}

	// Typed vectors and coercion.
	switch name {
	case "length", "xlength", "XLENGTH", "LENGTH":
		return int64(Length(runtimeValue(dynamicArg(values, 0))))
	case "TYPEOF":
		return cTypeOf(runtimeValue(dynamicArg(values, 0)))
	case "type2char", "type2str":
		return dynamicTypeName(dynamicArg(values, 0))
	case "R_typeToChar":
		return dynamicTypeName(cTypeOf(runtimeValue(dynamicArg(values, 0))))
	case "isString":
		_, ok := runtimeValue(dynamicArg(values, 0)).(*CharacterVector)
		return ok
	case "isInteger":
		_, ok := runtimeValue(dynamicArg(values, 0)).(*IntegerVector)
		return ok
	case "isReal":
		_, ok := runtimeValue(dynamicArg(values, 0)).(*DoubleVector)
		return ok
	case "isLogical":
		_, ok := runtimeValue(dynamicArg(values, 0)).(*LogicalVector)
		return ok
	case "isVector":
		return dynamicIsVector(runtimeValue(dynamicArg(values, 0)))
	case "isList", "isNewList":
		_, ok := runtimeValue(dynamicArg(values, 0)).(*List)
		return ok
	case "isNull", "ISNULL":
		_, ok := runtimeValue(dynamicArg(values, 0)).(Null)
		return ok || dynamicArg(values, 0) == nil
	case "isSymbol":
		_, ok := dynamicArg(values, 0).(dynamicSymbol)
		return ok
	case "isLanguage":
		_, ok := runtimeValue(dynamicArg(values, 0)).(*Language)
		return ok
	case "isFunction":
		_, ok := runtimeValue(dynamicArg(values, 0)).(*Closure)
		return ok
	case "isPrimitive":
		t := cTypeOf(runtimeValue(dynamicArg(values, 0)))
		return t == "SPECIALSXP" || t == "BUILTINSXP"
	case "isExpression":
		return false
	case "isByteCode":
		return cTypeOf(runtimeValue(dynamicArg(values, 0))) == "BCODESXP"
	case "isValidString", "isValidStringF":
		vector, ok := runtimeValue(dynamicArg(values, 0)).(*CharacterVector)
		return ok && len(vector.Data) > 0 && !dynamicCharacterMissing(vector, 0)
	case "isObject":
		return dynamicAttribute(runtimeValue(dynamicArg(values, 0)), "class") != NullValue
	case "inherits":
		return dynamicInherits(runtimeValue(dynamicArg(values, 0)), dynamicName(dynamicArg(values, 1)))
	case "STRING_ELT":
		return dynamicStringElement(runtimeValue(dynamicArg(values, 0)), dynamicIndex(dynamicArg(values, 1)))
	case "VECTOR_ELT":
		return dynamicListElement(runtimeValue(dynamicArg(values, 0)), dynamicIndex(dynamicArg(values, 1)))
	case "SET_STRING_ELT":
		dynamicSetStringElement(runtimeValue(dynamicArg(values, 0)), dynamicIndex(dynamicArg(values, 1)), dynamicArg(values, 2))
		return dynamicArg(values, 0)
	case "SET_VECTOR_ELT":
		if v, ok := runtimeValue(dynamicArg(values, 0)).(*List); ok {
			i := dynamicIndex(dynamicArg(values, 1))
			if i >= 0 && i < len(v.Data) {
				v.Data[i] = runtimeValue(dynamicArg(values, 2))
			}
		}
		return dynamicArg(values, 0)
	case "allocVector":
		return allocateDynamicVector(dynamicName(dynamicArg(values, 0)), dynamicIndex(dynamicArg(values, 1)))
	case "ScalarInteger":
		return dynamicScalarInteger(dynamicArg(values, 0))
	case "ScalarReal":
		return dynamicScalarReal(dynamicArg(values, 0))
	case "ScalarLogical":
		return dynamicScalarLogical(dynamicArg(values, 0))
	case "ScalarString":
		return dynamicScalarString(dynamicArg(values, 0))
	case "INTEGER", "REAL", "LOGICAL", "RAW", "COMPLEX":
		return runtimeValue(dynamicArg(values, 0))
	case "asLogical", "asRbool", "asBool2", "asLogicalNA":
		return dynamicAsLogical(runtimeValue(dynamicArg(values, 0)))
	case "asInteger":
		return dynamicAsInteger(runtimeValue(dynamicArg(values, 0)))
	case "asReal":
		return dynamicAsReal(runtimeValue(dynamicArg(values, 0)))
	case "asXLength":
		return int64(dynamicIndex(dynamicArg(values, 0)))
	case "asChar":
		return dynamicAsChar(runtimeValue(dynamicArg(values, 0)))
	case "coerceVector":
		return dynamicCoerceVector(runtimeValue(dynamicArg(values, 0)), dynamicName(dynamicArg(values, 1)))
	case "IS_LONG_VEC":
		return int64(Length(runtimeValue(dynamicArg(values, 0)))) > int64(math.MaxInt32)
	}

	// Character/symbol and C-memory compatibility.
	switch name {
	case "mkChar", "translateChar", "translateCharFP", "EncodeChar", "CHAR":
		return dynamicString(dynamicArg(values, 0))
	case "mkString":
		return &CharacterVector{Data: []string{dynamicString(dynamicArg(values, 0))}}
	case "install", "installTrChar":
		return dynamicSymbol(dynamicString(dynamicArg(values, 0)))
	case "PRINTNAME":
		return dynamicName(dynamicArg(values, 0))
	case "streql":
		return dynamicString(dynamicArg(values, 0)) == dynamicString(dynamicArg(values, 1))
	case "strcpy", "strncpy", "strcat":
		return dynamicCStringCopy(name, values)
	case "Rstrdup":
		return string([]byte(dynamicString(dynamicArg(values, 0))))
	case "R_alloc":
		count := maxDynamic(0, dynamicIndex(dynamicArg(values, 0)))
		if len(values) > 1 {
			count *= maxDynamic(0, dynamicIndex(dynamicArg(values, 1)))
		}
		return make([]byte, count)
	case "calloc":
		count := maxDynamic(0, dynamicIndex(dynamicArg(values, 0))) * maxDynamic(0, dynamicIndex(dynamicArg(values, 1)))
		return make([]byte, count)
	case "malloc":
		return make([]byte, maxDynamic(0, dynamicIndex(dynamicArg(values, 0))))
	case "free":
		return NullValue
	case "memset":
		return dynamicMemset(dynamicArg(values, 0), dynamicArg(values, 1), dynamicArg(values, 2))
	case "memcpy":
		return dynamicMemcpy(dynamicArg(values, 0), dynamicArg(values, 1), dynamicArg(values, 2))
	}

	// Attributes and copy-on-write.
	switch name {
	case "getAttrib":
		return dynamicAttribute(runtimeValue(dynamicArg(values, 0)), dynamicAttributeKey(dynamicArg(values, 1)))
	case "setAttrib":
		setDynamicAttribute(runtimeValue(dynamicArg(values, 0)), dynamicAttributeKey(dynamicArg(values, 1)), runtimeValue(dynamicArg(values, 2)))
		return runtimeValue(dynamicArg(values, 0))
	case "classgets":
		setDynamicAttribute(runtimeValue(dynamicArg(values, 0)), "class", runtimeValue(dynamicArg(values, 1)))
		return runtimeValue(dynamicArg(values, 0))
	case "duplicate", "lazy_duplicate", "shallow_duplicate":
		return cloneDynamicValue(runtimeValue(dynamicArg(values, 0)))
	case "SHALLOW_DUPLICATE_ATTRIB":
		dynamicCopyAttributes(runtimeValue(dynamicArg(values, 0)), runtimeValue(dynamicArg(values, 1)))
		return dynamicArg(values, 0)
	}

	// Environment/frame/binding operations use the existing host Environment;
	// there is no second environment implementation in this bridge.
	switch name {
	case "isEnvironment":
		return dynamicEnvironment(dynamicArg(values, 0)) != nil
	case "ENCLOS", "parent.frame":
		if environment := dynamicEnvironment(dynamicArg(values, 0)); environment != nil && environment.Parent != nil {
			return environment.Parent
		}
		return NullValue
	case "SET_ENCLOS":
		if environment := dynamicEnvironment(dynamicArg(values, 0)); environment != nil {
			environment.Parent = dynamicEnvironment(dynamicArg(values, 1))
			return environment
		}
		return NullValue
	case "FRAME":
		if environment := dynamicEnvironment(dynamicArg(values, 0)); environment != nil {
			return environment
		}
		return NullValue
	case "HASHTAB":
		// The host Environment supplies lookup semantics itself. Returning NULL
		// selects GNU R's unhashed frame paths instead of maintaining a duplicate
		// hash table beside it.
		return NullValue
	case "SET_HASHTAB", "R_HashFrame":
		return NullValue
	case "IS_HASHED":
		return false
	case "NewEnvironment":
		return b.dynamicNewEnvironment(dynamicArg(values, 1), dynamicArg(values, 2))
	case "simple_as_environment":
		if environment := dynamicEnvironment(dynamicArg(values, 0)); environment != nil {
			return environment
		}
		return b.dynamicNewEnvironment(dynamicArg(values, 0), b.env)
	case "defineVar", "setVar", "gsetVar":
		return b.dynamicDefineVar(dynamicArg(values, 0), dynamicArg(values, 1), dynamicArg(values, len(values)-1))
	case "findVarInFrame", "R_findVarInFrame":
		return b.dynamicFindVarInFrame(dynamicArg(values, 0), dynamicArg(values, len(values)-1))
	case "findFun", "findFun3", "R_findVar":
		return b.dynamicFindVar(dynamicArg(values, 0), dynamicArg(values, len(values)-1))
	case "R_EnvironmentIsLocked":
		return b.environmentLocked[dynamicEnvironment(dynamicArg(values, 0))]
	case "R_LockEnvironment":
		if environment := dynamicEnvironment(dynamicArg(values, 0)); environment != nil {
			b.environmentLocked[environment] = true
			if dynamicTruth(dynamicArg(values, 1)) {
				for _, binding := range environment.Names(false) {
					b.dynamicSetBindingLock(environment, binding, true)
				}
			}
		}
		return NullValue
	case "R_LockBinding":
		if environment := dynamicEnvironment(dynamicArg(values, 1)); environment != nil {
			b.dynamicSetBindingLock(environment, dynamicName(dynamicArg(values, 0)), true)
		}
		return NullValue
	case "R_unLockBinding":
		if environment := dynamicEnvironment(dynamicArg(values, 1)); environment != nil {
			b.dynamicSetBindingLock(environment, dynamicName(dynamicArg(values, 0)), false)
		}
		return NullValue
	case "R_BindingIsLocked":
		return b.dynamicBindingLocked(dynamicEnvironment(dynamicArg(values, 1)), dynamicName(dynamicArg(values, 0)))
	case "R_BindingIsActive", "IS_ACTIVE_BINDING":
		return b.dynamicBindingActive(dynamicEnvironment(dynamicArg(values, 1)), dynamicName(dynamicArg(values, 0)))
	case "R_MakeActiveBinding":
		return b.dynamicMakeActiveBinding(dynamicArg(values, 0), dynamicArg(values, 1), dynamicArg(values, 2))
	case "R_ActiveBindingFunction":
		return b.dynamicActiveBindingFunction(dynamicArg(values, 0), dynamicArg(values, 1))
	case "R_GetBindingType":
		return b.dynamicBindingType(dynamicArg(values, 0), dynamicArg(values, 1))
	case "R_lsInternal3":
		if environment := dynamicEnvironment(dynamicArg(values, 0)); environment != nil {
			return &CharacterVector{Data: environment.Names(false)}
		}
		return &CharacterVector{}
	case "FrameNames":
		if environment := dynamicEnvironment(dynamicArg(values, 0)); environment != nil {
			return &CharacterVector{Data: environment.Names(false)}
		}
		return &CharacterVector{}
	case "FrameSize":
		if environment := dynamicEnvironment(dynamicArg(values, 0)); environment != nil {
			return int64(len(environment.Names(false)))
		}
		return int64(0)
	case "FrameValues":
		return b.dynamicFrameValues(dynamicArg(values, 0))
	case "resolveDotsEnv":
		if environment := dynamicEnvironment(dynamicArg(values, 0)); environment != nil {
			return environment
		}
		return b.env
	}

	// The evaluator entry path eagerly converts ordinary arguments to host
	// Values. eval() can therefore be identity for already-materialized data,
	// but executable language objects must still go through a real evaluator.
	switch name {
	case "eval", "Rf_eval_with_gd":
		value := runtimeValue(dynamicArg(values, 0))
		switch value.(type) {
		case *Language:
			panic(fmt.Sprintf("GNU R runtime operation %q requires host language evaluation", name))
		default:
			return value
		}
	case "evalListKeepMissing":
		return dynamicArg(values, 0)
	}

	// Condition-handler entries are represented explicitly and kept on the
	// translated runtime's existing R_HandlerStack symbol.
	switch name {
	case "mkHandlerEntry":
		return &dynamicHandlerEntry{
			Class:   dynamicString(dynamicArg(values, 0)),
			Parent:  dynamicArg(values, 1),
			Handler: dynamicArg(values, 2),
			Target:  dynamicArg(values, 3),
			Result:  dynamicArg(values, 4),
			Calling: dynamicTruth(dynamicArg(values, 5)),
		}
	case "ENTRY_HANDLER":
		if entry, ok := dynamicArg(values, 0).(*dynamicHandlerEntry); ok {
			return entry.Handler
		}
		return NullValue
	case "IS_CALLING_ENTRY":
		if entry, ok := dynamicArg(values, 0).(*dynamicHandlerEntry); ok {
			return entry.Calling
		}
		return false
	case "findConditionHandler":
		return b.dynamicFindConditionHandler(runtimeValue(dynamicArg(values, 0)))
	case "CHECK_RESTART":
		if _, ok := runtimeValue(dynamicArg(values, 0)).(*List); !ok {
			panic(translatedRuntimeError{message: "invalid restart object"})
		}
		return dynamicArg(values, 0)
	}

	// External pointers are pure-Go containers. Native routine resolution is a
	// separate, explicitly unsupported ABI boundary handled below.
	switch name {
	case "R_MakeExternalPtr":
		return &dynamicExternalPtr{Addr: dynamicArg(values, 0), Tag: dynamicArg(values, 1), Protected: dynamicArg(values, 2)}
	case "R_ExternalPtrAddr":
		if pointer, ok := dynamicArg(values, 0).(*dynamicExternalPtr); ok {
			return pointer.Addr
		}
		return NullValue
	case "R_ExternalPtrTag":
		if pointer, ok := dynamicArg(values, 0).(*dynamicExternalPtr); ok {
			return pointer.Tag
		}
		return NullValue
	case "R_RegisterCFinalizerEx", "R_RegisterFinalizerEx":
		if pointer, ok := dynamicArg(values, 0).(*dynamicExternalPtr); ok {
			pointer.Finalizer = dynamicArg(values, 1)
			pointer.OnExit = dynamicTruth(dynamicArg(values, 2))
		}
		return NullValue
	}

	// Small pure-Go utilities shared by connection/save-load entries.
	switch name {
	case "NextConnection":
		return int64(b.dynamicNextConnection())
	case "getConnection", "getConnectionCheck":
		index := dynamicIndex(dynamicArg(values, 0))
		if b.connections != nil {
			if connection, ok := b.connections.values[index]; ok {
				return connection
			}
		}
		panic(translatedRuntimeError{message: fmt.Sprintf("invalid connection %d", index)})
	case "R_ExpandFileName", "filenameToWchar":
		return dynamicExpandFileName(dynamicString(dynamicArg(values, 0)))
	case "RC_fopen":
		file, openErr := dynamicOpenFile(dynamicString(dynamicArg(values, 0)), dynamicString(dynamicArg(values, 1)))
		if openErr != nil {
			panic(translatedRuntimeError{message: openErr.Error()})
		}
		return file
	case "fclose":
		if file, ok := dynamicArg(values, 0).(*os.File); ok && file != nil {
			if closeErr := file.Close(); closeErr != nil {
				return int64(-1)
			}
		}
		return int64(0)
	case "Rprintf":
		fmt.Print(dynamicFormatMessage(values, 0))
		return NullValue
	case "Rconn_printf":
		message := dynamicFormatMessage(values, 1)
		if writer, ok := dynamicArg(values, 0).(interface{ Write([]byte) (int, error) }); ok {
			_, _ = writer.Write([]byte(message))
		}
		return NullValue
	case "defaultSaveVersion", "defaultSerializeVersion":
		return int64(3)
	case "R_compute_identical":
		return reflect.DeepEqual(runtimeValue(dynamicArg(values, 0)), runtimeValue(dynamicArg(values, 1)))
	case "SET_RDEBUG":
		b.debugFlags[dynamicIdentity(dynamicArg(values, 0))] = dynamicTruth(dynamicArg(values, 1))
		return NullValue
	case "SET_RSTEP":
		b.stepFlags[dynamicIdentity(dynamicArg(values, 0))] = dynamicTruth(dynamicArg(values, 1))
		return NullValue
	case "RDEBUG":
		return b.debugFlags[dynamicIdentity(dynamicArg(values, 0))]
	case "HASHASH":
		return false
	case "HASHVALUE", "R_Newhashpjw":
		return int64(dynamicHashPJW(dynamicString(dynamicArg(values, 0))))
	case "BuiltinSize":
		return int64(len(rprimitive.PrimitiveTable))
	case "BuiltinNames":
		return &CharacterVector{Data: dynamicBuiltinNames()}
	case "BuiltinValues":
		names := dynamicBuiltinNames()
		list := &List{Data: make([]Value, len(names)), Names: append([]string(nil), names...)}
		for i, primitive := range names {
			list.Data[i] = dynamicSymbol(primitive)
		}
		return list
	case "PRIMNAME":
		if primitive, ok := dynamicArg(values, 0).(rprimitive.PrimitiveOp); ok {
			return primitive.Primitive
		}
		return ""
	}

	// Native ABI calls are never reported as successful by the pure-Go bridge.
	switch name {
	case "resolveNativeRoutine", "R_FindSymbol", "AddDLL", "DeleteDLL", "GetFullDLLPath", "R_getDllTable", "R_getRegisteredRoutines", "R_getSymbolInfo", "Rf_MakeDLLInfo", "R_doDotCall", "fun":
		panic(rprimitive.PluginBoundaryError{Operation: name})
	}

	return b.Runtime.Call(name, values...)
}

// DynamicRuntime's buffer macros are modelled as typed runtime vectors.  The
// translated C loops use Index/AssignIndex against these buffers, so keeping
// mutation here preserves one shared storage model instead of C-like slices.
type dynamicVectorRef struct {
	value Value
	index int
}

func (r dynamicVectorRef) Get() rprimitive.Value      { return dynamicVectorIndex(r.value, r.index) }
func (r dynamicVectorRef) Set(value rprimitive.Value) { dynamicVectorSet(r.value, r.index, value) }
func (b *translatedRuntimeBridge) Index(value, index rprimitive.Value) rprimitive.Value {
	if table, ok := value.(*dynamicConnectionTable); ok {
		return dynamicConnectionRef{table: table, index: dynamicIndex(index)}.Get()
	}
	if runtime := runtimeValue(value); runtime != NullValue {
		return dynamicVectorIndex(runtime, dynamicIndex(index))
	}
	return b.Runtime.Index(value, index)
}
func (b *translatedRuntimeBridge) IndexRef(value, index rprimitive.Value) rprimitive.Ref {
	if table, ok := value.(*dynamicConnectionTable); ok {
		return dynamicConnectionRef{table: table, index: dynamicIndex(index)}
	}
	if runtime := runtimeValue(value); runtime != NullValue {
		return dynamicVectorRef{runtime, dynamicIndex(index)}
	}
	return b.Runtime.IndexRef(value, index)
}
func (b *translatedRuntimeBridge) AssignIndex(target, index, value rprimitive.Value) rprimitive.Value {
	b.IndexRef(target, index).Set(value)
	return value
}
func dynamicVectorIndex(value Value, index int) rprimitive.Value {
	if index < 0 {
		return NullValue
	}
	switch v := value.(type) {
	case *RawVector:
		if index < len(v.Data) {
			return int64(v.Data[index])
		}
	case *LogicalVector:
		if index < len(v.Data) {
			return v.Data[index]
		}
	case *IntegerVector:
		if index < len(v.Data) {
			if dynamicIntegerMissing(v, index) {
				return int64(math.MinInt32)
			}
			return v.Data[index]
		}
	case *DoubleVector:
		if index < len(v.Data) {
			if dynamicDoubleMissing(v, index) {
				return dynamicNARealValue
			}
			return v.Data[index]
		}
	case *ComplexVector:
		if index < len(v.Data) {
			return v.Data[index]
		}
	case *CharacterVector:
		if index < len(v.Data) {
			if dynamicCharacterMissing(v, index) {
				return dynamicNAStringValue
			}
			return v.Data[index]
		}
	case *List:
		if index < len(v.Data) {
			return v.Data[index]
		}
	}
	return NullValue
}
func dynamicVectorSet(value Value, index int, item rprimitive.Value) {
	if index < 0 {
		return
	}
	switch v := value.(type) {
	case *RawVector:
		if index < len(v.Data) {
			v.Data[index] = byte(dynamicIndex(item))
		}
	case *LogicalVector:
		if index < len(v.Data) {
			if x, ok := item.(Logical); ok {
				v.Data[index] = x
			} else {
				v.Data[index] = Logical(dynamicIndex(item))
			}
		}
	case *IntegerVector:
		if index < len(v.Data) {
			dynamicEnsureMissing(&v.Missing, len(v.Data))
			if x, ok := item.(int64); ok && x == int64(math.MinInt32) {
				v.Data[index] = 0
				v.Missing[index] = true
			} else {
				v.Data[index] = int64(dynamicIndex(item))
				v.Missing[index] = false
			}
		}
	case *DoubleVector:
		if index < len(v.Data) {
			dynamicEnsureMissing(&v.Missing, len(v.Data))
			if _, missing := item.(*dynamicNAReal); missing {
				v.Data[index] = 0
				v.Missing[index] = true
			} else if x, ok := item.(float64); ok {
				v.Data[index] = x
				v.Missing[index] = false
			} else {
				v.Data[index] = float64(dynamicIndex(item))
				v.Missing[index] = false
			}
		}
	case *CharacterVector:
		dynamicSetStringElement(v, index, item)
	case *List:
		if index < len(v.Data) {
			v.Data[index] = runtimeValue(item)
		}
	}
}

func dynamicIndex(value rprimitive.Value) int {
	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		if math.IsNaN(v) {
			return 0
		}
		return int(v)
	case Logical:
		return int(v)
	case *IntegerVector:
		if len(v.Data) != 0 && !dynamicIntegerMissing(v, 0) {
			return int(v.Data[0])
		}
	case *DoubleVector:
		if len(v.Data) != 0 && !dynamicDoubleMissing(v, 0) && !math.IsNaN(v.Data[0]) {
			return int(v.Data[0])
		}
	case *LogicalVector:
		if len(v.Data) != 0 {
			return int(v.Data[0])
		}
	}
	return 0
}
func maxDynamic(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func dynamicArg(values []rprimitive.Value, index int) rprimitive.Value {
	if index < 0 || index >= len(values) {
		return nil
	}
	return values[index]
}

func dynamicMutableValue(value Value) bool {
	switch value.(type) {
	case *RawVector, *LogicalVector, *IntegerVector, *DoubleVector, *ComplexVector, *CharacterVector, *List, *Environment, *EnvironmentValue:
		return true
	default:
		return false
	}
}

func dynamicFormatMessage(values []rprimitive.Value, formatIndex int) string {
	if formatIndex < 0 || formatIndex >= len(values) {
		return "GNU R runtime error"
	}
	format := dynamicString(values[formatIndex])
	if formatIndex+1 >= len(values) {
		return format
	}
	args := make([]any, 0, len(values)-formatIndex-1)
	for _, value := range values[formatIndex+1:] {
		args = append(args, dynamicPrintable(value))
	}
	return fmt.Sprintf(format, args...)
}

func dynamicPrintable(value rprimitive.Value) any {
	switch value := value.(type) {
	case dynamicSymbol:
		return string(value)
	case *dynamicNAString:
		return "NA"
	case *dynamicNAReal:
		return "NA"
	case *dynamicUnbound:
		return "<unbound>"
	case *dynamicMissingArg:
		return "<missing>"
	default:
		return value
	}
}

func dynamicName(value rprimitive.Value) string {
	switch value := value.(type) {
	case dynamicSymbol:
		return string(value)
	case string:
		return value
	case *dynamicNAString:
		return "NA"
	case nil:
		return ""
	default:
		return fmt.Sprint(value)
	}
}

func dynamicString(value rprimitive.Value) string {
	switch value := value.(type) {
	case string:
		return value
	case dynamicSymbol:
		return string(value)
	case []byte:
		if i := strings.IndexByte(string(value), 0); i >= 0 {
			return string(value[:i])
		}
		return string(value)
	case *dynamicNAString, *dynamicNAReal:
		return "NA"
	case nil:
		return ""
	default:
		if host := runtimeValue(value); host != NullValue {
			switch x := host.(type) {
			case *CharacterVector:
				if len(x.Data) != 0 {
					if dynamicCharacterMissing(x, 0) {
						return "NA"
					}
					return x.Data[0]
				}
			case *IntegerVector:
				if len(x.Data) != 0 {
					if dynamicIntegerMissing(x, 0) {
						return "NA"
					}
					return strconv.FormatInt(x.Data[0], 10)
				}
			case *DoubleVector:
				if len(x.Data) != 0 {
					if dynamicDoubleMissing(x, 0) {
						return "NA"
					}
					return strconv.FormatFloat(x.Data[0], 'g', -1, 64)
				}
			}
		}
		return fmt.Sprint(value)
	}
}

func dynamicTruth(value rprimitive.Value) bool {
	switch value := value.(type) {
	case bool:
		return value
	case int:
		return value != 0
	case int64:
		return value != 0 && value != int64(math.MinInt32)
	case float64:
		return value != 0 && !math.IsNaN(value)
	case Logical:
		return value == True
	case *LogicalVector:
		return len(value.Data) != 0 && value.Data[0] == True
	case *IntegerVector:
		return len(value.Data) != 0 && !dynamicIntegerMissing(value, 0) && value.Data[0] != 0
	case *DoubleVector:
		return len(value.Data) != 0 && !dynamicDoubleMissing(value, 0) && value.Data[0] != 0 && !math.IsNaN(value.Data[0])
	case *dynamicNAReal, *dynamicNAString, *dynamicMissingArg:
		return false
	default:
		return value != nil
	}
}

func dynamicListElement(value Value, index int) rprimitive.Value {
	if index < 0 {
		return NullValue
	}
	if list, ok := value.(*List); ok && index < len(list.Data) {
		return list.Data[index]
	}
	return NullValue
}

func dynamicListTail(value Value, skip int) rprimitive.Value {
	list, ok := value.(*List)
	if !ok || skip < 0 {
		return NullValue
	}
	if skip == 0 {
		return list
	}
	if skip >= len(list.Data) {
		return NullValue
	}
	out := &List{Data: list.Data[skip:]}
	if skip < len(list.Names) {
		out.Names = list.Names[skip:]
	}
	if attrs := attributesOf(list); attrs != nil {
		applyAttributes(out, attrs)
	}
	return out
}

func dynamicSetListTail(head, tail Value) rprimitive.Value {
	list, ok := head.(*List)
	if !ok || len(list.Data) == 0 {
		return head
	}
	if _, nilTail := tail.(Null); nilTail || tail == nil {
		list.Data = list.Data[:1]
		if len(list.Names) > 1 {
			list.Names = list.Names[:1]
		}
		return list
	}
	other, ok := tail.(*List)
	if !ok {
		list.Data = append(list.Data[:1], tail)
		dynamicEnsureListNames(list)
		return list
	}
	first := list.Data[0]
	firstName := ""
	if len(list.Names) != 0 {
		firstName = list.Names[0]
	}
	list.Data = append([]Value{first}, other.Data...)
	list.Names = append([]string{firstName}, other.Names...)
	dynamicEnsureListNames(list)
	return list
}

func dynamicCons(head, tail Value) *List {
	list := &List{Data: []Value{head}, Names: []string{""}}
	if other, ok := tail.(*List); ok {
		list.Data = append(list.Data, other.Data...)
		list.Names = append(list.Names, other.Names...)
	}
	dynamicEnsureListNames(list)
	return list
}

func dynamicEnsureListNames(list *List) {
	if list == nil {
		return
	}
	if len(list.Names) < len(list.Data) {
		list.Names = append(list.Names, make([]string, len(list.Data)-len(list.Names))...)
	} else if len(list.Names) > len(list.Data) {
		list.Names = list.Names[:len(list.Data)]
	}
}

func dynamicFormals(values []rprimitive.Value) *List {
	result := &List{Data: make([]Value, len(values)), Names: make([]string, len(values))}
	for i, value := range values {
		result.Data[i] = runtimeValue(dynamicMissingValue)
		result.Names[i] = dynamicName(value)
	}
	return result
}

func dynamicVectorToList(value Value) *List {
	if list, ok := value.(*List); ok {
		return cloneDynamicValue(list).(*List)
	}
	n := Length(value)
	result := &List{Data: make([]Value, n), Names: make([]string, n)}
	for i := 0; i < n; i++ {
		result.Data[i] = runtimeValue(dynamicVectorIndex(value, i))
	}
	if attrs := attributesOf(value); attrs != nil {
		if names, ok := attrs["names"].(*CharacterVector); ok {
			copy(result.Names, names.Data)
		}
	}
	return result
}

func dynamicAppendLists(left, right Value) Value {
	if _, nilLeft := left.(Null); nilLeft || left == nil {
		return right
	}
	if _, nilRight := right.(Null); nilRight || right == nil {
		return left
	}
	a, ok := left.(*List)
	if !ok {
		return left
	}
	b, ok := right.(*List)
	if !ok {
		a.Data = append(a.Data, right)
		dynamicEnsureListNames(a)
		return a
	}
	a.Data = append(a.Data, b.Data...)
	a.Names = append(a.Names, b.Names...)
	dynamicEnsureListNames(a)
	return a
}

func dynamicTypeName(value rprimitive.Value) string {
	name := dynamicName(value)
	switch name {
	case "NILSXP":
		return "NULL"
	case "SYMSXP":
		return "symbol"
	case "LISTSXP":
		return "pairlist"
	case "CLOSXP":
		return "closure"
	case "ENVSXP":
		return "environment"
	case "PROMSXP":
		return "promise"
	case "LANGSXP":
		return "language"
	case "SPECIALSXP":
		return "special"
	case "BUILTINSXP":
		return "builtin"
	case "LGLSXP":
		return "logical"
	case "INTSXP":
		return "integer"
	case "REALSXP":
		return "double"
	case "CPLXSXP":
		return "complex"
	case "STRSXP":
		return "character"
	case "VECSXP":
		return "list"
	case "EXPRSXP":
		return "expression"
	case "BCODESXP":
		return "bytecode"
	case "EXTPTRSXP":
		return "externalptr"
	case "RAWSXP":
		return "raw"
	default:
		return name
	}
}

func dynamicIsVector(value Value) bool {
	switch value.(type) {
	case *RawVector, *LogicalVector, *IntegerVector, *DoubleVector, *ComplexVector, *CharacterVector, *List:
		return true
	default:
		return false
	}
}

func dynamicCharacterMissing(vector *CharacterVector, index int) bool {
	return vector != nil && index >= 0 && index < len(vector.Missing) && vector.Missing[index]
}

func dynamicIntegerMissing(vector *IntegerVector, index int) bool {
	return vector != nil && index >= 0 && index < len(vector.Missing) && vector.Missing[index]
}

func dynamicDoubleMissing(vector *DoubleVector, index int) bool {
	return vector != nil && index >= 0 && index < len(vector.Missing) && vector.Missing[index]
}

func dynamicEnsureMissing(missing *[]bool, n int) {
	if n <= len(*missing) {
		return
	}
	*missing = append(*missing, make([]bool, n-len(*missing))...)
}

func dynamicStringElement(value Value, index int) rprimitive.Value {
	vector, ok := value.(*CharacterVector)
	if !ok || index < 0 || index >= len(vector.Data) {
		return dynamicNAStringValue
	}
	if dynamicCharacterMissing(vector, index) {
		return dynamicNAStringValue
	}
	return vector.Data[index]
}

func dynamicSetStringElement(value Value, index int, item rprimitive.Value) {
	vector, ok := value.(*CharacterVector)
	if !ok || index < 0 || index >= len(vector.Data) {
		return
	}
	dynamicEnsureMissing(&vector.Missing, len(vector.Data))
	if _, missing := item.(*dynamicNAString); missing {
		vector.Data[index] = ""
		vector.Missing[index] = true
		return
	}
	vector.Data[index] = dynamicString(item)
	vector.Missing[index] = false
}

func dynamicScalarInteger(value rprimitive.Value) *IntegerVector {
	result := &IntegerVector{Data: []int64{0}, Missing: []bool{false}}
	switch value := value.(type) {
	case int64:
		if value == int64(math.MinInt32) {
			result.Missing[0] = true
		} else {
			result.Data[0] = value
		}
	case int:
		result.Data[0] = int64(value)
	case float64:
		if math.IsNaN(value) {
			result.Missing[0] = true
		} else {
			result.Data[0] = int64(value)
		}
	case *dynamicNAReal, *dynamicMissingArg:
		result.Missing[0] = true
	default:
		result.Data[0] = int64(dynamicIndex(value))
	}
	return result
}

func dynamicScalarReal(value rprimitive.Value) *DoubleVector {
	result := &DoubleVector{Data: []float64{0}, Missing: []bool{false}}
	switch value := value.(type) {
	case float64:
		result.Data[0] = value
	case int64:
		if value == int64(math.MinInt32) {
			result.Missing[0] = true
		} else {
			result.Data[0] = float64(value)
		}
	case int:
		result.Data[0] = float64(value)
	case *dynamicNAReal, *dynamicMissingArg:
		result.Missing[0] = true
	default:
		result.Data[0] = float64(dynamicIndex(value))
	}
	return result
}

func dynamicScalarLogical(value rprimitive.Value) *LogicalVector {
	logical := NA
	switch value := value.(type) {
	case Logical:
		logical = value
	case bool:
		if value {
			logical = True
		} else {
			logical = False
		}
	case int64:
		if value != int64(math.MinInt32) {
			if value == 0 {
				logical = False
			} else {
				logical = True
			}
		}
	case int:
		if value == 0 {
			logical = False
		} else {
			logical = True
		}
	}
	return &LogicalVector{Data: []Logical{logical}}
}

func dynamicScalarString(value rprimitive.Value) *CharacterVector {
	result := &CharacterVector{Data: []string{""}, Missing: []bool{false}}
	if _, missing := value.(*dynamicNAString); missing {
		result.Missing[0] = true
		return result
	}
	result.Data[0] = dynamicString(value)
	return result
}

func dynamicAsLogical(value Value) rprimitive.Value {
	switch value := value.(type) {
	case *LogicalVector:
		if len(value.Data) == 0 {
			return NA
		}
		return value.Data[0]
	case *IntegerVector:
		if len(value.Data) == 0 || dynamicIntegerMissing(value, 0) {
			return NA
		}
		if value.Data[0] == 0 {
			return False
		}
		return True
	case *DoubleVector:
		if len(value.Data) == 0 || dynamicDoubleMissing(value, 0) || math.IsNaN(value.Data[0]) {
			return NA
		}
		if value.Data[0] == 0 {
			return False
		}
		return True
	case *CharacterVector:
		if len(value.Data) == 0 || dynamicCharacterMissing(value, 0) {
			return NA
		}
		parsed, err := strconv.ParseBool(strings.TrimSpace(value.Data[0]))
		if err != nil {
			return NA
		}
		if parsed {
			return True
		}
		return False
	case Null:
		return NA
	case Logical:
		return value
	}
	return NA
}

func dynamicAsInteger(value Value) rprimitive.Value {
	switch value := value.(type) {
	case *IntegerVector:
		if len(value.Data) == 0 || dynamicIntegerMissing(value, 0) {
			return int64(math.MinInt32)
		}
		return value.Data[0]
	case *LogicalVector:
		if len(value.Data) == 0 || value.Data[0] == NA {
			return int64(math.MinInt32)
		}
		return int64(value.Data[0])
	case *DoubleVector:
		if len(value.Data) == 0 || dynamicDoubleMissing(value, 0) || math.IsNaN(value.Data[0]) {
			return int64(math.MinInt32)
		}
		return int64(value.Data[0])
	case *CharacterVector:
		if len(value.Data) == 0 || dynamicCharacterMissing(value, 0) {
			return int64(math.MinInt32)
		}
		if parsed, err := strconv.ParseInt(strings.TrimSpace(value.Data[0]), 10, 64); err == nil {
			return parsed
		}
		return int64(math.MinInt32)
	case Null:
		return int64(math.MinInt32)
	}
	return int64(dynamicIndex(value))
}

func dynamicAsReal(value Value) rprimitive.Value {
	switch value := value.(type) {
	case *DoubleVector:
		if len(value.Data) == 0 || dynamicDoubleMissing(value, 0) {
			return dynamicNARealValue
		}
		return value.Data[0]
	case *IntegerVector:
		if len(value.Data) == 0 || dynamicIntegerMissing(value, 0) {
			return dynamicNARealValue
		}
		return float64(value.Data[0])
	case *LogicalVector:
		if len(value.Data) == 0 || value.Data[0] == NA {
			return dynamicNARealValue
		}
		return float64(value.Data[0])
	case *CharacterVector:
		if len(value.Data) == 0 || dynamicCharacterMissing(value, 0) {
			return dynamicNARealValue
		}
		if parsed, err := strconv.ParseFloat(strings.TrimSpace(value.Data[0]), 64); err == nil {
			return parsed
		}
		return dynamicNARealValue
	case Null:
		return dynamicNARealValue
	}
	return float64(dynamicIndex(value))
}

func dynamicAsChar(value Value) rprimitive.Value {
	if vector, ok := value.(*CharacterVector); ok {
		if len(vector.Data) == 0 || dynamicCharacterMissing(vector, 0) {
			return dynamicNAStringValue
		}
		return vector.Data[0]
	}
	return dynamicString(value)
}

func dynamicCoerceVector(value Value, target string) Value {
	target = strings.ToUpper(strings.TrimSpace(target))
	if target == cTypeOf(value) {
		return cloneDynamicValue(value)
	}
	n := Length(value)
	var result Value
	switch target {
	case "STRSXP":
		out := &CharacterVector{Data: make([]string, n), Missing: make([]bool, n)}
		for i := 0; i < n; i++ {
			item := dynamicVectorIndex(value, i)
			switch item.(type) {
			case *dynamicNAString, *dynamicNAReal, *dynamicMissingArg:
				out.Missing[i] = true
			default:
				if x, ok := item.(int64); ok && x == int64(math.MinInt32) {
					out.Missing[i] = true
				} else {
					out.Data[i] = dynamicString(item)
				}
			}
		}
		result = out
	case "INTSXP":
		out := &IntegerVector{Data: make([]int64, n), Missing: make([]bool, n)}
		for i := 0; i < n; i++ {
			item := dynamicScalarInteger(dynamicVectorIndex(value, i))
			out.Data[i] = item.Data[0]
			out.Missing[i] = item.Missing[0]
		}
		result = out
	case "REALSXP":
		out := &DoubleVector{Data: make([]float64, n), Missing: make([]bool, n)}
		for i := 0; i < n; i++ {
			item := dynamicScalarReal(dynamicVectorIndex(value, i))
			out.Data[i] = item.Data[0]
			out.Missing[i] = item.Missing[0]
		}
		result = out
	case "LGLSXP":
		out := &LogicalVector{Data: make([]Logical, n)}
		for i := 0; i < n; i++ {
			item := dynamicAsLogical(runtimeValue(dynamicVectorIndex(value, i)))
			if logical, ok := item.(Logical); ok {
				out.Data[i] = logical
			} else {
				out.Data[i] = NA
			}
		}
		result = out
	case "RAWSXP":
		out := &RawVector{Data: make([]byte, n)}
		for i := 0; i < n; i++ {
			v := dynamicIndex(dynamicVectorIndex(value, i))
			if v < 0 {
				v = 0
			}
			if v > 255 {
				v = 255
			}
			out.Data[i] = byte(v)
		}
		result = out
	case "VECSXP", "LISTSXP":
		result = dynamicVectorToList(value)
	default:
		return value
	}
	if attrs := attributesOf(value); attrs != nil {
		applyAttributes(result, cloneDynamicAttributes(attrs))
	}
	return result
}

func dynamicInherits(value Value, class string) bool {
	classes, ok := dynamicAttribute(value, "class").(*CharacterVector)
	if !ok {
		return false
	}
	for i, current := range classes.Data {
		if !dynamicCharacterMissing(classes, i) && current == class {
			return true
		}
	}
	return false
}

func dynamicCStringCopy(name string, values []rprimitive.Value) rprimitive.Value {
	target := dynamicArg(values, 0)
	source := dynamicString(dynamicArg(values, 1))
	if name == "strcat" {
		source = dynamicString(target) + source
	}
	limit := len(source)
	if name == "strncpy" && len(values) > 2 {
		limit = dynamicIndex(values[2])
		if limit < 0 {
			limit = 0
		}
		if limit > len(source) {
			limit = len(source)
		}
	}
	if bytes, ok := target.([]byte); ok {
		copyLen := limit
		if copyLen > len(bytes) {
			copyLen = len(bytes)
		}
		copy(bytes[:copyLen], source[:copyLen])
		if copyLen < len(bytes) {
			bytes[copyLen] = 0
		}
		return bytes
	}
	return source[:limit]
}

func dynamicMemset(target, fill, count rprimitive.Value) rprimitive.Value {
	bytes, ok := target.([]byte)
	if !ok {
		return target
	}
	n := dynamicIndex(count)
	if n < 0 {
		n = 0
	}
	if n > len(bytes) {
		n = len(bytes)
	}
	for i := 0; i < n; i++ {
		bytes[i] = byte(dynamicIndex(fill))
	}
	return bytes
}

func dynamicMemcpy(target, source, count rprimitive.Value) rprimitive.Value {
	dst, ok := target.([]byte)
	if !ok {
		return target
	}
	src, ok := source.([]byte)
	if !ok {
		return target
	}
	n := dynamicIndex(count)
	if n < 0 {
		n = 0
	}
	if n > len(dst) {
		n = len(dst)
	}
	if n > len(src) {
		n = len(src)
	}
	copy(dst[:n], src[:n])
	return dst
}

func cloneDynamicAttributes(attrs map[string]Value) map[string]Value {
	if attrs == nil {
		return nil
	}
	copyAttrs := make(map[string]Value, len(attrs))
	for key, value := range attrs {
		copyAttrs[key] = cloneDynamicValue(value)
	}
	return copyAttrs
}

func dynamicCopyAttributes(target, source Value) {
	attrs := attributesOf(source)
	if attrs == nil {
		return
	}
	applyAttributes(target, cloneDynamicAttributes(attrs))
}

func (b *translatedRuntimeBridge) rootEnvironment() *Environment {
	environment := b.env
	if environment == nil {
		return nil
	}
	seen := map[*Environment]bool{}
	for environment.Parent != nil && !seen[environment] {
		seen[environment] = true
		environment = environment.Parent
	}
	return environment
}

func (b *translatedRuntimeBridge) dynamicNewEnvironment(frame, parentValue rprimitive.Value) *Environment {
	parent := dynamicEnvironment(parentValue)
	if parent == nil {
		parent = b.env
	}
	environment := &Environment{Parent: parent}
	if list, ok := runtimeValue(frame).(*List); ok {
		for i, value := range list.Data {
			name := ""
			if i < len(list.Names) {
				name = list.Names[i]
			}
			if name != "" {
				environment.Set(name, value)
			}
		}
	}
	return environment
}

func (b *translatedRuntimeBridge) dynamicDefineVar(symbol, value, environmentValue rprimitive.Value) rprimitive.Value {
	environment := dynamicEnvironment(environmentValue)
	if environment == nil {
		environment = b.env
	}
	if environment == nil {
		panic(translatedRuntimeError{message: "cannot define a variable without an environment"})
	}
	name := dynamicName(symbol)
	if b.dynamicBindingLocked(environment, name) {
		panic(translatedRuntimeError{message: fmt.Sprintf("cannot change value of locked binding for %q", name)})
	}
	if b.environmentLocked[environment] && !dynamicEnvironmentHasName(environment, name) {
		panic(translatedRuntimeError{message: fmt.Sprintf("cannot add bindings to a locked environment: %q", name)})
	}
	if b.dynamicBindingActive(environment, name) {
		panic(translatedRuntimeError{message: fmt.Sprintf("active binding %q requires host callback evaluation", name)})
	}
	hostValue := runtimeValue(value)
	environment.Set(name, hostValue)
	return hostValue
}

func (b *translatedRuntimeBridge) dynamicFindVar(symbol, environmentValue rprimitive.Value) rprimitive.Value {
	environment := dynamicEnvironment(environmentValue)
	if environment == nil {
		environment = b.env
	}
	if environment == nil {
		return dynamicUnboundValue
	}
	name := dynamicName(symbol)
	if b.dynamicBindingActive(environment, name) {
		panic(translatedRuntimeError{message: fmt.Sprintf("active binding %q requires host callback evaluation", name)})
	}
	value, err := environment.Get(b.ctx, name)
	if err != nil {
		return dynamicUnboundValue
	}
	return value
}

func (b *translatedRuntimeBridge) dynamicFindVarInFrame(symbol, environmentValue rprimitive.Value) rprimitive.Value {
	environment := dynamicEnvironment(environmentValue)
	if environment == nil {
		return dynamicUnboundValue
	}
	name := dynamicName(symbol)
	if !dynamicEnvironmentHasName(environment, name) {
		return dynamicUnboundValue
	}
	if b.dynamicBindingActive(environment, name) {
		panic(translatedRuntimeError{message: fmt.Sprintf("active binding %q requires host callback evaluation", name)})
	}
	value, err := environment.Get(b.ctx, name)
	if err != nil {
		return dynamicUnboundValue
	}
	return value
}

func dynamicEnvironmentHasName(environment *Environment, name string) bool {
	if environment == nil {
		return false
	}
	for _, current := range environment.Names(false) {
		if current == name {
			return true
		}
	}
	return false
}

func (b *translatedRuntimeBridge) dynamicSetBindingLock(environment *Environment, name string, locked bool) {
	if environment == nil {
		return
	}
	locks := b.bindingLocked[environment]
	if locks == nil {
		locks = map[string]bool{}
		b.bindingLocked[environment] = locks
	}
	if locked {
		locks[name] = true
	} else {
		delete(locks, name)
	}
}

func (b *translatedRuntimeBridge) dynamicBindingLocked(environment *Environment, name string) bool {
	return environment != nil && b.bindingLocked[environment] != nil && b.bindingLocked[environment][name]
}

func (b *translatedRuntimeBridge) dynamicBindingActive(environment *Environment, name string) bool {
	return environment != nil && b.activeBindings[environment] != nil && b.activeBindings[environment][name] != nil
}

func (b *translatedRuntimeBridge) dynamicMakeActiveBinding(symbol, function, environmentValue rprimitive.Value) rprimitive.Value {
	environment := dynamicEnvironment(environmentValue)
	if environment == nil {
		panic(translatedRuntimeError{message: "active binding requires an environment"})
	}
	name := dynamicName(symbol)
	if b.dynamicBindingLocked(environment, name) {
		panic(translatedRuntimeError{message: fmt.Sprintf("cannot change active binding for locked name %q", name)})
	}
	bindings := b.activeBindings[environment]
	if bindings == nil {
		bindings = map[string]rprimitive.Value{}
		b.activeBindings[environment] = bindings
	}
	bindings[name] = function
	return NullValue
}

func (b *translatedRuntimeBridge) dynamicActiveBindingFunction(symbol, environmentValue rprimitive.Value) rprimitive.Value {
	environment := dynamicEnvironment(environmentValue)
	if environment == nil || b.activeBindings[environment] == nil {
		return NullValue
	}
	if function := b.activeBindings[environment][dynamicName(symbol)]; function != nil {
		return function
	}
	return NullValue
}

func (b *translatedRuntimeBridge) dynamicBindingType(symbol, environmentValue rprimitive.Value) rprimitive.Value {
	environment := dynamicEnvironment(environmentValue)
	name := dynamicName(symbol)
	if environment == nil {
		return "R_BindingTypeUnbound"
	}
	if b.dynamicBindingActive(environment, name) {
		return "R_BindingTypeActive"
	}
	if dynamicEnvironmentHasName(environment, name) {
		return "R_BindingTypeValue"
	}
	return "R_BindingTypeUnbound"
}

func (b *translatedRuntimeBridge) dynamicFrameValues(environmentValue rprimitive.Value) Value {
	environment := dynamicEnvironment(environmentValue)
	if environment == nil {
		return &List{}
	}
	names := environment.Names(false)
	values := &List{Data: make([]Value, 0, len(names)), Names: append([]string(nil), names...)}
	for _, name := range names {
		value, err := environment.Get(b.ctx, name)
		if err != nil {
			values.Data = append(values.Data, NullValue)
		} else {
			values.Data = append(values.Data, value)
		}
	}
	return values
}

func (b *translatedRuntimeBridge) dynamicFindConditionHandler(condition Value) rprimitive.Value {
	stack, ok := runtimeValue(b.Runtime.Symbol("R_HandlerStack")).(*List)
	if !ok || len(stack.Data) == 0 {
		return NullValue
	}
	classes := []string{"condition"}
	if vector, ok := dynamicAttribute(condition, "class").(*CharacterVector); ok && len(vector.Data) != 0 {
		classes = append([]string(nil), vector.Data...)
	}
	for i, item := range stack.Data {
		entry, ok := item.(*dynamicHandlerEntry)
		if !ok {
			continue
		}
		for _, class := range classes {
			if entry.Class == class || entry.Class == "condition" {
				return dynamicListTail(stack, i)
			}
		}
	}
	return NullValue
}

func (b *translatedRuntimeBridge) dynamicNextConnection() int {
	if b.connections == nil {
		b.connections = &dynamicConnectionTable{values: map[int]rprimitive.Value{}}
	}
	for index := 3; index < 128; index++ {
		if _, used := b.connections.values[index]; !used {
			return index
		}
	}
	panic(translatedRuntimeError{message: "all connection slots are in use"})
}
func dynamicOpenFile(name, mode string) (*os.File, error) {
	name = dynamicExpandFileName(name)
	mode = strings.TrimSpace(mode)
	if mode == "" {
		mode = "r"
	}
	flags := os.O_RDONLY
	perm := os.FileMode(0o666)
	switch mode {
	case "r", "rb":
		flags = os.O_RDONLY
	case "r+", "r+b", "rb+":
		flags = os.O_RDWR
	case "w", "wb":
		flags = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	case "w+", "w+b", "wb+":
		flags = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	case "a", "ab":
		flags = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	case "a+", "a+b", "ab+":
		flags = os.O_RDWR | os.O_CREATE | os.O_APPEND
	default:
		return nil, fmt.Errorf("unsupported file open mode %q", mode)
	}
	return os.OpenFile(name, flags, perm)
}

func dynamicIdentity(value rprimitive.Value) string {
	if value == nil {
		return "<nil>"
	}
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Func, reflect.Chan, reflect.UnsafePointer:
		if rv.IsNil() {
			return fmt.Sprintf("%T:nil", value)
		}
		return fmt.Sprintf("%T:%x", value, rv.Pointer())
	default:
		return fmt.Sprintf("%T:%v", value, value)
	}
}

func dynamicHashPJW(text string) uint32 {
	var hash uint32
	for i := 0; i < len(text); i++ {
		hash = (hash << 4) + uint32(text[i])
		high := hash & 0xF0000000
		if high != 0 {
			hash ^= high >> 24
			hash &^= high
		}
	}
	return hash
}

func dynamicBuiltinNames() []string {
	names := make([]string, 0, len(rprimitive.PrimitiveTable))
	for name := range rprimitive.PrimitiveTable {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
func dynamicExpandFileName(name string) string {
	if name == "" || strings.Contains(name, "://") {
		return name
	}
	name = os.ExpandEnv(name)
	if name == "~" || strings.HasPrefix(name, "~/") || strings.HasPrefix(name, `~\\`) {
		if home, err := os.UserHomeDir(); err == nil {
			if name == "~" {
				name = home
			} else {
				name = filepath.Join(home, name[2:])
			}
		}
	}
	return filepath.Clean(name)
}
func dynamicAttributeKey(value rprimitive.Value) string {
	key := strings.Trim(fmt.Sprint(value), "{} ")
	switch key {
	case "R_ClassSymbol":
		return "class"
	case "R_NameSymbol":
		return "name"
	case "R_NamesSymbol":
		return "names"
	case "R_DimSymbol":
		return "dim"
	case "R_DimNamesSymbol":
		return "dimnames"
	case "R_LevelsSymbol":
		return "levels"
	case "R_RowNamesSymbol":
		return "row.names"
	}
	return key
}
func cTypeOf(value Value) string {
	if value == nil {
		return "NILSXP"
	}
	switch value.(type) {
	case Null:
		return "NILSXP"
	case dynamicSymbol:
		return "SYMSXP"
	case *dynamicExternalPtr:
		return "EXTPTRSXP"
	case *dynamicMissingArg:
		return "PROMSXP"
	case *RawVector:
		return "RAWSXP"
	case *LogicalVector:
		return "LGLSXP"
	case *IntegerVector:
		return "INTSXP"
	case *DoubleVector, *dynamicNAReal:
		return "REALSXP"
	case *ComplexVector:
		return "CPLXSXP"
	case *CharacterVector, *dynamicNAString:
		return "STRSXP"
	case *List:
		return "VECSXP"
	case *Environment, *EnvironmentValue:
		return "ENVSXP"
	case *Closure:
		return "CLOSXP"
	case *Language:
		return "LANGSXP"
	}
	return "NILSXP"
}
func allocateDynamicVector(kind string, n int) Value {
	if n < 0 {
		n = 0
	}
	switch strings.ToUpper(kind) {
	case "STRSXP", "CHARACTER":
		return &CharacterVector{Data: make([]string, n), Missing: make([]bool, n)}
	case "INTSXP", "INTEGER":
		return &IntegerVector{Data: make([]int64, n), Missing: make([]bool, n)}
	case "REALSXP", "DOUBLE", "NUMERIC":
		return &DoubleVector{Data: make([]float64, n), Missing: make([]bool, n)}
	case "LGLSXP", "LOGICAL":
		return &LogicalVector{Data: make([]Logical, n)}
	case "CPLXSXP", "COMPLEX":
		return &ComplexVector{Data: make([]complex128, n)}
	case "RAWSXP", "RAW":
		return &RawVector{Data: make([]byte, n)}
	default:
		return &List{Data: make([]Value, n), Names: make([]string, n)}
	}
}
func dynamicAttribute(value Value, key string) Value {
	attrs := attributesOf(value)
	if attrs != nil {
		if v, ok := attrs[key]; ok {
			return v
		}
	}
	return NullValue
}
func setDynamicAttribute(value Value, key string, attribute Value) {
	attrs := attributesOf(value)
	if attrs == nil {
		attrs = map[string]Value{}
		applyAttributes(value, attrs)
	}
	attrs[key] = attribute
}
func cloneDynamicValue(value Value) Value {
	var clone Value
	switch v := value.(type) {
	case *RawVector:
		clone = &RawVector{Data: append([]byte(nil), v.Data...)}
	case *LogicalVector:
		clone = &LogicalVector{Data: append([]Logical(nil), v.Data...)}
	case *IntegerVector:
		clone = &IntegerVector{Data: append([]int64(nil), v.Data...), Missing: append([]bool(nil), v.Missing...)}
	case *DoubleVector:
		clone = &DoubleVector{Data: append([]float64(nil), v.Data...), Missing: append([]bool(nil), v.Missing...)}
	case *ComplexVector:
		clone = &ComplexVector{Data: append([]complex128(nil), v.Data...)}
	case *CharacterVector:
		clone = &CharacterVector{Data: append([]string(nil), v.Data...), Missing: append([]bool(nil), v.Missing...)}
	case *List:
		clone = &List{Data: append([]Value(nil), v.Data...), Names: append([]string(nil), v.Names...)}
	default:
		return value
	}
	if attrs := attributesOf(value); attrs != nil {
		applyAttributes(clone, cloneDynamicAttributes(attrs))
	}
	return clone
}

type translatedRuntimeError struct{ message string }

func (e translatedRuntimeError) Error() string { return e.message }

func runtimeValue(value rprimitive.Value) Value {
	if v, ok := value.(Value); ok {
		return v
	}
	return NullValue
}
func dynamicEnvironment(value rprimitive.Value) *Environment {
	if environment, ok := value.(*Environment); ok {
		return environment
	}
	if wrapped, ok := value.(*EnvironmentValue); ok {
		return wrapped.Env
	}
	return nil
}
func firstInteger(value Value) int64 {
	if v, ok := value.(*IntegerVector); ok && len(v.Data) > 0 {
		return v.Data[0]
	}
	if v, ok := value.(*DoubleVector); ok && len(v.Data) > 0 {
		return int64(v.Data[0])
	}
	return 0
}
func firstDouble(value Value) float64 {
	if v, ok := value.(*DoubleVector); ok && len(v.Data) > 0 {
		return v.Data[0]
	}
	if v, ok := value.(*IntegerVector); ok && len(v.Data) > 0 {
		return float64(v.Data[0])
	}
	return 0
}

func (c *Context) executeTranslatedEntry(plan ExecutionPlan, args []syntax.Argument, env *Environment) (value Value, handled bool, err error) {
	if _, ok := rprimitive.PrimitiveTable[plan.Name]; !ok {
		return nil, false, nil
	}
	switch plan.Name {
	case ".C", ".Fortran", ".Call", ".Call.graphics", ".External", ".External2", ".External.graphics", "dyn.load", "dyn.unload", "getLoadedDLLs":
		return nil, true, rprimitive.PluginBoundaryError{Operation: plan.Name}
	}
	values := make([]Value, len(args))
	for i, argument := range args {
		if argument.Value == nil {
			return nil, true, fmt.Errorf("missing argument %d", i+1)
		}
		values[i], err = c.Eval(argument.Value, env)
		if err != nil {
			return nil, true, err
		}
	}
	bridge := newTranslatedRuntimeBridge(c, env)
	translatedRuntimeMu.Lock()
	previous := rprimitive.RT
	rprimitive.SetRuntime(bridge)
	defer func() {
		rprimitive.SetRuntime(previous)
		translatedRuntimeMu.Unlock()
		if recovered := recover(); recovered != nil {
			value = nil
			handled = true
			switch failure := recovered.(type) {
			case translatedRuntimeError:
				err = failure
			case rprimitive.PluginBoundaryError:
				err = failure
			case error:
				err = fmt.Errorf("translated %s runtime failure: %w", plan.Name, failure)
			default:
				err = fmt.Errorf("translated %s requires unimplemented GNU-R runtime operation: %v", plan.Name, recovered)
			}
		}
	}()
	result := rprimitive.InvokePrimitive(plan.Name, NullValue, NullValue, &List{Data: values}, env)
	if runtimeResult, ok := result.(Value); ok {
		return runtimeResult, true, nil
	}
	switch v := result.(type) {
	case nil:
		return NullValue, true, nil
	case bool:
		if v {
			return &LogicalVector{Data: []Logical{True}}, true, nil
		}
		return &LogicalVector{Data: []Logical{False}}, true, nil
	case int64:
		return &IntegerVector{Data: []int64{v}}, true, nil
	case float64:
		return &DoubleVector{Data: []float64{v}}, true, nil
	case string:
		return &CharacterVector{Data: []string{v}}, true, nil
	default:
		return nil, true, fmt.Errorf("translated %s returned unsupported dynamic value %T", plan.Name, result)
	}
}
