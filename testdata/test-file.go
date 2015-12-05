package testdata

//go:generate aspartame -name Foo -enumType FooEnumType -source $GOFILE
type FooEnumType int

const (
	bar  FooEnumType = 1<<10 + iota //render:BAZ
	baz              = 3 + 4 + 5
	quix             = 1234 + iota
	quux
)

const (
	other int = iota
	block
	of
	constValues
)

const (
	mixed      FooEnumType = iota
	collection int         = 23
)
