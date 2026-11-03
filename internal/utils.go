package commandergo

func IfEl[T any](res bool, trueVal, falseVal T) T {
	if res {
		return trueVal
	}
	return falseVal
}
