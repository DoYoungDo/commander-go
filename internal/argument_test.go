package commandergo

import "testing"

func TestNewArgument(t *testing.T) {
	arg, err := NewArgument("[name]")
	if err != nil {
		panic(err)
	}
	if arg.name != "name" {
		t.Fatalf("invalid argument name :%v", arg.name)
	}
	if arg.multiValue {
		t.Fatalf("argument expect single value, but parsed as multi value")
	}
	if arg.valueRequired {
		t.Fatalf("argument expect optional value, but parsed as required value")
	}

	arg, err = NewArgument("<name...>")
	if err != nil {
		panic(err)
	}
	if arg.name != "name" {
		t.Fatalf("invalid argument name :%v", arg.name)
	}
	if !arg.multiValue {
		t.Fatalf("argument expect multi value, but parsed as single value")
	}
	if !arg.valueRequired {
		t.Fatalf("argument expect required value, but parsed as optional value")
	}
}
