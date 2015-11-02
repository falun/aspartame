package generators

import (
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/falun/aspartame/types"
)

var enumTemplate string = `package {{ .Package }}

{{ . }}

// fin
`

func GenerateEnum(source *types.File, dest io.Writer) {
	templateParsed, templateError := template.New("enum").Parse(enumTemplate)

	if templateError != nil {
		fmt.Printf(templateError.Error())
		return
	}

	if dest == nil {
		dest = os.Stdout
	}

	templateParsed.Execute(dest, source)
}
