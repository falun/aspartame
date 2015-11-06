package generators

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/falun/aspartame/types"
)

var enumTemplate string = `package {{ .Package }}

import (
  "errors"
  "fmt"
)

const ({{ range $e := .Elements }}
    _{{ $e.Name }}{{if eq 0 $e.Index }} {{ $e.Type }}{{end}} {{ if eq $e.Index 0 }} = iota{{end}}{{ end }}
)

func (v {{ .EnumType }}) String() string {
  switch v { {{ range $e := .Elements }}
    case _{{ $e.Name }}: return "{{ $e.Type }}<{{ $e.Name | Cap }}>"{{end}}
    default: return fmt.Sprintf("{{ .EnumType }}<Unknown(%d)>", int(v))
  }
}

type {{ .EnumName }}Container struct { {{ range $e := .Elements }}
    {{ $e.Name | Cap }} {{ $e.Type }}{{ end }}
}

// uncertain if I want to return pointer or not
func (this *{{ .EnumName }}Container) ByValue(v int) (*{{ .EnumType }}, error) {
	value := new({{ .EnumType }})
  var err error = nil

  switch v { {{ range $e := .Elements }}
    case int(_{{ $e.Name }}): *value = _{{ $e.Name }}{{end}}
    default: err = errors.New(fmt.Sprintf("Unable to find {{ .EnumType }} associated with value %d", v))
  }

	return value, err
}

// uncertain if I want to return pointer or not
func (this *{{ .EnumName }}Container) ByName(s string) (*{{ .EnumType }}, error) {
  value := new({{ .EnumType }})
  var err error = nil

  switch s { {{ range $e := .Elements }}
    case "{{ $e.Name | Cap }}": *value = _{{ $e.Name }}{{end}}
    default: err = errors.New(fmt.Sprintf("Unable to find {{ .EnumType }} associated by name %s", s))
  }

  return value, err
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

// fin
`

func mkCap(s string) string {
	return fmt.Sprintf("%s%s", strings.ToUpper(string(s[0])), s[1:])
}

type ConstItem struct {
	Index        int
	ExportedName string
	Name         string
	Type         string
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
		ci := ConstItem{
			Index:        i,
			ExportedName: fmt.Sprintf("aoeu%s", v.Name),
			Name:         v.Name,
			Type:         v.Type,
		}
		ed.Elements = append(ed.Elements, ci)
	}

	ed.EnumType = ed.Elements[0].Type

	return ed
}

func GenerateEnum(
	source *types.File,
	enumName string,
	enumType string,
	dest io.Writer,
) {
	funcMap := template.FuncMap{
		"Cap": mkCap,
	}

	templateParsed, templateError := template.New("enum").Funcs(funcMap).Parse(enumTemplate)

	if templateError != nil {
		fmt.Printf(templateError.Error())
		return
	}

	if dest == nil {
		dest = os.Stdout
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
		templateParsed.Execute(dest, constBlockToEnumData(enumName, source, sourceConsts))
	} else {
		fmt.Println("Found sourceConsts:", sourceConsts)
	}
}
