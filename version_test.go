package commandergo

import "testing"

func TestVersionRegistersOption(t *testing.T) {
	cmd := New("app").Version("1.0.0")
	if !cmd._options.has("version") {
		t.Fatal("Version() should auto-register --version option")
	}
	opt := cmd._options["version"]
	if opt.alias != "V" {
		t.Fatalf("--version alias should be 'V', got %q", opt.alias)
	}
}

func TestNoVersionNoOption(t *testing.T) {
	cmd := New("app")
	if cmd._options.has("version") {
		t.Fatal("--version should not be registered without calling Version()")
	}
}
