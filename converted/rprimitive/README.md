# GNU R C-to-Go remaining translation

This archive is the Go translation of all **119 C batch files** supplied in `gnu-r-c-to-go-remaining-input.zip`.
The manifest contains **158 primitive mappings** to **119 translated C entry functions**. The three platform-conditional entries (`fifo`, `pipe`, `url/file`) are emitted as separate Windows and Unix Go files so both original preprocessor branches are retained.

## Structure

- `batch_XXX.go` — translated Go control flow for platform-neutral inputs.
- `batch_047_*`, `batch_087_*`, `batch_118_*` — Windows/Unix variants selected with Go build tags.
- `runtime.go` — C-free dynamic compatibility layer for GNU R internal runtime operations.
- `registry.go` — primitive name / C-entry / offset / source-file registry generated from `manifest.csv`.
- `manifest.csv` and `missing_for_converter.csv` — original mapping metadata.
- `translation_report.csv` — one row for each translated batch.
- `runtime_test.go` — inventory and runtime sanity tests.

## Runtime boundary

The supplied C snippets are not a complete copy of the GNU R runtime. They call hundreds of internal APIs and macros (`SEXP`, environments, connections, graphics, native routines, serialization, LAPACK, and more). The translated Go therefore routes those external runtime operations through the package-global `RT` interface in `runtime.go`. The **control flow itself is translated to Go**; no C compiler or cgo is required by this package.

An embedding GNU R Go port can install a custom `DynamicRuntime` with `SetRuntime`, or use `Runtime.Register` to implement individual GNU R operations.

## Source entries not present in the input

`missing_for_converter.csv` also lists 9 primitives for which the input explicitly says no C body was found. They cannot be translated from this archive because there is no function body in the supplied data:

`.Primitive`, `browser`, `deparseRd`, `format.info`, `formatC`, `parse`, `parseLatex`, `parse_Rd`, `~`

## Build check

Run:

```text
go test ./...
```

The package is intentionally C-free.
