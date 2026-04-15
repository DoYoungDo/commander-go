package commandergo

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
