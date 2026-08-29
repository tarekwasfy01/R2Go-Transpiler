package runtime

// ImplementedNames is the machine-readable surface currently recognized by
// the evaluator. Entries describe dispatch availability, not full GNU-R
// conformance; semantic status is established separately by differential tests.
var ImplementedNames = map[string]bool{}

func init() {
	for _, name := range []string{
		"if", "while", "for", "repeat", "break", "next", "return", "function", "<-", "=", "<<-", "(", "{",
		"+", "-", "*", "/", "^", "%%", "%/%", "%in%", ":", "==", "!=", "<", "<=", ">", ">=", "&", "|", "&&", "||", "!",
		"[", "[[", "$", "quote", "eval", "missing", "stop", "warning", "tryCatch", "UseMethod", "NextMethod",
		"c", "list", "length", "print", "typeof", "is.null", "is.na", "is.data.frame", "is.factor", "is.raw", "inherits",
		"as.logical", "as.integer", "as.double", "as.numeric", "as.complex", "as.raw", "as.character", "raw", "rawToChar", "charToRaw",
		"Re", "Im", "Mod", "Arg", "Conj", "sum", "prod", "min", "max", "mean", "seq", "seq.int", "seq_along", "rep",
		"any", "all", "which", "unique", "duplicated", "match", "paste", "paste0", "lapply", "sapply", "names", "levels",
		"dim", "nrow", "ncol", "attr", "attributes", "class", "unclass", "structure", "matrix", "array", "data.frame", "factor",
		"abs", "sqrt", "exp", "log", "log10", "sin", "cos", "tan", "floor", "ceiling", "trunc", "round", "rev", "head", "tail",
		"sort", "order", "rank", "cumsum", "cumprod", "cummin", "cummax", "diff", "pmin", "pmax", "nchar", "tolower", "toupper",
		"trimws", "substr", "substring", "startsWith", "endsWith", "strsplit", "getwd", "setwd", "file.exists", "dir.exists",
		"dir.create", "basename", "dirname", "normalizePath", "readLines", "writeLines", "list.files", "read.csv", "read.csv2",
		"write.csv", "write.csv2", "serialize", "unserialize", "conditionMessage",
	} {
		ImplementedNames[name] = true
	}
	// A dispatcher-family implementation covers every R_FunTab alias sharing
	// that C entry point.  This is the key completeness multiplier: adding one
	// family updates capability reporting without another name list.
	for _, descriptor := range PrimitiveTable {
		if plan, ok := ExecutionPlanByName[descriptor.Name]; ok && planExecutable(plan) {
			ImplementedNames[descriptor.Name] = true
		}
	}
}
