### Aspartame

A tool that makes some go things (currently only "enums") slightly less annoying to me. Probably not idiomatic.


##### Usage
Currently `aspartame` needs to be run as a command line tool as it only produces generated code to STDOUT:

	go get github.com/falun/aspartame
	go install github.com/falun/aspartame
	./bin/aspartame \
	   -name Foo \
	   -enumType FooEnumType \
	   -source ./src/github.com/falun/aspartame/testdata/

Eventually you should be able to use it as a standard `go:generate` tool:

	package main
	
	// go:generate aspartame -target enum -name Foo -source FooEnumType
	type FooEnumType int
	
	const (
		bar FooEnumType = iota
		baz
		quix
		quux
	)

The tool checks for two things when examining a const block declaration and determining if it meets the requirement for sweetening.

1. The values are _not_ exported. The reasoning is that we want to limit confusion about what is available in the package namespace.
2. The const block declares _only_ values of the same type.

##### Results
The generated code will be produced in the same package and currently provides the following convenience methods:

* `Values()`&mdash;produces an array of all values
* `String()`&mdash;readable output
* `ByValue(int)`&mdash;given an int find the corresponding enum value
* `ByName(string)`&mdash;given an enum name produce the corresponding value
* Enum value access via `$EnumName.$ValueName` (in our example above `Foo.Bar`, `Foo.Quix`, etc.)

##### Example
See [this play link](http://bit.ly/1RDyktS)

##### Limitations
They are legion. High on the list though is that `aspartame` currently only supports int-typed enums and doesn't parse out the values. Adding that should be pretty simple but didn't make the first cut.
