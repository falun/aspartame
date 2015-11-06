### Aspartame

A `go:generate` tool that makes some go things (currently only "enums") slightly less annoying to me. Probably not idiomatic.


##### Usage
This describes eventual usage as the tooling doesn't actually parse command line args.

	package main
	
	// go:generate aspartame -target enum -name Foo -source FooEnumType
	type FooEnumType int
	
	const (
		bar FooEnumType = iota
		baz
		quix
		quux
	)

##### Results
The generated code will be produced in the same package and currently provides the following convenience methods:

* `Values()`&mdash;produces an array of all values
* `String()`&mdash;readable output
* `ByValue(int)`&mdash;given an int find the corresponding enum value
* `ByName(string)`&mdash;given an enum name produce the corresponding value

##### Example
See [this play link](http://play.golang.org/p/ZTbQzepKTY)
