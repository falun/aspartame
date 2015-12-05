package generators

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/falun/aspartame/types"
)

var EnumGenerator *EnumGeneratorT = &EnumGeneratorT{}

type EnumGeneratorT struct {
	enumName string
	enumType string
}

func (eg *EnumGeneratorT) SetupFlags() {
	flag.StringVar(&(eg.enumName), "name", "", "[required:enum] What name should we export the sweetened enum as")
	flag.StringVar(&(eg.enumType), "enumType", "", "[required:enum] The type of the enum we'll be sweetening")
}

func (eg *EnumGeneratorT) LocateFile(inputPath string) *types.File {
	if !IsDirectory(inputPath) {
		f := types.NewFile(inputPath)
		if !eg.validate(f) {
			return nil
		}
		return f
	} else {
		pkg, err := ParseDir(inputPath)
		if err != nil {
			log.Fatal(err)
		}

		baseDir := pkg.Dir
		var possibilities []string
		for _, presumptiveFile := range pkg.GoFiles {
			possibilities = append(possibilities, filepath.Join(baseDir, presumptiveFile))
		}

		for _, path := range possibilities {
			f := types.NewFile(path)
			if eg.validate(f) {
				return f
			}
		}

		return nil
	}
}

func (eg *EnumGeneratorT) validate(f *types.File) bool {
	for _, consts := range f.Consts {
		t, err := consts.Type()

		if err == nil && t == eg.enumType && !consts.HasExported() {
			// check for leading _
			for _, c := range consts.Contents {
				if c.Name[0] == '_' {
					return false
				}
			}
			return true
		}
	}

	return false
}

func (eg *EnumGeneratorT) DoGenerate(source *types.File, dest *string) {
	err := false

	if eg.enumType == "" {
		fmt.Println("-enumType must be specified")
		err = true
	}

	if eg.enumName == "" {
		fmt.Println("-name must be specified")
		err = true
	}

	if err {
		return
	}

	GenerateEnum(source, eg.enumName, eg.enumType, dest)
}

var enumTemplate string = `package {{ .Package }}

import (
	"errors"
	"fmt"
)

const ({{ range $e := .Elements }}
	_{{ $e.Name }} {{ if eq $e.Index 0 }}{{ $e.Type }}{{ end }}{{ if eq $e.Value "" }}{{ if eq $e.Index 0}} = iota{{ end }}{{else}} = {{ $e.Value }}{{ end }}{{ end }}
)

func (v {{ .EnumType }}) String() string {
	switch v { {{ range $e := .Elements }}
		case _{{ $e.Name }}: return "{{ $e.RenderName }}"{{end}}
		default: return fmt.Sprintf("Unknown(%d)", int(v))
	}
}

type {{ .EnumName }}Container struct { {{ range $e := .Elements }}
	{{ $e.Name | Cap }} {{ $e.Type }}{{ end }}
}

// uncertain if I want to return pointer or not
func (this *{{ .EnumName }}Container) ByValue(v int) (*{{ .EnumType }}, error) {
	value := new({{ .EnumType }})

	switch v { {{ range $e := .Elements }}
		case int(_{{ $e.Name }}): *value = _{{ $e.Name }}{{end}}
		default: return nil, errors.New(fmt.Sprintf("Unable to find {{ .EnumType }} associated with value %d", v))
	}

	return value, nil
}

// uncertain if I want to return pointer or not
func (this *{{ .EnumName }}Container) ByName(s string) (*{{ .EnumType }}, error) {
	var value *{{ .EnumType }} = new({{ .EnumType }})

	switch s { {{ range $e := .Elements }}
		case "{{ $e.RenderName }}": *value = _{{ $e.Name }}{{end}}
		default: return nil, errors.New(fmt.Sprintf("Unable to find {{ .EnumType }} associated by name %s", s))
	}

	return value, nil
}

var valArray []{{ .EnumType }} = []{{ .EnumType }}{ {{ range $e := .Elements }}
	_{{ $e.Name }},{{ end }}
}

func (this *{{ .EnumName }}Container) Values() []{{ .EnumType }} {
	return valArray
}

var {{ .EnumName }} = &{{ .EnumName }}Container{ {{ range $e := .Elements }}
	{{ $e.Name | Cap }}: _{{ $e.Name }},{{ end }}
}
`

type ConstItem struct {
	Index      int
	Name       string
	RenderName string
	Type       string
	Value      string
}

type EnumData struct {
	Package  string
	EnumName string
	EnumType string
	Elements []ConstItem
}

func constBlockToEnumData(enumName string, f *types.File, cb *types.ConstBlock) *EnumData {
	ed := &EnumData{
		Package:  f.Package,
		EnumName: enumName,
		Elements: make([]ConstItem, 0),
	}

	for i, v := range cb.Contents {
		renderName := mkCap(v.Name)
		portions := strings.Split(v.Comment, ",")

		for _, p := range portions {
			if segments := strings.Split(p, ":"); len(segments) == 2 {
				if k := strings.Trim(segments[0], " \n"); k == "render" {
					renderName = strings.Trim(segments[1], " \n")
					break
				}
			}
		}

		ci := ConstItem{
			Index:      i,
			RenderName: renderName,
			Name:       v.Name,
			Type:       v.Type,
			Value:      v.Value,
		}
		ed.Elements = append(ed.Elements, ci)
	}

	ed.EnumType = ed.Elements[0].Type

	return ed
}

func mkCap(s string) string {
	return fmt.Sprintf("%s%s", strings.ToUpper(string(s[0])), s[1:])
}

func GenerateEnum(
	source *types.File,
	enumName string,
	enumType string,
	dest *string,
) {
	funcMap := template.FuncMap{
		"Cap": mkCap,
	}

	templateParsed, templateError := template.New("enum").Funcs(funcMap).Parse(enumTemplate)

	if templateError != nil {
		fmt.Printf(templateError.Error())
		return
	}

	destWriter := os.Stdout
	if dest != nil {
		basedir := filepath.Dir(source.Path)
		path := ""
		switch *dest {
		case "":
			path = filepath.Join(basedir, fmt.Sprintf("%s_enum.go", strings.ToLower(enumName)))
		default:
			path = filepath.Join(basedir, *dest)
		}
		fPtr, fErr := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)

		if fErr != nil {
			log.Fatal(fErr)
		}
		defer fPtr.Close()

		destWriter = fPtr
	}

	var sourceConsts *types.ConstBlock = nil

	for _, v := range source.Consts {
		enumT, tErr := v.Type()
		if tErr == nil && enumT == enumType {
			sourceConsts = &v
			break
		}
	}

	if sourceConsts != nil {
		templateParsed.Execute(destWriter, constBlockToEnumData(enumName, source, sourceConsts))
	} else {
		log.Fatal(fmt.Sprintf("Could not find enum source: %s", source.Path))
	}
}
