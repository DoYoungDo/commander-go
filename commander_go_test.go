package commandergo

import "testing"

func TestCommander(t *testing.T) {
	New("todo").
		Description("description").
		Version("0.0.1").
		Arguments("todo", "", nil).
		Command("add", "add todo").
		Arguments("todo", "todo", nil).
		Parent().
		Action(func(cmd *Command) {

		}).
		Parse()
}
