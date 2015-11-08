package testdata

// go:generate aspartame -target enum -name Foo -source FooEnumType
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
