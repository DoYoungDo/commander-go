package commandergo

import (
	"strings"
	"testing"
)

func TestNewRegistersHelp(t *testing.T) {
	cmd := New("app")
	if !cmd._options.has("help") {
		t.Fatal("New() should auto-register --help option")
	}
	opt, _ := cmd._options.get("help")
	if opt.alias != "h" {
		t.Fatalf("--help alias should be 'h', got %q", opt.alias)
	}
}

func TestHelpText(t *testing.T) {
	cmd := New("todo").
		Description("A simple todo CLI").
		Arguments("[filter]", "filter todos", nil).
		Options("--all", "show all todos", false)
	cmd.Command("add [todo]", "add a new todo").
		Options("--force", "force add", false)
	cmd.Command("list", "list todos")

	help := cmd.helpText()

	cases := []string{
		"Usage: todo",
		"[options]",
		"[command]",
		"[filter]",
		"A simple todo CLI",
		"Arguments:",
		"filter",
		"Options:",
		"--all",
		"-h, --help",
		"Commands:",
		"add [options] [todo]",
		"list",
	}
	for _, want := range cases {
		if !strings.Contains(help, want) {
			t.Errorf("helpText() missing %q\ngot:\n%s", want, help)
		}
	}

	if strings.Contains(help, "list [options]") {
		t.Errorf("helpText() should not show default help option as subcommand options\ngot:\n%s", help)
	}
}
