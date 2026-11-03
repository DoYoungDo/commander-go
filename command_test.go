package commandergo

import (
	"fmt"
	"testing"
)

func TestNewCommand(t *testing.T) {
	cmd, err := newCommandWithNameAndArg("cmd [arg...]")
	if err != nil {
		t.Fatal(err)
	}
	if cmd.name != "cmd" {
		t.Fatalf("command name is %s, want cmd", cmd.name)
	}
	if !cmd._arguments.has("arg") {
		t.Fatalf("argument arg already exists")
	}
	fmt.Println(cmd, err)
}
