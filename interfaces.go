package pkgcraft

type Iterator[T any] interface {
	HasNext() bool
	Next() T
}
