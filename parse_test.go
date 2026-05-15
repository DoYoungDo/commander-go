package commandergo

import (
	"io"
	"os"
	"strings"
	"testing"
)

func captureStderr(t *testing.T, fn func()) string {
	t.Helper()

	old := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe stderr: %v", err)
	}
	os.Stderr = w

	fn()

	w.Close()
	os.Stderr = old

	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read stderr: %v", err)
	}
	return string(out)
}

func TestParseBoolOption(t *testing.T) {
	var called bool
	New("app").
		Options("--verbose", "verbose", false).
		Action(func(ctx *Context) {
			called = true
			if !ctx.Opt("verbose").toBool() {
				t.Error("expected verbose=true")
			}
		}).
		Parse([]string{"--verbose"})
	if !called {
		t.Error("action not called")
	}
}

func TestParseOptionWithValue(t *testing.T) {
	var called bool
	New("app").
		Options("-o, --output <file>", "output file", "").
		Action(func(ctx *Context) {
			called = true
			if ctx.Opt("output").toString() != "out.txt" {
				t.Errorf("expected out.txt, got %v", ctx.Opt("output").toString())
			}
		}).
		Parse([]string{"--output", "out.txt"})
	if !called {
		t.Error("action not called")
	}
}

func TestParseOptionWithInlineValue(t *testing.T) {
	var called bool
	New("app").
		Options("-o, --output <file>", "output file", "").
		Action(func(ctx *Context) {
			called = true
			if ctx.Opt("output").toString() != "out.txt" {
				t.Errorf("expected out.txt, got %v", ctx.Opt("output").toString())
			}
		}).
		Parse([]string{"--output=out.txt"})
	if !called {
		t.Error("action not called")
	}
}

func TestParseShortAlias(t *testing.T) {
	var called bool
	New("app").
		Options("-v, --verbose", "verbose", false).
		Action(func(ctx *Context) {
			called = true
			if !ctx.Opt("verbose").toBool() {
				t.Error("expected verbose=true")
			}
		}).
		Parse([]string{"-v"})
	if !called {
		t.Error("action not called")
	}
}

func TestParsePositionalArg(t *testing.T) {
	var called bool
	New("app").
		Arguments("[name]", "name", nil).
		Action(func(ctx *Context) {
			called = true
			if ctx.Args()[0].toString() != "hello" {
				t.Errorf("expected hello, got %v", ctx.Args()[0].toString())
			}
		}).
		Parse([]string{"hello"})
	if !called {
		t.Error("action not called")
	}
}

func TestParseSubCommand(t *testing.T) {
	var called bool
	app := New("app")
	app.Command("add", "add item").
		Arguments("[item]", "item", nil).
		Action(func(ctx *Context) {
			called = true
			if ctx.Args()[0].toString() != "todo1" {
				t.Errorf("expected todo1, got %v", ctx.Args()[0].toString())
			}
		})
	app.Parse([]string{"add", "todo1"})
	if !called {
		t.Error("sub command action not called")
	}
}

func TestParseRequiredArgMissing(t *testing.T) {
	err := New("app").
		Arguments("<name>", "name", nil).
		Parse([]string{})
	if err == nil {
		t.Error("expected error for missing required argument")
	}
}

func TestParseUnknownOption(t *testing.T) {
	// 未知选项应该 warning 并继续，不返回 error，action 正常调用
	var called bool
	var err error
	stderr := captureStderr(t, func() {
		err = New("app").
			Action(func(ctx *Context) { called = true }).
			Parse([]string{"--unknown"})
	})
	if err != nil {
		t.Errorf("unknown option should warn not error, got: %v", err)
	}
	if !strings.Contains(stderr, "warning: unknown option: --unknown") {
		t.Errorf("expected unknown option warning, got: %q", stderr)
	}
	if !called {
		t.Error("action should still be called after unknown option warning")
	}
}

func TestParseValueTypeInference(t *testing.T) {
	var called bool
	New("app").
		Options("-n, --num <n>", "number", 0).
		Action(func(ctx *Context) {
			called = true
			if ctx.Opt("num").toInt() != 42 {
				t.Errorf("expected 42, got %v", ctx.Opt("num").toInt())
			}
		}).
		Parse([]string{"--num", "42"})
	if !called {
		t.Error("action not called")
	}
}

func TestParseFloatValueInference(t *testing.T) {
	var called bool
	New("app").
		Options("-r, --rate <n>", "rate", 0.0).
		Action(func(ctx *Context) {
			called = true
			if ctx.Opt("rate").toFloat() != 1.5 {
				t.Errorf("expected 1.5, got %v", ctx.Opt("rate").toFloat())
			}
		}).
		Parse([]string{"--rate", "1.5"})
	if !called {
		t.Error("action not called")
	}
}

func TestParseBoolValueInference(t *testing.T) {
	var called bool
	New("app").
		Options("-d, --debug <flag>", "debug flag", false).
		Action(func(ctx *Context) {
			called = true
			if !ctx.Opt("debug").isBool() {
				t.Error("expected bool value")
			}
		}).
		Parse([]string{"--debug", "true"})
	if !called {
		t.Error("action not called")
	}
}

func TestParseMultiAlias(t *testing.T) {
	var called bool
	New("app").
		Options("-v, --verbose", "verbose", false).
		Options("-d, --debug", "debug", false).
		Action(func(ctx *Context) {
			called = true
			if !ctx.Opt("verbose").toBool() {
				t.Error("expected verbose=true")
			}
			if !ctx.Opt("debug").toBool() {
				t.Error("expected debug=true")
			}
		}).
		Parse([]string{"-vd"})
	if !called {
		t.Error("action not called")
	}
}

func TestParseVersionOutput(t *testing.T) {
	err := New("app").Version("1.2.3").Parse([]string{"--version"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestParseHelpOutput(t *testing.T) {
	err := New("app").Description("test app").Parse([]string{"--help"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestParseMissingRequiredOptionValue(t *testing.T) {
	err := New("app").
		Options("-o, --output <file>", "output", nil).
		Parse([]string{"--output"})
	if err == nil {
		t.Error("expected error for missing required option value")
	}
}

func TestParseStrictPositionalArg(t *testing.T) {
	var called bool
	err := New("app").
		Arguments("[name]", "name", nil).
		Action(func(ctx *Context) {
			called = true
			if ctx.Args()[0].toString() != "bob" {
				t.Errorf("expected bob, got %v", ctx.Args()[0].toString())
			}
		}).
		ParseStrict([]string{"bob"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !called {
		t.Error("action not called")
	}
}

func TestParseStrictUnexpectedArg(t *testing.T) {
	err := New("app").
		Arguments("[name]", "name", nil).
		ParseStrict([]string{"bob", "extra"})
	if err == nil {
		t.Error("expected error for unexpected argument")
	}
}

func TestParseStrictOptionWithPositionalArg(t *testing.T) {
	var called bool
	err := New("app").
		Options("-o, --output <file>", "output file", nil).
		Arguments("[name]", "name", nil).
		Action(func(ctx *Context) {
			called = true
			if ctx.Opt("output").toString() != "out.txt" {
				t.Errorf("expected out.txt, got %v", ctx.Opt("output").toString())
			}
			if ctx.Args()[0].toString() != "bob" {
				t.Errorf("expected bob, got %v", ctx.Args()[0].toString())
			}
		}).
		ParseStrict([]string{"--output", "out.txt", "bob"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !called {
		t.Error("action not called")
	}
}

func TestParseStrictUnknownOption(t *testing.T) {
	err := New("app").ParseStrict([]string{"--unknown"})
	if err == nil {
		t.Error("expected error for unknown option")
	}
}

func TestParseStrictUnknownAlias(t *testing.T) {
	err := New("app").ParseStrict([]string{"-x"})
	if err == nil {
		t.Error("expected error for unknown alias")
	}
}

func TestParseStrictFlagEndTreatsOptionAsArg(t *testing.T) {
	var called bool
	err := New("app").
		Arguments("[name]", "name", nil).
		Action(func(ctx *Context) {
			called = true
			if ctx.Args()[0].toString() != "--unknown" {
				t.Errorf("expected --unknown, got %v", ctx.Args()[0].toString())
			}
		}).
		ParseStrict([]string{"--", "--unknown"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !called {
		t.Error("action not called")
	}
}

func TestParseStrictHelpSkipsRequiredArg(t *testing.T) {
	err := New("app").
		Arguments("<name>", "name", nil).
		ParseStrict([]string{"--help"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestParseStrictHelpRejectsUnexpectedArg(t *testing.T) {
	err := New("app").ParseStrict([]string{"--help", "extra"})
	if err == nil {
		t.Error("expected error for unexpected argument")
	}
}
