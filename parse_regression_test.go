package commandergo

import (
	"reflect"
	"testing"
)

func TestParseExistingSemantics(t *testing.T) {
	for _, strict := range []bool{false, true} {
		mode := "normal"
		if strict {
			mode = "strict"
		}
		t.Run(mode, func(t *testing.T) {
			parse := func(cmd *Command, args []string) error {
				if strict {
					return cmd.ParseStrict(args)
				}
				return cmd.Parse(args)
			}
			t.Run("option_presence", func(t *testing.T) {
				cmd := New("app").Options("-n, --number <value>", "", 42).
					Options("--enabled", "", true).Options("--disabled", "", false).
					Options("--unset", "", nil)
				for _, args := range [][]string{nil, {"--number", "7", "--disabled"}, {"-n"}, {"--number"}, nil} {
					called := false
					cmd.Action(func(ctx *Context) {
						called = true
						want := map[string]Varaint{}
						if len(args) > 0 {
							want["number"] = Varaint{value: 42}
						}
						if len(args) > 1 {
							want["number"] = Varaint{value: 7}
							want["disabled"] = Varaint{value: true}
						}
						if !reflect.DeepEqual(ctx.parsedOpts, want) {
							t.Errorf("options=%v, want %v", ctx.parsedOpts, want)
						}
						for _, name := range []string{"number", "enabled", "disabled", "unset", "help", "unknown"} {
							if ctx.Opt(name) != want[name] {
								t.Errorf("%s=%v, want %v", name, ctx.Opt(name), want[name])
							}
						}
					})
					if err := parse(cmd, args); err != nil || !called {
						t.Fatalf("args=%v err=%v called=%v", args, err, called)
					}
				}
			})
			t.Run("immediate_child_dispatch", func(t *testing.T) {
				for _, args := range [][]string{{"run"}, {"--help", "run"}, {"--version", "run"}, {"one", "two", "run"}} {
					called := false
					cmd := New("app").Version("1.2.3").Arguments("<source>", "", nil).
						Action(func(*Context) { t.Error("unexpected parent action") })
					cmd.Command("run", "").Action(func(*Context) { called = true })
					if err := parse(cmd, args); err != nil || !called {
						t.Fatalf("args=%v err=%v called=%v", args, err, called)
					}
				}
			})
			t.Run("tokens", func(t *testing.T) {
				for _, tc := range []struct {
					args        []string
					value       any
					positionals []Varaint
				}{
					{[]string{"--value="}, nil, []Varaint{{value: "--value="}}},
					{[]string{"-v="}, nil, []Varaint{{value: "-v="}}},
					{[]string{"--flag="}, nil, []Varaint{{value: "--flag="}}},
					{[]string{"--unknown1", "--bad_name", "-?"}, nil, []Varaint{{value: "--unknown1"}, {value: "--bad_name"}, {value: "-?"}}},
					{[]string{"--value", "--unknown1"}, "--unknown1", []Varaint{}},
					{[]string{"-v", "--bad_name"}, "--bad_name", []Varaint{}},
					{[]string{"--", "--", "--flag"}, nil, []Varaint{{value: "--flag"}}},
				} {
					called := false
					cmd := New("app").Options("-v, --value [value]", "", "fallback").
						Options("--flag", "", nil).Arguments("[items...]", "", nil).
						Action(func(ctx *Context) {
							called = true
							if ctx.Opt("value").value != tc.value || !reflect.DeepEqual(ctx.Args(), tc.positionals) {
								t.Errorf("input=%v value=%v args=%v", tc.args, ctx.Opt("value"), ctx.Args())
							}
						})
					if err := parse(cmd, tc.args); err != nil || !called {
						t.Fatalf("args=%v err=%v called=%v", tc.args, err, called)
					}
				}
			})
			t.Run("positional_defaults", func(t *testing.T) {
				called := false
				cmd := New("app").Arguments("<source>", "", nil).
					Arguments("[destination]", "", "backup").Action(func(ctx *Context) {
					called = true
					if !reflect.DeepEqual(ctx.Args(), []Varaint{{value: "backup"}}) {
						t.Errorf("args=%v", ctx.Args())
					}
				})
				if err := parse(cmd, nil); err != nil || !called {
					t.Fatalf("err=%v called=%v", err, called)
				}
			})
		})
	}
}
