# GNU R C-entry to Pure Go conversion

This directory contains one Go file for each of the 289 unique C entry points in `manifest.csv`, plus `runtime.go` with shared Pure-Go mechanisms. Multiple R primitives that share a C entry point map to the same Go function, matching the input manifest.

## Build

```sh
go test ./...
```

The package name is `rgo` and the module is `gnu-r-purego`. No cgo is used and no external R process is invoked.

## Conversion status

- `constraint_error_path`: 11
- `missing_helper`: 57
- `partial_constraint`: 8
- `partial_evaluator`: 36
- `partial_format`: 7
- `translated`: 105
- `translated_family`: 65

`conversion_manifest.csv` is authoritative per batch. `missing_helpers.md` explains constrained or corpus-external dependencies. Entry points that cannot be faithfully executed under the stated constraints return an explicit error `Value`; they are not reported as fully translated.

## Value model

`Value` represents logical, integer, double, complex, character, raw, list, symbol, environment, function, error and connection values. NA is tracked independently from floating-point NaN. Attributes (including names/dim/dimnames) are carried in `Attr`; clone/write helpers provide copy-on-write behavior for shared values.
