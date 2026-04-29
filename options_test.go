package commandergo

import (
	"testing"
)

func TestNewOption(t *testing.T) {
	option, err := NewOption("--f")
	if err != nil {
		t.Fatalf("NewOption failed: %v", err)
	}
	if option.name != "f" {
		t.Fatalf("NewOption expect name f, but got: %v", option.name)
	}

	if option.alias != "" {
		t.Fatalf("NewOption expect alias empty, but got: %v", option.alias)
	}

	if option.valueName != "" {
		t.Fatalf("NewOption expect valueName empty, but got: %v", option.valueName)
	}

	if option.valueRequired {
		t.Fatalf("NewOption expect optional value, but parsed as required value")
	}
}
