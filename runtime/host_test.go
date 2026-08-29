package runtime

import (
	"r2go/syntax"
	"testing"
	"time"
)

func TestMemoryHostKeepsEnvironmentAndClockDeterministic(t *testing.T) {
	host := NewMemoryHost()
	host.Time = time.Date(2001, 2, 3, 4, 5, 6, 0, time.UTC)
	ctx := NewContextWithHost(host)
	program, err := syntax.Parse(`{Sys.setenv("R2GO_MATRIX","ok");Sys.getenv("R2GO_MATRIX")}`)
	if err != nil {
		t.Fatal(err)
	}
	value, err := ctx.EvalProgram(program)
	if err != nil {
		t.Fatal(err)
	}
	if got := value.String(); got != `"ok"` {
		t.Fatalf("getenv=%s", got)
	}
	if host.Getenv("R2GO_MATRIX") != "ok" {
		t.Fatal("environment was not host-scoped")
	}
}

func TestMemoryHostRoutesFilePrimitivesWithoutOSWrites(t *testing.T) {
	host := NewMemoryHost()
	ctx := NewContextWithHost(host)
	program, err := syntax.Parse(`{dir.create("work");writeLines(c("a","b"),"work/x.txt");readLines("work/x.txt")}`)
	if err != nil {
		t.Fatal(err)
	}
	value, err := ctx.EvalProgram(program)
	if err != nil {
		t.Fatal(err)
	}
	if got := value.String(); got != `"a" "b"` {
		t.Fatalf("readLines=%s", got)
	}
	if _, exists := host.Files["/work/x.txt"]; !exists {
		t.Fatal("file was not stored in MemoryHost")
	}
}

func TestMemoryHostEffectMatrix(t *testing.T) {
	cases := []struct{ name, source, want string }{
		{"environment", `{Sys.setenv("A","x");Sys.unsetenv("A");Sys.getenv("A","missing")}`, `"missing"`},
		{"working-directory", `{dir.create("d");setwd("d");getwd()}`, `"/d"`},
		{"file-existence", `{dir.create("d");writeLines("x","d/a");c(dir.exists("d"),file.exists("d/a"))}`, `TRUE TRUE`},
		{"csv-roundtrip", `{dir.create("d");writeLines(c("id,label","1,a","2,b"),"d/a.csv");x<-read.csv("d/a.csv");c(nrow(x),x[[1]][2])}`, `2 2`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := syntax.Parse(tc.source)
			if err != nil {
				t.Fatal(err)
			}
			v, err := NewContextWithHost(NewMemoryHost()).EvalProgram(p)
			if err != nil {
				t.Fatal(err)
			}
			if got := v.String(); got != tc.want {
				t.Fatalf("got %s want %s", got, tc.want)
			}
		})
	}
}
