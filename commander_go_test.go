package commandergo

import "testing"

func TestCommander(t *testing.T) {
	New("todo").
		Description("description").
		Version("0.0.1").
		Arguments("[todo]", "", nil).
		Command("add", "add todo").
		Arguments("[todo]", "todo", nil).
		Parent().
		Action(func(ctx *Context) {

		}).
		Parse([]string{})
}

func TestChild(t *testing.T) {
	app := New("app")
	add := app.Command("add", "add item")
	tag := add.Command("tag", "tag item")
	create := tag.Command("create", "create tag")

	if app.Child("add") != add {
		t.Fatal("expected child command add")
	}

	if add.Child("tag") != tag {
		t.Fatal("expected child command tag")
	}

	if tag.Child("create") != create {
		t.Fatal("expected child command create")
	}

	if app.Child("tag") != nil {
		t.Fatal("expected nil when looking up grandchild from root command")
	}

	if add.Child("create") != nil {
		t.Fatal("expected nil when looking up grandchild from parent command")
	}

	if app.Child("missing") != nil {
		t.Fatal("expected nil for missing child command")
	}
}

func TestChildChain(t *testing.T) {
	app := New("app")
	create := app.
		Command("add", "add item").
		Command("tag", "tag item").
		Command("create", "create tag")

	child := app.Child("add").Child("tag").Child("create")
	if child != create {
		t.Fatal("expected chained child lookup to return create command")
	}
}

func TestParent(t *testing.T) {
	app := New("app")
	add := app.Command("add", "add item")
	tag := add.Command("tag", "tag item")
	create := tag.Command("create", "create tag")

	if app.Parent() != nil {
		t.Fatal("expected root command parent to be nil")
	}

	if add.Parent() != app {
		t.Fatal("expected add parent to be app")
	}

	if tag.Parent() != add {
		t.Fatal("expected tag parent to be add")
	}

	if create.Parent() != tag {
		t.Fatal("expected create parent to be tag")
	}

	if create.Parent().Parent().Parent() != app {
		t.Fatal("expected create to trace back to app through parents")
	}
}
