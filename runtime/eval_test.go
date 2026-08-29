package runtime

import (
	"bytes"
	"math"
	"path/filepath"
	"strconv"
	"testing"

	"r2go/syntax"
)

func evaluate(t *testing.T, src string) (Value, string) {
	t.Helper()
	p, err := syntax.Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	ctx := NewContext()
	var output bytes.Buffer
	ctx.Output = &output
	v, err := ctx.EvalProgram(p)
	if err != nil {
		t.Fatal(err)
	}
	return v, output.String()
}

func TestClosureAndDefaultPromise(t *testing.T) {
	v, _ := evaluate(t, "make <- function(x) function(y = x) x + y\nf <- make(4)\nf()")
	if got := v.String(); got != "8" {
		t.Fatalf("got %s", got)
	}
}

func TestUnusedPromiseIsNotForced(t *testing.T) {
	v, _ := evaluate(t, "f <- function(x) 42\nf(missing_object)")
	if got := v.String(); got != "42" {
		t.Fatalf("got %s", got)
	}
}

func TestVectorRecyclingAndLoop(t *testing.T) {
	v, _ := evaluate(t, "x <- c(1, 2, 3, 4)\ny <- x + c(10, 20)\ntotal <- 0\nfor (item in y) total <- total + item\ntotal")
	if got := v.String(); got != "70" {
		t.Fatalf("got %s", got)
	}
}

func TestQuoteEval(t *testing.T) {
	v, _ := evaluate(t, "x <- 9\neval(quote(x + 1))")
	if got := v.String(); got != "10" {
		t.Fatalf("got %s", got)
	}
}

func TestArgumentMatchingPasses(t *testing.T) {
	v, _ := evaluate(t, "f <- function(alpha, beta, ...) alpha * 10 + beta\nf(be = 2, al = 3, extra = missing_object)")
	if got := v.String(); got != "32" {
		t.Fatalf("got %s", got)
	}
}

func TestMissingObservesDefaultWithoutForcing(t *testing.T) {
	v, _ := evaluate(t, "f <- function(x = 7) missing(x)\nf()")
	if got := v.String(); got != "TRUE" {
		t.Fatalf("got %s", got)
	}
	v, _ = evaluate(t, "f <- function(x = 7) missing(x)\nf(9)")
	if got := v.String(); got != "FALSE" {
		t.Fatalf("got %s", got)
	}
}

func TestUnusedArgumentFails(t *testing.T) {
	p, err := syntax.Parse("f <- function(x) x\nf(1, y = 2)")
	if err != nil {
		t.Fatal(err)
	}
	_, err = NewContext().EvalProgram(p)
	if err == nil {
		t.Fatal("expected unused argument error")
	}
}

func TestRNumericSentinels(t *testing.T) {
	if !IsNAReal(NAReal()) || IsNaNButNotNA(NAReal()) {
		t.Fatal("R NA payload not preserved")
	}
	if !IsNaNButNotNA(math.NaN()) {
		t.Fatal("ordinary NaN misclassified")
	}
	want := []float64{-2, -2, 0, 0, 2, 2}
	input := []float64{-2.5, -1.5, -0.5, 0.5, 1.5, 2.5}
	for i, x := range input {
		if got := RoundZero(x); got != want[i] {
			t.Fatalf("RoundZero(%v)=%v", x, got)
		}
	}
}

func TestDistinctSubsettingSemantics(t *testing.T) {
	v, _ := evaluate(t, "x <- c(10, 20, 30)\nx[c(3, 1, 9)]")
	if got := v.String(); got != "30 10 NA" {
		t.Fatalf("got %s", got)
	}
	v, _ = evaluate(t, "x <- list(first = 10, second = 20)\nx[[2]] + x$first")
	if got := v.String(); got != "30" {
		t.Fatalf("got %s", got)
	}
}

func TestStructuredErrorConditionAndFinally(t *testing.T) {
	v, output := evaluate(t, "tryCatch(stop('boom'), error = function(e) conditionMessage(e), finally = print('cleanup'))")
	if got := v.String(); got != "\"boom\"" {
		t.Fatalf("got %s", got)
	}
	if output != "\"cleanup\"\n" {
		t.Fatalf("output=%q", output)
	}
}

func TestWarningIsRecorded(t *testing.T) {
	p, err := syntax.Parse("warning('careful')\n42")
	if err != nil {
		t.Fatal(err)
	}
	ctx := NewContext()
	v, err := ctx.EvalProgram(p)
	if err != nil {
		t.Fatal(err)
	}
	if v.String() != "42" || len(ctx.Warnings) != 1 || ctx.Warnings[0].Message != "careful" {
		t.Fatalf("value=%s warnings=%v", v.String(), ctx.Warnings)
	}
}

func TestIntegerAndTypedMissingValues(t *testing.T) {
	v, _ := evaluate(t, "typeof(1L)")
	if v.String() != "\"integer\"" {
		t.Fatalf("got %s", v.String())
	}
	v, _ = evaluate(t, "is.na(c(NA_integer_, 2L, NA_real_))")
	if v.String() != "TRUE FALSE TRUE" {
		t.Fatalf("got %s", v.String())
	}
}

func TestModuloIntegerDivisionAndMissingPropagation(t *testing.T) {
	v, _ := evaluate(t, "c(-5 %% 3, -5 %/% 3, NA_real_ + 1)")
	if v.String() != "1 -2 NA" {
		t.Fatalf("got %s", v.String())
	}
}

func TestLazyScalarLogic(t *testing.T) {
	v, _ := evaluate(t, "FALSE && stop('must not run')")
	if v.String() != "FALSE" {
		t.Fatalf("got %s", v.String())
	}
}

func TestAtomicVectorConstructors(t *testing.T) {
	v, _ := evaluate(t, "c(length(numeric(3)), sum(numeric(3)), length(integer(2)), length(logical(4)), length(character(1)), length(complex(2)))")
	if v.String() != "3 0 2 4 1 2" {
		t.Fatalf("constructors=%s", v.String())
	}
}

func TestNumericSummariesAndConversions(t *testing.T) {
	v, _ := evaluate(t, "sum(as.double(c(1L, 2L, 3L))) + mean(c(2, 4))")
	if v.String() != "9" {
		t.Fatalf("got %s", v.String())
	}
	v, _ = evaluate(t, "seq_along(c('a', 'b', 'c'))")
	if v.String() != "1L 2L 3L" {
		t.Fatalf("got %s", v.String())
	}
}

func TestNamesAttributesAndCopyOnModify(t *testing.T) {
	v, _ := evaluate(t, "x <- c(10, 20)\ny <- x\nnames(x) <- c('a', 'b')\nc(x['b'], y[1])")
	if v.String() != "20 10" {
		t.Fatalf("got %s", v.String())
	}
	v, _ = evaluate(t, "x <- structure(c(1, 2), class = 'demo', note = 'ok')\nc(class(x), attr(x, 'note'))")
	if v.String() != "\"demo\" \"ok\"" {
		t.Fatalf("got %s", v.String())
	}
}

func TestVectorListAndDollarReplacement(t *testing.T) {
	v, _ := evaluate(t, "x <- c(1, 2, 3)\nx[c(1, 3)] <- c(10, 30)\nx[5] <- 50\nx")
	if v.String() != "10 2 30 NA 50" {
		t.Fatalf("got %s", v.String())
	}
	v, _ = evaluate(t, "x <- list(a = 1)\nx$b <- 2\nx$a <- 9\nx$a + x$b")
	if v.String() != "11" {
		t.Fatalf("got %s", v.String())
	}
}

func TestMatrixColumnMajorIndexingAndDrop(t *testing.T) {
	v, _ := evaluate(t, "m <- matrix(1:6, nrow = 2, ncol = 3)\nc(m[2, 1], m[1, 3], nrow(m), ncol(m))")
	if v.String() != "2L 5L 2L 3L" {
		t.Fatalf("got %s", v.String())
	}
	v, _ = evaluate(t, "m <- matrix(1:6, 2, 3)\nm[, 2]")
	if v.String() != "3L 4L" {
		t.Fatalf("got %s", v.String())
	}
}

func TestMatrixInferenceByRowAndDimensionNames(t *testing.T) {
	v, _ := evaluate(t, "m <- matrix(1:6, nrow=2, byrow=TRUE)\nrownames(m) <- c('r1','r2')\ncolnames(m) <- c('a','b','c')\nc(dim(m), m[1,2], m[2,1], rownames(m), colnames(m))")
	if v.String() != `"2" "3" "2" "4" "r1" "r2" "a" "b" "c"` {
		t.Fatalf("matrix inference/byrow/dimnames=%s", v.String())
	}
}

func TestMatrixReplacementAndDimAssignment(t *testing.T) {
	v, _ := evaluate(t, "m <- matrix(1:6, 2, 3)\nm[2, c(1, 3)] <- c(20, 60)\nc(m[2, 1], m[2, 3], dim(m))")
	if v.String() != "20 60 2 3" {
		t.Fatalf("got %s", v.String())
	}
	v, _ = evaluate(t, "x <- 1:4\ndim(x) <- c(2L, 2L)\nc(x[2, 2], nrow(x))")
	if v.String() != "4L 2L" {
		t.Fatalf("got %s", v.String())
	}
}

func TestDataFrameConstructionAndColumnAccess(t *testing.T) {
	v, _ := evaluate(t, "d <- data.frame(id = 1:3, label = c('a', 'b', 'c'))\nc(nrow(d), ncol(d), d$id[2], is.data.frame(d))")
	// Numeric coercion turns TRUE into 1 in c(), matching R's common type.
	if v.String() != "3L 2L 2L 1L" {
		t.Fatalf("got %s", v.String())
	}
}

func TestDataFrameRowAndColumnSubsetting(t *testing.T) {
	v, _ := evaluate(t, "d <- data.frame(id = 1:3, score = c(10, 20, 30))\nc(d[2, 'score'], nrow(d[c(1,3), ]))")
	if v.String() != "20 2" {
		t.Fatalf("got %s", v.String())
	}
}

func TestThreeDimensionalArray(t *testing.T) {
	v, _ := evaluate(t, "a <- array(1:8, dim = c(2, 2, 2))\nc(a[2, 1, 2], dim(a))")
	if v.String() != "6L 2L 2L 2L" {
		t.Fatalf("got %s", v.String())
	}
}

func TestNestedReplacementCopyOnModify(t *testing.T) {
	v, _ := evaluate(t, "x <- list(a = c(1, 2, 3))\ny <- x\nx$a[2] <- 20\nc(x$a, y$a)")
	if v.String() != "1 20 3 1 2 3" {
		t.Fatalf("got %s", v.String())
	}
	v, _ = evaluate(t, "x <- list(inner = list(value = 1))\nx$inner$value <- 9\nx$inner$value")
	if v.String() != "9" {
		t.Fatalf("got %s", v.String())
	}
}

func TestFactorLevelsCodesAndSubsetting(t *testing.T) {
	v, _ := evaluate(t, "f <- factor(c('red', 'blue', 'red', NA_character_), levels = c('blue', 'red'))\nc(as.character(f[1:3]), levels(f), is.factor(f[2:3]))")
	if v.String() != "\"red\" \"blue\" \"red\" \"blue\" \"red\" \"TRUE\"" {
		t.Fatalf("got %s", v.String())
	}
}

func TestS3UseMethodAndBuiltinDispatch(t *testing.T) {
	v, _ := evaluate(t, "describe <- function(x) UseMethod('describe')\ndescribe.demo <- function(x) x$value + 1\ndescribe.default <- function(x) 0\na <- structure(list(value = 9), class = 'demo')\nc(describe(a), describe(3))")
	if v.String() != "10 0" {
		t.Fatalf("got %s", v.String())
	}
	v, output := evaluate(t, "print.demo <- function(x) print(x$value + 5)\na <- structure(list(value = 2), class = 'demo')\nprint(a)")
	if v.String() != "7" || output != "7\n" {
		t.Fatalf("value=%s output=%q", v.String(), output)
	}
}

func TestSequenceRepeatAndLogicalBaseFunctions(t *testing.T) {
	v, _ := evaluate(t, "c(seq(2, 8, by = 2), rep(c(9, 8), times = 2), which(c(FALSE, TRUE, NA, TRUE)))")
	if v.String() != "2 4 6 8 9 8 9 8 2 4" {
		t.Fatalf("got %s", v.String())
	}
	v, _ = evaluate(t, "c(any(FALSE, NA, na.rm = TRUE), all(TRUE, TRUE), all(TRUE, NA))")
	if v.String() != "FALSE TRUE NA" {
		t.Fatalf("got %s", v.String())
	}
}

func TestMatchingUniquenessAndPaste(t *testing.T) {
	v, _ := evaluate(t, "c(match(c('b','x'), c('a','b')), c('a','b') %in% c('b'), duplicated(c(1,1,2)))")
	if v.String() != "2L NA 0L 1L 0L 1L 0L" {
		t.Fatalf("got %s", v.String())
	}
	v, _ = evaluate(t, "paste(c('a','b'), 1:2, sep = '-', collapse = ',')")
	if v.String() != "\"a-1,b-2\"" {
		t.Fatalf("got %s", v.String())
	}
}

func TestVectorizedMathFunctions(t *testing.T) {
	v, _ := evaluate(t, "c(sqrt(c(4,9)), round(c(1.5,2.5)), floor(-1.2), abs(-3))")
	if v.String() != "2 3 2 2 -2 3" {
		t.Fatalf("got %s", v.String())
	}
}

func TestLapplyAndSapplyClosures(t *testing.T) {
	v, _ := evaluate(t, "sapply(1:4, function(x) x * x)")
	if v.String() != "1 4 9 16" {
		t.Fatalf("got %s", v.String())
	}
	v, _ = evaluate(t, "x <- lapply(c(2,3), function(x) c(x, x+1))\nx[[1]][2] + x[[2]][1]")
	if v.String() != "6" {
		t.Fatalf("got %s", v.String())
	}
	v, _ = evaluate(t, "classify <- function(x, lower, upper) if (x < lower) 'low' else if (x > upper) 'high' else 'mid'\nsapply(c(7,12,18), classify, lower=8, upper=17)")
	if v.String() != `"low" "mid" "high"` {
		t.Fatalf("sapply did not forward named arguments: got %s", v.String())
	}
}

func TestCompleteCasesAcrossDataFramesAndMatrices(t *testing.T) {
	v, _ := evaluate(t, "d <- data.frame(a=c(1,NA,3), b=c('x','y','z'))\ncomplete.cases(d)")
	if v.String() != "TRUE FALSE TRUE" {
		t.Fatalf("data-frame complete.cases: value=%v", v)
	}
	v, _ = evaluate(t, "m <- matrix(c(1,2,NA,4), nrow=2)\ncomplete.cases(m)")
	if v.String() != "FALSE TRUE" {
		t.Fatalf("matrix complete.cases: value=%v", v)
	}
}

func TestS3NextMethod(t *testing.T) {
	v, _ := evaluate(t, "describe <- function(x) UseMethod('describe')\ndescribe.child <- function(x) NextMethod() + 1\ndescribe.parent <- function(x) x$value\nx <- structure(list(value=7), class=c('child','parent'))\ndescribe(x)")
	if v.String() != "8" {
		t.Fatalf("got %s", v.String())
	}
}

func TestStatisticsOrderAndSequenceSemantics(t *testing.T) {
	v, _ := evaluate(t, "c(median(c(1,4,2,3)), var(c(1,2,3)), sd(c(1,2,3)))")
	if v.String() != "2.5 1 1" {
		t.Fatalf("statistics=%s", v.String())
	}
	v, _ = evaluate(t, "quantile(c(1,2,3,4), c(0,0.5,1))")
	if v.String() != "1 2.5 4" {
		t.Fatalf("quantile=%s", v.String())
	}
	v, _ = evaluate(t, "order(c(12,7,NA,18,4,25,9,16,NA,11), na.last=TRUE)")
	if v.String() != "5L 2L 7L 10L 1L 8L 4L 6L 3L 9L" {
		t.Fatalf("order=%s", v.String())
	}
	v, _ = evaluate(t, "seq(0.5, 5, length.out=10)")
	if v.String() != "0.5 1 1.5 2 2.5 3 3.5 4 4.5 5" {
		t.Fatalf("seq length.out=%s", v.String())
	}
}

func TestIfelseVectorizationAndMissingness(t *testing.T) {
	v, _ := evaluate(t, "ifelse(c(TRUE,FALSE,NA,TRUE), c(1,2), c(10,20))")
	if v.String() != "1 20 NA 2" {
		t.Fatalf("ifelse numeric=%s", v.String())
	}
	v, _ = evaluate(t, "ifelse(is.na(c(1,NA,3)), 'NA', as.character(c(1,NA,3)))")
	if v.String() != `"1" "NA" "3"` {
		t.Fatalf("ifelse character=%s", v.String())
	}
}

func TestTableOneWayCharacterCounts(t *testing.T) {
	v, _ := evaluate(t, "table(c('medium','low',NA,'high','medium','low'))")
	if v.String() != "1L 2L 2L" {
		t.Fatalf("table counts=%s", v.String())
	}
	dims, ok := dimensions(v)
	if !ok || len(dims) != 1 || dims[0] != 3 || !hasClass(v, "table") {
		t.Fatalf("table attributes missing: dims=%v class=%v", dims, classNames(v))
	}
}

func TestApplySplitAndGroupedLapply(t *testing.T) {
	v, _ := evaluate(t, "m <- matrix(1:6, nrow=2, byrow=TRUE)\napply(m, 1, function(row) sum(row^2))")
	if v.String() != "14 77" {
		t.Fatalf("apply rows=%s", v.String())
	}
	v, _ = evaluate(t, "m <- matrix(1:6, nrow=2, byrow=TRUE)\napply(m, 2, function(col) sum(col^2))")
	if v.String() != "17 29 45" {
		t.Fatalf("apply cols=%s", v.String())
	}
	v, _ = evaluate(t, "g <- split(c(10,20,30,40), c('b','a','b','a'))\nc(names(g), g$a[2], g$b[1])")
	if v.String() != `"a" "b" "40" "10"` {
		t.Fatalf("split=%s", v.String())
	}
	v, _ = evaluate(t, "g <- split(c(10,20,30,40), c('b','a','b','a'))\nlapply(g, function(z) mean(z))")
	if v.String() != "$a\n30\n\n$b\n20" {
		t.Fatalf("grouped lapply=%s", v.String())
	}
}

func TestAdditionalBaseVectorFamilies(t *testing.T) {
	v, _ := evaluate(t, "c(sort(c(3,1,2)), rev(1:3), cumsum(1:3), diff(c(2,5,9)))")
	if v.String() != "1 2 3 3 2 1 1 3 6 3 4" {
		t.Fatalf("got %s", v.String())
	}
	v, _ = evaluate(t, "paste(toupper(trimws(c(' a ','b '))), nchar(c('xx','yyy')), sep=':', collapse=',')")
	if v.String() != "\"A:2,B:3\"" {
		t.Fatalf("got %s", v.String())
	}
	v, _ = evaluate(t, "x <- strsplit(c('a,b','c,d'), ',')\npaste(x[[2]], collapse='-')")
	if v.String() != "\"c-d\"" {
		t.Fatalf("got %s", v.String())
	}
}

func TestComplexVectorsAndArithmetic(t *testing.T) {
	v, _ := evaluate(t, "z <- c(1+2i, 3-1i)\nc(Re(z), Im(z), Mod(3+4i), Re(Conj(2+3i)))")
	if v.String() != "1 3 2 -1 5 2" {
		t.Fatalf("got %s", v.String())
	}
	v, _ = evaluate(t, "c((1+2i)*(2-1i), NA_complex_ + 1)")
	if v.String() != "4+3i NA" {
		t.Fatalf("got %s", v.String())
	}
}

func TestDotsForwardingPreservesCallerAndLaziness(t *testing.T) {
	v, _ := evaluate(t, "inner <- function(a, b=1) a+b\nouter <- function(...) inner(...)\nx <- 9\nouter(a=x, b=3)")
	if v.String() != "12" {
		t.Fatalf("got %s", v.String())
	}
	v, _ = evaluate(t, "ignore <- function(...) 7\nouter <- function(...) ignore(...)\nouter(stop('not forced'))")
	if v.String() != "7" {
		t.Fatalf("got %s", v.String())
	}
}

func TestStopIfNotAndEnvironmentCore(t *testing.T) {
	v, _ := evaluate(t, "e <- new.env()\ne$x <- 7\ne$.hidden <- 9\nstopifnot(e$x == 7)\nc(ls(e), ls(e, all.names=TRUE))")
	if v.String() != `"x" ".hidden" "x"` {
		t.Fatalf("got %s", v.String())
	}

	p, err := syntax.Parse("stopifnot(1 == 2)")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = NewContext().EvalProgram(p); err == nil {
		t.Fatal("stopifnot accepted a false condition")
	}
}

func TestFractionalRecyclingRecordsWarning(t *testing.T) {
	p, err := syntax.Parse("1:3 + 1:2")
	if err != nil {
		t.Fatal(err)
	}
	ctx := NewContext()
	if _, err = ctx.EvalProgram(p); err != nil {
		t.Fatal(err)
	}
	if len(ctx.Warnings) != 1 || ctx.Warnings[0].Message != "longer object length is not a multiple of shorter object length" {
		t.Fatalf("warnings=%v", ctx.Warnings)
	}
}

func TestTextAndCSVFileIO(t *testing.T) {
	dir := t.TempDir()
	textPath := strconv.Quote(filepath.Join(dir, "lines.txt"))
	csvPath := strconv.Quote(filepath.Join(dir, "data.csv"))
	v, _ := evaluate(t, "p <- "+textPath+"\nwriteLines(c('alpha','beta'), p)\nc(file.exists(p), readLines(p))")
	if v.String() != "\"TRUE\" \"alpha\" \"beta\"" {
		t.Fatalf("got %s", v.String())
	}
	v, _ = evaluate(t, "p <- "+csvPath+"\nd <- data.frame(id=1:2, label=c('a','b'))\nwrite.csv(d,p)\nr <- read.csv(p)\nc(nrow(r),r$id[2],r$label[1])")
	if v.String() != "\"2\" \"2\" \"a\"" {
		t.Fatalf("got %s", v.String())
	}
}
