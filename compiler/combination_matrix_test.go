package compiler

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"r2go/runtime"
	"r2go/syntax"
)

func verifyCombination(t *testing.T, source string, requireMatrixLowering bool) runtime.Value {
	t.Helper()
	program, err := syntax.Parse(source)
	if err != nil {
		t.Fatalf("parse R %q: %v", source, err)
	}
	generated, err := GenerateMainWithSource(program, source)
	if err != nil {
		t.Fatalf("generate Go for %q: %v", source, err)
	}
	text := string(generated)
	if _, err := parser.ParseFile(token.NewFileSet(), "generated.go", generated, parser.AllErrors); err != nil {
		t.Fatalf("parse generated Go for %q: %v\n%s", source, err, text)
	}
	if strings.Contains(text, "RunEncodedProgramInContext") {
		t.Fatalf("opaque IR fallback for %q", source)
	}
	if requireMatrixLowering && strings.Contains(text, "compatibilitySource") {
		t.Fatalf("readable matrix lowering unexpectedly fell back for %q\n%s", source, text)
	}
	value, err := runtime.NewContext().EvalProgram(program)
	if err != nil {
		t.Fatalf("execute %q: %v", source, err)
	}
	return value
}

func canonicalRuntimeValue(value runtime.Value) string {
	switch x := value.(type) {
	case runtime.Null:
		return "NULL:"
	case *runtime.LogicalVector:
		parts := make([]string, len(x.Data))
		for i, item := range x.Data {
			parts[i] = item.String()
		}
		return "logical:" + strings.Join(parts, "|")
	case *runtime.IntegerVector:
		parts := make([]string, len(x.Data))
		for i, item := range x.Data {
			if i < len(x.Missing) && x.Missing[i] {
				parts[i] = "NA"
			} else {
				parts[i] = strconv.FormatInt(item, 10)
			}
		}
		return "integer:" + strings.Join(parts, "|")
	case *runtime.DoubleVector:
		parts := make([]string, len(x.Data))
		for i, item := range x.Data {
			if i < len(x.Missing) && x.Missing[i] {
				parts[i] = "NA"
			} else {
				parts[i] = strconv.FormatFloat(item, 'g', 17, 64)
			}
		}
		return "double:" + strings.Join(parts, "|")
	case *runtime.CharacterVector:
		parts := make([]string, len(x.Data))
		for i, item := range x.Data {
			if i < len(x.Missing) && x.Missing[i] {
				parts[i] = "NA"
			} else {
				parts[i] = strconv.Quote(item)
			}
		}
		return "character:" + strings.Join(parts, "|")
	default:
		return string(value.Kind()) + ":" + value.String()
	}
}

func canonicalRValue(t *testing.T, rscript, source string) string {
	t.Helper()
	rCode := `h <- commandArgs(TRUE)[[1]]; starts <- seq.int(1,nchar(h),by=2); src <- rawToChar(as.raw(strtoi(substring(h,starts,starts+1),16L))); x <- eval(parse(text=src), envir=new.env(parent=baseenv())); typ <- typeof(x); enc <- function(v) { if (length(v)==0) return(""); paste(vapply(seq_along(v), function(i) { if (is.na(v[[i]])) "NA" else if (typ=="character") encodeString(v[[i]], quote="\"") else if (typ=="logical") if (v[[i]]) "TRUE" else "FALSE" else as.character(v[[i]]) }, character(1)), collapse="|") }; cat(typ, ":", enc(x), sep="")`
	cmd := exec.Command(rscript, "--vanilla", "-e", rCode, hex.EncodeToString([]byte(source)))
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("GNU R reference failed for %q: %v: %s", source, err, stderr.String())
	}
	return strings.TrimSpace(string(out))
}

// TestNumericCombinationMatrix checks interfaces between whole function
// families. Its cartesian product catches bugs that isolated primitive rows do
// not expose (shape, type, attributes, NA and argument forwarding).
func TestNumericCombinationMatrix(t *testing.T) {
	producers := []string{
		"c(1,2,3,4)",
		"1:6",
		"seq(1,6)",
		"rep(c(1,2,3),2)",
		"sort(c(4,1,3,2))",
	}
	transforms := []string{"abs(%s)", "sqrt(%s)", "round(%s)", "rev(%s)", "sort(%s)", "unique(%s)", "cumsum(%s)", "as.numeric(%s)"}
	consumers := []string{"sum(%s)", "prod(%s)", "min(%s)", "max(%s)", "mean(%s)", "length(%s)"}
	for pi, producer := range producers {
		for ti, transform := range transforms {
			middle := fmt.Sprintf(transform, producer)
			for ci, consumer := range consumers {
				source := fmt.Sprintf(consumer, middle)
				name := fmt.Sprintf("p%02d_t%02d_c%02d", pi, ti, ci)
				t.Run(name, func(t *testing.T) { verifyCombination(t, source, true) })
			}
		}
	}
}

func TestStringCombinationMatrix(t *testing.T) {
	producers := []string{"c(' Alpha ','beta','ALPHA')", "paste(c('a','b'),1:2,sep='-')", "rep(c('x','Y'),2)"}
	transforms := []string{"tolower(%s)", "toupper(%s)", "trimws(%s)", "sort(%s)", "unique(%s)", "rev(%s)", "substr(%s,1,2)"}
	consumers := []string{"length(%s)", "nchar(%s)", "paste(%s,collapse=',')", "sort(%s)"}
	for pi, producer := range producers {
		for ti, transform := range transforms {
			middle := fmt.Sprintf(transform, producer)
			for ci, consumer := range consumers {
				source := fmt.Sprintf(consumer, middle)
				name := fmt.Sprintf("p%02d_t%02d_c%02d", pi, ti, ci)
				t.Run(name, func(t *testing.T) { verifyCombination(t, source, true) })
			}
		}
	}
}

func TestStructuralCombinationMatrix(t *testing.T) {
	type testCase struct{ source, expected string }
	cases := map[string]testCase{
		"matrix_transpose_subset": {"m <- matrix(1:9,3,3); sum(t(m)[2,])", "15L"},
		"matrix_dimnames_subset":  {"m <- matrix(1:4,2,2); rownames(m)<-c('a','b'); colnames(m)<-c('x','y'); m['b','y']", "4L"},
		"list_member_vector":      {"x <- list(a=1:4,b=c('x','y')); sum(x$a)", "10L"},
		"dataframe_member_subset": {"d <- data.frame(a=1:4,b=5:8); sum(d$a[d$b>6])", "7L"},
		"factor_table_names":      {"x <- factor(c('a','b','a')); names(table(x))", `"a" "b"`},
		"closure_matrix":          {"f <- function(x,n=2) { m <- matrix(x,n,n); sum(diag(m)) }; f(1:4)", "5L"},
		"closure_nested_apply":    {"f <- function(x) sapply(x,function(y) y*2); sum(f(1:4))", "20"},
		"control_subset_replace":  {"x <- 1:5; for(i in 2:4) { if(i%%2==0) x[i] <- x[i]*10 }; sum(x)", "69"},
		"missing_named_arguments": {"paste(rep(c('a','b'),times=2),collapse='|')", `"a|b|a|b"`},
		"logical_na_chain":        {"x <- c(TRUE,NA,FALSE); c(any(x,na.rm=TRUE),all(x,na.rm=TRUE),sum(is.na(x)))", "1L 0L 1L"},
	}
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			value := verifyCombination(t, test.source, false)
			if value.String() != test.expected {
				t.Fatalf("result mismatch for %q: got %q, want %q", test.source, value.String(), test.expected)
			}
		})
	}
}

func TestStructuralDifferentialAgainstGNU_R(t *testing.T) {
	rscript := `C:\Users\tarek\Desktop\GroundGIS\Release\Package\runtimes\R\bin\Rscript.exe`
	if _, err := os.Stat(rscript); err != nil {
		t.Skipf("GNU R reference is unavailable: %v", err)
	}
	cases := map[string]string{
		"matrix_transpose_subset": "m <- matrix(1:9,3,3); sum(t(m)[2,])",
		"matrix_dimnames_subset":  "m <- matrix(1:4,2,2); rownames(m)<-c('a','b'); colnames(m)<-c('x','y'); m['b','y']",
		"list_member_vector":      "x <- list(a=1:4,b=c('x','y')); sum(x$a)",
		"dataframe_member_subset": "d <- data.frame(a=1:4,b=5:8); sum(d$a[d$b>6])",
		"factor_table_names":      "x <- factor(c('a','b','a')); names(table(x))",
		"closure_matrix":          "f <- function(x,n=2) { m <- matrix(x,n,n); sum(diag(m)) }; f(1:4)",
		"closure_nested_apply":    "f <- function(x) sapply(x,function(y) y*2); sum(f(1:4))",
		"control_subset_replace":  "x <- 1:5; for(i in 2:4) { if(i%%2==0) x[i] <- x[i]*10 }; sum(x)",
		"missing_named_arguments": "paste(rep(c('a','b'),times=2),collapse='|')",
		"logical_na_chain":        "x <- c(TRUE,NA,FALSE); c(any(x,na.rm=TRUE),all(x,na.rm=TRUE),sum(is.na(x)))",
	}
	for name, source := range cases {
		t.Run(name, func(t *testing.T) {
			program, err := syntax.Parse(source)
			if err != nil {
				t.Fatal(err)
			}
			value, err := runtime.NewContext().EvalProgram(program)
			if err != nil {
				t.Fatal(err)
			}
			got, want := canonicalRuntimeValue(value), canonicalRValue(t, rscript, source)
			if got != want {
				t.Fatalf("GNU R mismatch for %q: Go=%q R=%q", source, got, want)
			}
		})
	}
}
