package commandergo

import "fmt"

type Varaint struct {
	value any
}

func (v Varaint) isString() bool {
	_, ok := v.value.(string)
	return ok
}
func (v Varaint) toString() string {
	str := v.value.(string)
	return str
}
func (v Varaint) isInt() bool {
	_, ok := v.value.(int)
	return ok
}
func (v Varaint) toInt() int {
	i := v.value.(int)
	return i
}
func (v Varaint) isBool() bool {
	_, ok := v.value.(bool)
	return ok
}
func (v Varaint) toBool() bool {
	b := v.value.(bool)
	return b
}
func (v Varaint) isFloat() bool {
	_, ok := v.value.(float64)
	return ok
}
func (v Varaint) toFloat() float64 {
	f := v.value.(float64)
	return f
}

// 导出方法供外部包使用
func (v Varaint) IsEmpty() bool {
	return v == Varaint{} || v.value == nil
}
func (v Varaint) IsString() bool   { return v.isString() }
func (v Varaint) ToString() string { return v.toString() }
func (v Varaint) IsInt() bool      { return v.isInt() }
func (v Varaint) ToInt() int       { return v.toInt() }
func (v Varaint) IsBool() bool     { return v.isBool() }
func (v Varaint) ToBool() bool     { return v.toBool() }
func (v Varaint) IsFloat() bool    { return v.isFloat() }
func (v Varaint) ToFloat() float64 { return v.toFloat() }
func (v Varaint) ForceToString() string {
	if v.isString() {
		return v.toString()
	}
	return fmt.Sprintf("%v", v.value)
}
