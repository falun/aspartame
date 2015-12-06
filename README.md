### Aspartame

A tool that makes some go things (currently only "enums") slightly less annoying to me. Probably not idiomatic; check [TODO](TODO.md) for potential changes / ways to contribute.


##### Usage
`aspartame` may be run via command line:

	go get github.com/falun/aspartame
	go install github.com/falun/aspartame
	./bin/aspartame \
	   -name Foo \
	   -enumType FooEnumType \
	   -source ./src/github.com/falun/aspartame/testdata/

or as a `go:generate` tool:

	package main
	
	//go:generate aspartame -name Foo -enumType FooEnumType -source $GOFILE
	type FooEnumType int
	
	const (
		bar FooEnumType = iota
		baz
		quix
		quux
	)

Note that `-source $GOFILE` is optional but prevents the tool from guessing the wrong file since a generated file will be a potential match if we're given only the type name to target.

The tool checks for two things when examining a const block declaration and determining if it meets the requirement for sweetening.

1. The values are _not_ exported. The reasoning is that we want to limit confusion about what is available in the package namespace;
2. The enum value names do not begin with `_`;
3. The const block declares _only_ values of the same type.

Additionally if you want to change how an enum is represented as a string value you can add a `render` directive as a comment to the declaration. An example of that would be:

	const (
		bar FooEnumType = iota
		baz
		quix // render: QuiX
		quux
	)

Now when converting `Foo.Quix` to or from a string "QuiX" is the expected name.  The const as accessed in code will just be a capitalized version of the initial definition and the rendered value does not come into play. The play example below makes use of the render tag.

##### Results
The generated code will be produced in the same package and currently provides the following convenience methods:

* Enum value access via `$EnumName.$ValueName` (in our example above `Foo.Bar`, `Foo.Quix`, etc.)
* `$EnumName.Values()`&mdash;produces an array of all values
* `$EnumName.String()`&mdash;readable output
* `$EnumName.ByValue(int)`&mdash;given an int find the corresponding enum value
* `$EnumName.ByName(string)`&mdash;given an enum name produce the corresponding value

##### Example
See [this play link](http://play.golang.org/p/hKnt5OLNvp)

##### Limitations
I'm sure there are a lot but thus far it's good enough for my uses. File an issue of something doesn't work the way you expect.
