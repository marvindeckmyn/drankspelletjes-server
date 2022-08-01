package types

func Opt[T interface{}](val T) *T {
	return &val
}

func Ptr[T interface{}](val T) *T {
	return &val
}
