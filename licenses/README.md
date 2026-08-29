# License material

The release build embeds license/notice texts into `rtogo.exe`.

`go generate` in `internal/licenses` runs `tools/licensepack`, which reads the complete Go module graph, collects conventional license/notice files from every dependency, adds rtogo's own `LICENSE`, and also adds the project-level notices in this `licenses` directory. It writes the generated bundle to:

`internal/licenses/generated/THIRD_PARTY_NOTICES.txt`

That generated file is compiled into the executable with `//go:embed` and can be printed with:

`rtogo --licenses`

The project-level notices currently include the R logo attribution/license information (`NOTICE_R_LOGO.txt`) and the Go Logo/trademark notice (`NOTICE_GO_LOGO.txt`). The three direct Go dependencies used by this starter are Gio, gvcode and Chroma; their dependency licenses are collected automatically from the module graph during a release build.
