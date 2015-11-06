package types

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

// Holds the information we need about a given Const
type Const struct {
	Type string
	Name string
}

func (c Const) String() string {
	return fmt.Sprintf("Const[%s, %s]", c.Name, c.Type)
}

// Logical grouping of consts within the file; this also represents the
// range over which an iota has the ability to define increasing values
type ConstBlock struct {
	Contents []Const
}

func (cb *ConstBlock) Length() int {
	return len(cb.Contents)
}

func (cb *ConstBlock) IsHomogenous() bool {
	if cb.Length() == 0 {
		return true
	}

	t := cb.Contents[0].Type
	for _, v := range cb.Contents {
		if t != v.Type {
			return false
		}
	}

	return true
}

var typeErrMany error = errors.New("ConstBlock has inconsistent Type")
var typeErrNone error = errors.New("ConstBlock has no contents (thus no type)")

func (cb *ConstBlock) Type() (string, error) {
	if cb.IsHomogenous() {
		if cb.Length() == 0 {
			return "", typeErrNone
		} else {
			return cb.Contents[0].Type, nil
		}
	} else {
		return "", typeErrMany
	}
}

type File struct {
	Path    string
	Package string
	Ast     *ast.File
	Consts  []ConstBlock
}

func (f *File) String() string {
	str := "File {\n"
	str += fmt.Sprintf("  path: %s\n", f.Path)
	str += fmt.Sprintf("  package: %s\n", f.Package)
	str += fmt.Sprintf("  ast: %p\n", f.Ast)
	str += fmt.Sprintf("  consts {\n")
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
			curConstBlock := ConstBlock{}
			curConstBlock.Contents = make([]Const, 0)

			curType := ""
			for _, spec := range decl.Specs {

				if vs, isValueSpec := spec.(*ast.ValueSpec); isValueSpec {
					if vs.Type != nil {
						// found a new type
						curType = fmt.Sprintf("%s", vs.Type)
					}

					for _, n := range vs.Names {
						curConstBlock.Contents = append(
							curConstBlock.Contents,
							Const{Type: curType, Name: n.Name})
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
		return nil
	}

	pkg := ""
	if fn := fAst.Name; fn != nil {
		pkg = fn.Name
	}

	f := File{
		Path:    path,
		Package: pkg,
		Ast:     fAst,
	}
	f.parseConsts()

	return &f
}
