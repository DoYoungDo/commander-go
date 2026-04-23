package commandergo

import "testing"

func TestVaraint(t *testing.T) {
	cases := []struct {
		v       Varaint
		isStr   bool
		isInt   bool
		isBool  bool
		isFloat bool
	}{
		{Varaint{value: "hello"}, true, false, false, false},
		{Varaint{value: 42}, false, true, false, false},
		{Varaint{value: true}, false, false, true, false},
		{Varaint{value: 3.14}, false, false, false, true},
	}
	for _, c := range cases {
		if c.v.isString() != c.isStr {
			t.Errorf("isString() = %v, want %v for %v", c.v.isString(), c.isStr, c.v.value)
		}
		if c.v.isInt() != c.isInt {
			t.Errorf("isInt() = %v, want %v for %v", c.v.isInt(), c.isInt, c.v.value)
		}
		if c.v.isBool() != c.isBool {
			t.Errorf("isBool() = %v, want %v for %v", c.v.isBool(), c.isBool, c.v.value)
		}
		if c.v.isFloat() != c.isFloat {
			t.Errorf("isFloat() = %v, want %v for %v", c.v.isFloat(), c.isFloat, c.v.value)
		}
	}

	f := Varaint{value: 3.14}
	if f.toFloat() != 3.14 {
		t.Error("toFloat() failed")
	}
}
