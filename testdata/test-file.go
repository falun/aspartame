package testdata

//go:generate aspartame -name Foo -enumType FooEnumType -source $GOFILE
type FooEnumType int

const (
	bar FooEnumType = iota
	baz
	quix
	quux
)

const (
	other int = iota
	block
	of
	constValues
)
