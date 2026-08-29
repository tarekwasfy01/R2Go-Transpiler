# R2Go Transpiler

**R2Go** is an experimental **Proof of Concept R-to-Go transpiler** written in Go.

The project explores how much **Base R** code can be parsed, executed, and translated into standalone Go code without requiring a normal R runtime.

> **Project status:** Proof of Concept / experimental.
> R2Go currently targets **Base R only**. **CRAN packages are not supported.**
> It is not intended to be a complete replacement for GNU R and not every R language construct can currently be translated.

## Download

**Windows executable:**

[Download R2Go.exe](https://github.com/tarekwasfy01/R2Go-Transpiler/releases/download/R2Go/R2Go.exe)

Repository:

https://github.com/tarekwasfy01/R2Go-Transpiler

## GUI

Start the executable without arguments:

```powershell
R2Go.exe
```

or explicitly start the GUI:

```powershell
R2Go.exe gui
```

The GUI provides an R source editor on the left and displays the generated Go code on the right.

## CLI Commands

| Command                                                           | Description                                                                              |
| ----------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `R2Go.exe`                                                        | Starts the graphical user interface.                                                     |
| `R2Go.exe gui`                                                    | Explicitly starts the GUI.                                                               |
| `R2Go.exe transpile input.R -o output.go`                         | Transpiles an R script into a Go source file.                                            |
| `R2Go.exe transpile --strict-native input.R -o output.go`         | Uses strict native Go lowering and rejects compatibility fallback.                       |
| `R2Go.exe transpile --source-comments=false input.R -o output.go` | Generates Go code without preserving original R source comments in compatibility output. |
| `R2Go.exe run input.R`                                            | Executes an R script using the Pure-Go R2Go runtime.                                     |
| `R2Go.exe ast input.R`                                            | Parses the R source and prints the generated syntax tree as JSON.                        |
| `R2Go.exe coverage`                                               | Displays Base R primitive coverage information.                                          |
| `R2Go.exe version`                                                | Displays the current R2Go version.                                                       |
| `R2Go.exe --licenses`                                             | Displays embedded third-party licenses and attribution notices.                          |
| `R2Go.exe licenses`                                               | Alternative command for displaying embedded licenses.                                    |
| `R2Go.exe help`                                                   | Displays command-line help.                                                              |

## Example

Transpile an R script:

```powershell
R2Go.exe transpile analysis.R -o analysis.go
```

Run the generated Go source:

```powershell
go run analysis.go
```

## Current Project Status

R2Go is currently a **Proof of Concept**.

The current implementation focuses on translating and executing functionality from **Base R** using Go-native components.

Current areas include:

* Base R syntax parsing
* Pure-Go R runtime components
* Native R-to-Go lowering where supported
* Scalar and vector operations
* Control flow
* User-defined functions where supported
* Base R primitive/function compatibility
* R syntax tree inspection
* GUI and CLI inside the same executable
* Standalone Go source generation

The project is still experimental and compatibility is not complete.

## Scope

### Supported target

**Base R**

R2Go is currently developed against Base R language functionality and Base R primitives.

### Not supported

**CRAN packages are currently not supported.**

This includes third-party packages installed through CRAN and arbitrary package-specific APIs.

R2Go should therefore currently be considered a:

> **Base R → Go Proof of Concept**

and not a complete replacement for GNU R or the complete R package ecosystem.

## Disclaimer

R2Go is experimental software under active development.

Some R constructs may not yet be translated correctly or may remain unsupported. Generated Go code should be reviewed and tested before being used in production environments.
