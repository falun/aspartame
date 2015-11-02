package main // import "github.com/falun/aspartame"

import (
	"flag"
	"fmt"
	"go/ast"
	// "go/build"
	// "go/format"
	"go/parser"
	"go/token"
)

type Const struct {
	Type string
	Name string
}

func (c Const) String() string {
	return fmt.Sprintf("Const[%s, %s]", c.Name, c.Type)
}

type ConstBlock struct {
	Contents []Const
}

type File struct {
	Path   string
	Ast    *ast.File
	Consts []ConstBlock
}

func (f *File) String() string {
	str := "File {\n"
	str += fmt.Sprintf("  path: %s\n", f.Path)
	str += fmt.Sprintf("  ast: %p\n", f.Ast)
	str += fmt.Sprintf("  Consts {\n")
	for i, cb := range f.Consts {
		str += fmt.Sprintf("    Const block %d\n", (i + 1))
		for _, c := range cb.Contents {
			str += fmt.Sprintf("      %s\n", c)
		}
	}
	str += fmt.Sprintf("  }\n")
	str += "}"

	return str
}

func (f *File) parseConsts() {
	blocks := make([]ConstBlock, 0)

	for _, d := range f.Ast.Decls {
		decl, isDecl := d.(*ast.GenDecl)
		if isDecl && (decl.Tok == token.CONST) {
			// in const block; how many value specs
			curConstBlock := ConstBlock{}
			curConstBlock.Contents = make([]Const, 0)

			var curType *string = nil
			for _, spec := range decl.Specs {
				if vs, isValueSpec := spec.(*ast.ValueSpec); isValueSpec {
					if vs.Type != nil {
						if curType == nil {
							curType = new(string)
						}
						*curType = fmt.Sprintf("%s", vs.Type)
					}

					for i, _ := range vs.Names {
						curConstBlock.Contents = append(
							curConstBlock.Contents,
							Const{Type: *curType, Name: vs.Names[i].Name})
					}
				}
			}

			blocks = append(blocks, curConstBlock)
		}
	}

	f.Consts = blocks
}

func NewFile(path string) *File {
	fset := token.NewFileSet()
	fAst, err := parser.ParseFile(fset, path, nil, parser.ParseComments)

	if err != nil {
		fmt.Println(err)
	}

	f := File{Path: path, Ast: fAst}
	f.parseConsts()

	return &f
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		return
	}

	filePath := flag.Args()[0]

	f := NewFile(filePath)
	fmt.Println(f)
}
