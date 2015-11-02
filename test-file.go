package main

// go:generate aspartame -enum Bar -type BarEnumType
type BarEnumType int

const (
	baz BarEnumType = 1234 + iota
	quix
	quux
)

const (
	other int = iota
	block
	of
	constValues
)
