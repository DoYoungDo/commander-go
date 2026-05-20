package commandergo

import "testing"

func TestNewFromJsonBuildsCommandTree(t *testing.T) {
	cmd := NewFromJson(`{
		"name": "todo",
		"version": "1.0.0",
		"description": "todo app",
		"arguments": [
			{
				"name": "filter",
				"desc": "filter todos",
				"defaultValue": "all",
				"valueRequired": false,
				"multiValue": false
			}
		],
		"options": [
			{
				"name": "output",
				"alias": "o",
				"desc": "output file",
				"value": {
					"valueName": "file",
					"defaultValue": "stdout",
					"valueRequired": true
				}
			}
		],
		"commands": [
			{
				"name": "add",
				"description": "add todo",
				"arguments": [
					{
						"name": "item",
						"valueRequired": true,
						"multiValue": false
					}
				]
			}
		]
	}`)

	if cmd.name != "todo" {
		t.Fatalf("expected command name todo, got %q", cmd.name)
	}
	if cmd.version != "1.0.0" {
		t.Fatalf("expected version 1.0.0, got %q", cmd.version)
	}
	if cmd.description != "todo app" {
		t.Fatalf("expected description todo app, got %q", cmd.description)
	}
	if len(cmd._arguments) != 1 || cmd._arguments[0].name != "filter" {
		t.Fatalf("expected root argument filter, got %+v", cmd._arguments)
	}
	opt, ok := cmd._options.get("output")
	if !ok {
		t.Fatal("expected output option to be registered")
	}
	if opt.alias != "o" {
		t.Fatalf("expected output alias o, got %q", opt.alias)
	}
	if opt.valueName != "file" || !opt.valueRequired {
		t.Fatalf("expected output option to require <file>, got valueName=%q required=%v", opt.valueName, opt.valueRequired)
	}
	child := cmd.Child("add")
	if child == nil {
		t.Fatal("expected add subcommand to exist")
	}
	if child.description != "add todo" {
		t.Fatalf("expected child description add todo, got %q", child.description)
	}
	if len(child._arguments) != 1 || child._arguments[0].name != "item" || !child._arguments[0].valueRequired {
		t.Fatalf("expected add subcommand required item argument, got %+v", child._arguments)
	}
}

func TestNewFromJsonNormalizesNumberDefaultsAndAppliesArgumentDefault(t *testing.T) {
	var called bool
	cmd := NewFromJson(`{
		"name": "app",
		"arguments": [
			{
				"name": "count",
				"defaultValue": 3,
				"valueRequired": false,
				"multiValue": false
			}
		],
		"options": [
			{
				"name": "rate",
				"alias": "r",
				"value": {
					"valueName": "n",
					"defaultValue": 1.5,
					"valueRequired": false
				}
			}
		]
	}`)

	cmd.Action(func(ctx *Context) {
		called = true
		if len(ctx.Args()) != 1 {
			t.Fatalf("expected 1 argument, got %d", len(ctx.Args()))
		}
		if !ctx.Args()[0].IsInt() {
			t.Fatalf("expected default argument to be int, got %#v", ctx.Args()[0])
		}
		if ctx.Args()[0].ToInt() != 3 {
			t.Fatalf("expected default argument value 3, got %d", ctx.Args()[0].ToInt())
		}
		if !ctx.Opt("rate").IsFloat() {
			t.Fatalf("expected default option to be float, got %#v", ctx.Opt("rate"))
		}
		if ctx.Opt("rate").ToFloat() != 1.5 {
			t.Fatalf("expected default option value 1.5, got %v", ctx.Opt("rate").ToFloat())
		}
	})

	if err := cmd.Parse([]string{"-r"}); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if !called {
		t.Fatal("expected action to be called")
	}
}

func TestNewFromJsonAllowsMissingOptionalFields(t *testing.T) {
	cmd := NewFromJson(`{
		"name": "app",
		"arguments": [
			{"name": "input"}
		],
		"options": [
			{"name": "verbose"}
		],
		"commands": [
			{"name": "list"}
		]
	}`)

	if cmd.name != "app" {
		t.Fatalf("expected command name app, got %q", cmd.name)
	}
	if cmd.version != "" {
		t.Fatalf("expected empty version, got %q", cmd.version)
	}
	if cmd.description != "" {
		t.Fatalf("expected empty description, got %q", cmd.description)
	}
	if len(cmd._arguments) != 1 {
		t.Fatalf("expected 1 argument, got %d", len(cmd._arguments))
	}
	arg := cmd._arguments[0]
	if arg.name != "input" || arg.desc != "" || !arg.defaultValue.IsEmpty() || arg.valueRequired || arg.multiValue {
		t.Fatalf("expected optional input argument with zero values, got %+v", arg)
	}
	opt, ok := cmd._options.get("verbose")
	if !ok {
		t.Fatal("expected verbose option to be registered")
	}
	if opt.alias != "" || opt.desc != "" || opt.valueName != "" || opt.valueRequired || !opt.defaultValue.IsEmpty() {
		t.Fatalf("expected bool verbose option with zero values, got %+v", opt)
	}
	if cmd.Child("list") == nil {
		t.Fatal("expected list subcommand to exist")
	}
	if len(cmd.Child("list")._arguments) != 0 || len(cmd.Child("list")._subCommands) != 0 {
		t.Fatalf("expected list subcommand to have no arguments or children, got %+v", cmd.Child("list"))
	}
}

func TestNewFromJsonWithEmptyCollections(t *testing.T) {
	cmd := NewFromJson(`{
		"name": "app",
		"arguments": [],
		"options": [],
		"commands": []
	}`)

	if len(cmd._arguments) != 0 {
		t.Fatalf("expected no arguments, got %d", len(cmd._arguments))
	}
	if len(cmd._subCommands) != 0 {
		t.Fatalf("expected no subcommands, got %d", len(cmd._subCommands))
	}
	if !cmd._options.has("help") {
		t.Fatal("expected default help option to remain registered")
	}
	if len(cmd._options) != 1 {
		t.Fatalf("expected only default help option, got %d options", len(cmd._options))
	}
}

func TestNewFromJsonBoolOptionWithMissingValue(t *testing.T) {
	var called bool
	cmd := NewFromJson(`{
		"name": "app",
		"options": [
			{"name": "verbose", "alias": "v"}
		]
	}`)

	cmd.Action(func(ctx *Context) {
		called = true
		if !ctx.Opt("verbose").IsBool() {
			t.Fatalf("expected verbose option to parse as bool, got %#v", ctx.Opt("verbose"))
		}
		if !ctx.Opt("verbose").ToBool() {
			t.Fatal("expected verbose option to be true")
		}
	})

	if err := cmd.ParseStrict([]string{"-v"}); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if !called {
		t.Fatal("expected action to be called")
	}
}

func TestNewFromJsonPanicsOnInvalidJson(t *testing.T) {
	expectPanic(t, func() {
		NewFromJson(`{"name":`)
	})
}

func TestNewFromJsonPanicsWhenRequiredFieldsAreMissing(t *testing.T) {
	tests := []struct {
		name string
		json string
	}{
		{
			name: "argument name",
			json: `{
				"name": "app",
				"arguments": [{}]
			}`,
		},
		{
			name: "option name",
			json: `{
				"name": "app",
				"options": [{}]
			}`,
		},
		{
			name: "option value name",
			json: `{
				"name": "app",
				"options": [
					{"name": "count", "value": {"valueRequired": true}}
				]
			}`,
		},
		{
			name: "subcommand name",
			json: `{
				"name": "app",
				"commands": [{}]
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectPanic(t, func() {
				NewFromJson(tt.json)
			})
		})
	}
}

func expectPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	fn()
}
