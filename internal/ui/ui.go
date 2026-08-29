package ui

import (
	"context"
	"fmt"
	"image/color"
	"io"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/oligo/gvcode"
	gvcolor "github.com/oligo/gvcode/color"
	"github.com/oligo/gvcode/textstyle/syntax"
	gvwidget "github.com/oligo/gvcode/widget"

	standalone "r2go"
	"r2go/internal/appcore"
	"r2go/internal/highlight"
	"r2go/internal/platform"
)

type conversionResult struct {
	generation uint64
	code       string
	tokens     []syntax.Token
	err        error
}

type saveResult struct {
	path string
	err  error
}

type App struct {
	window *app.Window
	theme  *material.Theme
	engine appcore.Engine
	pool   *appcore.Pool
	hl     *highlight.Service

	left  *gvcode.Editor
	right *gvcode.Editor

	convertBtn  widget.Clickable
	copyBtn     widget.Clickable
	saveBtn     widget.Clickable
	infoBtn     widget.Clickable
	copyInfoBtn widget.Clickable

	allowFallback  widget.Bool
	preserveSource widget.Bool

	showInfo bool
	status   string
	busy     bool

	convertGeneration atomic.Uint64
	leftGeneration    atomic.Uint64
	rightGeneration   atomic.Uint64
	cancelConvert     context.CancelFunc

	convertResults chan conversionResult
	saveResults    chan saveResult
}

func New(engine appcore.Engine) *App {
	if engine == nil {
		engine = appcore.R2GoEngine{}
	}

	// Give Go at least four scheduler Ps. rtogo's own CPU-bound helpers are kept
	// bounded; two Ps are intentionally not consumed by its worker budget.
	if runtime.GOMAXPROCS(0) < 4 {
		runtime.GOMAXPROCS(4)
	}

	th := material.NewTheme()
	w := &app.Window{}
	w.Option(
		app.Title("R2Go"),
		app.Size(unit.Dp(1280), unit.Dp(760)),
		app.MinSize(unit.Dp(900), unit.Dp(560)),
	)

	a := &App{
		window:         w,
		theme:          th,
		engine:         engine,
		pool:           appcore.NewPool(1), // exactly one active transpile request
		convertResults: make(chan conversionResult, 4),
		saveResults:    make(chan saveResult, 2),
		status:         "Ready",
		allowFallback:  widget.Bool{Value: true},
		preserveSource: widget.Bool{Value: true},
	}
	a.hl = highlight.NewService(w.Invalidate)
	a.left = newCodeEditor(th, false, codeColorScheme(true, true))
	a.right = newCodeEditor(th, true, codeColorScheme(true, true))
	const initialR = "# Enter R code here\nx <- c(1, 2, 3)\nprint(x * 2)\n"
	const initialGo = "// Go output will appear here.\n"
	a.left.SetText(initialR)
	a.right.SetText(initialGo)
	if tokens, err := highlight.Tokens(context.Background(), highlight.R, initialR); err == nil {
		a.left.SetSyntaxTokens(tokens...)
	}
	if tokens, err := highlight.Tokens(context.Background(), highlight.Go, initialGo); err == nil {
		a.right.SetSyntaxTokens(tokens...)
	}
	a.scheduleHighlight("left", a.left, highlight.R, a.leftGeneration.Add(1))
	a.scheduleHighlight("right", a.right, highlight.Go, a.rightGeneration.Add(1))
	return a
}

func (a *App) Close() {
	if a.cancelConvert != nil {
		a.cancelConvert()
	}
	a.hl.Close()
	a.pool.Close()
}

func (a *App) Run() error {
	// Pin the Gio frame/event loop to a dedicated OS thread. app.Main still owns
	// the platform main thread. Heavy work never runs on this goroutine.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	platform.BoostGUIThread()
	defer a.Close()

	var ops op.Ops
	for {
		e := a.window.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			a.applyBackgroundResults()
			a.handleEditorEvents(gtx)
			a.handleClicks(gtx)
			a.layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

func newCodeEditor(th *material.Theme, readOnly bool, scheme syntax.ColorScheme) *gvcode.Editor {
	ed := gvwidget.NewEditor(th)
	ed.WithOptions(
		gvcode.WithFont(font.Font{Typeface: "monospace", Weight: font.Bold}),
		gvcode.WithTextSize(unit.Sp(14)),
		gvcode.WithLineHeight(0, 1.35),
		gvcode.WithTabWidth(4),
		gvcode.WithSoftTab(true),
		gvcode.WrapLine(false),
		gvcode.WithDefaultGutters(),
		gvcode.WithGutterGap(unit.Dp(10)),
		gvcode.WithCornerRadius(unit.Dp(3)),
		gvcode.WithColorScheme(scheme),
		gvcode.ReadOnlyMode(readOnly),
	)
	return ed
}

func codeColorScheme(colors, blackText bool) syntax.ColorScheme {
	c := syntax.ColorScheme{Name: "r2go-high-contrast-light"}
	foreground := "#000000"
	if !blackText {
		foreground = "#4B5563"
	}
	c.Foreground = mustColor(foreground + "FF")
	c.Background = mustColor("#FFFFFFFF")
	c.SelectColor = mustColor("#94C5FFFF")
	c.LineColor = mustColor("#E5F0FFFF")
	c.LineNumberColor = mustColor(foreground + "FF")
	if !colors {
		for _, scope := range []string{"keyword", "name.function", "name.builtin", "name.class", "literal.string", "literal.number", "comment", "operator", "punctuation"} {
			c.AddStyle(syntax.StyleScope(scope), 0, mustColor(foreground+"FF"), gvcolor.Color{})
		}
		return c
	}
	c.AddStyle("keyword", syntax.Bold, mustColor("#003CFFFF"), gvcolor.Color{})
	c.AddStyle("name.function", syntax.Bold, mustColor("#7A00CCFF"), gvcolor.Color{})
	c.AddStyle("name.builtin", syntax.Bold, mustColor("#007A3DFF"), gvcolor.Color{})
	c.AddStyle("name.class", syntax.Bold, mustColor("#A03A00FF"), gvcolor.Color{})
	c.AddStyle("literal.string", syntax.Bold, mustColor("#008000FF"), gvcolor.Color{})
	c.AddStyle("literal.number", syntax.Bold, mustColor("#B000B0FF"), gvcolor.Color{})
	c.AddStyle("comment", syntax.Bold, mustColor("#000000FF"), gvcolor.Color{})
	c.AddStyle("operator", syntax.Bold, mustColor("#D00020FF"), gvcolor.Color{})
	c.AddStyle("punctuation", syntax.Bold, mustColor("#000000FF"), gvcolor.Color{})
	return c
}

func mustColor(hex string) gvcolor.Color {
	c, err := gvcolor.Hex2Color(hex)
	if err != nil {
		panic(err)
	}
	return c
}

func (a *App) handleEditorEvents(gtx layout.Context) {
	for {
		evt, ok := a.left.Update(gtx)
		if !ok {
			break
		}
		if _, changed := evt.(gvcode.ChangeEvent); changed {
			gen := a.leftGeneration.Add(1)
			a.scheduleHighlight("left", a.left, highlight.R, gen)
		}
	}
	for {
		_, ok := a.right.Update(gtx)
		if !ok {
			break
		}
	}
}

func (a *App) scheduleHighlight(tag string, ed *gvcode.Editor, lang highlight.Language, generation uint64) {
	// GetReader is gvcode's concurrency-safe way to expose the text buffer to
	// another goroutine. Chroma never runs on the GUI thread.
	a.hl.Submit(highlight.Request{
		Tag:        tag,
		Generation: generation,
		Language:   lang,
		Reader:     ed.GetReader(),
	})
}

func (a *App) applyBackgroundResults() {
	for {
		select {
		case res := <-a.convertResults:
			if res.generation != a.convertGeneration.Load() {
				continue
			}
			a.busy = false
			// The engine intentionally returns readable generated/error output
			// together with an error. Never hide that diagnostic from the editor.
			if strings.TrimSpace(res.code) != "" {
				a.right.SetText(res.code)
				a.rightGeneration.Add(1)
				a.right.SetSyntaxTokens(res.tokens...)
			}
			if res.err != nil {
				a.status = "Convert failed: " + res.err.Error()
				continue
			}
			a.status = "Converted"
		case res := <-a.saveResults:
			if res.err != nil {
				a.status = "Save failed: " + res.err.Error()
			} else if res.path != "" {
				a.status = "Saved: " + res.path
			} else {
				a.status = "Save cancelled"
			}
		case res := <-a.hl.Results():
			if res.Err != nil {
				continue
			}
			switch res.Tag {
			case "left":
				if res.Generation == a.leftGeneration.Load() {
					a.left.SetSyntaxTokens(res.Tokens...)
				}
			case "right":
				if res.Generation == a.rightGeneration.Load() {
					a.right.SetSyntaxTokens(res.Tokens...)
				}
			}
		default:
			return
		}
	}
}

func (a *App) handleClicks(gtx layout.Context) {
	if a.convertBtn.Clicked(gtx) {
		a.startConvert()
	}
	if a.copyBtn.Clicked(gtx) {
		// gvcode's reader is the concurrency-friendly path and avoids building a
		// second giant string in the GUI frame just to feed the clipboard.
		gtx.Execute(clipboard.WriteCmd{
			Type: "text/plain",
			Data: io.NopCloser(a.right.GetReader()),
		})
		a.status = "Copied Go code"
	}
	if a.saveBtn.Clicked(gtx) {
		a.startSaveAs(a.right.GetReader())
	}
	if a.infoBtn.Clicked(gtx) {
		a.showInfo = !a.showInfo
	}
	if a.copyInfoBtn.Clicked(gtx) {
		gtx.Execute(clipboard.WriteCmd{Type: "text/plain", Data: io.NopCloser(strings.NewReader(cliHelp))})
		a.status = "Copied CLI commands"
	}
}

func (a *App) startConvert() {
	if a.cancelConvert != nil {
		a.cancelConvert()
	}
	ctx, cancel := context.WithCancel(context.Background())
	a.cancelConvert = cancel
	generation := a.convertGeneration.Add(1)
	options := appcore.TranspileOptions{
		AllowIRFallback:  a.allowFallback.Value,
		PreserveOriginal: a.preserveSource.Value,
	}
	// Do not materialize a potentially large R document as a new string on the
	// GUI event loop. The reader is consumed by the conversion worker.
	inputReader := a.left.GetReader()
	a.busy = true
	a.status = "Converting…"

	err := a.pool.Submit(func() {
		var (
			code   string
			tokens []syntax.Token
			runErr error
		)
		defer func() {
			if r := recover(); r != nil {
				runErr = fmt.Errorf("transpiler panic: %v", r)
			}
			select {
			case a.convertResults <- conversionResult{generation: generation, code: code, tokens: tokens, err: runErr}:
			default:
			}
			a.window.Invalidate()
		}()
		input, err := io.ReadAll(inputReader)
		if err != nil {
			runErr = fmt.Errorf("read R input: %w", err)
			return
		}
		if configurable, ok := a.engine.(appcore.ConfigurableEngine); ok {
			code, runErr = configurable.TranspileWithOptions(ctx, string(input), options)
		} else {
			code, runErr = a.engine.Transpile(ctx, string(input))
		}
		if strings.TrimSpace(code) != "" {
			highlighted, highlightErr := highlight.Tokens(ctx, highlight.Go, code)
			if highlightErr == nil {
				tokens = highlighted
			} else if runErr == nil {
				runErr = highlightErr
			}
		}
	})
	if err != nil {
		a.busy = false
		a.status = "Convert not started: " + err.Error()
	}
}

func (a *App) startSaveAs(reader io.Reader) {
	a.status = "Choose save location…"
	go func() {
		// Snapshot before the modal picker opens. Holding an editor-backed reader
		// across a long-running dialog allowed later editor updates to invalidate it.
		data, err := io.ReadAll(reader)
		path := ""
		if err == nil {
			path, err = platform.SaveGoFileDialog("output.go")
		}
		if err == nil && path != "" {
			err = standalone.WriteProgram(path, data)
		}
		select {
		case a.saveResults <- saveResult{path: path, err: err}:
		default:
		}
		a.window.Invalidate()
	}()
}

func (a *App) layout(gtx layout.Context) layout.Dimensions {
	if a.busy {
		// Keep frames flowing while a long conversion is running. The actual work
		// is outside the event loop; this only updates responsive visual feedback.
		gtx.Execute(op.InvalidateCmd{At: gtx.Now.Add(250 * time.Millisecond)})
	}
	return layout.Inset{Top: 12, Bottom: 10, Left: 12, Right: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(a.layoutHeader),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions { return layout.Spacer{Height: 10}.Layout(gtx) }),
			layout.Flexed(1, a.layoutMain),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if !a.showInfo {
					return layout.Dimensions{}
				}
				return a.layoutInfo(gtx)
			}),
			layout.Rigid(a.layoutFooter),
		)
	})
}

func (a *App) layoutHeader(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions { return layout.Dimensions{Size: gtx.Constraints.Min} }),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			label := material.Body2(a.theme, a.status)
			label.Alignment = text.End
			label.Color = color.NRGBA{R: 87, G: 96, B: 106, A: 255}
			return label.Layout(gtx)
		}),
	)
}

func (a *App) layoutMain(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return a.layoutEditorPanel(gtx, "R", a.left, nil)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: 14, Right: 14}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(a.theme, &a.convertBtn, "Convert")
					btn.Background = color.NRGBA{R: 9, G: 105, B: 218, A: 255}
					btn.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
					btn.CornerRadius = unit.Dp(8)
					btn.Inset = layout.Inset{Top: 12, Bottom: 12, Left: 20, Right: 20}
					return btn.Layout(gtx)
				})
			})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			actions := func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions { return smallButton(gtx, a.theme, &a.copyBtn, "Copy") }),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions { return layout.Spacer{Width: 8}.Layout(gtx) }),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions { return smallButton(gtx, a.theme, &a.saveBtn, "Save As") }),
				)
			}
			return a.layoutEditorPanel(gtx, "Go", a.right, actions)
		}),
	)
}

func (a *App) layoutEditorPanel(gtx layout.Context, title string, ed *gvcode.Editor, actions layout.Widget) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					label := material.Body1(a.theme, title)
					label.Font.Weight = font.SemiBold
					return label.Layout(gtx)
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions { return layout.Dimensions{Size: gtx.Constraints.Min} }),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if actions == nil {
						return layout.Dimensions{}
					}
					return actions(gtx)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions { return layout.Spacer{Height: 8}.Layout(gtx) }),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return widget.Border{
				Color: color.NRGBA{R: 208, G: 215, B: 222, A: 255},
				Width: unit.Dp(1), CornerRadius: unit.Dp(8),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: 6, Bottom: 6, Left: 6, Right: 6}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return ed.Layout(gtx, a.theme.Shaper)
				})
			})
		}),
	)
}

func (a *App) layoutInfo(gtx layout.Context) layout.Dimensions {
	return layout.Inset{Top: 10, Bottom: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return widget.Border{Color: color.NRGBA{R: 208, G: 215, B: 222, A: 255}, Width: 1, CornerRadius: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: 10, Bottom: 10, Left: 12, Right: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						label := material.Body2(a.theme, cliHelp)
						label.Color = color.NRGBA{R: 31, G: 35, B: 40, A: 255}
						return label.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return smallButton(gtx, a.theme, &a.copyInfoBtn, "Copy CLI commands")
						})
					}),
				)
			})
		})
	})
}

const cliHelp = `CLI - same r2go.exe as this GUI
  r2go run input.R                 execute R source with the Pure-Go runtime
  r2go ast input.R                 print the parsed R syntax tree as JSON
  r2go transpile input.R -o out.go generate a standalone Go main package
  r2go transpile --strict-native input.R -o out.go
                                    reject every compatibility fallback
  r2go transpile --source-comments=false input.R -o out.go
                                    omit original R comments from fallback output
  r2go coverage                    print primitive coverage metadata
  r2go version                     print the version
  r2go --licenses                  print embedded third-party notices
  r2go gui                         start this graphical editor

Examples
  r2go transpile analysis.R -o analysis.go
  go run analysis.go

Convert uses the same parser and generator as "r2go transpile". The Copy
button copies this complete command reference to the Windows clipboard.`

func (a *App) layoutFooter(gtx layout.Context) layout.Dimensions {
	return layout.Inset{Top: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		oldFallback, oldSource := a.allowFallback.Value, a.preserveSource.Value
		dims := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions { return smallButton(gtx, a.theme, &a.infoBtn, "Info") }),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							cb := material.CheckBox(a.theme, &a.allowFallback, "Pure-Go compatibility blocks (uncheck for strict native)")
							cb.TextSize = unit.Sp(12)
							return cb.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions { return layout.Spacer{Width: 18}.Layout(gtx) }),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							cb := material.CheckBox(a.theme, &a.preserveSource, "Keep original R comments")
							cb.TextSize = unit.Sp(12)
							return cb.Layout(gtx)
						}),
					)
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.Caption(a.theme, fmt.Sprintf("GUI thread pinned | worker budget %d of %d Ps", appcore.RecommendedWorkers(), runtime.GOMAXPROCS(0)))
				label.Color = color.NRGBA{R: 110, G: 118, B: 129, A: 255}
				return label.Layout(gtx)
			}),
		)
		if oldFallback != a.allowFallback.Value || oldSource != a.preserveSource.Value {
			a.status = "Transpiler options updated"
			a.window.Invalidate()
		}
		return dims
	})
}

func smallButton(gtx layout.Context, th *material.Theme, click *widget.Clickable, text string) layout.Dimensions {
	btn := material.Button(th, click, text)
	btn.Background = color.NRGBA{R: 236, G: 242, B: 248, A: 255}
	btn.Color = color.NRGBA{R: 31, G: 35, B: 40, A: 255}
	btn.CornerRadius = 7
	btn.TextSize = unit.Sp(12)
	btn.Inset = layout.Inset{Top: 6, Bottom: 6, Left: 11, Right: 11}
	return btn.Layout(gtx)
}
